package domain

type SystemData struct {
	CpuTemp string `json:"cpu_temp"`
	Uptime  string `json:"uptime"`
	Cpu     string `json:"cpu"`
	Memory  string `json:"memory"`
}

type AverageLoad struct {
	Load1  float64 `json:"load_1"`
	Load5  float64 `json:"load_5"`
	Load15 float64 `json:"load_15"`
	Procs  string  `json:"procs"`
}

type ProcessesData struct{}

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
	CoreLoad struct {
		Average float64   `json:"average"`
		Cores   []float64 `json:"cores"`
	}

	CpuLoad struct {
		Frames map[string]CoreLoad `json:"frames"`
	}
)
