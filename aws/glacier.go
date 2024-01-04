package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/glacier"
	"github.com/jczornik/glacier_backup/config"
)

func UploadData(cfg *config.Config) error {
	awscfg, err := NewConfig(cfg)
	if err != nil {
		return err
	}

	client := glacier.NewFromConfig(awscfg)
	vaults, err := client.ListVaults(context.TODO(), &glacier.ListVaultsInput{AccountId: &cfg.AWS.AccountID})
	if err != nil {
		return err
	}

	fmt.Println("Vaults:")
	for _, element := range vaults.VaultList {
		fmt.Println(*element.VaultName)
	}

	return nil
}
