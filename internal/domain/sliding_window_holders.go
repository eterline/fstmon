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

// Slice - slice from sliding window.
func (sw SlidingWindow[T]) Slice() []T {
	return sw.Values
}

// Window - time window from sliding window.
func (sw SlidingWindow[T]) Window() time.Duration {
	return sw.Duration
}

// Len - len of sliding window.
func (sw SlidingWindow[T]) Len() int {
	return len(sw.Values)
}

// ============================ CPU domain structures ============================

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
NetworkingIOTimeline - instantaneous counters and speeds of a single network interface.

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
InterfacesIOTimeline - network interface names to their instantaneous NetworkingIOTimeline metrics.
*/
type InterfacesIOTimeline map[string]NetworkingIOTimeline

// ============================ System domain structures ============================

type AvgLoadTimeline struct {
	Load1        SlidingWindow[float64] `json:"load1"`         // Historical 1-minute average load
	Load5        SlidingWindow[float64] `json:"load5"`         // Historical 5-minute average load
	Load15       SlidingWindow[float64] `json:"load15"`        // Historical 15-minute average load
	RunningProcs SlidingWindow[int]     `json:"running_procs"` // Historical number of currently running processes
	TotalProcs   SlidingWindow[int]     `json:"total_procs"`   // Historical total number of processes
}

/*
MemoryMetricsTimeline - Historical system memory metrics.
*/
type MemoryMetricsTimeline struct {
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
ThermalMetricsTimeline - the historical thermal data of a sensor.

	Each field contains a sliding window of temperature values collected
	over a defined time interval.
*/
type ThermalMetricsTimeline struct {
	Current SlidingWindow[float64] `json:"current"`
	Minimum SlidingWindow[float64] `json:"minimum"`
	Maximum SlidingWindow[float64] `json:"maximum"`
}

/*
ThermalMetricMapTimeline - sensor identifiers to their corresponding thermal timeline structures.
*/
type ThermalMetricMapTimeline map[string]ThermalMetricsTimeline

// ============================ Storage domain structures ============================

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
