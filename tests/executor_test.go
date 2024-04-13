package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/isavita/codeexec/internal/executor"
)

const timeout = 200 * time.Millisecond

func TestDockerExecutor(t *testing.T) {
	// Test case 1: Successful execution
	t.Run("SuccessfulExecution", func(t *testing.T) {
		code := "print('Hello, World!')"
		language := "python"
		exec, err := executor.NewDockerExecutor()
		if err != nil {
			t.Fatalf("Failed to create Docker executor: %v", err)
		}

		result, err := exec.Execute(code, language, timeout)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expected := "Hello, World!"
		if result != expected {
			t.Errorf("Expected output: %s, but got: %s", expected, result)
		}
	})

	// Test case 2: Syntax error
	t.Run("SyntaxError", func(t *testing.T) {
		code := "print('Hello, World!"
		language := "python"
		exec, err := executor.NewDockerExecutor()
		if err != nil {
			t.Fatalf("Failed to create Docker executor: %v", err)
		}
		_, err = exec.Execute(code, language, timeout)
		if err == nil {
			t.Error("Expected a syntax error, but got nil")
		}
	})

	// Test case 3: Print and return value
	t.Run("PrintAndReturnValue", func(t *testing.T) {
		code := `
def add(a, b):
	return a + b

a = add(2, 3)
print(a)
a
`
		language := "python"
		exec, err := executor.NewDockerExecutor()
		if err != nil {
			t.Fatalf("Failed to create Docker executor: %v", err)
		}
		result, err := exec.Execute(code, language, timeout)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expected := "5"
		if result != expected {
			t.Errorf("Expected output: %s, but got: %s", expected, result)
		}
	})

	// Test case 4: Execution timeout
	t.Run("ExecutionTimeout", func(t *testing.T) {
		code := `
import time

print("Starting execution")
time.sleep(5)
print("Execution finished")
`
		language := "python"
		exec, err := executor.NewDockerExecutor()
		if err != nil {
			t.Fatalf("Failed to create Docker executor: %v", err)
		}
		_, err = exec.Execute(code, language, timeout)
		if err == nil {
			t.Error("Expected an execution timeout error, but got nil")
		} else if !strings.Contains(err.Error(), "container execution timed out") {
			t.Errorf("Expected a timeout error, but got: %v", err)
		}
	})
	// Test case 5: Exceed memory limit
	t.Run("ExceedMemoryLimit", func(t *testing.T) {
		code := `
import numpy as np

print("Allocating large array")
large_array = np.zeros((10240000, 1024000, 1024000), dtype=np.int8)
print("Allocation finished")
`
		language := "python"
		exec, err := executor.NewDockerExecutor()
		if err != nil {
			t.Fatalf("Failed to create Docker executor: %v", err)
		}
		_, err = exec.Execute(code, language, 10*timeout)
		if err == nil {
			t.Error("Expected a resource limitation error, but got nil")
			// TODO: this check should be "container exceeded memory limit" instead of "container exited with non-zero status code: 1"
		} else if !strings.Contains(err.Error(), "container exited with non-zero status code: 1") {
			t.Errorf("Expected a resource limitation error, but got: %v", err)
		}
	})

	// Test case for JavaScript execution
	t.Run("JavaScriptExecution", func(t *testing.T) {
		code := "console.log('Hello, JavaScript!');"
		language := "javascript"
		exec, err := executor.NewDockerExecutor()
		if err != nil {
			t.Fatalf("Failed to create Docker executor: %v", err)
		}

		result, err := exec.Execute(code, language, timeout)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expected := "Hello, JavaScript!"
		if result != expected {
			t.Errorf("Expected output: %s, but got: %s", expected, result)
		}
	})

	// Add more test cases for different scenarios
}

func TestJavaScriptSyntaxCheck(t *testing.T) {
	// Test case: JavaScript syntax error
	t.Run("JavaScriptSyntaxError", func(t *testing.T) {
		// Deliberately broken JavaScript code
		code := "console.log('Hello, JavaScript!';" // Missing closing parenthesis
		language := "javascript"
		exec, err := executor.NewDockerExecutor()
		if err != nil {
			t.Fatalf("Failed to create Docker executor: %v", err)
		}

		_, err = exec.Execute(code, language, timeout)
		if err == nil {
			t.Error("Expected a syntax error, but got nil")
		} else if !strings.Contains(err.Error(), "syntax check failed") {
			t.Errorf("Expected a syntax error, but got: %v", err)
		}
	})
}
