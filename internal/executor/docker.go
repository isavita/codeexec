package executor

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerExecutor struct {
	client *client.Client
}

func NewDockerExecutor() (*DockerExecutor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerExecutor{client: cli}, nil
}

func (e *DockerExecutor) Execute(code, language string, timeout time.Duration) (string, error) {
	if err := e.syntaxCheck(code, language); err != nil {
		return "", fmt.Errorf("syntax check failed: %v", err)
	}

	containerID, err := e.createContainer(code, language)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %v", err)
	}
	defer e.removeContainer(containerID)

	if err := e.startContainer(containerID); err != nil {
		return "", fmt.Errorf("failed to start container: %v", err)
	}

	exitCode, err := e.waitForContainer(containerID, timeout)
	if err != nil {
		return "", err
	}

	if err := e.checkContainerStatus(containerID, exitCode); err != nil {
		return "", err
	}

	output, err := e.getContainerOutput(containerID)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (e *DockerExecutor) createContainer(code, language string) (string, error) {
	ctx := context.Background()
	resp, err := e.client.ContainerCreate(ctx, &container.Config{
		Image: getImageForLanguage(language),
		Cmd:   []string{"code." + getFileExtensionForLanguage(language)},
	}, &container.HostConfig{
		Resources: container.Resources{
			Memory:     64 * 1024 * 1024,
			MemorySwap: 64 * 1024 * 1024,
			CPUQuota:   50000,
		},
		Binds: []string{
			fmt.Sprintf("%s:/app/code.%s", createCodeFile(code, language), getFileExtensionForLanguage(language)),
		},
	}, nil, nil, "")
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (e *DockerExecutor) startContainer(containerID string) error {
	ctx := context.Background()
	return e.client.ContainerStart(ctx, containerID, container.StartOptions{})
}

func (e *DockerExecutor) waitForContainer(containerID string, timeout time.Duration) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	statusCh, errCh := e.client.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return 0, fmt.Errorf("failed to wait for container: %v", err)
		}
	case result := <-statusCh:
		return result.StatusCode, nil
	case <-ctx.Done():
		return 0, fmt.Errorf("container execution timed out after %s", timeout)
	}
	return 0, nil
}

func (e *DockerExecutor) checkContainerStatus(containerID string, exitCode int64) error {
	ctx := context.Background()
	containerInfo, err := e.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return fmt.Errorf("failed to inspect container: %v", err)
	}

	if containerInfo.State.OOMKilled {
		return fmt.Errorf("container exceeded memory limit")
	} else if exitCode != 0 {
		return fmt.Errorf("container exited with non-zero status code: %d", exitCode)
	}
	return nil
}

func (e *DockerExecutor) getContainerOutput(containerID string) (string, error) {
	ctx := context.Background()
	out, err := e.client.ContainerLogs(ctx, containerID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", fmt.Errorf("failed to retrieve container logs: %v", err)
	}
	defer out.Close()

	var stdout, stderr strings.Builder
	if _, err := stdcopy.StdCopy(&stdout, &stderr, out); err != nil {
		return "", fmt.Errorf("failed to read container logs: %v", err)
	}

	if stderr.Len() > 0 {
		return "", fmt.Errorf("execution error: %s", stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func (e *DockerExecutor) removeContainer(containerID string) {
	ctx := context.Background()
	e.client.ContainerRemove(ctx, containerID, container.RemoveOptions{})
}

func (e *DockerExecutor) syntaxCheck(code, language string) error {
	// Implement syntax check logic based on the language
	switch language {
	case "python":
		// Use a lightweight Python container to check the syntax
		ctx := context.Background()
		resp, err := e.client.ContainerCreate(ctx, &container.Config{
			Image: "python:3.11-alpine",
			Cmd:   []string{"python", "-c", "import ast; ast.parse('''" + code + "''')"},
		}, nil, nil, nil, "")
		if err != nil {
			return err
		}
		defer e.client.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})

		if err := e.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
			return err
		}

		statusCh, errCh := e.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				return err
			}
		case status := <-statusCh:
			if status.StatusCode != 0 {
				return fmt.Errorf("syntax check failed")
			}
		}
	case "javascript":
		ctx := context.Background()
		codeFile := createCodeFile(code, language) // Ensure this function returns a valid path
		defer os.Remove(codeFile)                  // Clean up after checking

		// Make sure the container path is specified correctly
		containerPath := "/app/" + filepath.Base(codeFile) // Use a specific folder instead of root

		resp, err := e.client.ContainerCreate(ctx, &container.Config{
			Image: "node:19-alpine",
			Cmd:   []string{"node", "--check", containerPath},
		}, &container.HostConfig{
			Binds: []string{fmt.Sprintf("%s:%s", codeFile, containerPath)},
		}, nil, nil, "")
		if err != nil {
			return err
		}
		defer e.client.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})

		if err := e.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
			return err
		}

		statusCh, errCh := e.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				return err
			}
		case status := <-statusCh:
			if status.StatusCode != 0 {
				// Retrieve logs to see the syntax error
				out, logErr := e.client.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
				if logErr != nil {
					return fmt.Errorf("syntax check failed with status code %d and unable to retrieve logs: %v", status.StatusCode, logErr)
				}
				defer out.Close()
				var stderr strings.Builder
				if _, err := stdcopy.StdCopy(nil, &stderr, out); err != nil {
					return fmt.Errorf("failed to read container logs: %v", err)
				}
				return fmt.Errorf("syntax check failed: status code %d, error: %s", status.StatusCode, stderr.String())
			}
		}
	default:
		return fmt.Errorf("unsupported language: %s", language)
	}
	return nil
}

func createCodeFile(code, language string) string {
	tempDir, err := os.MkdirTemp("", "code")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}

	codeFile := filepath.Join(tempDir, "code."+getFileExtensionForLanguage(language))
	err = os.WriteFile(codeFile, []byte(code), 0644)
	if err != nil {
		log.Fatalf("Failed to write code to file: %v", err)
	}

	return codeFile
}

func getImageForLanguage(language string) string {
	switch language {
	case "python":
		return "python-exec"
	case "javascript":
		return "javascript-exec"
	default:
		return ""
	}
}

func getFileExtensionForLanguage(language string) string {
	switch language {
	case "python":
		return "py"
	case "javascript":
		return "js"
	default:
		return ""
	}
}
