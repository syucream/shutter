package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	ErrInvalidAPIResponse = errors.New("AWS API response is invalid")
	ErrAlreadyTerminated  = errors.New("The instance has already terminated")
)

type awsClient struct {
	ec2Client *ec2.EC2
	asClient  *autoscaling.AutoScaling
}

func NewAwsClient(config *Config) (*awsClient, error) {
	creds := credentials.NewEnvCredentials()
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      &config.AwsRegion,
	})
	if err != nil {
		return nil, err
	}

	return &awsClient{
		ec2Client: ec2.New(sess),
		asClient:  autoscaling.New(sess),
	}, nil
}

func (c *awsClient) DescribeAutoscalingGroup(name string) (*autoscaling.Group, error) {
	res, err := c.asClient.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{&name},
	})
	if err != nil {
		return nil, err
	}

	if len(res.AutoScalingGroups) != 1 {
		return nil, ErrInvalidAPIResponse
	}

	return res.AutoScalingGroups[0], nil
}

func (c *awsClient) DescribeInstance(instanceId string) (*ec2.Instance, error) {
	name := "instance-id"

	res, err := c.ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   &name,
				Values: []*string{&instanceId},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(res.Reservations) != 1 || len(res.Reservations[0].Instances) != 1 {
		return nil, ErrInvalidAPIResponse
	}

	return res.Reservations[0].Instances[0], nil
}

func (c *awsClient) DescribeInstanceDetails(instanceId string) (*autoscaling.InstanceDetails, error) {
	res, err := c.asClient.DescribeAutoScalingInstances(&autoscaling.DescribeAutoScalingInstancesInput{
		InstanceIds: []*string{&instanceId},
	})
	if err != nil {
		return nil, err
	}

	if len(res.AutoScalingInstances) != 1 {
		return nil, ErrInvalidAPIResponse
	}

	return res.AutoScalingInstances[0], nil
}

func (c *awsClient) CompleteLifecycleAction(instanceId string, lifecycleActionResult, lifecycleHook string) error {
	details, err := c.DescribeInstanceDetails(instanceId)
	if err != nil {
		return err
	}

	if *details.LifecycleState == "Terminated" {
		return ErrAlreadyTerminated
	}

	_, err = c.asClient.CompleteLifecycleAction(&autoscaling.CompleteLifecycleActionInput{
		AutoScalingGroupName:  details.AutoScalingGroupName,
		InstanceId:            &instanceId,
		LifecycleActionResult: &lifecycleActionResult,
		LifecycleHookName:     &lifecycleHook,
	})

	return err
}
