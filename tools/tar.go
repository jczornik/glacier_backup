package tools

import (
	"os/exec"

	"github.com/jczornik/glacier_backup/config"
)

func NewBackupCmd(src config.BackupSrc, snapshotFile string) *exec.Cmd {
	return exec.Command(tar, "-czg", snapshotFile, src)
}
