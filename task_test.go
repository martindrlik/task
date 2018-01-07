package task_test

import (
	"testing"
	"time"

	"github.com/martindrlik/task"
)

func TestParseKey(t *testing.T) {
	tt := []struct {
		name      string
		keyString string
		keyInt    int64
		withError bool
	}{
		{"parsing ff should pass", "ff", 15*36 + 15, false},
		{"parsing č should fail", "č", 0, true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			k, err := task.ParseKey(tc.keyString)
			if !tc.withError && err != nil {
				t.Errorf("expected to pass with no error, got error %v", err)
			}
			keyInt := int64(k)
			if keyInt != tc.keyInt {
				t.Errorf("expected to parse to int value %v, got %v", tc.keyInt, keyInt)
			}
		})
	}
}

func TestKeyString(t *testing.T) {
	k := task.Key(35)
	if ks := k.String(); ks != "z" {
		t.Errorf("expected key int value 35 to be \"z\" string value, got %q", ks)
	}
}

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

func TestErrNoTask(t *testing.T) {
	tm := task.NewManager()
	var k task.Key
	if _, err := tm.Task(k); err != task.ErrNoTask {
		t.Errorf("there should be no task and error %v, got %v", task.ErrNoTask, err)
	}
}

func TestFromAddedToFinished(t *testing.T) {
	tm := task.NewManager()
	k := tm.Add(time.Now(), "write more unit tests")
	tm.Finish(time.Now(), k)
	if finished := tm.Finished(); len(finished) != 1 {
		t.Errorf("expected one finished task, got %d", len(finished))
	}
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

func TestMust(t *testing.T) {
	defer func() {
		switch x := recover().(type) {
		case error:
			if errorString := x.Error(); errorString != task.ErrNoTask.Error() {
				t.Errorf("expected panic with error to be %v, got %v", task.ErrNoTask, errorString)
			}
		default:
			t.Errorf("expected panic with error %v, got %T %v", task.ErrNoTask, x, x)
		}
	}()
	tm := task.NewManager()
	var k task.Key
	task.Must(tm.Task(k))
	t.Errorf("expected panic if there is no task for given key")
}
