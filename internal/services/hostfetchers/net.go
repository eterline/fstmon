package hostfetchers

import (
	"context"
	"sync"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/utils/output"
	psnet "github.com/shirou/gopsutil/v4/net"
)

type NetworkMon struct {
	data domain.InterfacesData
	mu   sync.RWMutex
	err  error
}

func InitNetworkMon(ctx context.Context) *NetworkMon {
	self := new(NetworkMon)
	go self.updates(ctx)
	return self
}

func (mon *NetworkMon) updates(ctx context.Context) {

	const (
		poolDur = 5 * time.Second
	)

	tm := time.NewTicker(poolDur)
	defer tm.Stop()

	for {
		stat1, err := psnet.IOCountersWithContext(ctx, true)
		if err != nil {
			mon.mu.Lock()
			mon.err = err
			mon.mu.Unlock()
			return
		}

		tm.Reset(poolDur)

		select {
		case <-ctx.Done():
			return
		case <-tm.C:
		}

		stat2, err := psnet.IOCountersWithContext(ctx, true)
		if err != nil {
			mon.mu.Lock()
			mon.err = err
			mon.mu.Unlock()
			return
		}

		mon.mu.Lock()
		mon.data = make(domain.InterfacesData, len(stat2))
		for _, cntr2 := range stat2 {

			if cntr1, ok := findCounter(cntr2.Name, stat1); ok {

				rx, tx := netSpeed(cntr1, cntr2, poolDur)

				mon.data[cntr2.Name] = domain.NetworkingData{
					FullRX:  output.SizeString(cntr2.BytesRecv),
					FullTX:  output.SizeString(cntr2.BytesSent),
					SpeedRX: output.SpeedString(rx),
					SpeedTX: output.SpeedString(tx),
				}
			}
		}
		mon.mu.Unlock()
	}
}

func (mon *NetworkMon) Fetch() (data domain.InterfacesData, err error) {
	mon.mu.RLock()
	defer mon.mu.RUnlock()

	if mon.err != nil {
		return nil, err
	}

	return mon.data, nil
}

func findCounter(name string, stat []psnet.IOCountersStat) (psnet.IOCountersStat, bool) {
	for _, st := range stat {
		if name == st.Name {
			return st, true
		}
	}
	return psnet.IOCountersStat{}, true
}

func netSpeed(c1, c2 psnet.IOCountersStat, dur time.Duration) (rx, tx uint64) {
	rx = (c2.BytesRecv - c1.BytesRecv) / uint64(dur.Seconds())
	tx = (c2.BytesSent - c1.BytesSent) / uint64(dur.Seconds())
	return rx, tx
}
