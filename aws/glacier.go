package aws

import (
	"context"
	"os"

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

func UploadData(cfg aws.Config, account string, vault string, archive string) error {
	file, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer file.Close()

	client := glacier.NewFromConfig(cfg)
	input := glacier.UploadArchiveInput{
		AccountId: &account,
		VaultName: &vault,
		Body:      file,
	}

	_, err = client.UploadArchive(context.TODO(), &input)

	return err
}
