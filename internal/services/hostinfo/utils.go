// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package hostinfo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/utils/output"
	"github.com/eterline/fstmon/pkg/procf"
	psnet "github.com/shirou/gopsutil/v4/net"
)

func cpuTemperature(ctx context.Context) (float64, error) {

	data, err := procf.Temperatures(ctx)
	if err != nil {
		return 0, err
	}

	count := 0
	temp := float64(0)
	contains := func(value string) bool {
		return strings.Contains(value, "coretemp") ||
			strings.Contains(value, "k10temp")
	}

	for key, value := range data {
		if !contains(key) {
			continue
		}
		count++
		temp += value.Current
	}

	if count == 0 {
		return 0, nil
	}

	return (temp / float64(count)), nil
}

func calcRxTx(ctx context.Context, data *domain.InterfacesData) error {

	counters1, err := psnet.IOCountersWithContext(ctx, true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("measurement cancelled by context")
	case <-time.After(1 * time.Second):
	}

	counters2, err := psnet.IOCountersWithContext(ctx, true)
	if err != nil {
		return err
	}

	var find = func(name string) *psnet.IOCountersStat {
		for _, c := range counters2 {
			if c.Name == name {
				return &c
			}
		}
		return nil
	}

	dataTmp := make(domain.InterfacesData, len(counters2))

	for _, counter := range counters1 {
		counter2 := find(counter.Name)
		if counter2 == nil {
			continue
		}

		dataTmp[counter.Name] = domain.NetworkingData{
			FullRX:  output.SizeString(counter2.BytesRecv),
			FullTX:  output.SizeString(counter2.BytesSent),
			SpeedRX: output.SpeedString(counter2.BytesRecv - counter.BytesRecv),
			SpeedTX: output.SpeedString(counter2.BytesSent - counter.BytesSent),
		}
	}

	*data = dataTmp
	return nil
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
