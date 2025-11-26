// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package hostfetchers

import (
	"context"
	"sync"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/utils/output"
	pshost "github.com/shirou/gopsutil/v4/host"
	psmem "github.com/shirou/gopsutil/v4/mem"
)

type SystemMon struct {
	data domain.SystemData
	mu   sync.RWMutex
	err  error
}

type CpuDataExporter interface {
	Data() (data domain.CpuLoad, err error)
}

func InitSystemMon(ctx context.Context, cpu CpuDataExporter, poolDuration time.Duration) *SystemMon {
	self := new(SystemMon)
	go self.updates(ctx, cpu, poolDuration)
	return self
}

func (mon *SystemMon) updates(ctx context.Context, cpu CpuDataExporter, poolDuration time.Duration) {
	var loadCpu = float64(0.0)

	tm := time.NewTicker(poolDuration)
	defer tm.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-tm.C:
			mon.mu.Lock()

			if ld, err := cpu.Data(); err == nil {
				loadCpu = ld.Average
			}

			ramData, err := psmem.VirtualMemoryWithContext(ctx)
			if err != nil {
				mon.err = err
				mon.mu.Unlock()
				continue
			}

			uptime, err := pshost.UptimeWithContext(ctx)
			if err != nil {
				mon.err = err
				mon.mu.Unlock()
				continue
			}

			mon.data = domain.SystemData{
				Cpu:           loadCpu,
				RAM:           output.UsageSizes(ramData.Used, ramData.Total),
				RAMUsage:      ramData.UsedPercent,
				Uptime:        output.FmtTime(time.Duration(uptime) * time.Second),
				UptimeSeconds: uptime,
			}

			mon.mu.Unlock()
		}
	}
}

func (mon *SystemMon) Fetch() (data domain.SystemData, err error) {
	mon.mu.RLock()
	defer mon.mu.RUnlock()
	return mon.data, nil
}
