package aws

import (
	"AutoScaler/internal/models"
	"AutoScaler/internal/provider"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Provider struct {
	clients *Clients
}

func NewEC2Provider(clients *Clients) *EC2Provider {
	return &EC2Provider{
		clients: clients,
	}
}

func (p *EC2Provider) Launch(ctx context.Context, req provider.LaunchRequest) (models.Instance, error) {
	tags := []types.Tag{
		{
			Key:   aws.String("asg:group-name"),
			Value: aws.String(req.GroupName),
		},
		{
			Key:   aws.String("asg:managed-by"),
			Value: aws.String("autoscaler"),
		},
	}

	for key, value := range req.LaunchSpec.Tags {
		tags = append(tags, types.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}

	securityGroupIds := make([]string, len(req.LaunchSpec.SecurityGroups))

	copy(securityGroupIds, req.LaunchSpec.SecurityGroups)

	subnetIds := []string{}
	for _, subnetID := range req.LaunchSpec.SubnetIDs {
		subnetIds = append(subnetIds, subnetID)
	}

	input := &ec2.RunInstancesInput{
		ImageId:          aws.String(req.LaunchSpec.ImageID),
		InstanceType:     types.InstanceType(req.LaunchSpec.InstanceType),
		MinCount:         aws.Int32(1),
		MaxCount:         aws.Int32(1),
		KeyName:          aws.String(req.LaunchSpec.KeyName),
		SecurityGroupIds: securityGroupIds,
		SubnetId:         aws.String(subnetIds[0]),
		UserData:         aws.String(req.LaunchSpec.UserData),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags:         tags,
			},
		},
		ClientToken: aws.String(req.ClientToken),
	}

	if req.LaunchSpec.IAMProfile != "" {
		input.IamInstanceProfile = &types.IamInstanceProfileSpecification{
			Name: aws.String(req.LaunchSpec.IAMProfile),
		}
	}

	result, err := p.clients.EC2.RunInstances(ctx, input)
	if err != nil {
		return models.Instance{}, fmt.Errorf("failed to run instances: %w", err)
	}

	if len(result.Instances) == 0 {
		return models.Instance{}, fmt.Errorf("no instances returned")
	}

	instance := result.Instances[0]
	return models.Instance{
		ID:        aws.ToString(instance.InstanceId),
		GroupName: req.GroupName,
		Zone:      aws.ToString(instance.Placement.AvailabilityZone),
		State:     MapEC2State(string(instance.State.Name)),
		Health:    models.HealthStatus("unknown"),
		LaunchAt:  *instance.LaunchTime,
	}, nil
}

func (p *EC2Provider) Terminate(ctx context.Context, ids []string) error {
	instanceIds := make([]string, len(ids))
	copy(instanceIds, ids)

	input := &ec2.TerminateInstancesInput{
		InstanceIds: instanceIds,
	}

	_, err := p.clients.EC2.TerminateInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to terminate instances: %w", err)
	}

	return nil
}

func (p *EC2Provider) List(ctx context.Context, groupName string) ([]models.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:asg:group-name"),
				Values: []string{groupName},
			},
		},
	}

	var instances []models.Instance
	paginator := ec2.NewDescribeInstancesPaginator(p.clients.EC2, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to describe instances: %w", err)
		}

		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				instances = append(instances, models.Instance{
					ID:        aws.ToString(instance.InstanceId),
					GroupName: groupName,
					Zone:      aws.ToString(instance.Placement.AvailabilityZone),
					State:     MapEC2State(string(instance.State.Name)),
					Health:    models.HealthStatus("unknown"),
					LaunchAt:  *instance.LaunchTime,
				})
			}
		}
	}

	return instances, nil
}

func (p *EC2Provider) Zones(ctx context.Context) ([]string, error) {
	input := &ec2.DescribeAvailabilityZonesInput{}

	result, err := p.clients.EC2.DescribeAvailabilityZones(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe availability zones: %w", err)
	}

	var zones []string
	for _, zone := range result.AvailabilityZones {
		zones = append(zones, aws.ToString(zone.ZoneName))
	}

	return zones, nil
}
