package task_test

import (
	"fmt"
	"time"

	"github.com/martindrlik/task"
)

func ExampleTask() {
	tm := task.NewManager()
	k := tm.Add(time.Now(), "write more unit tests")
	for _, k := range tm.Added() {
		t := task.Must(tm.Task(k))
		fmt.Println(t.Title, "added")
	}
	tm.Start(time.Now(), k)
	for _, k := range tm.Started() {
		t := task.Must(tm.Task(k))
		fmt.Println(t.Title, "started")
	}
	tm.Finish(time.Now(), k)
	for _, k := range tm.Finished() {
		t := task.Must(tm.Task(k))
		fmt.Println(t.Title, "finished")
	}
	// Output:
	// write more unit tests added
	// write more unit tests started
	// write more unit tests finished
}
