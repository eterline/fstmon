package domain

import "time"

// SlidingWindow represents a generic time-based window of values.
// Duration specifies the time span of the window, and Values stores
// the measured values over time.
type SlidingWindow[T any] struct {
	Duration time.Duration `json:"duration"`
	Values   []T           `json:"values"`
}

func CreateSlidingWindow[T any](l []T, intervalDuration time.Duration) SlidingWindow[T] {
	return SlidingWindow[T]{
		Duration: intervalDuration * time.Duration(len(l)),
		Values:   l,
	}
}

func CreateSlidingWindowSeconds[T any](l []T, sec int) SlidingWindow[T] {
	if sec < 0 {
		sec = 0
	}

	intervalDuration := time.Duration(sec) * time.Second
	return SlidingWindow[T]{
		Duration: intervalDuration * time.Duration(len(l)),
		Values:   l,
	}
}

// Add adds a new value to the sliding window.
// Optionally trims old values if needed (not implemented here; see note below).
func (sw *SlidingWindow[T]) Add(value T) {
	sw.Values = append(sw.Values, value)
}

func (sw SlidingWindow[T]) Slice() []T {
	return sw.Values
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

// =======

/*
CpuCoreTimeline - historical metrics of a single CPU core .
*/
type CpuCoreTimeline struct {
	Load      SlidingWindow[float64] `json:"load"`      // Historical load %
	Frequency SlidingWindow[float64] `json:"frequnecy"` // Historical frequency in MHz
}

/*
CpuTimeline - historical metrics of all CPU cores, including average CPU load across all cores.
*/
type CpuTimeline struct {
	Average CpuCoreTimeline   `json:"average"` // Historical average metrics
	Cores   []CpuCoreTimeline `json:"cores"`   // Timeline per individual core
}

// ============================ Networking domain structures ============================

/*
NetworkingIO - instantaneous counters and speeds of a single

	network interface, including received/sent bytes, packets, and errors.
*/
type NetworkingIO struct {
	BytesFullRX uint64 `json:"bytes_full_rx"` // Total bytes received
	BytesFullTX uint64 `json:"bytes_full_tx"` // Total bytes sent

	BytesSpeedRX uint64 `json:"bytes_speed_rx"` // RX speed in bytes per second
	BytesSpeedTX uint64 `json:"bytes_speed_tx"` // TX speed in bytes per second

	PacketsRx uint64 `json:"packets_rx"` // Total packets received
	PacketsTx uint64 `json:"packets_tx"` // Total packets sent

	ErrPacketsRx uint64 `json:"err_packets_rx"` // Total RX packet errors
	ErrPacketsTx uint64 `json:"err_packets_tx"` // Total TX packet errors
}

/*
InterfacesIO - network interface names to their instantaneous.
*/
type InterfacesIO map[string]NetworkingIO

/*
NetworkingIO - instantaneous counters and speeds of a single network interface.

	Including received/sent bytes, packets, and errors.
*/
type NetworkingIOTimeline struct {
	BytesFullRX SlidingWindow[uint64] `json:"bytes_full_rx"` // Historical Total bytes received
	BytesFullTX SlidingWindow[uint64] `json:"bytes_full_tx"` // Historical Total bytes sent

	BytesSpeedRX SlidingWindow[uint64] `json:"bytes_speed_rx"` // Historical RX speed in bytes per second
	BytesSpeedTX SlidingWindow[uint64] `json:"bytes_speed_tx"` // Historical TX speed in bytes per second

	PacketsRx SlidingWindow[uint64] `json:"packets_rx"` // Historical Total packets received
	PacketsTx SlidingWindow[uint64] `json:"packets_tx"` // Historical Total packets sent

	ErrPacketsRx SlidingWindow[uint64] `json:"err_packets_rx"` // Historical Total RX packet errors
	ErrPacketsTx SlidingWindow[uint64] `json:"err_packets_tx"` // Historical Total TX packet errors
}

/*
InterfacesIO - network interface names to their instantaneous NetworkingIO metrics.
*/
type InterfacesIOTimeline map[string]NetworkingIOTimeline

// ============================ System domain structures ============================

type SystemInfo struct {
	Uptime  time.Duration `json:"uptime"`
	Idle    time.Duration `json:"idle"`
	AvgLoad AvgLoad       `json:"average_load"`
}

type AvgLoad struct {
	Load1        float64 `json:"load1"`         // 1-minute average load
	Load5        float64 `json:"load5"`         // 5-minute average load
	Load15       float64 `json:"load15"`        // 15-minute average load
	RunningProcs int     `json:"running_procs"` // number of currently running processes
	TotalProcs   int     `json:"total_procs"`   // total number of processes
}

type AvgLoadTimeline struct {
	Load1        SlidingWindow[float64] `json:"load1"`         // Historical 1-minute average load
	Load5        SlidingWindow[float64] `json:"load5"`         // Historical 5-minute average load
	Load15       SlidingWindow[float64] `json:"load15"`        // Historical 15-minute average load
	RunningProcs SlidingWindow[int]     `json:"running_procs"` // Historical number of currently running processes
	TotalProcs   SlidingWindow[int]     `json:"total_procs"`   // Historical total number of processes
}

/*
MemoryMetric - system memory metrics.
*/
type MemoryMetric struct {
	Total       uint64  `json:"total"`        // total physical RAM (bytes)
	Available   uint64  `json:"available"`    // memory available for allocation (bytes)
	Used        uint64  `json:"used"`         // memory actively used by applications (bytes)
	UsedPercent float64 `json:"used_percent"` // percentage of used memory
	Free        uint64  `json:"free"`         // unallocated physical memory (bytes)

	SwapTotal uint64 `json:"swap_total"` // total swap space (bytes)
	SwapUsed  uint64 `json:"swap_used"`  // used swap space (bytes)
	SwapFree  uint64 `json:"swap_free"`  // free swap space (bytes)
}

/*
MemoryMetricTimeline - Historical system memory metrics.
*/
type MemoryMetricTimeline struct {
	Total       SlidingWindow[uint64]  `json:"total"`        // Historical total physical RAM (bytes)
	Available   SlidingWindow[uint64]  `json:"available"`    // Historical memory available for allocation (bytes)
	Used        SlidingWindow[uint64]  `json:"used"`         // Historical memory actively used by applications (bytes)
	UsedPercent SlidingWindow[float64] `json:"used_percent"` // Historical percentage of used memory
	Free        SlidingWindow[uint64]  `json:"free"`         // Historical unallocated physical memory (bytes)

	SwapTotal SlidingWindow[uint64] `json:"swap_total"` // Historical total swap space (bytes)
	SwapUsed  SlidingWindow[uint64] `json:"swap_used"`  // Historical used swap space (bytes)
	SwapFree  SlidingWindow[uint64] `json:"swap_free"`  // Historical free swap space (bytes)
}

// =======

/*
ThermalMetric - an instantaneous thermal state of a sensor.
It includes current, minimum, and maximum observed temperatures.
*/
type ThermalMetric struct {
	Current float64 `json:"current"`
	Minimum float64 `json:"minimum"`
	Maximum float64 `json:"maximum"`
}

/*
ThermalMetricMap - sensor identifiers to their instantaneous thermal metrics.
*/
type ThermalMetricMap map[string]ThermalMetric

/*
ThermalMetricTimeline - the historical thermal data of a sensor.

	Each field contains a sliding window of temperature values collected
	over a defined time interval.
*/
type ThermalMetricTimeline struct {
	Current SlidingWindow[float64] `json:"current"`
	Minimum SlidingWindow[float64] `json:"minimum"`
	Maximum SlidingWindow[float64] `json:"maximum"`
}

/*
ThermalMetricMapTimeline - sensor identifiers to their corresponding thermal timeline structures.
*/
type ThermalMetricMapTimeline map[string]ThermalMetricTimeline

// ============================ Storage domain structures ============================

/*
PartitionInfo - static information about a filesystem partition.

	These attributes describe the identity and configuration of the partition, and do not change over time.
*/
type PartitionInfo struct {
	Device     string   `json:"device"`      // Device path, e.g. "/dev/sda1"
	MountPoint string   `json:"mount_point"` // Mount point, e.g. "/"
	FsType     string   `json:"fs_type"`     // Filesystem type, e.g. "ext4"
	Opts       []string `json:"opts"`        // Mount options
}

/*
PartitionsInfo - a map of mount point to static partition info.
*/
type PartitionsInfo map[string]PartitionInfo

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

/*
PartitionMetricsTimeline - historical metrics for a filesystem partition captured over a sliding time window.

	It tracks changes in space and inode usage over time.
*/
type PartitionMetricsTimeline struct {
	TotalBytes  SlidingWindow[uint64]  `json:"total"`        // Total bytes
	UsedBytes   SlidingWindow[uint64]  `json:"used_bytes"`   // Historical used bytes
	FreeBytes   SlidingWindow[uint64]  `json:"free_bytes"`   // Historical free bytes
	UsedPercent SlidingWindow[float64] `json:"used_percent"` // Historical space usage %

	InodesTotal       SlidingWindow[uint64]  `json:"inodes_total"`        // Total inodes
	InodesUsed        SlidingWindow[uint64]  `json:"inodes_used"`         // Historical used inodes
	InodesFree        SlidingWindow[uint64]  `json:"inodes_free"`         // Historical free inodes
	InodesUsedPercent SlidingWindow[float64] `json:"inodes_used_percent"` // Historical inode usage %
}

/*
PartitionsMetricsTimeline - a map of mount point to historical timeline metrics.
*/
type PartitionsMetricsTimeline map[string]PartitionMetricsTimeline
