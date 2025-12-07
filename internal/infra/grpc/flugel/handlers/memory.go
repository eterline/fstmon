package handlers

import (
	"context"
	"log/slog"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/common"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/convert"
)

type memoryHandlers struct {
	common.UnimplementedMemoryServiceServer
	handlersGroup
}

func NewMemoryHandlers(l *slog.Logger, s ActualStateStore) *memoryHandlers {
	return &memoryHandlers{
		handlersGroup: newHandlersGroup(l, s),
	}
}

func (nh *memoryHandlers) GetMemoryMetrics(context.Context, *common.GetMemoryMetricsRequest) (*common.MemoryMetricsResponse, error) {
	data, err := GetMetric[domain.MemoryMetrics](nh.store, "memory")
	if err != nil {
		nh.log.Error("failed get memory metrics", "error", err)
		return nil, err
	}

	res := convert.MemoryMetricsToResponse(&data)
	return res, nil
}
