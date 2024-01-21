package main

import (
	"log"
	"os"

	"github.com/jczornik/glacier_backup/config"
	"github.com/jczornik/glacier_backup/persistent"
	"github.com/jczornik/glacier_backup/tools"
	"github.com/jczornik/glacier_backup/workflow"
)

const (
	ecRequiredTool = iota
	ecConfigPath
	ecReadingConf
	ecCreatingBackup
	ecUploadingBackup
	ecCreatingWorkflow
)

func main() {
	if err := tools.CheckRequiredTools(); err != nil {
		log.Printf("Required tool does not exist. Err: %s\n", err)
		os.Exit(ecRequiredTool)
	}

	if len(os.Args) < 2 {
		log.Println("Please specify path to configuration")
		os.Exit(ecConfigPath)
	}

	configPath := os.Args[1]

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Printf("Cannot read configuration. Err: %s\n", err)
		os.Exit(ecReadingConf)
	}

	db := persistent.NewDBClient(cfg.Db)
	if err := persistent.CheckAndUpdateSchema(db); err != nil {
		log.Println(err)
	}

	var workflows = make([]workflow.Workflow, len(cfg.Backups))

	for i, cbackup := range cfg.Backups {
		workflows[i], err = workflow.NewEncryptedBackup(cbackup.Src, cbackup.Dst, cbackup.Pass, cfg.AWS.AccountID, cbackup.Vault, cfg.AWS.Profile, cbackup.Keep, db)
		if err != nil {
			log.Println(err)
			os.Exit(ecCreatingWorkflow)
		}
	}

	for _, w := range workflows {
		if err := w.Exec(); err != nil {
			log.Fatal(err)
		} else {
			log.Println("Workflow OK!")
		}
	}
}
