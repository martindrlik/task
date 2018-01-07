package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/martindrlik/task/task"
)

var tasks = task.NewManager()

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		proc(s)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "could not read stdin: %v\n", err)
		os.Exit(1)
	}
}

func proc(s string) {
	if keys, ok := list(s); ok {
		printTasks(s, keys)
	} else if title := strings.TrimPrefix(s, "todo "); title != s {
		tasks.Add(time.Now(), title)
	} else if keyString := strings.TrimPrefix(s, "start "); keyString != s {
		if key, ok := parseKey(keyString); ok {
			tasks.Start(time.Now(), key)
		}
	} else if keyString := strings.TrimPrefix(s, "done "); keyString != s {
		if key, ok := parseKey(keyString); ok {
			tasks.Finish(time.Now(), key)
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
		fmt.Println(prefix, k, t.Title)
	}
}

func parseKey(s string) (task.Key, bool) {
	key, err := task.ParseKey(s)
	if err == nil {
		return key, true
	}
	fmt.Fprintf(os.Stderr, "could not parse task key: %v\n", err)
	return key, false
}
