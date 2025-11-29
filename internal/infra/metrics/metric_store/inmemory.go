package metricstore

import "sync"

type MetricMemoryStore struct {
	mu   sync.Mutex
	pool map[string]*ValueStore[any]
}

func NewMetricMemoryStore() *MetricMemoryStore {
	return &MetricMemoryStore{
		pool: map[string]*ValueStore[any]{},
	}
}

func (mms *MetricMemoryStore) Save(v any, key string) {
	vs, ok := mms.pool[key]
	if !ok {
		vs = NewValueStore[any]()
		mms.mu.Lock()
		mms.pool[key] = vs
		mms.mu.Unlock()
	}
	vs.Save(v)
}

func (mms *MetricMemoryStore) Get(key string) (v any, storeExists bool) {
	vs, ok := mms.pool[key]
	if !ok {
		return nil, false
	}
	return vs.Get(), true
}

func (mms *MetricMemoryStore) Clear() {
	for k, vs := range mms.pool {
		vs.Clear()
		delete(mms.pool, k)
	}
}
