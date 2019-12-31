package shutter

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type fakeAwsClient struct {
	AwsClient

	ResDescribeAutoscalingGroup *autoscaling.Group
	ErrDescribeAutoscalingGroup error

	ResDescribeInstance *ec2.Instance
	ErrDescribeInstance error

	ResDescribeInstanceDetails *autoscaling.InstanceDetails
	ErrDescribeInstanceDetails error

	ErrCompleteLifecycleAction error
}

func (c *fakeAwsClient) DescribeAutoscalingGroup(name string) (*autoscaling.Group, error) {
	return c.ResDescribeAutoscalingGroup, c.ErrDescribeAutoscalingGroup
}

func (c *fakeAwsClient) DescribeInstance(instanceId string) (*ec2.Instance, error) {
	return c.ResDescribeInstance, c.ErrDescribeInstance
}

func (c *fakeAwsClient) DescribeInstanceDetails(instanceId string) (*autoscaling.InstanceDetails, error) {
	return c.ResDescribeInstanceDetails, c.ErrDescribeInstanceDetails
}

func (c *fakeAwsClient) CompleteLifecycleAction(instanceId string, lifecycleActionResult, lifecycleHook string) error {
	return c.ErrCompleteLifecycleAction
}
