// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package monitors

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/utils/output"
	"github.com/eterline/fstmon/internal/utils/process"
	pscpu "github.com/shirou/gopsutil/v4/cpu"
)

type CpuLoadMonitoring struct {
	data domain.CpuLoad
	mu   sync.RWMutex
	err  process.ErrorHolder
}

func InitCpuLoadMon(ctx context.Context, poolDur time.Duration) *CpuLoadMonitoring {
	self := new(CpuLoadMonitoring)
	go self.updates(ctx, poolDur)
	return self
}

func (mon *CpuLoadMonitoring) updates(ctx context.Context, poolDur time.Duration) {
	for {
		if ctx.Err() != nil {
			return
		}

		loads, err := pscpu.PercentWithContext(ctx, poolDur, true)
		if err != nil {
			mon.err.SetError(err)
			continue
		}

		cpuInfo, err := pscpu.InfoWithContext(ctx)
		if err != nil {
			mon.err.SetError(err)
			continue
		}

		mon.err.ClearError()

		mon.mu.Lock()
		mon.data.Average = output.AverageFloat(loads)
		mon.data.Cores = sortCores(cpuInfo, loads)
		mon.mu.Unlock()
	}
}

func (mon *CpuLoadMonitoring) Data() (data domain.CpuLoad, err error) {

	if err, ok := mon.err.Err(); ok {
		return domain.CpuLoad{}, err
	}

	mon.mu.RLock()
	defer mon.mu.RUnlock()
	return mon.data, nil
}

func sortCores(stats []pscpu.InfoStat, loads []float64) []domain.CpuCore {

	var (
		statsLen = len(stats)
		statsMap = make(map[int]pscpu.InfoStat, statsLen)
		cores    = make([]domain.CpuCore, statsLen)
	)

	for _, v := range stats {
		id, err := strconv.Atoi(v.CoreID)
		if err == nil {
			statsMap[id] = v
		}
	}

	for i := range statsLen {
		cores[i] = domain.CpuCore{
			Load:      loads[i],
			Frequency: statsMap[i].Mhz,
		}
	}

	return cores
}
