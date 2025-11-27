package system

import (
	"context"
	"errors"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	psnet "github.com/shirou/gopsutil/v4/net"
)

/*
hardwareMetricNetwork  - provides network interface metrics for the host.

	It allows fetching per-interface counters, including bytes sent/received,
	packet counts, and error counts. Supports measuring speed over a short interval.
*/
type hardwareMetricNetwork struct{}

/*
NewHardwareMetricNetwork  - creates a new hardwareMetricNetwork instance.
*/
func NewHardwareMetricNetwork() *hardwareMetricNetwork {
	return &hardwareMetricNetwork{}
}

/*
ScrapeInterfacesIO - collects network metrics for all interfaces.

	It performs two snapshots of per-interface I/O counters with a 1-second interval
	to calculate approximate network speed.
*/
func (hmn *hardwareMetricNetwork) ScrapeInterfacesIO(ctx context.Context) (domain.InterfacesIO, error) {
	io0, err := psnet.IOCountersWithContext(ctx, true)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return domain.InterfacesIO{}, errors.New("measurement cancelled by context")
	case <-time.After(1 * time.Second):
	}

	io1, err := psnet.IOCountersWithContext(ctx, true)
	if err != nil {
		return nil, err
	}

	ioCounterMap := make(domain.InterfacesIO, len(io1))

	/* Initialize map with full counters and assume speed = full counters initially */
	for _, io := range io1 {
		ioCounterMap[io.Name] = domain.NetworkingIO{
			BytesFullRX:  io.BytesRecv,
			BytesFullTX:  io.BytesSent,
			BytesSpeedRX: io.BytesRecv,
			BytesSpeedTX: io.BytesSent,
			PacketsRx:    io.PacketsRecv,
			PacketsTx:    io.PacketsSent,
			ErrPacketsRx: io.Errin,
			ErrPacketsTx: io.Errout,
		}
	}

	/* Subtract first snapshot to compute actual speed during the interval */
	for i := range io0 {
		io := &io0[i]
		c, ok := ioCounterMap[io.Name]
		if !ok {
			continue
		}

		c.BytesSpeedRX -= io.BytesRecv
		c.BytesSpeedTX -= io.BytesSent
		ioCounterMap[io.Name] = c
	}

	return ioCounterMap, nil
}
