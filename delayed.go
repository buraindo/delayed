package delayed

import (
	"errors"
	"math"
	"sync"
	"time"
)

type manager struct {
	ticker *time.Ticker
	tasks  map[int64][]*task

	startTime int64
	bucket    int64
	duration  int64

	stopCh chan struct{}
	mutex  sync.Mutex
}

func NewTaskManager(d time.Duration) (*manager, error) {
	if d < time.Second {
		return nil, errors.New("too small duration, must be at least 1s")
	}

	tm := &manager{
		ticker:    time.NewTicker(d),
		tasks:     map[int64][]*task{},
		startTime: time.Now().UnixNano(),
		bucket:    1,
		duration:  d.Nanoseconds(),
		stopCh:    make(chan struct{}),
	}
	tm.start()
	return tm, nil
}

func (m *manager) Shutdown() {
	close(m.stopCh)
	m.ticker.Stop()
}

func (m *manager) Run(action action, in time.Duration) (*future, error) {
	return m.RunAt(action, time.Now().Add(in))
}

func (m *manager) RunAt(action action, at time.Time) (*future, error) {
	if at.Before(time.Now()) {
		return nil, errors.New("can't run tasks in the past")
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()

	f := &future{}
	f.mutex.Lock()
	t := &task{
		action: action,
		future: f,
	}

	b := int64(math.Ceil(float64(at.UnixNano()-m.startTime) / float64(m.duration)))
	if m.tasks[b] == nil {
		m.tasks[b] = make([]*task, 0)
	}
	m.tasks[b] = append(m.tasks[b], t)
	return f, nil
}

func (m *manager) start() {
	go m.process()
}

func (m *manager) process() {
	for {
		select {
		case <-m.stopCh:
			return
		case <-m.ticker.C:
			m.mutex.Lock()
			bucket := m.bucket
			m.bucket++
			for b := bucket - 1; b <= bucket; b++ {
				if m.tasks[b] == nil {
					continue
				}
				for _, t := range m.tasks[b] {
					go func(t *task) {
						result, err := t.action()
						t.future.result = result
						t.future.err = err
						t.future.mutex.Unlock()
					}(t)
				}
				delete(m.tasks, b)
			}
			m.mutex.Unlock()
		}
	}
}

type action func() (any, error)

type future struct {
	result any
	err    error
	mutex  sync.RWMutex
}

func (f *future) Get() any {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.result
}

func (f *future) HasError() bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.err != nil
}

func (f *future) Error() error {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.err
}

type task struct {
	action action
	future *future
}
