// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package httphomepage

import (
	"context"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/infra/http/common/api"
	"github.com/eterline/fstmon/internal/log"
)

type ActualStateStore interface {
	ActualMetric(key string) (value any, scheduleExists bool, stateExists bool, retryIn time.Duration)
}

type HomepageHandlerGroup struct {
	actualStore ActualStateStore
}

func New(ass ActualStateStore) *HomepageHandlerGroup {
	return &HomepageHandlerGroup{
		actualStore: ass,
	}
}

// GetMetric – generic wrapper for retrieving actual metrics by key.
// – 404 if schedule does not exist
// – 503 if state does not exist (adds Retry-In header)
// – 500 if the metric type is not assignable to T
func GetMetric[T any](ctx context.Context, ass ActualStateStore, w http.ResponseWriter, key string) (T, bool) {
	log := log.MustLoggerFromContext(ctx)

	var zero T

	value, scheduleExists, stateExists, retryIn := ass.ActualMetric(key)
	if !scheduleExists {
		api.NewResponse().
			SetCode(http.StatusNotFound).
			SetMessage("worker not exists").
			AddStringError("metric not found").
			Write(w)

		log.Error("invalid metric key request", "metric_key", key)
		return zero, false
	}

	if !stateExists {
		retry := strconv.Itoa(int(retryIn.Seconds()))
		w.Header().Set("Retry-In", retry)
		api.NewResponse().
			SetCode(http.StatusServiceUnavailable).
			SetMessage("metric not available").
			AddStringError("metric did not scraped yet").
			Write(w)

		log.Warn("metric not exists yet", "metric_key", key, "retry_in_sec", retry)
		return zero, false
	}

	casted, ok := value.(T)
	if !ok {
		api.InternalErrorResponse().
			SetMessage("internal type mismatch").
			Write(w)

		expectedType := reflect.TypeOf(zero)
		actualType := reflect.TypeOf(value)

		log.Error(
			"store type mismatch", "metric_key", key,
			"expected", expectedType, "got", actualType,
		)

		return zero, false
	}

	return casted, true
}

func (hhg *HomepageHandlerGroup) HandleThermal(w http.ResponseWriter, r *http.Request) {
	log := log.MustLoggerFromContext(r.Context())
	m, ok := GetMetric[domain.ThermalMetricsMap](r.Context(), hhg.actualStore, w, "thermal")
	if !ok {
		return
	}

	err := api.NewResponse().WrapData(m).Write(w)
	if err != nil {
		log.Error("response error", "error", err)
	}
}

func (hhg *HomepageHandlerGroup) HandleSystem(w http.ResponseWriter, r *http.Request) {
	log := log.MustLoggerFromContext(r.Context())

	m, ok := GetMetric[domain.SystemInfo](r.Context(), hhg.actualStore, w, "system")
	if !ok {
		return
	}

	dto := Domain2DTOSystem(m)

	err := api.NewResponse().WrapData(dto).Write(w)
	if err != nil {
		log.Error("response error", "error", err)
	}
}

func (hhg *HomepageHandlerGroup) HandleNetwork(w http.ResponseWriter, r *http.Request) {
	log := log.MustLoggerFromContext(r.Context())

	m, ok := GetMetric[domain.InterfacesIOMap](r.Context(), hhg.actualStore, w, "net_io")
	if !ok {
		return
	}

	dto := Domain2DTONetworkInterfaceIO(m)

	err := api.NewResponse().WrapData(dto).Write(w)
	if err != nil {
		log.Error("response error", "error", err)
	}
}

func (hhg *HomepageHandlerGroup) HandleMemory(w http.ResponseWriter, r *http.Request) {
	log := log.MustLoggerFromContext(r.Context())

	m, ok := GetMetric[domain.MemoryMetrics](r.Context(), hhg.actualStore, w, "memory")
	if !ok {
		return
	}

	dto := Domain2DTOMemory(m)

	err := api.NewResponse().WrapData(dto).Write(w)
	if err != nil {
		log.Error("response error", "error", err)
	}
}

func (hhg *HomepageHandlerGroup) HandleCpu(w http.ResponseWriter, r *http.Request) {
	log := log.MustLoggerFromContext(r.Context())

	pkg, ok := GetMetric[domain.CpuPackage](r.Context(), hhg.actualStore, w, "cpu")
	if !ok {
		return
	}

	mtrcs, ok := GetMetric[domain.CpuMetrics](r.Context(), hhg.actualStore, w, "cpu_usage")
	if !ok {
		return
	}

	dto := Domain2DTOCpu(pkg, mtrcs)

	err := api.NewResponse().WrapData(dto).Write(w)
	if err != nil {
		log.Error("response error", "error", err)
	}
}

func (hhg *HomepageHandlerGroup) HandlePartitions(w http.ResponseWriter, r *http.Request) {
	log := log.MustLoggerFromContext(r.Context())

	data, ok := GetMetric[domain.Partitions](r.Context(), hhg.actualStore, w, "partitions")
	if !ok {
		return
	}

	dto := Domain2DTOPartitions(data)

	err := api.NewResponse().WrapData(dto).Write(w)
	if err != nil {
		log.Error("response error", "error", err)
	}
}

func (hhg *HomepageHandlerGroup) HandleDiskIO(w http.ResponseWriter, r *http.Request) {
	log := log.MustLoggerFromContext(r.Context())

	disks, ok := GetMetric[domain.DiskIOMap](r.Context(), hhg.actualStore, w, "disk_io")
	if !ok {
		return
	}

	dto := Domain2DTODiskIOs(disks)

	err := api.NewResponse().WrapData(dto).Write(w)
	if err != nil {
		log.Error("response error", "error", err)
	}
}
