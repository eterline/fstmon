package system

import (
	"context"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/prometheus/procfs"
)

/*
hardwareMetricNetwork  - provides network interface metrics for the host.

	It allows fetching per-interface counters, including bytes sent/received,
	packet counts, and error counts. Supports measuring speed over a short interval.
*/
type hardwareMetricNetwork struct {
	fs procfs.FS
}

/*
NewHardwareMetricNetwork  - creates a new hardwareMetricNetwork instance.
*/
func NewHardwareMetricNetwork(fs procfs.FS) *hardwareMetricNetwork {
	return &hardwareMetricNetwork{
		fs: fs,
	}
}

/*
ScrapeInterfacesIO - collects network metrics for all interfaces.

	It performs two snapshots of per-interface I/O counters with a 1-second interval
	to calculate approximate network speed.
*/
func (hmn *hardwareMetricNetwork) ScrapeInterfacesIO(ctx context.Context) (domain.MetricWrapper[domain.InterfacesIO], error) {
	io0, err := hmn.fs.NetDev()
	if err != nil {
		return domain.EmptyWrapMetric[domain.InterfacesIO](),
			ErrScrapeInterfacesIO.Wrap(err)
	}

	select {
	case <-ctx.Done():
		return domain.EmptyWrapMetric[domain.InterfacesIO](),
			ErrScrapeInterfacesIO.Wrap(ctx.Err())
	case <-time.After(1 * time.Second):
	}

	io1, err := hmn.fs.NetDev()
	if err != nil {
		return domain.EmptyWrapMetric[domain.InterfacesIO](),
			ErrScrapeInterfacesIO.Wrap(err)
	}

	select {
	case <-ctx.Done():
		return domain.EmptyWrapMetric[domain.InterfacesIO](),
			ErrScrapeMemoryMetrics.Wrap(ctx.Err())
	default:
	}

	ioCounterMap := make(domain.InterfacesIO, len(io1))

	/* Initialize map with full counters and assume speed = full counters initially */
	for _, v := range io1 {
		ioCounterMap[v.Name] = domain.NetworkingIO{
			BytesFullRX:   v.RxBytes,
			BytesFullTX:   v.TxBytes,
			BytesPerSec:   domain.NewSpeedIO(v.RxBytes, v.TxBytes),
			PacketsRx:     v.RxPackets,
			PacketsTx:     v.TxPackets,
			PacketsPerSec: domain.NewSpeedIO(v.RxPackets, v.TxPackets),
			ErrPacketsRx:  v.RxErrors,
			ErrPacketsTx:  v.TxErrors,
		}
	}

	/* Subtract first snapshot to compute actual speed during the interval */
	for _, v := range io0 {
		c, ok := ioCounterMap[v.Name]
		if !ok {
			continue
		}

		c.BytesPerSec.RX -= v.RxBytes
		c.BytesPerSec.TX -= v.TxBytes

		c.PacketsPerSec.RX -= v.RxPackets
		c.PacketsPerSec.TX -= v.TxPackets

		ioCounterMap[v.Name] = c
	}

	return domain.WrapMetric(ioCounterMap), nil
}
