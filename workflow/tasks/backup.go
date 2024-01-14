package tasks

import (
	"os"

	"github.com/jczornik/glacier_backup/backup"
)

type EncryptedBackupTask struct {
	src       string
	pass      string
	artifacts backup.Artifacts
}

func NewEncryptedBackupTask(src string, artifacts backup.Artifacts, pass string) EncryptedBackupTask {
	return EncryptedBackupTask{src, pass, artifacts}
}

func (t EncryptedBackupTask) Exec() error {
	return backup.CreateEncryptedBackup(t.src, t.artifacts, t.pass)
}

func (t EncryptedBackupTask) Rollback() error {
	if err := os.Remove(t.artifacts.Snapshot); !os.IsNotExist(err) {
		return err
	}

	if err := os.Remove(t.artifacts.Archive); !os.IsNotExist(err) {
		return err
	}

	return nil
}
