package handlers

import (
	"context"
	"log/slog"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/common"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/convert"
)

type thernalHandlers struct {
	common.UnimplementedThermalServiceServer
	handlersGroup
}

func NewThernalHandlers(l *slog.Logger, s ActualStateStore) *thernalHandlers {
	return &thernalHandlers{
		handlersGroup: newHandlersGroup(l, s),
	}
}

func (nh *thernalHandlers) GetThermal(context.Context, *common.GetThermalRequest) (*common.ThermalResponse, error) {
	data, err := GetMetric[domain.ThermalMetricsMap](nh.store, "thermal")
	if err != nil {
		nh.log.Error("failed get thernal metrics", "error", err)
		return nil, err
	}

	res := convert.ThermalMetricsMapToResponse(data)
	return res, nil
}
