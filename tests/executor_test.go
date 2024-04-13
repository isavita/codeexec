package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/isavita/codeexec/internal/executor"
)

const (
	timeout         = 200 * time.Millisecond
	extendedTimeout = 10 * timeout
)

func TestDockerExecutor(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		language string
		expected string
		err      string
		timeout  time.Duration
	}{
		{
			name:     "SuccessfulExecution",
			code:     "print('Hello, World!')",
			language: "python",
			expected: "Hello, World!",
			timeout:  timeout,
		},
		{
			name:     "SyntaxError",
			code:     "print('Hello, World!",
			language: "python",
			err:      "syntax check failed",
			timeout:  timeout,
		},
		{
			name: "PrintAndReturnValue",
			code: `
def add(a, b):
	return a + b

a = add(2, 3)
print(a)
a
`,
			language: "python",
			expected: "5",
			timeout:  timeout,
		},
		{
			name: "ExecutionTimeout",
			code: `
import time

print("Starting execution")
time.sleep(5)
print("Execution finished")
`,
			language: "python",
			err:      "container execution timed out",
			timeout:  timeout,
		},
		{
			name: "ExceedMemoryLimit",
			code: `
import numpy as np

print("Allocating large array")
large_array = np.zeros((10240000, 1024000, 1024000), dtype=np.int8)
print("Allocation finished")
`,
			language: "python",
			// TODO: return "container exceeded memory limit",
			err:     "container exited with non-zero status code: 1",
			timeout: extendedTimeout,
		},
		{
			name:     "JavaScriptExecution",
			code:     "console.log('Hello, JavaScript!');",
			language: "javascript",
			expected: "Hello, JavaScript!",
			timeout:  timeout,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exec, err := executor.NewDockerExecutor()
			if err != nil {
				t.Fatalf("Failed to create Docker executor: %v", err)
			}

			result, err := exec.Execute(tc.code, tc.language, tc.timeout)
			if tc.err != "" {
				if err == nil {
					t.Errorf("Expected an error containing '%s', but got nil", tc.err)
				} else if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("Expected an error containing '%s', but got: %v", tc.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tc.expected {
					t.Errorf("Expected output: %s, but got: %s", tc.expected, result)
				}
			}
		})
	}
}

func TestJavaScriptSyntaxCheck(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		language string
		err      string
	}{
		{
			name:     "JavaScriptSyntaxError",
			code:     "console.log('Hello, JavaScript!';",
			language: "javascript",
			err:      "syntax check failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exec, err := executor.NewDockerExecutor()
			if err != nil {
				t.Fatalf("Failed to create Docker executor: %v", err)
			}

			_, err = exec.Execute(tc.code, tc.language, timeout)
			if err == nil {
				t.Errorf("Expected an error containing '%s', but got nil", tc.err)
			} else if !strings.Contains(err.Error(), tc.err) {
				t.Errorf("Expected an error containing '%s', but got: %v", tc.err, err)
			}
		})
	}
}
