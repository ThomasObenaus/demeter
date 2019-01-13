package nomadConnector

import (
	nomadApi "github.com/hashicorp/nomad/api"
)

// queryOptions sets sokars default QueryOptions for making GET calls to
// the API.
func (nc *connectorImpl) queryOptions() (queryOptions *nomadApi.QueryOptions) {
	return &nomadApi.QueryOptions{AllowStale: true}
}

func (nc *connectorImpl) ScaleBy(amount int) error {
	nc.log.Info().Str("job", nc.jobName).Int("amount", amount).Msg("Scaling ...")

	// In order to scale the job, we need information on the current status of the
	// running job from Nomad.
	jobInfo, _, err := nc.nomad.Jobs().Info(nc.jobName, nc.queryOptions())

	if err != nil {
		nc.log.Error().Err(err).Msg("Unable to determine job info")
		return err
	}

	// Use the current task count in order to determine whether or not a scaling
	// event will violate the min/max job policy.
	for i, _ := range jobInfo.TaskGroups {
		//if group.ScaleDirection == ScalingDirectionOut && *taskGroup.Count >= group.Max ||
		//	group.ScaleDirection == ScalingDirectionIn && *taskGroup.Count <= group.Min {
		//	logging.Debug("client/job_scaling: scale %v not permitted due to constraints on job \"%v\" and group \"%v\"",
		//		group.ScaleDirection, *jobInfo.ID, group.GroupName)
		//	return
		//}

		//logging.Info("client/job_scaling: scale %v will now be initiated against job \"%v\" and group \"%v\"",
		//	group.ScaleDirection, jobName, group.GroupName)

		// Depending on the scaling direction decrement/incrament the count;
		// currently replicator only supports addition/subtraction of 1.
		//if *taskGroup.Name == group.GroupName && group.ScaleDirection == ScalingDirectionOut {
		//	*jobResp.TaskGroups[i].Count++
		//	state.ScaleOutRequests++
		//}
		//
		//if *taskGroup.Name == group.GroupName && group.ScaleDirection == ScalingDirectionIn {
		//	*jobResp.TaskGroups[i].Count--
		//	state.ScaleInRequests++
		//}

		*jobInfo.TaskGroups[i].Count++
	}

	// Submit the job to the Register API endpoint with the altered count number
	// and check that no error is returned.
	_, _, err = nc.nomad.Jobs().Register(jobInfo, &nomadApi.WriteOptions{})

	if err != nil {
		nc.log.Error().Err(err).Msg("Unable to scale")
		return err
	}

	//// Track the scaling submission time.
	//state.LastScalingEvent = time.Now()
	//if err != nil {
	//	logging.Error("client/job_scaling: issue submitting job %s for scaling action: %v", jobName, err)
	//	return
	//}
	//
	//// Setup our metric scaling direction namespace.
	//m := fmt.Sprintf("scale_%s", strings.ToLower(group.ScaleDirection))
	//
	//success := nc.scaleConfirmation(resp.EvalID)

	//if !success {
	//	metrics.IncrCounter([]string{"job", jobName, group.GroupName, m, "failure"}, 1)
	//	state.FailureCount++
	//
	//	return
	//}
	//
	//metrics.IncrCounter([]string{"job", jobName, group.GroupName, m, "success"}, 1)
	//logging.Info("client/job_scaling: scaling of job \"%v\" and group \"%v\" successfully completed",
	//	jobName, group.GroupName)

	nc.log.Info().Str("job", nc.jobName).Int("amount", amount).Msg("Scaling ...done")
	return nil
}
