package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	glaciercfg "github.com/jczornik/glacier_backup/config"
)

func NewConfig(cfg *glaciercfg.Config) (aws.Config, error) {
	profile := cfg.AWS.Profile
	return awscfg.LoadDefaultConfig(context.TODO(), awscfg.WithSharedConfigProfile(profile))
}
