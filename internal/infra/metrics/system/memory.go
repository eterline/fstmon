// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package system

import (
	"context"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/prometheus/procfs"
)

type hardwareMetricMemory struct {
	fs procfs.FS
}

func NewHardwareMetricMemory(fs procfs.FS) *hardwareMetricMemory {
	return &hardwareMetricMemory{
		fs: fs,
	}
}

func (hmm *hardwareMetricMemory) ScrapeMemoryMetrics(ctx context.Context) (domain.MemoryMetrics, error) {
	mem, err := hmm.fs.Meminfo()
	if err != nil {
		return domain.MemoryMetrics{},
			ErrScrapeMemoryMetrics.Wrap(err)
	}

	if err := ctx.Err(); err != nil {
		return domain.MemoryMetrics{},
			ErrScrapeMemoryMetrics.Wrap(ctx.Err())
	}

	total := uwPtr(mem.MemTotal)
	avail := uwPtr(mem.MemAvailable)
	used := total - avail

	stotal := uwPtr(mem.SwapTotalBytes)
	sfree := uwPtr(mem.SwapFreeBytes)
	sused := stotal - sfree

	data := domain.MemoryMetrics{
		Total:       total,
		Available:   avail,
		Used:        used,
		Free:        uwPtr(mem.MemFree),
		UsedPercent: usedPercent[uint64, float64](avail, total),

		SwapTotal:       stotal,
		SwapFree:        sfree,
		SwapUsed:        sused,
		SwapUsedPercent: usedPercent[uint64, float64](sused, stotal),
	}

	return data, nil
}
