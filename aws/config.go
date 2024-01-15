package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func NewConfig(profile string) (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
}
