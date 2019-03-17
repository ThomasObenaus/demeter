package sokar

import (
	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"
)

func (sk *Sokar) scaleEventProcessor(scaleEventChannel <-chan sokarIF.ScaleEvent) {
	sk.wg.Add(1)
	defer sk.wg.Done()

	sk.logger.Info().Msg("ScaleEventProcessor started.")

	for {
		select {
		case <-sk.stopChan:
			sk.logger.Info().Msg("ScaleEventProcessor stopped.")
			return
		case se := <-scaleEventChannel:
			sk.handleScaleEvent(se)
		}
	}
}

func (sk *Sokar) handleScaleEvent(scaleEvent sokarIF.ScaleEvent) {
	sk.logger.Info().Msgf("Scale Event received: %v", scaleEvent)

	sk.metrics.scaleEventsTotal.Inc()

	currentCount, err := sk.scaler.GetCount()
	if err != nil {
		sk.metrics.failedScalingTotal.Inc()
		sk.logger.Error().Err(err).Msg("Scaling ignored. Failed to obtain current count.")
		return
	}
	sk.metrics.currentCount.Set(float64(currentCount))

	// plan
	plannedCount := sk.capacityPlanner.Plan(scaleEvent.ScaleFactor, currentCount)
	sk.metrics.plannedCount.Set(float64(plannedCount))
	err = sk.scaler.ScaleTo(plannedCount)

	// HACK: For now we ignore all rejected scaling tickets
	if err != nil {
		sk.metrics.failedScalingTotal.Inc()
		sk.logger.Error().Err(err).Msg("Failed to scale.")
	}
}
