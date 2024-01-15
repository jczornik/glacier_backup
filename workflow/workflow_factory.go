package workflow

import (
	"github.com/jczornik/glacier_backup/backup"
	"github.com/jczornik/glacier_backup/workflow/tasks"
)

func NewEncryptedBackup(src string, dst string, pass string, rmLocalCopy bool) Workflow {
	preserveManifest := tasks.NewPreserveTask(src, dst)

	artifacts := backup.NewArtifactNames(src, dst)
	encBackup := tasks.NewEncryptedBackupTask(src, artifacts, pass)

	cleanup := tasks.NewCleanupTask(src, dst, artifacts, rmLocalCopy)

	return Workflow{[]task{preserveManifest, encBackup, cleanup}}

}
