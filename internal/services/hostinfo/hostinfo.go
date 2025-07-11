package hostinfo

import (
	"context"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/pkg/procf"
)

type HostInfo struct{}

func InitHostInfo() *HostInfo {
	return new(HostInfo)
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
