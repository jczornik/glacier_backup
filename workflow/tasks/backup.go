package tasks

import (
	"log"
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
	log.Printf("Starting encrypted backup for %s\n", t.src)
	err := backup.CreateEncryptedBackup(t.src, t.artifacts, t.pass)

	if err != nil {
		log.Printf("Error while creating encrypted backup for %s.\n", t.src)
	} else {
		log.Printf("Successfully created encrypted backup for %s.\n Archive: %s\n", t.src, t.artifacts.Archive)
	}

	return err
}

func (t EncryptedBackupTask) Rollback() error {
	log.Printf("Rollback for creating backup for %s\n", t.src)
	if err := os.Remove(t.artifacts.Snapshot); !os.IsNotExist(err) {
		log.Printf("Error while removing snapshot %s\n", t.artifacts.Snapshot)
		return err
	}

	if err := os.Remove(t.artifacts.Archive); !os.IsNotExist(err) {
		log.Printf("Error while removing archive %s\n", t.artifacts.Archive)
		return err
	}

	return nil
}