package shutter

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/ec2"
	"go.uber.org/zap"
	"testing"
)

func TestFinisher_terminateHandler(t *testing.T) {
	publicIpAddr := "8.8.8.8"
	pricateIpAddr := "192.168.0.1"

	validClient := &fakeAwsClient{
		ResDescribeInstance: &ec2.Instance{
			PublicIpAddress:  &publicIpAddr,
			PrivateIpAddress: &pricateIpAddr,
		},
	}
	invalidClient := &fakeAwsClient{
		ResDescribeInstance: nil,
		ErrDescribeInstance: errors.New("error"),
	}

	validConfig := &Config{
		Finisher: Finisher{
			Terminate: Terminate{
				Command: "echo test",
			},
		},
	}
	invalidConfig := &Config{
		Finisher: Finisher{
			Terminate: Terminate{
				Command: "invalidcommand",
			},
		},
	}

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		finisher      *finisher
		expectedValue state
		isError       bool
	}{
		// no error
		{
			finisher: &finisher{
				config: validConfig,
				client: validClient,
				logger: logger,
			},
			expectedValue: waitState,
			isError:       false,
		},

		// errors
		{
			finisher: &finisher{
				config: invalidConfig,
				client: validClient,
				logger: logger,
			},
			expectedValue: abortedState,
			isError:       true,
		},

		{
			finisher: &finisher{
				config: validConfig,
				client: invalidClient,
				logger: logger,
			},
			expectedValue: abortedState,
			isError:       true,
		},
	}

	for _, c := range cases {
		v, err := c.finisher.terminateHandler()

		if v != c.expectedValue {
			t.Fatalf("expected: %v, but actual: %v\n", c.expectedValue, v)
		}
		if (err == nil && c.isError) || (err != nil && !c.isError) {
			t.Fatalf("expected: %v, but actual: %v\n", c.isError, err)
		}
	}
}

func TestFinisher_waitHandler(t *testing.T) {
	publicIpAddr := "8.8.8.8"
	pricateIpAddr := "192.168.0.1"

	validClient := &fakeAwsClient{
		ResDescribeInstance: &ec2.Instance{
			PublicIpAddress:  &publicIpAddr,
			PrivateIpAddress: &pricateIpAddr,
		},
	}
	invalidClient := &fakeAwsClient{
		ResDescribeInstance: nil,
		ErrDescribeInstance: errors.New("error"),
	}

	validConfig := &Config{
		Finisher: Finisher{
			Wait: Wait{
				Command:     "echo test",
				MaxTries:    1,
				IntervalSec: 0,
			},
		},
	}
	invalidConfig := &Config{
		Finisher: Finisher{
			Wait: Wait{
				Command:     "invalidcommand",
				MaxTries:    1,
				IntervalSec: 0,
			},
		},
	}

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		finisher      *finisher
		expectedValue state
		isError       bool
	}{
		// no error
		{
			finisher: &finisher{
				config: validConfig,
				client: validClient,
				logger: logger,
			},
			expectedValue: completedState,
			isError:       false,
		},

		// errors
		{
			finisher: &finisher{
				config: invalidConfig,
				client: validClient,
				logger: logger,
			},
			expectedValue: abortedState,
			isError:       true,
		},

		{
			finisher: &finisher{
				config: validConfig,
				client: invalidClient,
				logger: logger,
			},
			expectedValue: abortedState,
			isError:       true,
		},
	}

	for _, c := range cases {
		v, err := c.finisher.waitHandler()

		if v != c.expectedValue {
			t.Fatalf("expected: %v, but actual: %v\n", c.expectedValue, v)
		}
		if (err == nil && c.isError) || (err != nil && !c.isError) {
			t.Fatalf("expected: %v, but actual: %v\n", c.isError, err)
		}
	}
}

func TestFinisher_completeHandler(t *testing.T) {
	validClient := &fakeAwsClient{
		ErrCompleteLifecycleAction: nil,
	}
	invalidClient := &fakeAwsClient{
		ErrCompleteLifecycleAction: errors.New("error"),
	}

	config := &Config{
		Finisher: Finisher{
			LifecycleHookName:     "test",
			LifecycleActionResult: "ABANDON",
		},
	}

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		finisher      *finisher
		expectedValue state
		isError       bool
	}{
		// no error
		{
			finisher: &finisher{
				config: config,
				client: validClient,
				logger: logger,
			},
			expectedValue: finishedState,
			isError:       false,
		},

		// error
		{
			finisher: &finisher{
				config: config,
				client: invalidClient,
				logger: logger,
			},
			expectedValue: abortedState,
			isError:       true,
		},
	}

	for _, c := range cases {
		v, err := c.finisher.completeHandler()

		if v != c.expectedValue {
			t.Fatalf("expected: %v, but actual: %v\n", c.expectedValue, v)
		}
		if (err == nil && c.isError) || (err != nil && !c.isError) {
			t.Fatalf("expected: %v, but actual: %v\n", c.isError, err)
		}
	}
}

func TestFinisher_IsFinished(t *testing.T) {
	cases := []struct {
		finisher *finisher
		expected bool
	}{
		// should be finished
		{
			finisher: &finisher{
				state: finishedState,
			},
			expected: true,
		},
		{
			finisher: &finisher{
				state: abortedState,
			},
			expected: true,
		},

		{
			finisher: &finisher{
				state: initState,
			},
			expected: false,
		},
		{
			finisher: &finisher{
				state: terminateState,
			},
			expected: false,
		},
		{
			finisher: &finisher{
				state: waitState,
			},
			expected: false,
		},
		{
			finisher: &finisher{
				state: completedState,
			},
			expected: false,
		},
	}

	for _, c := range cases {
		actual := c.finisher.IsFinished()

		if actual != c.expected {
			t.Fatalf("expected: %v, but actual: %v\n", c.expected, actual)
		}
	}
}
