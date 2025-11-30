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

func (hmp *hardwareMetricPartitions) ScrapePartitionsInfo(ctx context.Context) (domain.PartitionsInfo, error) {
	prts, err := diskPs.PartitionsWithContext(ctx, true)
	if err != nil {
		return domain.PartitionsInfo{},
			ErrScrapePartitionsInfo.Wrap(err)
	}

	data := make(domain.PartitionsInfo, len(prts))

	for i, v := range prts {
		data[i] = domain.PartitionInfo{
			Device:     v.Device,
			Mount:      v.Mountpoint,
			Filesystem: v.Fstype,
			Options:    v.Opts,
		}
	}

	return data, nil
}

func (hmp *hardwareMetricPartitions) ScrapePartitionIO(ctx context.Context) (domain.PartitionIOMap, error) {
	ct0, err := diskPs.IOCountersWithContext(ctx)
	if err != nil {
		return nil, ErrScrapePartitionsIO.Wrap(err)
	}

	select {
	case <-ctx.Done():
		return domain.PartitionIOMap{},
			ErrScrapePartitionsIO.Wrap(ctx.Err())
	case <-time.After(1 * time.Second):
	}

	ct1, err := diskPs.IOCountersWithContext(ctx)
	if err != nil {
		return nil, ErrScrapePartitionsIO.Wrap(err)
	}

	data := make(domain.PartitionIOMap, len(ct0))

	for iface, io := range ct1 {
		rxTime := usecase.MsToDuration(io.ReadTime)
		txTime := usecase.MsToDuration(io.WriteTime)

		data[iface] = domain.PartitionIO{
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
