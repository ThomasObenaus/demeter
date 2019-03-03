package scaler

import (
	"fmt"

	"github.com/thomasobenaus/sokar/helper"
	sokar "github.com/thomasobenaus/sokar/sokar/iface"
)

func amountToScaleType(amount int) string {
	scaleTypeStr := "UP"
	if amount < 0 {
		scaleTypeStr = "DOWN"
	}
	return scaleTypeStr
}

// checkScalingPolicy verifies if the desired
func checkScalingPolicy(desiredCount uint, min uint, max uint) policyCheckResult {

	result := policyCheckResult{minPolicyViolated: false, maxPolicyViolated: false}

	result.desiredCount = desiredCount
	result.validCount = desiredCount

	// check if desiredCount exceeds minimum
	if desiredCount < min {
		result.validCount = min
		result.minPolicyViolated = true
	}

	// check if desiredCount exceeds maximum
	if desiredCount > max {
		result.validCount = max
		result.maxPolicyViolated = true
	}

	return result
}

// trueIfNil returns a ScaleResult filled in with an appropriate error message in case the given scaler is nil
func trueIfNil(s *Scaler) (result sokar.ScaleResult, ok bool) {
	ok = false
	result = sokar.ScaleResult{State: sokar.ScaleUnknown}

	if s == nil {
		ok = true
		result = sokar.ScaleResult{
			State:            sokar.ScaleFailed,
			StateDescription: "Scaler is nil",
			NewCount:         0,
		}
	}
	return result, ok
}

// scale scales the job from currentCount to desiredCount.
// Internally it is checked if a scaling is needed and if the scaling policy is valid.
func (s *Scaler) scale(desiredCount uint, currentCount uint) sokar.ScaleResult {
	if r, ok := trueIfNil(s); ok {
		return r
	}

	jobName := s.job.jobName
	min := s.job.minCount
	max := s.job.maxCount

	s.logger.Info().Str("job", jobName).Msgf("Request to scale job from %d to %d.", currentCount, desiredCount)

	dead, err := s.scalingTarget.IsJobDead(jobName)
	if err != nil {
		return sokar.ScaleResult{
			State:            sokar.ScaleFailed,
			StateDescription: fmt.Sprintf("Error obtaining if job is dead: %s.", err.Error()),
		}
	}

	if dead {
		return sokar.ScaleResult{
			State:            sokar.ScaleIgnored,
			StateDescription: fmt.Sprintf("Job '%s' is dead. Can't scale", jobName),
		}
	}

	chkResult := checkScalingPolicy(desiredCount, min, max)
	newCount := chkResult.validCount
	if chkResult.minPolicyViolated {
		s.logger.Info().Str("job", jobName).Msgf("Job.MinCount (%d) policy violated (wanted %d, have %d). Scale limited to %d.", min, chkResult.desiredCount, currentCount, newCount)
	}
	if chkResult.maxPolicyViolated {
		s.logger.Info().Str("job", jobName).Msgf("Job.MaxCount (%d) policy violated (wanted %d, have %d). Scale limited to %d.", max, chkResult.desiredCount, currentCount, newCount)
	}

	diff := helper.SubUint(newCount, currentCount)
	scaleNeeded := (diff != 0)

	if !scaleNeeded {
		s.logger.Info().Str("job", jobName).Msg("No scaling needed/ possible.")
		return sokar.ScaleResult{
			State:            sokar.ScaleIgnored,
			StateDescription: "No scaling needed/ possible.",
			NewCount:         newCount,
		}
	}

	scaleTypeStr := amountToScaleType(diff)
	s.logger.Info().Str("job", jobName).Msgf("Scale %s by %d to %d.", scaleTypeStr, diff, newCount)

	// Set the new job count
	err = s.scalingTarget.SetJobCount(s.job.jobName, newCount)
	if err != nil {
		return sokar.ScaleResult{
			State:            sokar.ScaleFailed,
			StateDescription: fmt.Sprintf("Error adjusting job count to %d: %s.", newCount, err.Error()),
		}
	}

	return sokar.ScaleResult{
		State:            sokar.ScaleDone,
		StateDescription: "Scaling successfully done.",
		NewCount:         newCount,
	}
}

// ScaleBy Scales the target component by the given amount of instances
func (s *Scaler) ScaleBy(amount int) sokar.ScaleResult {
	if r, ok := trueIfNil(s); ok {
		return r
	}

	jobName := s.job.jobName
	count, err := s.scalingTarget.GetJobCount(jobName)
	if err != nil {
		return sokar.ScaleResult{
			State:            sokar.ScaleFailed,
			StateDescription: fmt.Sprintf("Error obtaining job count: %s.", err.Error()),
		}
	}

	desiredCount := helper.IncUint(count, amount)

	return s.scale(desiredCount, count)
}