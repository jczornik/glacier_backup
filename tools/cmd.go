package tools

import (
	"errors"
	"io"
	"os/exec"
)

type Cmd struct {
	cmd              *exec.Cmd
	successExitCodes []int
}

func NewCmd(cmd *exec.Cmd, allowedExitCodes []int) Cmd {
	return Cmd{cmd, allowedExitCodes}
}

func (c Cmd) Start() error {
	return c.cmd.Start()
}

func (c Cmd) Wait() error {
	err := c.cmd.Wait()
	if c.isSuccessExitCode(err) {
		return nil
	} else {
		return err
	}
}

func (c Cmd) isSuccessExitCode(err error) bool {
	if err == nil {
		return true
	}

	exitError, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}

	exitCode := exitError.ExitCode()

	if exitCode == 0 {
		return true
	}

	if c.successExitCodes == nil {
		return false
	}

	for _, code := range c.successExitCodes {
		if exitCode == code {
			return true
		}
	}

	return false
}

type CmdChain struct {
	cmds []Cmd
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
	errs := make([]error, idx+1)
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
	c.cmds[0].cmd.Stdout = out
	return nil
}

func (c CmdChain) SetStderr(out io.Writer) error {
	if len(c.cmds) == 0 {
		return errors.New("Cannot set output for empty chain")
	}
	c.cmds[len(c.cmds)-1].cmd.Stderr = out
	return nil
}

func Pipe(cmd1 Cmd, cmd2 Cmd) (CmdChain, error) {
	var chain CmdChain
	var err error

	if cmd1.cmd == nil || cmd2.cmd == nil {
		return chain, errors.New("Cannot pipe if any command is nil")
	}

	cmd2.cmd.Stdin, err = cmd1.cmd.StdoutPipe()
	if err != nil {
		return chain, err
	}

	arr := []Cmd{cmd2, cmd1}
	return CmdChain{arr}, nil
}
