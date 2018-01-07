# task

Package task provides way to manage work described by its title and added-started-finished state.

```
$ go get github.com/martindrlik/task
```

Command task uses task package to implement simple task manager for terminal.

```
$ go get github.com/martindrlik/task/cmd/task
$ task
todo read book
todo
(prints) todo 1 read book
start 1
todo
(prints) todo none
doing
(prints) doing 1 read book
```