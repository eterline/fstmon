package system

import (
	"context"

	"github.com/eterline/fstmon/internal/domain"
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

func (hmp *hardwareMetricPartitions) ScrapePartitionsInfo(ctx context.Context) (domain.MetricWrapper[domain.PartitionsInfo], error) {
	prts, err := diskPs.PartitionsWithContext(ctx, true)
	if err != nil {
		return domain.EmptyWrapMetric[domain.PartitionsInfo](),
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

	return domain.WrapMetric(data), nil
}

// func (hmp *hardwareMetricPartitions) ScrapePartitionMetrics(ctx context.Context) (domain.PartitionMetricsMap, error) {
// 	counters, err := diskPs.IOCountersWithContext(ctx)
// 	if err != nil {
// 		return nil, ErrScrapePartitionsMetrics.Wrap(err)
// 	}

// 	return nil, nil
// }
