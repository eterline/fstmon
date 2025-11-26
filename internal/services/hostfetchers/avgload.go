// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package hostfetchers

import (
	"context"
	"sync"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/pkg/procf"
)

type AverageLoadMon struct {
	data domain.AverageLoad
	mu   sync.RWMutex
	err  error
}

func InitAverageLoadMon(ctx context.Context, poolDuration time.Duration) *AverageLoadMon {
	self := new(AverageLoadMon)
	go self.updates(ctx, poolDuration)
	return self
}

func (mon *AverageLoadMon) updates(ctx context.Context, poolDuration time.Duration) {
	var update = func() {
		mon.mu.Lock()

		data, err := procf.FetchProcLoadAvg()
		if err != nil {
			mon.err = err
			mon.mu.Unlock()
			return
		}

		mon.data = domain.AverageLoad{
			Load1:  data.Load1,
			Load5:  data.Load5,
			Load15: data.Load15,
			Procs:  data.RunningProcs,
		}
		mon.mu.Unlock()
	}

	tm := time.NewTicker(poolDuration)
	defer tm.Stop()
	update()

	for {
		select {
		case <-ctx.Done():
			return

		case <-tm.C:
			update()
		}
	}
}

func (mon *AverageLoadMon) Fetch() (data domain.AverageLoad, err error) {
	mon.mu.RLock()
	defer mon.mu.RUnlock()

	if mon.err != nil {
		return domain.AverageLoad{}, err
	}

	return mon.data, nil
}
