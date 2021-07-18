package workr

import (
	"errors"
	"fmt"
)

// SuccessfulTaskIDs returns the list of successful Task identifiers.
func SuccessfulTaskIDs(tasks []*Task) []interface{} {
	n := len(tasks)
	if n == 0 {
		return nil
	}
	res := make([]interface{}, 0, n)
	for _, t := range tasks {
		if t.Err == nil {
			res = append(res, t.ID)
		}
	}
	return res
}

func newTask(opts ...Option) *Task {
	t := new(Task)
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Task represents a Task.
type Task struct {
	ID        interface{}
	Err       error
	ErrToSkip []error
	Func      func() error
}

func (t *Task) do() {
	defer func() {
		r := recover()
		switch {
		case r != nil:
			// Do not panic: throw an error.
			t.Err = fmt.Errorf("%w: task ID(%v): %v", ErrPanic, t.ID, r)
		case t.Err != nil:
			t.Err = fmt.Errorf("worker pool: task ID(%v): %w", t.ID, t.Err)
		}
	}()
	t.Err = t.Func()
}

func (t Task) isErrSkipped() bool {
	for _, target := range t.ErrToSkip {
		if errors.Is(t.Err, target) {
			return true
		}
	}
	return false
}
