package workflow

import (
	"fmt"

	"github.com/jczornik/glacier_backup/backup"
	"github.com/jczornik/glacier_backup/config"
	"github.com/jczornik/glacier_backup/persistent"
	"github.com/jczornik/glacier_backup/workflow/tasks"
)

func NewEncryptedBackup(conf config.BackupConfig, awsConf config.AWSConfig, db persistent.DBClient) (Workflow, error) {
	preserveManifest := tasks.NewPreserveTask(conf.Src, conf.Dst)
	artifacts := backup.NewArtifactNames(conf.Src, conf.Dst)
	encBackup := tasks.NewEncryptedBackupTask(conf.Src, artifacts, conf.Pass, conf.CanChange)
	upload := tasks.NewUploadToGlacierTask(artifacts.Archive, awsConf.AccountID, conf.Vault, awsConf.Profile)
	cleanup := tasks.NewCleanupTask(conf.Src, conf.Dst, artifacts, conf.Keep)

	tasks := []task{preserveManifest, encBackup, upload, cleanup}

	name := fmt.Sprintf("Backup for %s", conf.Src)
	return NewPWorkflow(name, tasks, db)
}
