package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/martindrlik/task"
)

var (
	tasks = task.NewManager()
	now   = time.Now
)

func main() {
	err := ProcFile(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "task could not read stdin: %v\n", err)
		os.Exit(1)
	}
}

// ProcFile evaluates line by line of reader r.
func ProcFile(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := scanner.Text()
		Proc(s)
	}
	return scanner.Err()
}

// Proc evaluates string s.
func Proc(s string) {
	if keys, ok := list(s); ok {
		printTasks(s, keys)
	} else if title := strings.TrimPrefix(s, "todo "); title != s {
		tasks.Add(now(), title)
	} else if keyString := strings.TrimPrefix(s, "start "); keyString != s {
		if key, ok := parseKey(keyString); ok {
			tasks.Start(now(), key)
		}
	} else if keyString := strings.TrimPrefix(s, "done "); keyString != s {
		if key, ok := parseKey(keyString); ok {
			tasks.Finish(now(), key)
		}
	}
}

func list(s string) ([]task.Key, bool) {
	switch s {
	case "todo":
		return tasks.Added(), true
	case "doing":
		return tasks.Started(), true
	case "done":
		return tasks.Finished(), true
	}
	return nil, false
}

func printTasks(prefix string, keys []task.Key) {
	if len(keys) == 0 {
		fmt.Println(prefix, "none")
		return
	}
	for _, k := range keys {
		t := task.Must(tasks.Task(k))
		updated := maxTime(t.Added, t.Started)
		if t.Finished.After(updated) {
			fmt.Printf("%s %v %s took %v\n", prefix, k, t.Title, t.Finished.Sub(updated))
		} else {
			fmt.Println(prefix, k, t.Title)
		}
	}
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func parseKey(s string) (task.Key, bool) {
	key, err := task.ParseKey(s)
	if err == nil {
		return key, true
	}
	fmt.Fprintf(os.Stderr, "could not parse task key: %v\n", err)
	return key, false
}
