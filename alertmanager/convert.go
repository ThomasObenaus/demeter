package alertmanager

import (
	"strings"

	saa "github.com/thomasobenaus/sokar/scaleAlertAggregator"
)

func genEmitterName(name string) string {

	result := "AM"
	if len(name) > 0 {
		result += "." + name
	}

	return result
}

// amResponseToScalingAlerts extracts alerts from the response of the alertmanager
func amResponseToScalingAlerts(resp response) (emitter string, packet saa.ScaleAlertPacket) {
	result := saa.ScaleAlertPacket{}
	for _, alert := range resp.Alerts {
		result.ScaleAlerts = append(result.ScaleAlerts, amAlertToScalingAlert(alert))
	}

	return genEmitterName(resp.Receiver), result
}

func amAlertToScalingAlert(alert alert) saa.ScaleAlert {

	name, ok := alert.Labels["alertname"]
	if !ok {
		name = "NO_NAME"
	}

	return saa.ScaleAlert{
		Name:      name,
		Firing:    isFiring(alert.Status),
		StartedAt: alert.StartsAt,
	}
}

func isFiring(status string) bool {
	status = strings.ToLower(status)
	status = strings.TrimSpace(status)
	return status == "firing"
}
