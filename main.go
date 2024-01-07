package main

import (
	"log"
	"os"

	"github.com/jczornik/glacier_backup/backup"
	"github.com/jczornik/glacier_backup/config"
	"github.com/jczornik/glacier_backup/tools"
)

func main() {
	if err := tools.CheckRequiredTools(); err != nil {
		log.Printf("Required prog does not exist. Err: %s\n", err)
		os.Exit(2)
	}

	if len(os.Args) < 2 {
		log.Println("Please specify path to configuration")
		os.Exit(1)
	}

	configPath := os.Args[1]

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Printf("Cannot read configuration. Err: %s\n", err)
		os.Exit(3)
	}

	for _, cbackup := range cfg.Backups {
		if err := backup.CreateEncryptedBackup(cbackup.Src, cbackup.Dst, "1234"); err != nil {
			log.Println(err)
			os.Exit(4)
		}
	}

	//aws.UploadData(cfg)
}
