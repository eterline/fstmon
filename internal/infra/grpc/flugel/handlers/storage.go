package handlers

import (
	"context"
	"log/slog"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/common"
	"github.com/eterline/fstmon/internal/infra/grpc/flugel/convert"
)

type storageHandlers struct {
	common.UnimplementedCpuServiceServer
	handlersGroup
}

func NewStorageHandlers(l *slog.Logger, s ActualStateStore) *storageHandlers {
	return &storageHandlers{
		handlersGroup: newHandlersGroup(l, s),
	}
}

func (nh *storageHandlers) GetPartitions(context.Context, *common.GetPartitionsRequest) (*common.PartitionsResponse, error) {
	data, err := GetMetric[domain.Partitions](nh.store, "partitions")
	if err != nil {
		nh.log.Error("failed get partitions", "error", err)
		return nil, err
	}

	res := convert.PartitionsToMessage(data)
	return res, nil
}

func (nh *storageHandlers) GetDiskIO(context.Context, *common.GetDiskIORequest) (*common.DiskIOMapResponse, error) {
	data, err := GetMetric[domain.DiskIOMap](nh.store, "disk_io")
	if err != nil {
		nh.log.Error("failed get disk io", "error", err)
		return nil, err
	}

	res := convert.DiskIOMapResponseToMessage(data)
	return res, nil
}
