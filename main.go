package main

import (
	"fmt"
	"os"

	"github.com/jczornik/glacier_backup/config"
	"github.com/jczornik/glacier_backup/incrbackup"
	"github.com/jczornik/glacier_backup/tools"
)

func main() {
	if err := tools.CheckRequiredTools(); err != nil {
		fmt.Printf("Required prog does not exist. Err: %s\n", err)
		os.Exit(2)
	}

	if len(os.Args) < 2 {
		fmt.Println("Please specify path to configuration")
		os.Exit(1)
	}

	configPath := os.Args[1]

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Printf("Cannot read configuration. Err: %s\n", err)
		os.Exit(3)
	}

	for src, dst := range cfg.Local.Paths {
		if err := incrbackup.CreateBackup(src, dst); err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
	}
}
