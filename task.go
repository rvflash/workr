// Copyright (c) 2023 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package workr

import (
	"errors"
	"fmt"
)

// FailedResult only returns the list of failed tasks.
func FailedResult(list Result) Result {
	return tasks(list, true)
}

// SuccessfulResult only returns the list of successful tasks.
func SuccessfulResult(list Result) Result {
	return tasks(list, false)
}

func tasks(list Result, failed bool) Result {
	res := make(Result, 0, len(list))
	for _, t := range list {
		if t.failed() == failed {
			res = append(res, t)
		}
	}
	return res
}

// Result is a list of Task.
type Result []*Task

// IDList returns the list of task's identifiers.
func (r Result) IDList() []interface{} {
	n := len(r)
	if n == 0 {
		return nil
	}
	res := make([]interface{}, n)
	for k, t := range r {
		res[k] = t.ID
	}
	return res
}

// Error returns all errors occurred.
func (r Result) Error() error {
	var err error
	for _, t := range r {
		err = errors.Join(err, t.err)
	}
	return err
}

// Metadata returns the aggregation of all Metadata.
func (r Result) Metadata() []interface{} {
	n := len(r)
	if n == 0 {
		return nil
	}
	// Best effort of allocations.
	res := make([]interface{}, 0, len(r[0].Metadata))
	for _, t := range r {
		res = append(res, t.Metadata...)
	}
	return res
}

func newTask(f func() error, opts ...Option) *Task {
	t := &Task{f: f}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Task represents a Task.
type Task struct {
	// ID uniquely identified a Task.
	ID interface{}
	// Metadata is the list of Task's metadata.
	Metadata []interface{}

	skipped []error
	f       func() error
	err     error
}

// Error returns the result of the task.
func (t *Task) Error() error {
	if t.err != nil {
		if t.ID != nil {
			return fmt.Errorf("workr task ID(%v): %w", t.ID, t.err)
		}
		return t.err
	}
	return nil
}

// ErrorSkipped returns true if the error must be ignored.
func (t *Task) ErrorSkipped() bool {
	if t.err == nil || len(t.skipped) == 0 {
		return false
	}
	for _, target := range t.skipped {
		if errors.Is(t.err, target) {
			return true
		}
	}
	return false
}

func (t *Task) do() (ok bool) {
	defer func() {
		r := recover()
		if r != nil {
			// Do not panic: throw an error.
			t.err = fmt.Errorf("%w: %v", ErrPanic, r)
		}
		ok = !t.failed()
	}()
	t.err = t.f()

	return
}

func (t *Task) failed() bool {
	return t == nil || (t.err != nil && !t.ErrorSkipped())
}
