package tasks

import (
	"fmt"
	"log"
	"os"

	"github.com/jczornik/glacier_backup/backup"
)

const cleanupName = "Cleanup"

type CleanupTask struct {
	bckSrc           string
	bckDst           string
	artifacts        backup.Artifacts
	keepLocalArchive bool
}

func NewCleanupTask(bckSrc string, bckDst string, artifacts backup.Artifacts, rmLocalArchive bool) CleanupTask {
	return CleanupTask{bckSrc, bckDst, artifacts, rmLocalArchive}
}

func (t CleanupTask) cleanupLocalArchive() error {
	if t.keepLocalArchive {
		log.Printf("Keeping local archive %s", t.artifacts.Archive)
		return nil
	}

	if err := os.Remove(t.artifacts.Archive); err != nil {
		log.Printf("Fail cleanup for %s\n", t.artifacts.Archive)
		return err
	}

	log.Printf("Successfully cleanuped archive %s\n", t.artifacts.Archive)
	return nil
}

func (t CleanupTask) cleanupPreservedManifest() error {
	preserved := fmt.Sprintf("%s.%s", backup.MakeManifestName(t.bckSrc, t.bckDst), PreservedExt)

	if err := os.Remove(preserved); err != nil && !os.IsNotExist(err) {
		log.Printf("Fail cleanup for %s\n", t.artifacts.Archive)
		return err
	} else if err != nil && os.IsNotExist(err) {
		log.Printf("There is no preserved manifest for %s", t.bckSrc)
	} else {
		log.Printf("Successfully cleanuped preserved manifest %s\n", preserved)
	}

	return nil
}

func (t CleanupTask) Exec() error {
	log.Printf("Starting cleanup for %s\n", t.bckSrc)

	if err := t.cleanupLocalArchive(); err != nil {
		return err
	}

	if err := t.cleanupPreservedManifest(); err != nil {
		return err
	}

	log.Printf("Cleanup for %s was successful\n", t.bckSrc)

	return nil
}

func (t CleanupTask) Rollback() error {
	log.Printf("Starting rollback for cleanup for %s\n", t.artifacts.Archive)
	if _, err := os.Stat(t.artifacts.Archive); os.IsNotExist(err) {
		log.Printf("File %s already deleted - nothing we can do\n", t.artifacts.Archive)
	}

	log.Printf("Successfully rolbacked cleanup for %s\n", t.artifacts.Archive)
	return nil
}

func (t CleanupTask) Name() string {
	return cleanupName
}
