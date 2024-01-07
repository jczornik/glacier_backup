package tools

import (
	"os/exec"
	"strings"
	"testing"
)

func TestPipe(t *testing.T) {
	// Given
	expected := "Hello!"
	cmd1 := exec.Command("echo", "-n", expected)
	cmd2 := exec.Command("cat")
	out := new(strings.Builder)
	cmd2.Stdout = out

	// When
	chain, _ := Pipe(cmd1, cmd2)
	chain.Run()

	// Then
	if out.String() != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, out.String())
	}
}

func TestPipeFirstNotExisting(t *testing.T) {
	// Given
	expected := "Hello!"
	cmd1 := exec.Command("nonExistingCmd", expected)
	cmd2 := exec.Command("cat")

	// When
	chain, _ := Pipe(cmd1, cmd2)
	err := chain.Run()

	// Then
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestPipeAndRedirectOut(t *testing.T) {
	// Given
	expected := "Hello!"
	cmd1 := exec.Command("echo", "-n", expected)
	cmd2 := exec.Command("cat")
	out := new(strings.Builder)

	// When
	chain, _ := Pipe(cmd1, cmd2)
	chain.SetStdout(out)
	chain.Run()

	// Then
	if out.String() != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, out.String())
	}
}
