package workflow

import (
	"github.com/jczornik/glacier_backup/backup"
	"github.com/jczornik/glacier_backup/workflow/tasks"
)

func NewEncryptedBackup(src string, dst string, pass string) (*Workflow, error) {
	prevManifest, err := backup.GetManifestForSrc(src, dst)
	if err != nil {
		return nil, err
	}

	artifacts := backup.NewArtifactNames(src, dst)
	encBackup := tasks.NewEncryptedBackupTask(src, artifacts, pass)

	if prevManifest != nil {
		preserveManifest := tasks.NewPreserveTask(*prevManifest)
		return &Workflow{[]task{preserveManifest, encBackup}}, nil
	}

	return &Workflow{[]task{encBackup}}, nil

}
