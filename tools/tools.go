package tools

import (
	"errors"
	"io"
	"os/exec"
	"strings"
)

const (
	tar       = "tar"
	gpg       = "gpg"
	separator = " "
	required  = tar + separator + gpg
)

func checkIfToolExists(prog string) error {
	_, err := exec.LookPath(prog)
	return err
}

func CheckRequiredTools() error {
	for _, element := range strings.Split(required, separator) {
		if err := checkIfToolExists(element); err != nil {
			return err
		}
	}

	return nil
}

type CmdChain struct {
	cmds []*exec.Cmd
}

func (c CmdChain) Run() error {
	for i, cmd := range c.cmds {
		if err := cmd.Start(); err != nil {
			c.waitForN(i)
			return err
		}
	}

	return c.waitForN(len(c.cmds) - 1)
}

func (c CmdChain) waitForN(idx int) error {
	errs := make([]error, idx + 1)
	for i := 0; i <= idx; i++ {
		errs[i] = c.cmds[i].Wait()
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func (c CmdChain) SetStdout(out io.Writer) error {
	if len(c.cmds) == 0 {
		return errors.New("Cannot set output for empty chain")
	}
	c.cmds[0].Stdout = out
	return nil
}

func Pipe(cmd1 *exec.Cmd, cmd2 *exec.Cmd) (CmdChain, error) {
	var chain CmdChain
	var err error

	if cmd1 == nil || cmd2 == nil {
		return chain, errors.New("Cannot pipe if any command is nil")
	}

	cmd2.Stdin, err = cmd1.StdoutPipe()
	if err != nil {
		return chain, err
	}

	arr := []*exec.Cmd{cmd2, cmd1}
	return CmdChain{arr}, nil
}
