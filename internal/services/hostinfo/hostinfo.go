package hostinfo

import (
	"context"
	"fmt"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/utils/output"
	"github.com/eterline/fstmon/pkg/procf"
	pscpu "github.com/shirou/gopsutil/v4/cpu"
	psmem "github.com/shirou/gopsutil/v4/mem"
)

type HostInfo struct{}

func InitHostInfo() *HostInfo {
	return &HostInfo{}
}

func (hi *HostInfo) System(ctx context.Context) (domain.SystemData, error) {

	loads, err := pscpu.PercentWithContext(ctx, 5*time.Second, true)
	if err != nil {
		return domain.SystemData{}, err
	}

	stat, err := psmem.VirtualMemory()
	if err != nil {
		return domain.SystemData{}, err
	}

	uptime, err := procf.FetchProcUptime()
	if err != nil {
		return domain.SystemData{}, err
	}

	temp, err := cpuTemperature(ctx)
	if err != nil {
		return domain.SystemData{}, err
	}

	data := domain.SystemData{
		CpuTemp: output.CelsiusString(temp),
		Uptime:  output.FmtTime(uptime.Uptime),
		Cpu:     fmt.Sprintf("%.0f%%", output.AverageFloat(loads)),
		Memory:  fmt.Sprintf("%.0f%%", stat.UsedPercent),
	}

	return data, nil
}

func (hi *HostInfo) Processes(ctx context.Context) (domain.ProcessesData, error) {
	return domain.ProcessesData{}, nil
}

func (hi *HostInfo) Networking(ctx context.Context) (domain.InterfacesData, error) {
	data := domain.InterfacesData{}

	err := calcRxTx(ctx, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (hi *HostInfo) PartUse(ctx context.Context) (domain.PartsUsages, error) {

	partArr, err := procf.FetchPartitions()
	if err != nil {
		return nil, err
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

	return useMap, nil
}
