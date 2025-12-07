package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/eterline/fstmon/internal/infra/grpc/flugel/common"
	"github.com/eterline/fstmon/internal/log"
	"google.golang.org/grpc"
)

type ActualStateStore interface {
	ActualMetric(key string) (value any, scheduleExists bool, stateExists bool, retryIn time.Duration)
}

func GetMetric[T any](ass ActualStateStore, key string) (T, error) {
	var zero T

	value, scheduleExists, stateExists, _ := ass.ActualMetric(key)
	if !scheduleExists {
		return zero, fmt.Errorf("worker not exists under key: '%s'", key)
	}

	if !stateExists {
		return zero, fmt.Errorf("metric not exists yet: '%s'", key)
	}

	casted, ok := value.(T)
	if !ok {
		return zero, errors.New("store type mismatch")
	}

	return casted, nil
}

type handlersGroup struct {
	log   *slog.Logger
	store ActualStateStore
}

func newHandlersGroup(l *slog.Logger, s ActualStateStore) handlersGroup {
	return handlersGroup{
		log:   l,
		store: s,
	}
}

// TODO: make another app instance for grpc agent
func RegisterToGrpcServer(ctx context.Context, s *grpc.Server, a ActualStateStore) {
	log := log.MustLoggerFromContext(ctx)
	log.Info("init grpc server handlers")

	pb.RegisterCpuServiceServer(s, NewCpuHandlers(log, a))
	pb.RegisterNetworkServiceServer(s, NewNetworkHandlers(log, a))
	pb.RegisterSystemServiceServer(s, NewSystemHandlers(log, a))
	pb.RegisterMemoryServiceServer(s, NewMemoryHandlers(log, a))
	pb.RegisterThermalServiceServer(s, NewThernalHandlers(log, a))
	pb.RegisterStorageServiceServer(s, NewStorageHandlers(log, a))
}
