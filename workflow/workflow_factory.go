package workflow

import (
	"github.com/jczornik/glacier_backup/backup"
	"github.com/jczornik/glacier_backup/persistent"
	"github.com/jczornik/glacier_backup/workflow/tasks"
)

func NewEncryptedBackup(src string, dst string, pass string, accountId string, vault string, profile string, rmLocalCopy bool, client persistent.DBClient) (Workflow, error) {
	preserveManifest := tasks.NewPreserveTask(src, dst)
	artifacts := backup.NewArtifactNames(src, dst)
	encBackup := tasks.NewEncryptedBackupTask(src, artifacts, pass)
	upload := tasks.NewUploadToGlacierTask(artifacts.Archive, accountId, vault, profile)
	cleanup := tasks.NewCleanupTask(src, dst, artifacts, rmLocalCopy)

	tasks := []task{preserveManifest, encBackup, upload, cleanup}

	return NewWorkflow(tasks, client)
}
