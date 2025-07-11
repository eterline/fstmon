package monitors

import (
	"context"
	"sync"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/utils/output"
	pscpu "github.com/shirou/gopsutil/v4/cpu"
)

type CpuLoadMonitoring struct {
	data domain.CpuLoad
	mu   sync.RWMutex
	err  error
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
		mon.mu.Lock()
		if err != nil {
			mon.err = err
			mon.mu.Unlock()
			continue
		}

		mon.err = nil
		mon.data.Cores = loads
		mon.data.Average = output.AverageFloat(loads)
		mon.mu.Unlock()
	}
}

func (mon *CpuLoadMonitoring) Data() (data domain.CpuLoad, err error) {
	mon.mu.RLock()
	defer mon.mu.RUnlock()

	if mon.err != nil {
		return domain.CpuLoad{}, err
	}

	return mon.data, nil
}
