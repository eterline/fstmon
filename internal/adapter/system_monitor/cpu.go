package systemmonitor

import (
	"context"
	"errors"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/utils/output"
	"github.com/eterline/fstmon/pkg/procf"
	pscpu "github.com/shirou/gopsutil/v4/cpu"
)

/*
hardwareMetricCPU - provides CPU hardware metrics, both static and dynamic.

	It can fetch information about CPU package and per-core metrics.
	Interval defines the averaging duration for dynamic metrics.
*/
type hardwareMetricCPU struct {
	interval time.Duration
}

/*
NewHardwareMetricCPU - creates a new hardwareMetricCPU instance.

	Interval specifies the time duration for CPU load averaging in CpuMetrics.
*/
func NewHardwareMetricCPU(interval time.Duration) *hardwareMetricCPU {
	return &hardwareMetricCPU{
		interval: interval,
	}
}

/*
CpuPackage - returns static CPU package information, including:

  - vendor

  - model name

  - microcode version

  - CPU flags

  - per-core static details (CoreID, PhysicalID, Cache)

    Returns an error if fetching CPU info fails or no cores are detected.
*/
func (hmc *hardwareMetricCPU) CpuPackage() (domain.CpuPackage, error) {
	info, err := procf.FetchCpuInfo()
	if err != nil {
		return domain.CpuPackage{}, err
	}

	if len(info.Cores) == 0 {
		return domain.CpuPackage{}, errors.New("failed to get cpu package data")
	}

	c0 := info.Cores[0]

	pkg := domain.CpuPackage{
		Vendor:    c0.VendorID,
		ModelName: c0.ModelName,
		Microcode: c0.Microcode,
		Flags:     c0.Flags,
		Cores:     make([]domain.CpuCoreInfo, 0, len(info.Cores)),
	}

	for _, c := range info.Cores {
		core := domain.CpuCoreInfo{
			CoreID:     int(c.CoreID),
			PhysicalID: int(c.PhysicalID),
			Siblings:   int(c.Siblings),
			CacheKB:    int(c.CacheSize),
		}
		pkg.Cores = append(pkg.Cores, core)
	}

	return pkg, nil
}

/*
CpuMetrics - returns dynamic CPU metrics, including:

  - per-core load

  - per-core frequency

  - average load/frequency across all cores

    Uses the configured interval for calculating CPU load percentages.
    Returns an error if fetching CPU info or CPU load fails,
    or if the number of cores and load entries do not match.
*/
func (hmc *hardwareMetricCPU) CpuMetrics(ctx context.Context) (domain.CpuMetrics, error) {
	cpuInfo, err := pscpu.InfoWithContext(ctx)
	if err != nil {
		return domain.CpuMetrics{}, err
	}

	loads, err := pscpu.PercentWithContext(ctx, hmc.interval, true)
	if err != nil {
		return domain.CpuMetrics{}, err
	}

	if len(cpuInfo) != len(loads) {
		return domain.CpuMetrics{}, errors.New("unequal load stats count and core information count")
	}

	metrics := domain.CpuMetrics{
		Cores: make([]domain.CpuCoreMetrics, len(cpuInfo)),
	}

	{
		f := func(idx int) float64 { return cpuInfo[idx].Mhz }
		metrics.Average.Load = output.AverageFloat(loads)
		metrics.Average.Frequency = output.AverageFloatArr(cpuInfo, f)
	}

	for i := range metrics.Cores {
		metrics.Cores[i].Load = loads[i]
		metrics.Cores[i].Frequency = cpuInfo[i].Mhz
	}

	return metrics, nil
}
