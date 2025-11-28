package system

import (
	"context"

	"github.com/eterline/fstmon/internal/domain"
	sensorPs "github.com/shirou/gopsutil/v4/sensors"
)

type hardwareThermalMetrics struct{}

func NewHardwareThermalMetrics() *hardwareThermalMetrics {
	return &hardwareThermalMetrics{}
}

func (hms *hardwareMetricSystem) ScrapeThermalMetrics(ctx context.Context) (domain.MetricWrapper[domain.ThermalMetricsMap], error) {
	stat, err := sensorPs.TemperaturesWithContext(ctx)
	if err != nil {
		return domain.EmptyWrapMetric[domain.ThermalMetricsMap](),
			ErrScrapeThermalMetrics.Wrap(err)
	}

	data := make(domain.ThermalMetricsMap, len(stat))

	for _, s := range stat {
		data[s.SensorKey] = domain.ThermalMetrics{
			Current: s.Temperature,
			Max:     s.High,
			Crit:    s.Critical,
		}
	}

	return domain.WrapMetric(data), nil
}
