// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package domain

type SystemData struct {
	Cpu           float64 `json:"cpu"`
	RAM           string  `json:"ram"`
	RAMUsage      float64 `json:"ram_usage"`
	Uptime        string  `json:"uptime"`
	UptimeSeconds uint64  `json:"uptime_seconds"`
}

type AverageLoad struct {
	Load1  float64 `json:"load_1"`
	Load5  float64 `json:"load_5"`
	Load15 float64 `json:"load_15"`
	Procs  string  `json:"procs"`
}

type (
	InterfacesData map[string]NetworkingData

	NetworkingData struct {
		FullRX  string `json:"full_rx"`
		FullTX  string `json:"full_tx"`
		SpeedRX string `json:"speed_rx"`
		SpeedTX string `json:"speed_tx"`
	}
)

type (
	PartsUsages map[string]PartUse

	PartUse struct {
		Name    string `json:"name"`
		Size    string `json:"size"`
		Use     string `json:"use"`
		Percent int32  `json:"percent"`
	}
)

type TemperatureMap map[string]float64

type (
	CpuCore struct {
		Load      float64 `json:"load"`
		Frequency float64 `json:"frequency"`
	}

	CpuLoad struct {
		Average float64   `json:"average"`
		Cores   []CpuCore `json:"cores"`
	}
)
