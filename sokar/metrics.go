package sokar

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	m "github.com/thomasobenaus/sokar/metrics"
)

// Metrics represents the collection of metrics internally set by sokar.
type Metrics struct {
	scaleEventsTotal   m.Counter
	failedScalingTotal m.Counter
	preScaleJobCount   m.Gauge
	plannedJobCount    m.Gauge
	postScaleJobCount  m.Gauge
	scaleFactor        m.Gauge
}

// NewMetrics returns the metrics collection needed for the SAA.
func NewMetrics() Metrics {

	scaleEventsTotal := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "sokar",
		Name:      "scale_events_total",
		Help:      "Number of received ScaleEvents in total.",
	})

	failedScalingTotal := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "sokar",
		Name:      "failed_scaling_total",
		Help:      "Number of failed scaling actions in total.",
	})

	preScaleJobCount := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "cap",
		Name:      "pre_scale_job_count",
		Help:      "The job count before the scaling action. Based on this count sokar does the planning.",
	})

	plannedJobCount := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "cap",
		Name:      "planned_job_count",
		Help:      "The count planned by the CapacityPlanner for the current scale action.",
	})

	postScaleJobCount := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "cap",
		Name:      "post_scale_job_count",
		Help:      "The job count after the scaling action.",
	})

	scaleFactor := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "cap",
		Name:      "scale_factor",
		Help:      "The scale factor (gradient) as it was received with a ScalingEvent.",
	})

	return Metrics{
		scaleEventsTotal:   scaleEventsTotal,
		failedScalingTotal: failedScalingTotal,
		preScaleJobCount:   preScaleJobCount,
		plannedJobCount:    plannedJobCount,
		postScaleJobCount:  postScaleJobCount,
		scaleFactor:        scaleFactor,
	}
}
