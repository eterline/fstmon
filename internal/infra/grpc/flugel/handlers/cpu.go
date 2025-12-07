package handlers

import (
	"context"
	"log/slog"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/common"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/convert"
)

type cpuHandlers struct {
	common.UnimplementedCpuServiceServer
	handlersGroup
}

func NewCpuHandlers(l *slog.Logger, s ActualStateStore) *cpuHandlers {
	return &cpuHandlers{
		handlersGroup: newHandlersGroup(l, s),
	}
}

func (cs *cpuHandlers) GetCpuInfo(ctx context.Context, r *common.GetCpuInfoRequest) (*common.CpuPackageResponse, error) {
	data, err := GetMetric[domain.CpuPackage](cs.store, "cpu")
	if err != nil {
		cs.log.Error("failed get cpu info", "error", err)
		return nil, err
	}

	res := convert.CpuPackageToResponse(&data)
	return res, nil
}

func (cs *cpuHandlers) GetCpuMetrics(ctx context.Context, r *common.GetCpuMetricsRequest) (*common.CpuMetricsResponse, error) {
	data, err := GetMetric[domain.CpuMetrics](cs.store, "cpu_usage")
	if err != nil {
		cs.log.Error("failed get cpu metrics", "error", err)
		return nil, err
	}

	res := convert.CpuMetricsToResponse(&data)
	return res, nil
}
