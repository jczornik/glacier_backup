package main

import (
	"log"
	"os"

	// "github.com/jczornik/glacier_backup/aws"
	"github.com/jczornik/glacier_backup/backup"
	"github.com/jczornik/glacier_backup/config"
	"github.com/jczornik/glacier_backup/tools"
)

const (
	ecRequiredTool = iota
	ecConfigPath
	ecReadingConf
	ecCreatingBackup
	ecUploadingBackup
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

	for _, cbackup := range cfg.Backups {
		if _, err := backup.CreateEncryptedBackup(cbackup.Src, cbackup.Dst, "1234"); err != nil {
			log.Println(err)
			os.Exit(ecCreatingBackup)
		}
	}

	// if err := aws.UploadData(cfg); err != nil {
	// 	log.Printf("Error while uploading backup, Err: %s\n", err)
	// 	os.Exit(ecUploadingBackup)
	// }
}
