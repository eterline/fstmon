package handlers

import (
	"context"
	"log/slog"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/common"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/convert"
)

type machineInfohandlers struct {
	common.UnimplementedMachineInfoServiceServer
	log   *slog.Logger
	store ActualStateStore
}

func NewMachineInfohandlers(l *slog.Logger, s ActualStateStore) *machineInfohandlers {
	return &machineInfohandlers{
		log:   l,
		store: s,
	}
}

// ==========================

func (cs *machineInfohandlers) GetCpuInfo(ctx context.Context, r *common.GetCpuInfoRequest) (*common.CpuPackageResponse, error) {
	data, err := GetMetric[domain.CpuPackage](cs.store, "cpu")
	if err != nil {
		cs.log.Error("failed get cpu info", "error", err)
		return nil, err
	}

	res := convert.CpuPackageToResponse(&data)
	return res, nil
}

func (cs *machineInfohandlers) GetCpuMetrics(ctx context.Context, r *common.GetCpuMetricsRequest) (*common.CpuMetricsResponse, error) {
	data, err := GetMetric[domain.CpuMetrics](cs.store, "cpu_usage")
	if err != nil {
		cs.log.Error("failed get cpu metrics", "error", err)
		return nil, err
	}

	res := convert.CpuMetricsToResponse(&data)
	return res, nil
}

// ==========================

func (nh *machineInfohandlers) GetMemoryMetrics(context.Context, *common.GetMemoryMetricsRequest) (*common.MemoryMetricsResponse, error) {
	data, err := GetMetric[domain.MemoryMetrics](nh.store, "memory")
	if err != nil {
		nh.log.Error("failed get memory metrics", "error", err)
		return nil, err
	}

	res := convert.MemoryMetricsToResponse(&data)
	return res, nil
}

// ==========================

func (nh *machineInfohandlers) GetInterfacesIO(context.Context, *common.GetInterfacesIORequest) (*common.InterfacesIOResponse, error) {
	data, err := GetMetric[domain.InterfacesIOMap](nh.store, "net_io")
	if err != nil {
		nh.log.Error("failed get interfaces io", "error", err)
		return nil, err
	}

	res := convert.InterfacesIOMapToResponse(data)
	return res, nil
}

// ==========================

func (nh *machineInfohandlers) GetPartitions(context.Context, *common.GetPartitionsRequest) (*common.PartitionsResponse, error) {
	data, err := GetMetric[domain.Partitions](nh.store, "partitions")
	if err != nil {
		nh.log.Error("failed get partitions", "error", err)
		return nil, err
	}

	res := convert.PartitionsToMessage(data)
	return res, nil
}

func (nh *machineInfohandlers) GetDiskIO(context.Context, *common.GetDiskIORequest) (*common.DiskIOMapResponse, error) {
	data, err := GetMetric[domain.DiskIOMap](nh.store, "disk_io")
	if err != nil {
		nh.log.Error("failed get disk io", "error", err)
		return nil, err
	}

	res := convert.DiskIOMapResponseToMessage(data)
	return res, nil
}

// ==========================

func (nh *machineInfohandlers) GetSystemInfo(context.Context, *common.GetSystemInfoRequest) (*common.SystemInfoResponse, error) {
	data, err := GetMetric[domain.SystemInfo](nh.store, "system")
	if err != nil {
		nh.log.Error("failed get system info", "error", err)
		return nil, err
	}

	res := convert.SystemInfoToResponse(&data)
	return res, nil
}

// ==========================

func (nh *machineInfohandlers) GetThermal(context.Context, *common.GetThermalRequest) (*common.ThermalResponse, error) {
	data, err := GetMetric[domain.ThermalMetricsMap](nh.store, "thermal")
	if err != nil {
		nh.log.Error("failed get thernal metrics", "error", err)
		return nil, err
	}

	res := convert.ThermalMetricsMapToResponse(data)
	return res, nil
}
