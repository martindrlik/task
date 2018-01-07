package task

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	// ErrNoTask indicates that given key points
	// to task which not exists in task manager.
	ErrNoTask = errors.New("no task for given key")
)

// Key is type used to identify tasks.
type Key int64

func (key Key) String() string {
	return strconv.FormatInt(int64(key), 36)
}

// ParseKey parses task key from string s.
func ParseKey(s string) (Key, error) {
	i, err := strconv.ParseInt(s, 36, 64)
	if err != nil {
		return 0, fmt.Errorf("task could not parse key: %v", err)
	}
	return Key(i), nil
}

// Task represents a task data.
type Task struct {
	Title    string
	Added    time.Time
	Started  time.Time
	Finished time.Time
}

// Manager manages tasks.
type Manager struct {
	all      map[Key]Task
	added    []Key
	started  []Key
	finished []Key
}

// NewManager initializes ready to used task manager.
func NewManager() *Manager {
	tm := Manager{
		all: make(map[Key]Task),
	}
	return &tm
}

// Add appends task with title in added state and returns its key.
func (tm *Manager) Add(time time.Time, title string) Key {
	key := Key(len(tm.all) + 1)
	t := Task{
		Title: title,
		Added: time,
	}
	tm.all[key] = t
	tm.added = append(tm.added, key)
	return key
}

// Start sets started state to task given by its key.
func (tm *Manager) Start(time time.Time, key Key) {
	added := tm.added
	i, ok := index(added, key)
	if ok {
		tm.added = append(added[:i], added[i+1:]...)
		tm.started = append(tm.started, key)
	}
	t := tm.all[key]
	t.Started = time
	tm.all[key] = t
}

// Finish sets finished state to task given by its key.
func (tm *Manager) Finish(time time.Time, key Key) {
	i, ok := index(tm.started, key)
	if ok {
		tm.started = append(tm.started[:i], tm.started[i+1:]...)
		tm.finished = append(tm.finished, key)
		t := tm.all[key]
		t.Finished = time
		tm.all[key] = t
		return
	}
	i, ok = index(tm.added, key)
	if ok {
		tm.added = append(tm.added[:i], tm.added[i+1:]...)
		tm.finished = append(tm.finished, key)
	}
	t := tm.all[key]
	t.Finished = time
	tm.all[key] = t
}

func index(keys []Key, key Key) (int, bool) {
	for i, k := range keys {
		if key == k {
			return i, true
		}
	}
	return 0, false
}

// Added returns all keys of tasks in added state.
func (tm *Manager) Added() []Key {
	return tm.added
}

// Started returns all keys of tasks in started state.
func (tm *Manager) Started() []Key {
	return tm.started
}

// Finished returns all keys of tasks in finished state.
func (tm *Manager) Finished() []Key {
	return tm.finished
}

// Task returns task for given key.
func (tm *Manager) Task(key Key) (Task, error) {
	task, ok := tm.all[key]
	if !ok {
		return Task{}, ErrNoTask
	}
	return task, nil
}

// Must is a helper that wraps a call to a function returning (Task, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//	var t = task.Must(tm.Task(key))
func Must(t Task, err error) Task {
	if err != nil {
		panic(err)
	}
	return t
}

func (tm *Manager) tasks(keys []Key) []Task {
	tasks := make([]Task, len(keys))
	for i, key := range keys {
		tasks[i] = tm.all[key]
	}
	return tasks
}
