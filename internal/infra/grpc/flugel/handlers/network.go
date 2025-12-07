package handlers

import (
	"context"
	"log/slog"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/common"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/convert"
)

type networkHandlers struct {
	common.UnimplementedCpuServiceServer
	handlersGroup
}

func NewNetworkHandlers(l *slog.Logger, s ActualStateStore) *networkHandlers {
	return &networkHandlers{
		handlersGroup: newHandlersGroup(l, s),
	}
}

func (nh *networkHandlers) GetInterfacesIO(context.Context, *common.GetInterfacesIORequest) (*common.InterfacesIOResponse, error) {
	data, err := GetMetric[domain.InterfacesIOMap](nh.store, "net_io")
	if err != nil {
		nh.log.Error("failed get interfaces io", "error", err)
		return nil, err
	}

	res := convert.InterfacesIOMapToResponse(data)
	return res, nil
}
