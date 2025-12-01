// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package domain

import "time"

type MetricState struct {
	Value      any
	Available  bool
	LastUpdate time.Time
}

// ============================ CPU domain structures ============================

/*
CpuCoreInfo – static information about a single CPU core

	such as its physical package, core ID, number of siblings, and cache size.
*/
type CpuCoreInfo struct {
	PhysicalID int `json:"physical_id"` // ID of the physical CPU package
	CoreID     int `json:"core_id"`     // ID of the core within the physical package
	Siblings   int `json:"siblings"`    // Number of sibling cores
	CacheKB    int `json:"cache_kb"`    // L3 cache size in KB
}

/*
CpuPackage – static information about the CPU package, including
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
CpuCoreMetrics – instantaneous dynamic metrics of a single CPU core,

	including load percentage and current frequency in MHz.
*/
type CpuCoreMetrics struct {
	Load      float64 `json:"load"`      // Current core load in percent
	Frequency float64 `json:"frequnecy"` // Current core frequency in MHz
}

/*
CpuMetrics – current dynamic metrics for all CPU cores, including the average CPU load across all cores.
*/
type CpuMetrics struct {
	Average CpuCoreMetrics   `json:"average"` // Average metrics across all cores
	Cores   []CpuCoreMetrics `json:"cores"`   // Metrics per individual core
}

// ============================ Networking domain structures ============================

/*
NetworkingIO – instantaneous counters and speeds of a single

	network interface, including received/sent bytes, packets, and errors.
*/
type InterfaceIO struct {
	BytesTotal       IO[uint64] `json:"bytes_total"`         // Total bytes
	PacketsTotal     IO[uint64] `json:"packets_total"`       // Total packets
	ErrPacketsTotal  IO[uint64] `json:"error_packets_total"` // Total packet errors
	DropPacketsTotal IO[uint64] `json:"drop_packets_total"`  // Total packet drops

	BytesPerSec   IO[uint64] `json:"bytes_per_sec"`   // Per second bytes
	PacketsPerSec IO[uint64] `json:"packets_per_sec"` // Per second packets
}

/*
InterfacesIO – network interface names to their instantaneous.
*/
type InterfacesIOMap map[string]InterfaceIO

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
MemoryMetric – system memory metrics.
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
ThermalMetric – an instantaneous thermal state of a sensor.
It includes current, critical, and maximum observed temperatures.
*/
type ThermalMetrics struct {
	Current float64 `json:"current"`
	Max     float64 `json:"max"`
	Crit    float64 `json:"crit"`
}

/*
ThermalMetricMap – sensor identifiers to their instantaneous thermal metrics.
*/
type ThermalMetricsMap map[string]ThermalMetrics

// ============================ Storage domain structures ============================

/*
PartitionUsage – dynamic, time-sensitive metrics of a filesystem partition such as space usage and inode usage.
*/
type PartitionUsage struct {
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
PartitionInfo – static information about a filesystem partition.

	These attributes describe the identity and configuration of the partition, and do not change over time.
*/
type Partition struct {
	Device     string          `json:"device"`     // Device path, e.g. "/dev/sda1"
	Mount      string          `json:"mount"`      // Mount point, e.g. "/"
	Filesystem string          `json:"filesystem"` // Filesystem type, e.g. "ext4"
	Options    []string        `json:"options"`    // Mount options
	Usage      *PartitionUsage `json:"usage"`      // Mount options
}

/*
PartitionsInfo – a map of mount point to static partition info.
*/
type Partitions []Partition

// =======

// DiskIO – I/O statistics for a partition.
type DiskIO struct {
	IopsInProgress uint64 `json:"iops_in_progress"` // Number of I/O operations currently in progress

	Ops       IO[uint64] `json:"ops"`
	MergedOps IO[uint64] `json:"merged_ops"` // RX = merged_read_count, TX = merged_write_count
	Bytes     IO[uint64] `json:"bytes"`      // RX = read_bytes, TX = write_bytes

	OpsPerSec       IO[uint64] `json:"ops_per_sec"`
	MergedOpsPerSec IO[uint64] `json:"merged_ops_per_sec"`
	BytesPerSec     IO[uint64] `json:"bytes_per_sec"`

	Time       IO[time.Duration] `json:"time"`        // RX = read_time, TX = write_time
	IoTime     time.Duration     `json:"io_time"`     // Total I/O time
	WeightedIO time.Duration     `json:"weighted_io"` // Weighted I/O time
}

type DiskIOMap map[string]DiskIO

// SetReadTime – sets the read time from time.Duration
func (p *DiskIO) SetReadTime(ms uint64) {
	p.Time.RX = time.Millisecond * time.Duration(ms)
}

// SetWriteTime – sets the write time from time.Duration
func (p *DiskIO) SetWriteTime(ms uint64) {
	p.Time.TX = time.Millisecond * time.Duration(ms)
}

// SetTotalTime – sets the total I/O time from time.Duration
func (p *DiskIO) SetTotalTime(ms uint64) {
	p.IoTime = time.Millisecond * time.Duration(ms)
}

// SetWeightedIOTime – sets the weighted I/O time from time.Duration
func (p *DiskIO) SetWeightedIOTime(ms uint64) {
	p.WeightedIO = time.Millisecond * time.Duration(ms)
}
