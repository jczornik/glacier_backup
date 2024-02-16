package workflow

import (
	"fmt"

	"github.com/jczornik/glacier_backup/backup"
	"github.com/jczornik/glacier_backup/config"
	"github.com/jczornik/glacier_backup/persistent"
	"github.com/jczornik/glacier_backup/workflow/tasks"
)

func NewEncryptedBackup(conf config.BackupConfig, awsConf config.AWSConfig, db persistent.DBClient) (Workflow, error) {
	name := fmt.Sprintf("Backup for %s", conf.Src)
	workflow, err := NewPWorkflow(name, db)
	if err != nil {
		return &BasicWorkflow{}, err
	}

	preserveManifest := tasks.NewPreserveTask(conf.Src, conf.Dst)
	artifacts := backup.NewArtifactNames(conf.Src, conf.Dst)
	encBackup := tasks.NewEncryptedBackupTask(conf.Src, artifacts, conf.Pass, conf.CanChange)
	upload := tasks.NewUploadToGlacierTask(artifacts.Archive, awsConf.AccountID, conf.Vault, awsConf.Profile)
	saveArtifacts := tasks.NewSaveArtifactsTask(artifacts.Snapshot, artifacts.Archive, db, workflow.Id())
	cleanup := tasks.NewCleanupTask(conf.Src, conf.Dst, artifacts, conf.Keep)

	tasks := []task{preserveManifest, encBackup, upload, saveArtifacts, cleanup}

	err = workflow.AddTasks(tasks)
	return workflow, err
}
