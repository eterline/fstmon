package system

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/prometheus/procfs"
)

type hardwareMetricSystem struct {
	fs procfs.FS
}

func NewHardwareMetricSystem(fs procfs.FS) *hardwareMetricSystem {
	return &hardwareMetricSystem{
		fs: fs,
	}
}

func (hms *hardwareMetricSystem) ScrapeSystemInfo(ctx context.Context) (domain.SystemInfo, error) {
	avg, err := hms.fs.LoadAvg()
	if err != nil {
		return domain.SystemInfo{},
			ErrScrapeSystemInfo.Wrap(err)
	}

	select {
	case <-ctx.Done():
		return domain.SystemInfo{},
			ErrScrapeSystemInfo.Wrap(ctx.Err())
	default:
	}

	uptime, err := hms.readProfsUptime()
	if err != nil {
		return domain.SystemInfo{},
			ErrScrapeSystemInfo.Wrap(err)
	}

	select {
	case <-ctx.Done():
		return domain.SystemInfo{},
			ErrScrapeSystemInfo.Wrap(ctx.Err())
	default:
	}

	r, t := hms.readRunProcs()

	data := domain.SystemInfo{
		Uptime:       uptime.Uptime,
		Idle:         uptime.Idle,
		Load1:        avg.Load1,
		Load5:        avg.Load5,
		Load15:       avg.Load15,
		RunningProcs: r,
		TotalProcs:   t,
	}

	return data, nil
}

type Uptime struct {
	Uptime time.Duration
	Idle   time.Duration
}

func (hms *hardwareMetricSystem) readProfsUptime() (Uptime, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return Uptime{}, fmt.Errorf("failed read uptime: %w", err)
	}

	var u, i float64

	_, err = fmt.Sscanf(string(data), "%f %f", &u, &i)
	if err != nil {
		return Uptime{}, fmt.Errorf("failed scan fields uptime: %w", err)
	}

	return Uptime{
		Uptime: time.Duration(u * float64(time.Second)),
		Idle:   time.Duration(i * float64(time.Second)),
	}, nil
}

func (hms *hardwareMetricSystem) readRunProcs() (run, total int) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, 0
	}

	var load1, load5, load15 float64
	var r, t int

	_, err = fmt.Sscanf(
		string(data), "%f %f %f %d/%d",
		&load1, &load5, &load15, &r, &t,
	)
	if err != nil {
		return 0, 0
	}

	return r, t
}
