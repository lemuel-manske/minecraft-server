package ec2

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var _ Client = (*EC2Impl)(nil)

type EC2Impl struct {
	svc *awsec2.Client
}

func New(region string) (*EC2Impl, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	return &EC2Impl{svc: awsec2.NewFromConfig(cfg)}, nil
}

func (c *EC2Impl) GetInstance(ctx context.Context, name string) (*Instance, error) {
	out, err := c.svc.DescribeInstances(ctx, &awsec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{name},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"pending", "running", "stopping", "stopped"},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("describe instances: %w", err)
	}

	for _, r := range out.Reservations {
		for _, i := range r.Instances {
			inst := &Instance{
				ID:       aws.ToString(i.InstanceId),
				PublicIP: aws.ToString(i.PublicIpAddress),
				State:    string(i.State.Name),
			}
			for _, t := range i.Tags {
				if aws.ToString(t.Key) == "Profile" {
					inst.Profile = aws.ToString(t.Value)
				}
			}
			return inst, nil
		}
	}

	return nil, fmt.Errorf("instance %q not found", name)
}

func (c *EC2Impl) StartInstance(ctx context.Context, id string) error {
	_, err := c.svc.StartInstances(ctx, &awsec2.StartInstancesInput{
		InstanceIds: []string{id},
	})
	return err
}

func (c *EC2Impl) StopInstance(ctx context.Context, id string) error {
	_, err := c.svc.StopInstances(ctx, &awsec2.StopInstancesInput{
		InstanceIds: []string{id},
	})
	return err
}

func (c *EC2Impl) WaitRunning(ctx context.Context, id string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	for {
		out, err := c.svc.DescribeInstances(ctx, &awsec2.DescribeInstancesInput{
			InstanceIds: []string{id},
		})
		if err != nil {
			return "", fmt.Errorf("describe instance: %w", err)
		}

		for _, r := range out.Reservations {
			for _, i := range r.Instances {
				if i.State.Name == types.InstanceStateNameRunning && i.PublicIpAddress != nil {
					return aws.ToString(i.PublicIpAddress), nil
				}
			}
		}

		select {
		case <-ctx.Done():
			return "", fmt.Errorf("instance %s did not reach running state in time", id)
		case <-time.After(5 * time.Second):
		}
	}
}
