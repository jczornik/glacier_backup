package tasks

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/jczornik/glacier_backup/aws"
)

const uploadName = "UploadToGlacier"

type UploadToGlacierTask struct {
	archive string
	account string
	vault   string
	profile string
}

func NewUploadToGlacierTask(archive string, account string, vault string, profile string) UploadToGlacierTask {
	return UploadToGlacierTask{archive, account, vault, profile}
}

func (t UploadToGlacierTask) Exec() error {
	log.Printf("Startting uploading archive %s to vault %s\n", t.archive, t.vault)
	cfg, err := aws.NewConfig(t.profile)
	if err != nil {
		return err
	}

	vaults, err := aws.GetVaults(cfg, t.account)
	if err != nil {
		return err
	}

	if !slices.Contains(vaults, t.vault) {
		return errors.New(fmt.Sprintf("Provided vault name %s is not available", t.vault))
	}

	if err := aws.UploadData(cfg, t.account, t.vault, t.archive); err != nil {
		log.Printf("Error while uploading archive %s to glacier\n", t.archive)
		return err
	}

	log.Printf("Successfully uploaded archive %s to vault %s\n", t.archive, t.vault)
	return nil
}

func (t UploadToGlacierTask) Rollback() error {
	return nil
}

func (t UploadToGlacierTask) Name() string {
	return uploadName
}
