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
	"github.com/eterline/fstmon/pkg/procf"
)

type PartUseMon struct {
	data domain.PartsUsages
	mu   sync.RWMutex
	err  error
}

func InitPartUseMon(ctx context.Context, poolDuration time.Duration) *PartUseMon {
	self := new(PartUseMon)
	go self.updates(ctx, poolDuration)
	return self
}

func (mon *PartUseMon) updates(ctx context.Context, poolDuration time.Duration) {
	var update = func() {
		mon.mu.Lock()
		partArr, err := procf.FetchPartitions()
		if err != nil {
			mon.err = err
			mon.mu.Unlock()
			return
		}

		useMap := domain.PartsUsages{}

		for _, use := range extractUsage(partArr.Main) {
			useMap[use.Name] = use
		}

		for _, use := range extractUsage(partArr.LoopBack) {
			useMap[use.Name] = use
		}

		for _, use := range extractUsage(partArr.Mount) {
			useMap[use.Name] = use
		}

		for _, use := range extractUsage(partArr.Docker) {
			useMap[use.Name] = use
		}

		mon.data = useMap
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

func (mon *PartUseMon) Fetch() (data domain.PartsUsages, err error) {
	mon.mu.RLock()
	defer mon.mu.RUnlock()

	if mon.err != nil {
		return nil, err
	}

	return mon.data, nil
}

func extractUsage(parts []procf.Partition) []domain.PartUse {
	useArr := make([]domain.PartUse, 0)

	for _, part := range parts {
		if part.Usage == nil {
			continue
		}

		use := domain.PartUse{
			Name:    part.Name(),
			Size:    output.SizeString(part.Usage.Total),
			Use:     output.SizeString(part.Usage.Used),
			Percent: int32(part.Usage.UsedPercent),
		}

		useArr = append(useArr, use)
	}

	return useArr
}
