package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
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
