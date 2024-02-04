package tools

import (
	"os/exec"
	"strings"
	"testing"
)

func TestPipe(t *testing.T) {
	// Given
	expected := "Hello!"
	cmd1 := Cmd{exec.Command("echo", "-n", expected), nil}
	cmd2 := Cmd{exec.Command("cat"), nil}
	out := new(strings.Builder)
	cmd2.cmd.Stdout = out

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
	cmd1 := Cmd{exec.Command("nonExistingCmd", expected), nil}
	cmd2 := Cmd{exec.Command("cat"), nil}

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
	cmd1 := Cmd{exec.Command("echo", "-n", expected), nil}
	cmd2 := Cmd{exec.Command("cat"), nil}
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
