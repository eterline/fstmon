// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package httphomepage

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/eterline/fstmon/internal/domain"
)

// ============================ CPU dto ============================

// DTOCpuCore – simplified per-core dynamic metrics.
type DTOCpuCore struct {
	Load      string `json:"load"`      // "12.5%"
	Frequency string `json:"frequency"` // "3200MHz"
}

// DTOCpu – aggregated CPU info for homepage.
type DTOCpu struct {
	Vendor      string `json:"vendor"`       // "Intel", "AMD"
	Model       string `json:"model"`        // "Ryzen 5 5600X"
	CoreCount   int    `json:"core_count"`   // "6"
	ThreadCount int    `json:"thread_count"` // "12"

	Load      string       `json:"load"`      // "15.4%"
	Frequency string       `json:"frequency"` // "3250MHz"`
	Cores     []DTOCpuCore `json:"cores"`     // per-core metrics
}

func Domain2DTOCpu(pkg domain.CpuPackage, m domain.CpuMetrics) *DTOCpu {
	if len(pkg.Cores) == 0 {
		return &DTOCpu{}
	}

	dto := DTOCpu{
		Vendor:      pkg.Vendor,
		Model:       pkg.ModelName,
		CoreCount:   len(pkg.Cores),
		ThreadCount: pkg.Cores[0].Siblings,
		Load:        fmt.Sprintf("%.1f%%", m.Average.Load),
		Frequency:   fmt.Sprintf("%.1fMhz", m.Average.Frequency),
		Cores:       make([]DTOCpuCore, len(m.Cores)),
	}

	for i, core := range m.Cores {
		dto.Cores[i] = DTOCpuCore{
			Load:      fmt.Sprintf("%.1f%%", core.Load),
			Frequency: fmt.Sprintf("%.1fMhz", core.Frequency),
		}
	}

	return &dto
}

// ============================ RAM dto ============================

type DTOMemoryT struct {
	UsedPercent   string `json:"used_percent"`    // "52.5%"
	Total         string `json:"total"`           // "15.62GB"
	Used          string `json:"used"`            // "8.22GB"
	Free          string `json:"free"`            // "7.32GB"
	UsedTotal     string `json:"used_total"`      // "8.22GB/15.62GB"
	UsedTotalFree string `json:"used_total_free"` // "8.22GB/15.62GB/7.32GB"
}

type DTOMemory struct {
	RAM  DTOMemoryT `json:"ram"`
	Swap DTOMemoryT `json:"swap"`
}

func Domain2DTOMemory(v domain.MemoryMetrics) *DTOMemory {
	return &DTOMemory{
		RAM: DTOMemoryT{
			Total:         NewQBBSBuilder(0).Add(v.Total).Build(),
			Used:          NewQBBSBuilder(0).Add(v.Used).Build(),
			Free:          NewQBBSBuilder(0).Add(v.Free).Build(),
			UsedPercent:   fmt.Sprintf("%.1f%%", v.UsedPercent),
			UsedTotal:     NewQBBSBuilder('/').Add(v.Used).Add(v.Total).Build(),
			UsedTotalFree: NewQBBSBuilder('/').Add(v.Used).Add(v.Total).Add(v.Free).Build(),
		},
		Swap: DTOMemoryT{
			Total:         NewQBBSBuilder(0).Add(v.SwapTotal).Build(),
			Used:          NewQBBSBuilder(0).Add(v.SwapUsed).Build(),
			Free:          NewQBBSBuilder(0).Add(v.SwapFree).Build(),
			UsedPercent:   fmt.Sprintf("%.1f%%", v.SwapUsedPercent),
			UsedTotal:     NewQBBSBuilder('/').Add(v.SwapUsed).Add(v.SwapTotal).Build(),
			UsedTotalFree: NewQBBSBuilder('/').Add(v.SwapUsed).Add(v.SwapTotal).Add(v.SwapFree).Build(),
		},
	}
}

// ============================ System dto ============================

type DTOSystem struct {
	Uptime       string `json:"uptime"`        // "3h-22m-5s"
	Idle         string `json:"idle"`          // "3h-22m-5s"
	Load1        string `json:"load1"`         // "0.42"
	Load5        string `json:"load5"`         // "0.38"
	Load15       string `json:"load15"`        // "0.32"
	RunningProcs int    `json:"running_procs"` // "3"
	TotalProcs   int    `json:"total_procs"`   // "187"
	SummaryProcs string `json:"summary_procs"` // "3/187"
}

func Domain2DTOSystem(v domain.SystemInfo) *DTOSystem {
	return &DTOSystem{
		Uptime:       formatDuration(v.Uptime, ":", true),
		Idle:         formatDuration(v.Idle, ":", true),
		Load1:        fmt.Sprintf("%.2f", v.Load1),
		Load5:        fmt.Sprintf("%.2f", v.Load5),
		Load15:       fmt.Sprintf("%.2f", v.Load15),
		RunningProcs: v.RunningProcs,
		TotalProcs:   v.TotalProcs,
		SummaryProcs: fmt.Sprintf("%d/%d", v.RunningProcs, v.TotalProcs),
	}
}

// ============================ Network dto ============================

type DTONetworkInterfaceIO struct {
	Name string `json:"name"` // "eth0"

	Bytes   IO[uint64] `json:"bytes"`   // "12.3GB"
	Packets IO[uint64] `json:"packets"` // "1.2M"

	BytesSpeed   IO[uint64] `json:"bytes_speed"`   // "120KB/s"
	PacketsSpeed IO[uint64] `json:"packets_speed"` // "133 pkt/s"

	Errors uint64 `json:"errors"` // "0"
	Drops  uint64 `json:"drops"`  // "0"
}

type DTONetworkIO struct {
	BytesTotal   IO[uint64]                       `json:"bytes_total"`
	PacketsTotal IO[uint64]                       `json:"packets_total"`
	ErrorsTotal  uint64                           `json:"errors_total"`
	DropsTotal   uint64                           `json:"drops_total"`
	Interfaces   map[string]DTONetworkInterfaceIO `json:"interfaces"`
}

func NewDTONetworkIO(ifCount int) *DTONetworkIO {
	return &DTONetworkIO{
		Interfaces: make(map[string]DTONetworkInterfaceIO, ifCount),
	}
}

func Domain2DTONetworkInterfaceIO(v domain.InterfacesIOMap) *DTONetworkIO {
	if len(v) == 0 {
		return NewDTONetworkIO(0)
	}

	dto := NewDTONetworkIO(len(v))

	var (
		pktTotalRx   = uint64(0)
		pktTotalTx   = uint64(0)
		bytesTotalRx = uint64(0)
		bytesTotalTx = uint64(0)
		errorsTotal  = uint64(0)
		dropTotal    = uint64(0)
	)

	for name, io := range v {
		pktTotalRx += io.PacketsTotal.RX
		pktTotalTx += io.PacketsTotal.TX
		bytesTotalRx += io.BytesTotal.RX
		bytesTotalTx += io.BytesTotal.TX
		errorsTotal += io.ErrPacketsTotal.Summary
		dropTotal += io.DropPacketsTotal.Summary

		iface := DTONetworkInterfaceIO{
			Name: name,

			Bytes:   NewIOBuilder(io.BytesTotal.RX, io.BytesTotal.TX).AutoUnits().Build(),
			Packets: NewIOBuilder(io.PacketsTotal.RX, io.PacketsTotal.TX).AutoMetricUnits().WithPostfix("Pkts").Build(),

			BytesSpeed:   NewIOBuilder(io.BytesPerSec.RX, io.BytesPerSec.TX).AutoUnitsPerSec().Build(),
			PacketsSpeed: NewIOBuilder(io.PacketsPerSec.RX, io.PacketsPerSec.TX).AutoMetricUnits().WithPostfix("Pkts/s").Build(),

			Errors: io.ErrPacketsTotal.Summary,
			Drops:  io.DropPacketsTotal.Summary,
		}

		dto.Interfaces[name] = iface
	}

	dto.BytesTotal = NewIOBuilder(bytesTotalRx, bytesTotalTx).AutoUnits().Build()
	dto.PacketsTotal = NewIOBuilder(pktTotalRx, pktTotalTx).AutoMetricUnits().WithPostfix("Pkts").Build()
	dto.ErrorsTotal = errorsTotal
	dto.DropsTotal = dropTotal

	return dto
}

// ============================ Disk systme and FS dto ============================

type DTOPartitionUsage struct {
	UsedPercent                  string `json:"used_percent"` // Used space percentage
	TotalBytes                   string `json:"total"`        // Total bytes
	UsedBytes                    string `json:"used"`         // Used bytes
	FreeBytes                    string `json:"free"`         // Free bytes
	UsedBytesTotalBytes          string `json:"used_total"`
	UsedBytesTotalBytesFreeBytes string `json:"used_total_free"`

	InodesUsedPercent string `json:"inodes_used_percent"` // Inode usage percentage
	InodesTotal       string `json:"inodes_total"`        // Total inodes
	InodesUsed        string `json:"inodes_used"`         // Used inodes
	InodesFree        string `json:"inodes_free"`         // Free inodes
}

type DTOPartition struct {
	Device        string             `json:"device"`         // Device path, e.g. "/dev/sda1"
	Mount         string             `json:"mount"`          // Mount point, e.g. "/"
	Filesystem    string             `json:"filesystem"`     // Filesystem type, e.g. "ext4"
	OptionsString string             `json:"options_string"` // Mount options
	Options       []string           `json:"options"`        // Mount options
	Usage         *DTOPartitionUsage `json:"usage"`
}

type DTOPartitions map[string]DTOPartition

func Domain2DTOPartitions(v domain.Partitions) *DTOPartitions {
	if len(v) == 0 {
		return &DTOPartitions{}
	}

	dto := make(DTOPartitions, len(v))

	for _, p := range v {
		part := DTOPartition{
			Device:        p.Device,
			Mount:         p.Mount,
			Filesystem:    p.Filesystem,
			OptionsString: strings.Join(p.Options, " "),
			Options:       p.Options,
		}

		if p.Usage != nil {
			part.Usage = &DTOPartitionUsage{
				UsedPercent:                  fmt.Sprintf("%.1f%%", p.Usage.UsedPercent),
				TotalBytes:                   NewQBBSBuilder(0).Add(p.Usage.TotalBytes).Build(),
				UsedBytes:                    NewQBBSBuilder(0).Add(p.Usage.UsedBytes).Build(),
				FreeBytes:                    NewQBBSBuilder(0).Add(p.Usage.FreeBytes).Build(),
				UsedBytesTotalBytes:          NewQBBSBuilder('/').Add(p.Usage.UsedBytes).Add(p.Usage.TotalBytes).Build(),
				UsedBytesTotalBytesFreeBytes: NewQBBSBuilder('/').Add(p.Usage.UsedBytes).Add(p.Usage.TotalBytes).Add(p.Usage.FreeBytes).Build(),
				InodesUsedPercent:            fmt.Sprintf("%.1f%%", p.Usage.InodesUsedPercent),
				InodesTotal:                  strconv.FormatUint(p.Usage.InodesTotal, 10),
				InodesUsed:                   strconv.FormatUint(p.Usage.InodesUsed, 10),
				InodesFree:                   strconv.FormatUint(p.Usage.InodesFree, 10),
			}
		}

		dto[strings.ReplaceAll(p.Device, "/", "&")] = part
	}

	return &dto
}

// DTODiskIO – DTO representation of DiskIO for output formatting layers.
// Contains string-formatted values and raw numeric fields.
type DTODiskIO struct {
	IopsInProgress uint64 `json:"iops_in_progress"`

	Ops       IO[uint64] `json:"ops"`
	MergedOps IO[uint64] `json:"merged_ops"`
	Bytes     IO[uint64] `json:"bytes"`

	OpsPerSec       IO[uint64] `json:"ops_per_sec"`
	MergedOpsPerSec IO[uint64] `json:"merged_ops_per_sec"`
	BytesPerSec     IO[uint64] `json:"bytes_per_sec"`

	Time IO[uint64] `json:"time"` // RX/TX in nanoseconds, string versions formatted in seconds

	IoTime     string `json:"io_time"`     // Total I/O time formatted
	WeightedIO string `json:"weighted_io"` // Weighted I/O time formatted
}

type DTODiskIOs map[string]DTODiskIO

/*
Domain2 - converts domain-level DiskIO to human-readable DTO (DTODiskIO).

Performs:
- Direct copy of raw counters.
- Automatic unit selection for bytes and metrics (K/M/G scaling).
- Optional “/s” postfix for per-second metrics.
- Conversion of time.Duration fields into readable "X.XXs" string form.
- Converts Duration into raw nanoseconds for DTO raw fields.
*/
func Domain2DTODiskIOs(d domain.DiskIOMap) *DTODiskIOs {
	s := make(DTODiskIOs, len(d))

	for dev, d := range d {
		dto := DTODiskIO{}

		// IopsInProgress - direct transfer of the in-flight operations counter
		dto.IopsInProgress = d.IopsInProgress

		// Ops - auto metric units (K/M/G)
		dto.Ops = NewIOBuilder(d.Ops.RX, d.Ops.TX).
			AutoMetricUnits().
			Build()

		// MergedOps - same as Ops
		dto.MergedOps = NewIOBuilder(d.MergedOps.RX, d.MergedOps.TX).
			AutoMetricUnits().
			Build()

		// Bytes - auto byte units (KiB / MiB / GiB)
		dto.Bytes = NewIOBuilder(d.Bytes.RX, d.Bytes.TX).
			AutoUnits().
			Build()

		// OpsPerSec - metric units + "/s"
		dto.OpsPerSec = NewIOBuilder(d.OpsPerSec.RX, d.OpsPerSec.TX).
			AutoMetricUnits().
			WithPostfix("/s").
			Build()

		// MergedOpsPerSec - same logic
		dto.MergedOpsPerSec = NewIOBuilder(d.MergedOpsPerSec.RX, d.MergedOpsPerSec.TX).
			AutoMetricUnits().
			WithPostfix("/s").
			Build()

		// BytesPerSec - auto byte units + "/s"
		dto.BytesPerSec = NewIOBuilder(d.BytesPerSec.RX, d.BytesPerSec.TX).
			AutoUnitsPerSec().
			Build()

		// Time - convert Duration → raw nanoseconds + "X.XXs" formatted text
		{
			rawRX := uint64(d.Time.RX) // nanoseconds
			rawTX := uint64(d.Time.TX)

			dto.Time.RawRX = rawRX
			dto.Time.RawTX = rawTX

			dto.Time.RX = fmt.Sprintf("%.2fs", d.Time.RX.Seconds())
			dto.Time.TX = fmt.Sprintf("%.2fs", d.Time.TX.Seconds())
		}

		// IoTime - formatted total I/O time ("X.XXs")
		dto.IoTime = fmt.Sprintf("%.2fs", d.IoTime.Seconds())

		// WeightedIO - formatted weighted time
		dto.WeightedIO = fmt.Sprintf("%.2fs", d.WeightedIO.Seconds())

		s[dev] = dto
	}

	return &s
}
