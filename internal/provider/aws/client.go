package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type Clients struct {
	EC2    *ec2.Client
	CW     *cloudwatch.Client
	ELBv2  *elasticloadbalancingv2.Client
	Config aws.Config
}

func NewClients(ctx context.Context, region, profile string) (*Clients, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region), config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}

	return &Clients{
		EC2:    ec2.NewFromConfig(cfg),
		CW:     cloudwatch.NewFromConfig(cfg),
		ELBv2:  elasticloadbalancingv2.NewFromConfig(cfg),
		Config: cfg,
	}, nil
}
