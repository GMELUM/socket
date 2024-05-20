package pool

import (
	"testing"
	"time"
)

func TestPoolSchedule(t *testing.T) {
	pool := New(3, 2)

	err := pool.Schedule(func() {})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	taskExecuted := make(chan bool)
	pool.Schedule(func() {
		taskExecuted <- true
	})
	select {
	case <-taskExecuted:
	case <-time.After(time.Second):
		t.Error("Task did not execute within time limit")
	}
}

func TestQueueOverflow(t *testing.T) {
	pool := New(3, 1)

	for i := 0; i < 10; i++ {
		err := pool.Schedule(func() {
			time.Sleep(time.Second * 10)
		})
		if err != nil {
			return
		}
	}

	t.Error("The task has not been abandoned")
}

func TestDisableQueue(t *testing.T) {
	pool := New(3, 0)

	for i := 0; i < 10; i++ {
		err := pool.Schedule(func() {
			time.Sleep(time.Second * 1)
		})
		if err != nil {
			t.Error(err)
		}
	}

	
}
