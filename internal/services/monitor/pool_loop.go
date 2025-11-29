package monitor

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/eterline/fstmon/internal/log"
)

/*
MetricSaverGeter — abstraction for storing and retrieving metric values.

Save — stores a value under a specific key.
Get  — retrieves the stored value by key.
*/
type MetricSaverGeter interface {
	Save(v any, key string)
	Get(key string) (any, bool)
}

/*
UpdateWorker — metric update function.

	Accepts a context.
	Returns the updated metric value or an error.
*/
type UpdateWorker func(context.Context) (any, error)

// updateJobContext — internal structure containing worker configuration.
type updateJobContext struct {
	job      UpdateWorker
	interval time.Duration
	cycles   uint64
}

func (ujc *updateJobContext) incCycle() {
	atomic.AddUint64(&ujc.cycles, 1)
}

func (ujc *updateJobContext) actualCycle() uint64 {
	return atomic.LoadUint64(&ujc.cycles)
}

/*
ServicePooler — manages a pool of metric update workers.

	Each worker periodically calls its UpdateWorker function
	and stores the result in the MetricSaverGeter.
*/
type ServicePooler struct {
	metricPool MetricSaverGeter
	jobPool    map[string]*updateJobContext
	wg         sync.WaitGroup
}

/*
NewServicePooler — creates a new worker pool instance.

s — implementation of MetricSaverGeter used to store metric results.
*/
func NewServicePooler(s MetricSaverGeter) *ServicePooler {
	return &ServicePooler{
		metricPool: s,
		jobPool:    map[string]*updateJobContext{},
	}
}

/*
AddMetricPooling — registers a periodic metric update worker.

uj       — update function to be executed.
key      — unique identifier for the metric.
interval — time interval between updates.
*/
func (sp *ServicePooler) AddMetricPooling(uj UpdateWorker, key string, interval time.Duration) {
	sp.jobPool[key] = &updateJobContext{
		job:      uj,
		interval: interval,
		cycles:   0,
	}
}

/*
RunPooling — starts all registered metric update workers.

Each worker runs in its own goroutine.
Updates are saved to the MetricSaverGeter under their respective keys.
*/
func (sp *ServicePooler) RunPooling(ctx context.Context) {
	log := log.MustLoggerFromContext(ctx)

	if len(sp.jobPool) == 0 {
		log.Warn("metric update workers do not exist")
		return
	}

	log.InfoContext(ctx, "metric update workers starting", "workers", len(sp.jobPool))

	// Start each worker in its own goroutine.
	for key, workerCtx := range sp.jobPool {

		log := log.With("worker_key", key)

		log.InfoContext(ctx,
			"worker start",
			"pool_interval", workerCtx.interval,
		)

		sp.wg.Add(1)

		go func(key string, workerCtx *updateJobContext, log *slog.Logger) {
			defer sp.wg.Done()

			t := time.NewTicker(workerCtx.interval)
			defer t.Stop()

			for {
				select {
				case <-ctx.Done():
					return

				case <-t.C:
					workerCtx.incCycle()
					log.DebugContext(ctx, "metric update request", "worker_cycle", workerCtx.actualCycle())
					update, err := workerCtx.job(ctx)
					if err != nil {
						log.ErrorContext(ctx, "metric update request error", "error", err)
						continue
					}

					sp.metricPool.Save(update, key)
				}
			}
		}(key, workerCtx, log)
	}
}

/*
Await — waits for all metric update workers to finish.

	Typically used after RunPooling when graceful shutdown is required.
*/
func (sp *ServicePooler) Await() {
	sp.wg.Wait()
}

/*
ActualMetric — retrieves the most recent value of a metric and
reports whether the worker and metric exist.

	value         — metric value (may be nil)
	workerExists  — true if worker with such key exists
	metricExists  — true if metric value was already saved
*/
func (sp *ServicePooler) ActualMetric(key string) (value any, workerExists, metricExists bool) {
	_, ok := sp.jobPool[key]
	if !ok {
		return nil, false, false
	}

	m, ok := sp.metricPool.Get(key)
	if !ok || m == nil {
		return nil, true, false
	}

	return m, true, true
}

func (sp *ServicePooler) MetricInterval(key string) (time.Duration, bool) {
	ctx, ok := sp.jobPool[key]
	if !ok {
		return 0, false
	}
	return ctx.interval, true
}
