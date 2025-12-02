// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package system

import (
	"context"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/utils/usecase"
	"github.com/prometheus/procfs"
	diskPs "github.com/shirou/gopsutil/v4/disk"
)

type hardwareMetricPartitions struct {
	fs procfs.FS
}

func NewHardwareMetricPartitions(fs procfs.FS) *hardwareMetricPartitions {
	return &hardwareMetricPartitions{
		fs: fs,
	}
}

func (hmp *hardwareMetricPartitions) ScrapeDiskIO(ctx context.Context) (domain.DiskIOMap, error) {
	ct0, err := diskPs.IOCountersWithContext(ctx)
	if err != nil {
		return nil, ErrScrapeDiskIO.Wrap(err)
	}

	select {
	case <-ctx.Done():
		return domain.DiskIOMap{}, ErrScrapeDiskIO.Wrap(ctx.Err())
	case <-time.After(1 * time.Second):
	}

	ct1, err := diskPs.IOCountersWithContext(ctx)
	if err != nil {
		return nil, ErrScrapeDiskIO.Wrap(err)
	}

	data := make(domain.DiskIOMap, len(ct0))

	for iface, io := range ct1 {
		rxTime := usecase.MsToDuration(io.ReadTime)
		txTime := usecase.MsToDuration(io.WriteTime)

		data[iface] = domain.DiskIO{
			IopsInProgress: io.IopsInProgress,

			Ops:       domain.NewIO(io.ReadCount, io.WriteCount),
			MergedOps: domain.NewIO(io.MergedReadCount, io.MergedWriteCount),
			Bytes:     domain.NewIO(io.ReadBytes, io.WriteBytes),

			OpsPerSec:       domain.NewIO(io.ReadCount, io.WriteCount),
			MergedOpsPerSec: domain.NewIO(io.MergedReadCount, io.MergedWriteCount),
			BytesPerSec:     domain.NewIO(io.ReadBytes, io.WriteBytes),

			Time:       domain.NewIO(rxTime, txTime),
			IoTime:     usecase.MsToDuration(io.IoTime),
			WeightedIO: usecase.MsToDuration(io.WeightedIO),
		}
	}

	/* Subtract first snapshot to compute actual speed during the interval */
	for iface, io := range ct0 {
		c, ok := data[iface]
		if !ok {
			continue
		}

		c.OpsPerSec.DecRX(io.ReadCount)
		c.OpsPerSec.DecTX(io.WriteCount)

		c.MergedOpsPerSec.DecRX(io.MergedReadCount)
		c.MergedOpsPerSec.DecTX(io.MergedWriteCount)

		c.BytesPerSec.DecRX(io.ReadBytes)
		c.BytesPerSec.DecTX(io.WriteBytes)

		data[iface] = c
	}

	return data, nil
}

func (hmp *hardwareMetricPartitions) ScrapePartitions(ctx context.Context) (domain.Partitions, error) {
	prts, err := diskPs.PartitionsWithContext(ctx, true)
	if err != nil {
		return domain.Partitions{}, ErrScrapePartitions.Wrap(err)
	}

	data := make(domain.Partitions, len(prts))

	for i, v := range prts {
		part := domain.Partition{
			Device:     v.Device,
			Mount:      v.Mountpoint,
			Filesystem: v.Fstype,
			Options:    v.Opts,
			Usage:      nil,
		}

		usage, err := diskPs.UsageWithContext(ctx, v.Mountpoint)
		if err == nil {
			part.Usage = &domain.PartitionUsage{
				TotalBytes:        usage.Total,
				FreeBytes:         usage.Free,
				UsedBytes:         usage.Used,
				UsedPercent:       usage.UsedPercent,
				InodesTotal:       usage.InodesTotal,
				InodesFree:        usage.InodesFree,
				InodesUsed:        usage.InodesUsed,
				InodesUsedPercent: usage.InodesUsedPercent,
			}
		}

		data[i] = part
	}

	return data, nil
}
