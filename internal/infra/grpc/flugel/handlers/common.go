package handlers

import (
	"context"
	"errors"
	"fmt"
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

// TODO: make another app instance for grpc agent
func RegisterToGrpcServer(ctx context.Context, s *grpc.Server, a ActualStateStore) {
	log := log.MustLoggerFromContext(ctx)

	log.Info("init grpc server handlers")
	pb.RegisterMachineInfoServiceServer(s, NewMachineInfohandlers(log, a))
}
