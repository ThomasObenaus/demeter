package aws

import (
	"fmt"
	"time"

	aws "github.com/aws/aws-sdk-go/service/autoscaling"
	iface "github.com/thomasobenaus/sokar/aws/iface"
)

// MonitorInstanceScaling will block until the instance is scaled up/ down
func MonitorInstanceScaling(autoScaling iface.AutoScaling, autoScalingGroupName string, activityID string) error {
	const timeout = time.Second * 180

	start := time.Now()

	for {
		state, err := getCurrentScalingState(autoScaling, autoScalingGroupName, activityID)
		if err != nil {
			return err
		}

		if state.progress >= 100 {
			// scaling completed
			return nil
		}

		if time.Since(start) >= timeout {
			return fmt.Errorf("MonitorInstanceScaling timed out after %v", timeout)
		}
	}
	return nil
}

type scalingState struct {
	status   string
	progress int64
}

func getCurrentScalingState(autoScaling iface.AutoScaling, autoScalingGroupName string, activityID string) (*scalingState, error) {

	activityIDs := make([]*string, 0)
	activityIDs = append(activityIDs, &activityID)
	input := aws.DescribeScalingActivitiesInput{AutoScalingGroupName: &autoScalingGroupName, ActivityIds: activityIDs}
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// First create the request
	req, output := autoScaling.DescribeScalingActivitiesRequest(&input)
	if req == nil {
		return nil, fmt.Errorf("Request generated by DescribeScalingActivitiesInput is nil")
	}

	// Now send the request
	if err := req.Send(); err != nil {
		return nil, err
	}

	if output == nil {
		return nil, fmt.Errorf("DescribeScalingActivitiesOutput is invalid")
	}

	if len(output.Activities) == 0 || output.Activities[0].StatusCode == nil || output.Activities[0].Progress == nil {
		return nil, fmt.Errorf("DescribeScalingActivitiesOutput contains no valid activities")
	}

	state := &scalingState{status: *output.Activities[0].StatusCode, progress: *output.Activities[0].Progress}
	return state, nil
}
