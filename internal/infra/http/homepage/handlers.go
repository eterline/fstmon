package httphomepage

import "time"

type ActualStateStore interface {
	ActualMetric(key string) (value any, scheduleExists bool, stateExists bool, retryIn time.Duration)
}

type StatsStateStore interface {
	ActualStats(key string) (value any, scheduleExists bool, stateExists bool)
}

type HomepageHandlerGroup struct {
	actualStore ActualStateStore
	statslStore StatsStateStore
}

func New(ass ActualStateStore, sss StatsStateStore) *HomepageHandlerGroup {
	return &HomepageHandlerGroup{
		actualStore: ass,
		statslStore: sss,
	}
}
