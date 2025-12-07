package handlers

import (
	"context"
	"log/slog"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/common"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/convert"
)

type systemHandlers struct {
	common.UnimplementedSystemServiceServer
	handlersGroup
}

func NewSystemHandlers(l *slog.Logger, s ActualStateStore) *systemHandlers {
	return &systemHandlers{
		handlersGroup: newHandlersGroup(l, s),
	}
}

func (nh *systemHandlers) GetSystemInfo(context.Context, *common.GetSystemInfoRequest) (*common.SystemInfoResponse, error) {
	data, err := GetMetric[domain.SystemInfo](nh.store, "system")
	if err != nil {
		nh.log.Error("failed get system info", "error", err)
		return nil, err
	}

	res := convert.SystemInfoToResponse(&data)
	return res, nil
}
