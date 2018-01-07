package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/martindrlik/task/task"
)

func ExampleProcFile() {
	tasks = task.NewManager()
	now = mockNow("2006", "2018", time.Second)
	file := &bytes.Buffer{}
	fmt.Fprint(file, "todo write more unit tests\n")
	fmt.Fprint(file, "todo")
	err := ProcFile(file)
	if err != nil {
		panic(err)
	}
	// Output:
	// todo 1 write more unit tests
}

func ExampleProc() {
	tasks = task.NewManager()
	now = mockNow("2006", "2018", time.Hour)
	Proc("todo write more unit tests")
	Proc("todo write even more unit tests")
	Proc("todo")
	Proc("doing")
	Proc("done 1")
	Proc("start 2")
	Proc("doing")
	Proc("done")
	// Output:
	// todo 1 write more unit tests
	// todo 2 write even more unit tests
	// doing none
	// doing 2 write even more unit tests
	// done 1 write more unit tests took 2h0m0s
}

func mockNow(layout, value string, add time.Duration) func() time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return func() time.Time {
		defer func() { t = t.Add(add) }()
		return t
	}
}
