# shutter

A single daemon based graceful shutdown manager for EC2 instances under an Autoscaling Group and a Lifecycle Hook.

## How to use

It accepts basic settings via a yaml file. You can specify the file path by `-file`.
And you can also run shutter as a daemon by using `-daemon`.

```sh
$ ./shutter -h
Usage of ./shutter:
  -daemon
        do as daemon
  -file string
        a config file path
  -instanceid string
        EC2 instance id (optional, used if daemon = false)
```

The yaml configuration layout is here:
Especially about `command` part, you can use these replacement will be replaced by actual instance values on runtime.

- `${INSTANCE_ID}`, an EC2 instance id
- `${PUBLIC_IP_ADDRESS}`, a public IP address for an EC2 instance
- `${PRIVATE_IP_ADDRESS}`, a private IP address for an EC2 instance

```yaml
aws_region: us-east-1
watcher:
  autoscaling_group_name: test
  interval_sec: 60
finisher:
  lifecycle_hook_name: test
  lifecycle_action_result: ABANDON
  terminate:
    # it should have idempotency
    command: "ping -c 5 ${PUBLIC_IP_ADDRESS}"
  wait:
    command: "ping -c 5 ${PUBLIC_IP_ADDRESS}"
    interval_sec: 60
    max_tries: 30
```

## Required minimum IAM policy

```
autoscaling:CompleteLifecycleAction
autoscaling:DescribeAutoScalingInstances
autoscaling:DescribeAutoScalingGroups
ec2:DescribeInstances
```

## shutter internals

### State management

shutter manages terminating instance by a structure like FSM. It has below states:

- init

  - noop
  
- terminate

  - Run a command to start termination. e.g. trigger graceful shutdown scripts
  
- wait

  - Wait until a command is succeeded. e.g. trigger script checks there's no active client
  
  
- complete

  - Complete an instance on an autoscaling group
  
- finished

  - noop, all tasks are succeeded
  
- aborted

  - noop, any task failed
  
