package tools

import (
	"os/exec"
	"strings"
	"testing"
	"io"
)

func catCtor(output io.Writer) CmdConstructor {
	ctor := func() *exec.Cmd {
		cat := exec.Command("cat")
		cat.Stdout = output
		return cat
	}

	return ctor
}

func ShouldSucceedWhenBufferEqMaxBuffer(t *testing.T) {
	// Given
	out1 := new(strings.Builder)
	catCtor := catCtor(out1)

	s := NewSplitter(1, 5, catCtor)
	input := []byte("Hello")

	// When
	wrote, err := s.Write([]byte(input))

	// Then
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if wrote != len(input) {
		t.Errorf("Expected to write %d bytes but wrote %d", len(input), wrote)
	}

	if out1.String() != string(input[:]) {
		t.Errorf("Expected '%s' but got '%s'", string(input[:]), out1.String())
	}
}

func ShouldSucceedWhenBufferSmallerThenMaxBuffer(t *testing.T) {
	// Given
	out1 := new(strings.Builder)
	catCtor := catCtor(out1)

	s := NewSplitter(1, 6, catCtor)
	input := []byte("Hello")

	// When
	wrote, err := s.Write([]byte(input))

	// Then
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if wrote != len(input) {
		t.Errorf("Expected to write %d bytes but wrote %d", len(input), wrote)
	}

	if out1.String() != string(input[:]) {
		t.Errorf("Expected '%s' but got '%s'", string(input[:]), out1.String())
	}
}
