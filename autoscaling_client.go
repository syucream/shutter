package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

var (
	ErrInvalidAPIResponse = errors.New("AWS API response is invalid")
	ErrAlreadyTerminated  = errors.New("The instance has already terminated")

	// TODO make it configurable
	lifecycleActionResult = "ABANDON"
	region                = "us-east-1"
)

type autoscalingClient struct {
	client *autoscaling.AutoScaling
}

func NewAutoscalingClient() (*autoscalingClient, error) {
	creds := credentials.NewEnvCredentials()
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      &region,
	})
	if err != nil {
		return nil, err
	}

	return &autoscalingClient{
		client: autoscaling.New(sess),
	}, nil
}

func (c *autoscalingClient) DescribeAutoscalingGroup(name string) (*autoscaling.Group, error) {
	res, err := c.client.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
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

func (c *autoscalingClient) DescribeInstance(instanceId string) (*autoscaling.InstanceDetails, error) {
	res, err := c.client.DescribeAutoScalingInstances(&autoscaling.DescribeAutoScalingInstancesInput{
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

func (c *autoscalingClient) CompleteLifecycleAction(instanceId string, lifecycleHook string) error {
	details, err := c.DescribeInstance(instanceId)
	if err != nil {
		return err
	}

	if *details.LifecycleState == "Terminated" {
		return ErrAlreadyTerminated
	}

	_, err = c.client.CompleteLifecycleAction(&autoscaling.CompleteLifecycleActionInput{
		AutoScalingGroupName:  details.AutoScalingGroupName,
		InstanceId:            &instanceId,
		LifecycleActionResult: &lifecycleActionResult,
		LifecycleHookName:     &lifecycleHook,
	})

	return err
}
