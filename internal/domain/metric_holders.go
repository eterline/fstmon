package domain

import "time"

// Numerable is a constraint that matches all numeric types in Go.
type Numerable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func NewSpeedIO[T Numerable](rx, tx T) SpeedIO[T] {
	return SpeedIO[T]{
		RX: rx,
		TX: tx,
	}
}

type SpeedIO[T Numerable] struct {
	RX T `json:"rx"`
	TX T `json:"tx"`
}

// ============================ CPU domain structures ============================

/*
CpuCoreInfo - static information about a single CPU core

	such as its physical package, core ID, number of siblings, and cache size.
*/
type CpuCoreInfo struct {
	PhysicalID int `json:"physical_id"` // ID of the physical CPU package
	CoreID     int `json:"core_id"`     // ID of the core within the physical package
	Siblings   int `json:"siblings"`    // Number of sibling cores
	CacheKB    int `json:"cache_kb"`    // L3 cache size in KB
}

/*
CpuPackage - static information about the CPU package, including
	vendor, model name, microcode version, CPU flags, and its cores.
*/

type CpuPackage struct {
	Vendor    string        `json:"vendor"`     // CPU vendor ID
	ModelName string        `json:"model_name"` // Human-readable model name
	Microcode string        `json:"microcode"`  // Microcode version
	Flags     []string      `json:"flags"`      // CPU feature flags
	Cores     []CpuCoreInfo `json:"cores"`      // Slice of CPU cores
}

// =======

/*
CpuCoreMetrics - instantaneous dynamic metrics of a single CPU core,

	including load percentage and current frequency in MHz.
*/
type CpuCoreMetrics struct {
	Load      float64 `json:"load"`      // Current core load in percent
	Frequency float64 `json:"frequnecy"` // Current core frequency in MHz
}

/*
CpuMetrics - current dynamic metrics for all CPU cores, including the average CPU load across all cores.
*/
type CpuMetrics struct {
	Average CpuCoreMetrics   `json:"average"` // Average metrics across all cores
	Cores   []CpuCoreMetrics `json:"cores"`   // Metrics per individual core
}

// ============================ Networking domain structures ============================

/*
NetworkingIO - instantaneous counters and speeds of a single

	network interface, including received/sent bytes, packets, and errors.
*/
type NetworkingIO struct {
	BytesFullRX uint64          `json:"bytes_full_rx"` // Total bytes received
	BytesFullTX uint64          `json:"bytes_full_tx"` // Total bytes sent
	BytesPerSec SpeedIO[uint64] `json:"bytes_per_sec"`

	PacketsRx     uint64          `json:"packets_rx"` // Total packets received
	PacketsTx     uint64          `json:"packets_tx"` // Total packets sent
	PacketsPerSec SpeedIO[uint64] `json:"packets_per_sec"`

	ErrPacketsRx uint64 `json:"err_packets_rx"` // Total RX packet errors
	ErrPacketsTx uint64 `json:"err_packets_tx"` // Total TX packet errors
}

/*
InterfacesIO - network interface names to their instantaneous.
*/
type InterfacesIO map[string]NetworkingIO

// ============================ System domain structures ============================

type SystemInfo struct {
	Uptime       time.Duration `json:"uptime"`
	Idle         time.Duration `json:"idle"`
	Load1        float64       `json:"load1"`         // 1-minute average load
	Load5        float64       `json:"load5"`         // 5-minute average load
	Load15       float64       `json:"load15"`        // 15-minute average load
	RunningProcs int           `json:"running_procs"` // number of currently running processes
	TotalProcs   int           `json:"total_procs"`   // total number of processes
}

// =======

/*
MemoryMetric - system memory metrics.
*/
type MemoryMetrics struct {
	Total     uint64 `json:"total"`     // total physical RAM (bytes)
	Available uint64 `json:"available"` // memory available for allocation (bytes)
	Used      uint64 `json:"used"`      // memory actively used by applications (bytes)
	Free      uint64 `json:"free"`      // unallocated physical memory (bytes)

	SwapTotal     uint64 `json:"swap_total"`     // total swap space (bytes)
	SwapAvailable uint64 `json:"swap_available"` // available swap space (bytes)
	SwapUsed      uint64 `json:"swap_used"`      // used swap space (bytes)
	SwapFree      uint64 `json:"swap_free"`      // free swap space (bytes)

	UsedPercent     float64 `json:"used_percent"`      // percentage of used memory
	SwapUsedPercent float64 `json:"swap_used_percent"` // percentage of used swap
}

// =======

/*
ThermalMetric - an instantaneous thermal state of a sensor.
It includes current, critical, and maximum observed temperatures.
*/
type ThermalMetrics struct {
	Current float64 `json:"current"`
	Max     float64 `json:"max"`
	Crit    float64 `json:"crit"`
}

/*
ThermalMetricMap - sensor identifiers to their instantaneous thermal metrics.
*/
type ThermalMetricsMap map[string]ThermalMetrics

// ============================ Storage domain structures ============================

/*
PartitionInfo - static information about a filesystem partition.

	These attributes describe the identity and configuration of the partition, and do not change over time.
*/
type PartitionInfo struct {
	Device     string   `json:"device"`     // Device path, e.g. "/dev/sda1"
	Mount      string   `json:"mount"`      // Mount point, e.g. "/"
	Filesystem string   `json:"filesystem"` // Filesystem type, e.g. "ext4"
	Options    []string `json:"options"`    // Mount options
}

/*
PartitionsInfo - a map of mount point to static partition info.
*/
type PartitionsInfo []PartitionInfo

// =======

/*
PartitionMetrics - dynamic, time-sensitive metrics of a filesystem partition such as space usage and inode usage.
*/
type PartitionMetrics struct {
	TotalBytes  uint64  `json:"total"`        // Total bytes
	UsedBytes   uint64  `json:"used"`         // Used bytes
	FreeBytes   uint64  `json:"free"`         // Free bytes
	UsedPercent float64 `json:"used_percent"` // Used space percentage

	InodesTotal       uint64  `json:"inodes_total"`        // Total inodes
	InodesUsed        uint64  `json:"inodes_used"`         // Used inodes
	InodesFree        uint64  `json:"inodes_free"`         // Free inodes
	InodesUsedPercent float64 `json:"inodes_used_percent"` // Inode usage percentage
}

type PartitionMetricsMap map[string]PartitionMetrics

type PartitionIO struct{}

type PartitionsIO map[string]PartitionIO
