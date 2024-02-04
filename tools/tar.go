package tools

import (
	"os/exec"

	"github.com/jczornik/glacier_backup/config"
)

const fileChangedExitCode = 1

func NewBackupCmd(src config.BackupSrc, snapshotFile string, ignoreFileChange bool) Cmd {
	var exitCodes []int

	if ignoreFileChange {
		exitCodes = []int{fileChangedExitCode}
	}

	cmd := exec.Command(tar, "-czg", snapshotFile, src)
	return NewCmd(cmd, exitCodes)
}
