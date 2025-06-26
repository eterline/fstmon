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

func (hi *HostInfo) AverageLoad() (domain.AverageLoad, error) {

	res := domain.AverageLoad{
		Load1:  0.0,
		Load5:  0.0,
		Load15: 0.0,
		Procs:  "0/0",
	}

	data, err := procf.FetchProcLoadAvg()
	if err == nil {
		res = domain.AverageLoad{
			Load1:  data.Load1,
			Load5:  data.Load5,
			Load15: data.Load15,
			Procs:  data.RunningProcs,
		}
	}

	return res, nil
}

func (hi *HostInfo) TemperatureMap(ctx context.Context) (domain.TemperatureMap, error) {

	data, err := procf.Temperatures(ctx)
	if err != nil {
		return nil, err
	}

	temps := make(domain.TemperatureMap, len(data))
	for name, values := range data {
		temps[name] = values.Current
	}

	return temps, nil
}
