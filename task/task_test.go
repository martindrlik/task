package task_test

import (
	"testing"
	"time"

	"github.com/martindrlik/task/task"
)

func TestTask(t *testing.T) {
	tm := task.NewManager()
	k := tm.Add(time.Now(), "write more unit tests")
	t.Run("add", func(t *testing.T) {
		if added := tm.Added(); len(added) != 1 {
			t.Errorf("expected one task added, got %d", len(added))
		}
	})
	t.Run("start", func(t *testing.T) {
		tm.Start(time.Now(), k)
		if started := tm.Started(); len(started) != 1 {
			t.Errorf("expected one task started, got %d", len(started))
		}
		if added := tm.Added(); len(added) != 0 {
			t.Errorf("expected no task added, got %d", len(added))
		}
	})
	t.Run("finish", func(t *testing.T) {
		tm.Finish(time.Now(), k)
		if finished := tm.Finished(); len(finished) != 1 {
			t.Errorf("expected one task finished, got %d", len(finished))
		}
		if started := tm.Started(); len(started) != 0 {
			t.Errorf("expected no task started, got %d", len(started))
		}
		if added := tm.Added(); len(added) != 0 {
			t.Errorf("expected no task added, got %d", len(added))
		}
	})
}

func TestAddMultiple(t *testing.T) {
	tm := task.NewManager()
	tm.Add(time.Now(), "write more unit tests")
	tm.Add(time.Now(), "write even more unit tests")
	added := tm.Added()
	if len(added) != 2 {
		t.Fatalf("expected two added, got %d added", len(added))
	}
	if title := task.Must(tm.Task(added[0])).Title; title != "write more unit tests" {
		t.Errorf("expected first task to have title %q, got %q", "write more unit tests", title)
	}
	if title := task.Must(tm.Task(added[1])).Title; title != "write even more unit tests" {
		t.Errorf("expected second task to have title %q, got %q", "write more unit tests", title)
	}
}

func TestTime(t *testing.T) {
	added := time.Now()
	started := added.Add(time.Second)
	finished := started.Add(time.Second)
	tm := task.NewManager()

	k := tm.Add(added, "write unit tests")
	tm.Start(started, k)
	tm.Finish(finished, k)

	task := task.Must(tm.Task(k))
	if added != task.Added {
		t.Errorf("expected added time to be %v, got %v", added, task.Added)
	}
	if started != task.Started {
		t.Errorf("expected started time to be %v, got %v", started, task.Started)
	}
	if finished != task.Finished {
		t.Errorf("expected finished time to be %v, got %v", finished, task.Finished)
	}
}
