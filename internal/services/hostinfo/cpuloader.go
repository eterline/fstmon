package hostinfo

import (
	"context"
	"sync"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	pscpu "github.com/shirou/gopsutil/v4/cpu"
)

func runLoad(ctx context.Context, dur time.Duration) Load {
	stats, err := pscpu.PercentWithContext(ctx, dur, true)
	if err != nil {
		return Load{}
	}

	return Load{
		Cores: stats,
	}
}

type Load struct {
	Cores []float64
}

func (l Load) Average() float64 {

	count := len(l.Cores)
	if count == 0 {
		return 0.0
	}

	sum := float64(0)
	for _, value := range l.Cores {
		sum += value
	}

	return sum / float64(count)
}

type CpuLoader struct {
	loadsMap map[time.Duration]Load
	mu       sync.RWMutex
}

func InitCpuLoader(ctx context.Context) *CpuLoader {

	self := &CpuLoader{
		loadsMap: map[time.Duration]Load{
			5 * time.Second:  {},
			10 * time.Second: {},
			30 * time.Second: {},
		},
	}

	go self.monitoring(ctx)

	return self
}

func (ld *CpuLoader) monitoring(ctx context.Context) {
	for duration := range ld.loadsMap {
		d := duration

		go func() {
			ticker := time.NewTicker(d)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					load := runLoad(ctx, d)
					ld.mu.Lock()
					ld.loadsMap[d] = load
					ld.mu.Unlock()
				}
			}
		}()
	}
}

func (ld *CpuLoader) CpuLoad() (domain.CpuLoad, error) {
	ld.mu.RLock()
	defer ld.mu.RUnlock()

	frames := make(map[string]domain.CoreLoad, len(ld.loadsMap))

	for frame, load := range ld.loadsMap {
		frames[frame.String()] = domain.CoreLoad{
			Average: load.Average(),
			Cores:   load.Cores,
		}
	}

	return domain.CpuLoad{
		Frames: frames,
	}, nil
}
