package metricstore

import (
	"reflect"
	"sync"
	"time"

	"github.com/eterline/fstmon/internal/domain"
)

/*
MetricInMemoryStore - in-memory storage, implements MetricRepository.

	Stores a per-key preallocated pointer to a concrete value and reuses it on updates
	to avoid repeated allocations.
*/
type MetricInMemoryStore struct {
	mu sync.RWMutex
	db map[string]*metricEntry
}

type metricEntry struct {
	// per-entry mutex protects fields below
	mu sync.RWMutex

	// valuePtr stores a pointer to the concrete value (always a pointer type).
	// We keep a pointer so we can update its Elem() without reallocating pointer/interface.
	// Example: for a value of type T, valuePtr has type *T (stored as interface{}).
	valuePtr any

	lastUpdate int64 // milliseconds since epoch
	available  bool
}

func (me *metricEntry) lastUpdateTime() time.Time {
	return time.UnixMilli(me.lastUpdate)
}

func (me *metricEntry) setLastUpdateTime(t time.Time) {
	me.lastUpdate = t.UnixMilli()
}

// NewMetricInMemoryStore - create new in-memory storage.
func NewMetricInMemoryStore() *MetricInMemoryStore {
	return &MetricInMemoryStore{
		db: make(map[string]*metricEntry),
	}
}

// =========================
// MetricRepository methods
// =========================

/*
SaveValue - saves value for the given key and timestamp.

	Behavior:
	- If this is the first time a value for this key is seen, allocate a pointer of the concrete type
		and store the incoming value into it.
	- On subsequent saves, if incoming value's concrete type matches the stored pointer's element type,
		copy the new value into the existing storage (no new allocation for the pointer/interface).
	- If the type differs, replace storage with a newly allocated pointer for the new type.

	Thread-safety:
	- Map-level locking for creating/accessing entries.
	- Per-entry locking for updating value/metadata.
*/
func (r *MetricInMemoryStore) SaveValue(key string, value any, ts time.Time) {

	// ensure entry exists
	r.mu.RLock()
	entry, ok := r.db[key]
	r.mu.RUnlock()

	if !ok {
		// create new entry under write lock
		r.mu.Lock()
		entry, ok = r.db[key]
		if !ok {
			entry = &metricEntry{}
			r.db[key] = entry
		}
		r.mu.Unlock()
	}

	// lock entry for write
	entry.mu.Lock()
	defer entry.mu.Unlock()

	nv := reflect.ValueOf(value)
	if !nv.IsValid() {
		// unset value
		entry.available = false
		entry.setLastUpdateTime(ts)
		return
	}

	// If we have no storage pointer yet -> allocate pointer to concrete type and set
	if entry.valuePtr == nil {
		var storage reflect.Value
		if nv.Kind() == reflect.Ptr {
			storage = reflect.New(nv.Elem().Type()) // *T
			storage.Elem().Set(nv.Elem())
		} else {
			storage = reflect.New(nv.Type())
			storage.Elem().Set(nv)
		}
		entry.valuePtr = storage.Interface()
		entry.available = true
		entry.setLastUpdateTime(ts)
		return
	}

	// we have existing storage pointer; try to reuse it
	ev := reflect.ValueOf(entry.valuePtr) // should be Ptr
	if ev.Kind() == reflect.Ptr && ev.Elem().CanSet() {

		target := ev.Elem()

		if nv.Kind() == reflect.Ptr {

			if nv.Elem().Type() == target.Type() {
				target.Set(nv.Elem())
				entry.available = true
				entry.setLastUpdateTime(ts)
				return
			}
		} else {
			if nv.Type() == target.Type() {
				target.Set(nv)
				entry.available = true
				entry.setLastUpdateTime(ts)
				return
			}
		}
	}

	// fallback: types differ or cannot set -> allocate new storage for new type
	var newStorage reflect.Value
	if nv.Kind() == reflect.Ptr {
		newStorage = reflect.New(nv.Elem().Type())
		newStorage.Elem().Set(nv.Elem())
	} else {
		newStorage = reflect.New(nv.Type())
		newStorage.Elem().Set(nv)
	}
	entry.valuePtr = newStorage.Interface()
	entry.available = true
	entry.setLastUpdateTime(ts)
}

/*
GetState - returns domain.MetricState and whether key exists.

	We return a copy of the stored value (dereferenced) so caller receives normal
	value (not pointer to internal buffer).
*/
func (r *MetricInMemoryStore) GetState(key string) (domain.MetricState, bool) {
	r.mu.RLock()
	entry, ok := r.db[key]
	r.mu.RUnlock()

	if !ok {
		return domain.MetricState{}, false
	}

	entry.mu.RLock()
	defer entry.mu.RUnlock()

	var v any
	if entry.available && entry.valuePtr != nil {
		ev := reflect.ValueOf(entry.valuePtr)
		if ev.Kind() == reflect.Ptr {
			v = ev.Elem().Interface()
		} else {
			v = entry.valuePtr
		}
	}

	return domain.MetricState{
		Value:      v,
		Available:  entry.available,
		LastUpdate: entry.lastUpdateTime(),
	}, true
}

/*
Close - clears all entries and releases internal storage.

	Behavior:
	- For each entry we zero the underlying storage (if pointer), mark unavailable and delete map entries.
*/
func (r *MetricInMemoryStore) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for k, entry := range r.db {
		entry.mu.Lock()
		if entry.valuePtr != nil {
			ev := reflect.ValueOf(entry.valuePtr)
			if ev.Kind() == reflect.Ptr && ev.Elem().CanSet() {
				ev.Elem().Set(reflect.Zero(ev.Elem().Type()))
			}
			entry.valuePtr = nil
		}

		entry.available = false
		entry.lastUpdate = 0
		entry.mu.Unlock()

		delete(r.db, k)
	}
}
