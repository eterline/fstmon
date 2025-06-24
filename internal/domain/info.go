package domain

type SystemData struct {
	CpuTemp string `json:"cpu_temp"`
	Uptime  string `json:"uptime"`
	Cpu     string `json:"cpu"`
	Memory  string `json:"memory"`
}

type ProcessesData struct{}

type NetworkingData struct {
	FullRX  string `json:"full_rx"`
	FullTX  string `json:"full_tx"`
	SpeedRX string `json:"speed_rx"`
	SpeedTX string `json:"speed_tx"`
}
