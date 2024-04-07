package tests

import (
	"testing"

	"github.com/isavita/codeexec/internal/executor"
)

func TestDockerExecutor(t *testing.T) {
	// Test case 1: Successful execution
	t.Run("SuccessfulExecution", func(t *testing.T) {
		code := "print('Hello, World!')"
		language := "python"
		exec, err := executor.NewDockerExecutor()
		if err != nil {
			t.Fatalf("Failed to create Docker executor: %v", err)
		}
		result, err := exec.Execute(code, language)
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
		_, err = exec.Execute(code, language)
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
		result, err := exec.Execute(code, language)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expected := "5"
		if result != expected {
			t.Errorf("Expected output: %s, but got: %s", expected, result)
		}
	})

	// Add more test cases for different scenarios
}
