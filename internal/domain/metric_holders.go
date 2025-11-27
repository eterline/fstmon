package domain

import "time"

// SlidingWindow represents a generic time-based window of values.
// Duration specifies the time span of the window, and Values stores
// the measured values over time.
type SlidingWindow[T any] struct {
	Duration time.Duration `json:"duration"`
	Values   []T           `json:"values"`
}

//
// CPU domain structures
//

// CpuCoreInfo contains static information about a single CPU core
// such as its physical package, core ID, number of siblings, and cache size.
type CpuCoreInfo struct {
	PhysicalID int `json:"physical_id"` // ID of the physical CPU package
	CoreID     int `json:"core_id"`     // ID of the core within the physical package
	Siblings   int `json:"siblings"`    // Number of sibling cores
	CacheKB    int `json:"cache_kb"`    // L1/L2/L3 cache size in KB
}

// CpuPackage holds static information about the CPU package, including
// vendor, model name, microcode version, CPU flags, and its cores.
type CpuPackage struct {
	Vendor    string        `json:"vendor"`     // CPU vendor ID
	ModelName string        `json:"model_name"` // Human-readable model name
	Microcode string        `json:"microcode"`  // Microcode version
	Flags     []string      `json:"flags"`      // CPU feature flags
	Cores     []CpuCoreInfo `json:"cores"`      // Slice of CPU cores
}

//
// Dynamic CPU metrics
//

// CpuCoreMetrics represents instantaneous dynamic metrics of a single CPU core,
// including load percentage and current frequency in MHz.
type CpuCoreMetrics struct {
	Load      float64 `json:"load"`      // Current core load in percent
	Frequency float64 `json:"frequnecy"` // Current core frequency in MHz
}

// CpuMetrics represents current dynamic metrics for all CPU cores, including
// the average CPU load across all cores.
type CpuMetrics struct {
	Average CpuCoreMetrics   `json:"average"` // Average metrics across all cores
	Cores   []CpuCoreMetrics `json:"cores"`   // Metrics per individual core
}

//
// CPU timeline metrics
//

// CpuCoreTimeline represents historical metrics of a single CPU core using
// a sliding time window for load and frequency.
type CpuCoreTimeline struct {
	Load      SlidingWindow[float64] `json:"load"`      // Historical load %
	Frequency SlidingWindow[float64] `json:"frequnecy"` // Historical frequency in MHz
}

// CpuTimeline represents historical metrics of all CPU cores, including
// average CPU load across all cores.
type CpuTimeline struct {
	Average CpuCoreTimeline   `json:"average"` // Historical average metrics
	Cores   []CpuCoreTimeline `json:"cores"`   // Timeline per individual core
}

//
// Networking domain structures
//

// NetworkingIO represents instantaneous counters and speeds of a single
// network interface, including received/sent bytes, packets, and errors.
type NetworkingIO struct {
	FullRX uint64 `json:"full_rx"` // Total bytes received
	FullTX uint64 `json:"full_tx"` // Total bytes sent

	SpeedRX uint64 `json:"speed_rx"` // RX speed in bytes per second
	SpeedTX uint64 `json:"speed_tx"` // TX speed in bytes per second

	PacketsRx uint64 `json:"packets_rx"` // Total packets received
	PacketsTx uint64 `json:"packets_tx"` // Total packets sent

	ErrPacketsRx uint64 `json:"err_packets_rx"` // Total RX packet errors
	ErrPacketsTx uint64 `json:"err_packets_tx"` // Total TX packet errors
}

// InterfacesIO maps network interface names to their instantaneous
// NetworkingIO metrics.
type InterfacesIO map[string]NetworkingIO
