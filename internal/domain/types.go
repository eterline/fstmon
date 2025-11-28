package domain

import "time"

// ==========================

type metric interface {
	CpuPackage | CpuMetrics |
		InterfacesIO |
		MemoryMetrics |
		PartitionsInfo | PartitionsIO |
		SystemInfo |
		ThermalMetricsMap
}

func WrapMetric[M metric](metric M) MetricWrapper[M] {
	return MetricWrapper[M]{
		CreatedAt: time.Now(),
		Metric:    metric,
	}
}

func EmptyWrapMetric[M metric]() MetricWrapper[M] {
	return MetricWrapper[M]{}
}

type MetricWrapper[M metric] struct {
	CreatedAt time.Time `json:"created_at"`
	Metric    M         `json:"metric"`
}

// ==========================
