package tools

import (
	"os/exec"
	"strings"
)

type TarOption = string

const (
	CreateNewArchive TarOption = "c"
	Gzip TarOption = "z"
	NewGnuIncremental = "g"
)

func RunTar(options []TarOption) error {
	optStr := "-" + strings.Join(options, "")
	return exec.Command(Tar, optStr).Run()
}
