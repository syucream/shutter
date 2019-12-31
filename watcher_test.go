package shutter

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"go.uber.org/zap"
	"testing"
)

func TestWatcher_Watch(t *testing.T) {
	client := &fakeAwsClient{
		ResDescribeAutoscalingGroup: &autoscaling.Group{
			Instances: []*autoscaling.Instance{
				// Its not target
				{
					LifecycleState: &inServiceLifecycleState,
				},
				// Its a target
				{
					LifecycleState: &terminatingLifecycleState,
				},
			},
		},
	}

	config := &Config{
		Watcher: Watcher{
			AutoscalingGroupName: "test",
		},
	}

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	watcher := NewWatcher(client, config, logger)

	actual, err := watcher.Watch()
	if err != nil {
		t.Fatal(err)
	}

	if len(actual) != 1 {
		t.Fatal("watch result is invalid", actual)
	}
}
