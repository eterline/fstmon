package httphomepage

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/infra/http/common/api"
)

type ActualStateStore interface {
	ActualMetric(key string) (value any, scheduleExists bool, stateExists bool, retryIn time.Duration)
}

type HomepageHandlerGroup struct {
	actualStore ActualStateStore
	log         *slog.Logger
}

func New(ass ActualStateStore, log *slog.Logger) *HomepageHandlerGroup {
	return &HomepageHandlerGroup{
		actualStore: ass,
		log:         log,
	}
}

// GetMetric - generic wrapper for retrieving actual metrics by key.
// - 404 if schedule does not exist
// - 503 if state does not exist (adds Retry-In header)
// - 500 if the metric type is not assignable to T
func GetMetric[T any](ass ActualStateStore, w http.ResponseWriter, key string) (T, bool) {
	var zero T

	value, scheduleExists, stateExists, retryIn := ass.ActualMetric(key)
	if !scheduleExists {
		api.NewResponse().
			SetCode(http.StatusNotFound).
			SetMessage("worker not exists").
			AddStringError("metric not found").
			Write(w)

		return zero, false
	}

	if !stateExists {
		w.Header().Set("Retry-In", strconv.Itoa(int(retryIn.Seconds())))
		api.NewResponse().
			SetCode(http.StatusServiceUnavailable).
			SetMessage("metric not available").
			AddStringError("metric did not scraped yet").
			Write(w)

		return zero, false
	}

	casted, ok := value.(T)
	if !ok {
		api.InternalErrorResponse().
			SetMessage("internal type mismatch").
			Write(w)

		return zero, false
	}

	return casted, true
}

func (hhg *HomepageHandlerGroup) HandleThermal(w http.ResponseWriter, r *http.Request) {
	m, ok := GetMetric[domain.ThermalMetricsMap](hhg.actualStore, w, "thermal")
	if !ok {
		return
	}

	err := api.NewResponse().WrapData(m).Write(w)
	if err != nil {
		hhg.log.Error("response error", "error", err)
	}
}

func (hhg *HomepageHandlerGroup) HandleSystem(w http.ResponseWriter, r *http.Request) {
	m, ok := GetMetric[domain.SystemInfo](hhg.actualStore, w, "system")
	if !ok {
		return
	}

	dto := Domain2DTOSystem(m)

	err := api.NewResponse().WrapData(dto).Write(w)
	if err != nil {
		hhg.log.Error("response error", "error", err)
	}
}

func (hhg *HomepageHandlerGroup) HandleNetwork(w http.ResponseWriter, r *http.Request) {
	m, ok := GetMetric[domain.InterfacesIOMap](hhg.actualStore, w, "net_io")
	if !ok {
		return
	}

	dto := Domain2DTONetworkInterfaceIO(m)

	err := api.NewResponse().WrapData(dto).Write(w)
	if err != nil {
		hhg.log.Error("response error", "error", err)
	}
}

func (hhg *HomepageHandlerGroup) HandleMemory(w http.ResponseWriter, r *http.Request) {
	m, ok := GetMetric[domain.MemoryMetrics](hhg.actualStore, w, "memory")
	if !ok {
		return
	}

	dto := Domain2DTOMemory(m)

	err := api.NewResponse().WrapData(dto).Write(w)
	if err != nil {
		hhg.log.Error("response error", "error", err)
	}
}

func (hhg *HomepageHandlerGroup) HandleCpu(w http.ResponseWriter, r *http.Request) {
	pkg, ok := GetMetric[domain.CpuPackage](hhg.actualStore, w, "cpu")
	if !ok {
		return
	}

	mtrcs, ok := GetMetric[domain.CpuMetrics](hhg.actualStore, w, "cpu_usage")
	if !ok {
		return
	}

	dto := Domain2DTOCpu(pkg, mtrcs)

	err := api.NewResponse().WrapData(dto).Write(w)
	if err != nil {
		hhg.log.Error("response error", "error", err)
	}
}
