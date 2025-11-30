package monitor

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/log"
	"github.com/eterline/fstmon/internal/utils/usecase"
)

/*
UpdateWorker - metric update function.

	Accepts a context.
	Returns the updated metric value or an error.
*/
type UpdateWorker func(context.Context) (any, error)

type WorkerConfig struct {
	Interval time.Duration
	Worker   UpdateWorker
}

type MetricStore interface {
	SaveValue(key string, value any, timestamp time.Time)
	GetState(key string) (domain.MetricState, bool)
}

// ServicePooler - manages a pool of metric update workers.
type ServicePooler struct {
	metricSt  MetricStore
	jobPool   map[string]*WorkerConfig
	workersWg sync.WaitGroup
}

// NewServicePooler - creates a new worker pool instance.
func NewServicePooler(ms MetricStore) *ServicePooler {
	return &ServicePooler{
		metricSt: ms,
		jobPool:  make(map[string]*WorkerConfig),
	}
}

// AddMetricPooling - registers a periodic metric update worker.
func (sp *ServicePooler) AddMetricPooling(w UpdateWorker, key string, wInterval time.Duration) {
	sp.jobPool[key] = &WorkerConfig{
		Interval: wInterval,
		Worker:   w,
	}
}

/*
RunPooling - starts all registered metric update workers.

Each worker runs in its own goroutine.
Updates are saved to the MetricSaverGeter under their respective keys.
*/
func (sp *ServicePooler) RunPooling(ctx context.Context) {
	logger := log.MustLoggerFromContext(ctx)

	metricKeys, _, n := usecase.MapSlicesLen(sp.jobPool)
	if n == 0 {
		logger.Warn("no metric workers registered")
		return
	}

	logger.Info("metric metric workers starting", "count", n, "workers", metricKeys)

	for key, cfg := range sp.jobPool {

		wlog := logger.With(
			"worker_key", key,
			"worker_interval", cfg.Interval,
		)

		wlog.Info("metric worker start")
		sp.workersWg.Add(1)

		go func(key string, cfg *WorkerConfig, wlog *slog.Logger) {

			defer sp.workersWg.Done()

			job := func() {
				wlog.Debug("worker start update metric")

				value, err := cfg.Worker(ctx)
				if err != nil {
					wlog.Error("worker update metric error", "error", err)
					return
				}

				now := time.Now()
				wlog.Debug("worker updated metric", "update_time", now.Format(time.RFC1123))

				sp.metricSt.SaveValue(key, value, now)
			}

			job() // first start

			ticker := time.NewTicker(cfg.Interval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					wlog.Info("metric worker shutdown")
					return
				case <-ticker.C:
					job()
				}
			}

		}(key, cfg, wlog)
	}

}

/*
Await - waits for all metric update workers to finish.

	Typically used after RunPooling when graceful shutdown is required.
*/
func (sp *ServicePooler) Wait() {
	sp.workersWg.Wait()
}

/*
ActualMetric - retrieves the last known state of a metric.

	value          - metric value (may be nil)
	scheduleExists - worker for this key was registered
	stateExists    - repository has at least one saved state
*/
func (sp *ServicePooler) ActualMetric(key string) (value any, scheduleExists, stateExists bool, retryIn time.Duration) {
	job, ok := sp.jobPool[key]
	if !ok {
		return nil, false, false, 0
	}

	state, ok := sp.metricSt.GetState(key)
	if !ok {
		return nil, true, false, job.Interval
	}

	nextUpdate := state.LastUpdate.Add(job.Interval)
	remaining := time.Until(nextUpdate)

	if remaining < 0 {
		remaining = 0
	}

	return state.Value, true, true, remaining
}
