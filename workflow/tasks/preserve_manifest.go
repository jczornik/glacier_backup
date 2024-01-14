package tasks

import (
	"fmt"
	"log"
	"os"

	"github.com/jczornik/glacier_backup/backup"
)

type PreserveTask struct {
	backupSrc string
	backupDst string
	original  *string
	preserved *string
}

func NewPreserveTask(backupSrc string, backupDst string) *PreserveTask {
	return &PreserveTask{backupSrc, backupDst, nil, nil}
}

func moveManifest(src string, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}

	return os.Rename(src, dst)
}

func (t *PreserveTask) Exec() error {
	log.Printf("Starting presering old manifest for backup %s\n", t.backupSrc)

	manifest, err := backup.GetManifestForSrc(t.backupSrc, t.backupDst)

	if err != nil {
		log.Printf("Error while searching for manifest for backup %s\n", t.backupSrc)
	}

	if manifest == nil {
		log.Printf("Manifest file for %s not found", t.backupSrc)
		return nil
	}

	t.original = manifest
	manifestDst := fmt.Sprintf("%s.old", *manifest)
	err = moveManifest(*manifest, manifestDst)

	if err != nil {
		log.Printf("Error while preserving manifest %s", *manifest)
	} else {
		log.Printf("Successfully preserved manifest %s.\nSaved to %s\n", *manifest, manifestDst)
		t.preserved = &manifestDst
	}

	return err
}

func (t *PreserveTask) Rollback() error {
	log.Printf("Rollback preserve for %s\n", t.backupSrc)

	if t.original == nil || t.preserved == nil {
		log.Printf("There is nothing to revert for preserve task for %s\n", t.backupSrc)
		return nil
	}

	err := moveManifest(*t.preserved, *t.original)

	if err != nil {
		log.Printf("Error while rollbacking preserve %s\n", *t.preserved)
	} else {
		log.Printf("Successfully rollbacked preserve %s\n", *t.preserved)
	}

	return err
}
