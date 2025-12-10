// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package convert

import (
	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/interface/grpc/flugel/common"
	"google.golang.org/protobuf/types/known/durationpb"
)

// ============================ CPU structures ============================

func cpuPackageDomainToDTO(d *domain.CpuPackage) *common.CpuPackage {
	if d == nil {
		return nil
	}

	cores := make([]*common.CpuCoreInfo, len(d.Cores))
	for i, c := range d.Cores {
		cores[i] = &common.CpuCoreInfo{
			PhysicalId: int32(c.PhysicalID),
			CoreId:     int32(c.CoreID),
			Siblings:   int32(c.Siblings),
			CacheKb:    int32(c.CacheKB),
		}
	}

	return &common.CpuPackage{
		Vendor:    d.Vendor,
		ModelName: d.ModelName,
		Microcode: d.Microcode,
		Flags:     append([]string(nil), d.Flags...), // копия
		Cores:     cores,
	}
}

func CpuPackageToResponse(d *domain.CpuPackage) *common.CpuPackageResponse {
	return &common.CpuPackageResponse{
		Cpu: cpuPackageDomainToDTO(d),
	}
}

func cpuCoreMetricsDomainToDTO(d domain.CpuCoreMetrics) *common.CpuCoreMetrics {
	return &common.CpuCoreMetrics{
		Load:      d.Load,
		Frequency: d.Frequency,
	}
}

func cpuMetricsDomainToDTO(d *domain.CpuMetrics) *common.CpuMetrics {
	if d == nil {
		return nil
	}

	cores := make([]*common.CpuCoreMetrics, len(d.Cores))
	for i, c := range d.Cores {
		cores[i] = cpuCoreMetricsDomainToDTO(c)
	}

	return &common.CpuMetrics{
		Average: cpuCoreMetricsDomainToDTO(d.Average),
		Cores:   cores,
	}
}

func CpuMetricsToResponse(d *domain.CpuMetrics) *common.CpuMetricsResponse {
	return &common.CpuMetricsResponse{
		Metrics: cpuMetricsDomainToDTO(d),
	}
}

// ============================ Networking structures ============================

func interfaceIODomainToDTO(i domain.InterfaceIO) *common.InterfaceIO {
	return &common.InterfaceIO{
		BytesTotal:        toIOUint64(i.BytesTotal),
		PacketsTotal:      toIOUint64(i.PacketsTotal),
		ErrorPacketsTotal: toIOUint64(i.ErrPacketsTotal),
		DropPacketsTotal:  toIOUint64(i.DropPacketsTotal),
		BytesPerSec:       toIOUint64(i.BytesPerSec),
		PacketsPerSec:     toIOUint64(i.PacketsPerSec),
	}
}

func interfacesIOMapDomainToDTO(m domain.InterfacesIOMap) *common.InterfacesIO {
	ma := make(map[string]*common.InterfaceIO, len(m))

	for name, i := range m {
		ma[name] = interfaceIODomainToDTO(i)
	}

	return &common.InterfacesIO{Interfaces: ma}
}

func InterfacesIOMapToResponse(m domain.InterfacesIOMap) *common.InterfacesIOResponse {
	return &common.InterfacesIOResponse{
		Data: interfacesIOMapDomainToDTO(m),
	}
}

// ============================ System structures ============================

func systemInfoToDTO(src *domain.SystemInfo) *common.SystemInfo {
	if src == nil {
		return nil
	}

	return &common.SystemInfo{
		Uptime:       durationpb.New(src.Uptime),
		Idle:         durationpb.New(src.Idle),
		Load1:        src.Load1,
		Load5:        src.Load5,
		Load15:       src.Load15,
		RunningProcs: int32(src.RunningProcs),
		TotalProcs:   int32(src.TotalProcs),
	}
}

func SystemInfoToResponse(src *domain.SystemInfo) *common.SystemInfoResponse {
	return &common.SystemInfoResponse{
		System: systemInfoToDTO(src),
	}
}

// ============================ Memory structures ============================

func MemoryMetricsToResponse(m *domain.MemoryMetrics) *common.MemoryMetricsResponse {
	return &common.MemoryMetricsResponse{
		Memory: &common.MemoryMetrics{
			Total:           m.Total,
			Available:       m.Available,
			Used:            m.Used,
			Free:            m.Free,
			SwapTotal:       m.SwapTotal,
			SwapAvailable:   m.SwapAvailable,
			SwapUsed:        m.SwapUsed,
			SwapFree:        m.SwapFree,
			UsedPercent:     m.UsedPercent,
			SwapUsedPercent: m.SwapUsedPercent,
		},
	}
}

// ============================ Thermal structures ============================

func ThermalMetricsMapToResponse(domainMap domain.ThermalMetricsMap) *common.ThermalResponse {
	sensors := make(map[string]*common.ThermalMetrics, len(domainMap))
	for name, m := range domainMap {
		sensors[name] = &common.ThermalMetrics{
			Current: m.Current,
			Max:     m.Max,
			Crit:    m.Crit,
		}
	}

	return &common.ThermalResponse{
		Metrics: &common.ThermalMetricsMap{
			Sensors: sensors,
		},
	}
}

// ============================ Storage structures ============================

func PartitionUsageToMessage(u *domain.PartitionUsage) *common.PartitionUsage {
	if u == nil {
		return nil
	}
	return partitionUsageToMessage(u)
}

func PartitionToMessage(p *domain.Partition) *common.Partition {
	if p == nil {
		return nil
	}
	return partitionToMessage(p)
}

func PartitionsToDTO(ps domain.Partitions) *common.Partitions {
	return partitionsToMessage(ps)
}

func PartitionsToMessage(ps domain.Partitions) *common.PartitionsResponse {
	return &common.PartitionsResponse{
		Partitions: PartitionsToDTO(ps),
	}
}

func partitionUsageToMessage(u *domain.PartitionUsage) *common.PartitionUsage {
	return &common.PartitionUsage{
		Total:             u.TotalBytes,
		Used:              u.UsedBytes,
		Free:              u.FreeBytes,
		UsedPercent:       u.UsedPercent,
		InodesTotal:       u.InodesTotal,
		InodesUsed:        u.InodesUsed,
		InodesFree:        u.InodesFree,
		InodesUsedPercent: u.InodesUsedPercent,
	}
}

func partitionToMessage(p *domain.Partition) *common.Partition {
	var usageMsg *common.PartitionUsage
	if p.Usage != nil {
		usageMsg = partitionUsageToMessage(p.Usage)
	}
	return &common.Partition{
		Device:     p.Device,
		Mount:      p.Mount,
		Filesystem: p.Filesystem,
		Options:    p.Options,
		Usage:      usageMsg,
	}
}

func partitionsToMessage(ps domain.Partitions) *common.Partitions {
	msgs := make([]*common.Partition, 0, len(ps))
	for _, p := range ps {
		msgs = append(msgs, partitionToMessage(&p))
	}
	return &common.Partitions{
		Partitions: msgs,
	}
}

func DiskIOToMessage(d domain.DiskIO) *common.DiskIO {
	return diskIOToMessage(d)
}

func DiskIOMapToMessage(m domain.DiskIOMap) *common.DiskIOMap {
	return diskIOMapToMessage(m)
}

func DiskIOMapResponseToMessage(m domain.DiskIOMap) *common.DiskIOMapResponse {
	return &common.DiskIOMapResponse{
		Io: diskIOMapToMessage(m),
	}
}

func diskIOToMessage(d domain.DiskIO) *common.DiskIO {
	return &common.DiskIO{
		IopsInProgress:  d.IopsInProgress,
		Ops:             toIOUint64(d.Ops),
		MergedOps:       toIOUint64(d.MergedOps),
		Bytes:           toIOUint64(d.Bytes),
		OpsPerSec:       toIOUint64(d.OpsPerSec),
		MergedOpsPerSec: toIOUint64(d.MergedOpsPerSec),
		BytesPerSec:     toIOUint64(d.BytesPerSec),
		Time:            toIODuration(d.Time),
		IoTime:          durationpb.New(d.IoTime),
		WeightedIo:      durationpb.New(d.WeightedIO),
	}
}

func diskIOMapToMessage(m domain.DiskIOMap) *common.DiskIOMap {
	disks := make(map[string]*common.DiskIO, len(m))
	for k, v := range m {
		disks[k] = diskIOToMessage(v)
	}
	return &common.DiskIOMap{
		Disks: disks,
	}
}
