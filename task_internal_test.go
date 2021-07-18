// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package workr

import (
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestNewTask(t *testing.T) {
	var (
		oops  = errors.New("oops")
		again = errors.New("again")
		dt    = map[string]struct {
			in  []Option
			out *Task
		}{
			"Default":       {out: &Task{}},
			"Identifier":    {in: []Option{SetID(math.MaxInt8)}, out: &Task{ID: math.MaxInt8}},
			"Error to skip": {in: []Option{AddErrToSkip(oops)}, out: &Task{ErrToSkip: []error{oops}}},
			"Errors to skip": {
				in:  []Option{AddErrToSkip(oops), AddErrToSkip(again)},
				out: &Task{ErrToSkip: []error{oops, again}},
			},
		}
	)
	t.Parallel()
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := newTask(tt.in...)
			is.New(t).Equal(out, tt.out)
		})
	}
}

func TestTask_Do(t *testing.T) {
	var (
		oops = errors.New("oops")
		dt   = map[string]struct {
			in  *Task
			err error
			msg string
		}{
			"Default": {
				in:  new(Task),
				err: ErrPanic,
				msg: "worker pool: panic: task ID(<nil>): runtime error",
			},
			"Error": {
				in: &Task{Func: func() error {
					return oops
				}},
				err: oops,
				msg: "worker pool: task ID(<nil>): oops",
			},
			"OK": {
				in: &Task{Func: func() error {
					return nil
				}},
			},
		}
	)
	t.Parallel()
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tt.in.do()
			are := is.New(t)
			are.True(errors.Is(tt.in.Err, tt.err))                                     // mismatch error
			are.True(tt.in.Err == nil || strings.HasPrefix(tt.in.Err.Error(), tt.msg)) // mismatch error message
		})
	}
}

func TestTask_Do2(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("panic expected")
		}
	}()
	var job *Task
	job.do()
}

func TestTask_IsErrSkipped(t *testing.T) {
	var (
		oops = errors.New("oops")
		dt   = map[string]struct {
			in  *Task
			out bool
		}{
			"Default": {in: new(Task)},
			"Blank":   {in: &Task{Err: oops, ErrToSkip: []error{}}},
			"OK":      {in: &Task{Err: oops, ErrToSkip: []error{oops}}, out: true},
		}
	)
	t.Parallel()
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := tt.in.isErrSkipped()
			is.New(t).Equal(tt.out, out)
		})
	}
}
