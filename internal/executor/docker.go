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
	// Perform syntax check
	if err := e.syntaxCheck(code, language); err != nil {
		return "", fmt.Errorf("syntax check failed: %v", err)
	}

	// Create a new Docker container
	ctx := context.Background()
	resp, err := e.client.ContainerCreate(ctx, &container.Config{
		Image: getImageForLanguage(language),
		Cmd:   []string{"code." + getFileExtensionForLanguage(language)},
	}, &container.HostConfig{
		Resources: container.Resources{
			Memory:     64 * 1024 * 1024, // 64MB
			MemorySwap: 64 * 1024 * 1024, // 64MB (same as Memory to disable swap)
			CPUQuota:   50000,            // 50% CPU
		},
		Binds: []string{
			fmt.Sprintf("%s:/app/code.%s", createCodeFile(code, language), getFileExtensionForLanguage(language)),
		},
	}, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %v", err)
	}
	defer e.client.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})

	// Start the container
	if err := e.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %v", err)
	}

	// Wait for the container to finish or timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	statusCh, errCh := e.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	var exitCode int64
	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("failed to wait for container: %v", err)
		}
	case result := <-statusCh:
		exitCode = result.StatusCode
	case <-ctx.Done():
		return "", fmt.Errorf("container execution timed out after %s", timeout)
	}

	// Check the container exit code and status
	containerInfo, err := e.client.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %v", err)
	}

	if containerInfo.State.OOMKilled {
		return "", fmt.Errorf("container exceeded memory limit")
	} else if exitCode != 0 {
		return "", fmt.Errorf("container exited with non-zero status code: %d", exitCode)
	}

	// Retrieve the container logs
	out, err := e.client.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", fmt.Errorf("failed to retrieve container logs: %v", err)
	}
	defer out.Close()

	// Read the logs and return the output
	var stdout, stderr strings.Builder
	if _, err := stdcopy.StdCopy(&stdout, &stderr, out); err != nil {
		return "", fmt.Errorf("failed to read container logs: %v", err)
	}

	if stderr.Len() > 0 {
		return "", fmt.Errorf("execution error: %s", stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
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
	// Add more cases for other supported languages
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
	// Add more cases for other supported languages
	default:
		return ""
	}
}

func getRunCommandForLanguage(language, fileName string) string {
	// Return the appropriate run command for the given language and file name
	switch language {
	case "python":
		return fmt.Sprintf("python %s", fileName)
	// Add more cases for other supported languages
	default:
		return ""
	}
}

func getFileExtensionForLanguage(language string) string {
	switch language {
	case "python":
		return "py"
	// Add more cases for other supported languages
	default:
		return ""
	}
}
