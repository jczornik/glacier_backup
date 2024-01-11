package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glacier"
)

func GetVaults(cfg aws.Config, account string) ([]string, error) {
	client := glacier.NewFromConfig(cfg)
	vaults, err := client.ListVaults(context.TODO(), &glacier.ListVaultsInput{AccountId: &account})
	if err != nil {
		return nil, err
	}

	vaultNames := make([]string, len(vaults.VaultList))
	for i, element := range vaults.VaultList {
		vaultNames[i] = *element.VaultName
	}

	return vaultNames, nil
}

func UploadData(cfg aws.Config, account string) error {
	glacier.NewFromConfig(cfg)
	return nil
}
