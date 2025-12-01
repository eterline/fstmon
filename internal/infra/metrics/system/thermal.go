// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
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

func (hms *hardwareThermalMetrics) ScrapeThermalMetrics(ctx context.Context) (domain.ThermalMetricsMap, error) {
	stat, err := sensorPs.TemperaturesWithContext(ctx)
	if err != nil {
		return domain.ThermalMetricsMap{},
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

	return data, nil
}
