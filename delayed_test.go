package delayed_test

import (
	"errors"
	"github.com/buraindo/delayed"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestCreateManager(t *testing.T) {
	m, err := delayed.NewTaskManager(time.Second)
	if err != nil {
		t.Fatalf("creating task manager failed: %v", err)
	}
	m.Shutdown()
}

func TestCreateManagerError(t *testing.T) {
	_, err := delayed.NewTaskManager(500 * time.Millisecond)
	if err == nil {
		t.Fatalf("error should not be nil")
	}
}

func TestRunError(t *testing.T) {
	m, err := delayed.NewTaskManager(time.Second)
	if err != nil {
		t.Fatalf("creating task manager failed: %v", err)
	}
	f, err := m.Run(func() (any, error) {
		return nil, nil
	}, -time.Second)
	if err == nil {
		t.Fatalf("error should not be nil")
	}
	if f != nil {
		t.Fatalf("future should be nil")
	}
	m.Shutdown()
}

func TestRunSimple(t *testing.T) {
	m, err := delayed.NewTaskManager(time.Second)
	if err != nil {
		t.Fatalf("creating task manager failed: %v", err)
	}
	f, err := m.Run(func() (any, error) {
		return "hello", nil
	}, time.Second)
	if err != nil {
		t.Fatalf("running task failed: %v", err)
	}
	if f == nil {
		t.Fatalf("future should not be nil")
	}
	if f.HasError() {
		t.Fatalf("hasError should be false")
	}
	if f.Error() != nil {
		t.Fatalf("error should be nil")
	}
	if f.Get() != "hello" {
		t.Fatalf("result should be 'hello'")
	}
	m.Shutdown()
}

func TestRunSimpleError(t *testing.T) {
	m, err := delayed.NewTaskManager(time.Second)
	if err != nil {
		t.Fatalf("creating task manager failed: %v", err)
	}
	f, err := m.Run(func() (any, error) {
		return nil, errors.New("some error")
	}, time.Second)
	if err != nil {
		t.Fatalf("running task failed: %v", err)
	}
	if f == nil {
		t.Fatalf("future should not be nil")
	}
	if !f.HasError() {
		t.Fatalf("hasError should be true")
	}
	if f.Error() == nil {
		t.Fatalf("error should not be nil")
	}
	if f.Get() != nil {
		t.Fatalf("result should be nil")
	}
	m.Shutdown()
}

func TestRun(t *testing.T) {
	m, err := delayed.NewTaskManager(time.Second)
	if err != nil {
		t.Fatalf("creating task manager failed: %v", err)
	}
	cases := []struct {
		hasError bool
		result   any
	}{
		{
			result: 42,
		},
		{
			result: errors.New("result"),
		},
		{
			hasError: true,
		},
		{
			result: "hi",
		},
		{
			result: struct{}{},
		},
		{
			result: false,
		},
		{
			hasError: true,
		},
	}
	for _, c := range cases {
		f, err := m.Run(func() (any, error) {
			var e error
			if c.hasError {
				e = errors.New("some error")
			}
			return c.result, e
		}, 50*time.Millisecond)
		if err != nil {
			t.Fatalf("running task failed: %v", err)
		}
		if f == nil {
			t.Fatalf("future should not be nil")
		}
		if c.hasError {
			if !f.HasError() {
				t.Fatalf("hasError should be true")
			}
			if f.Error() == nil {
				t.Fatalf("error should not be nil")
			}
			if f.Get() != nil {
				t.Fatalf("result should be nil")
			}
		} else {
			if f.HasError() {
				t.Fatalf("hasError should be false")
			}
			if f.Error() != nil {
				t.Fatalf("error should be nil")
			}
			if f.Get() != c.result {
				t.Fatalf("results don't match")
			}
		}
	}
	m.Shutdown()
}

func TestRunRandom(t *testing.T) {
	m, err := delayed.NewTaskManager(1200 * time.Millisecond)
	if err != nil {
		t.Fatalf("creating task manager failed: %v", err)
	}
	results := make([]any, 0)
	for i := 0; i < 1000000; i++ {
		results = append(results, rand.Int())
	}
	var wg sync.WaitGroup
	for _, r := range results {
		wg.Add(1)
		go func(result any) {
			defer wg.Done()

			d := rand.Int()%500 + 800
			f, err := m.Run(func() (any, error) {
				return result, nil
			}, time.Duration(d)*time.Millisecond)
			if err != nil {
				t.Errorf("running task failed: %v", err)
				return
			}
			if f == nil {
				t.Errorf("future should not be nil")
				return
			}
			if f.HasError() {
				t.Errorf("hasError should be false")
				return
			}
			if f.Error() != nil {
				t.Errorf("error should be nil")
				return
			}
			if f.Get() != result {
				t.Errorf("results don't match")
				return
			}
		}(r)
	}
	wg.Wait()
	m.Shutdown()
}
