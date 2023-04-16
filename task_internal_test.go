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
			// inputs
			opts []Option
			// outputs
			out *Task
		}{
			"Default":       {out: &Task{}},
			"Identifier":    {opts: []Option{ID(math.MaxInt8)}, out: &Task{ID: math.MaxInt8}},
			"Error to skip": {opts: []Option{SkipError(oops)}, out: &Task{skipped: []error{oops}}},
			"Errors to skip": {
				opts: []Option{SkipError(oops), SkipError(again)},
				out:  &Task{skipped: []error{oops, again}},
			},
		}
	)
	t.Parallel()
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := newTask(nil, tt.opts...)
			is.New(t).Equal(out, tt.out)
		})
	}
}

func TestTask_Do(t *testing.T) {
	var (
		oops = errors.New("oops")
		dt   = map[string]struct {
			// inputs
			in *Task
			// outputs
			ok  bool
			err error
			msg string
		}{
			"Default": {
				in:  new(Task),
				err: ErrPanic,
				msg: "runtime error",
			},
			"Error": {
				in: &Task{f: func() error {
					return oops
				}},
				err: oops,
			},
			"OK": {
				in: &Task{f: func() error {
					return nil
				}},
				ok: true,
			},
		}
	)
	t.Parallel()
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var (
				ok  = tt.in.do()
				are = is.New(t)
			)
			are.Equal(tt.ok, ok)                                                      // mismatch result
			are.True(errors.Is(tt.in.err, tt.err))                                    // mismatch error
			are.True(tt.in.err == nil || strings.Contains(tt.in.err.Error(), tt.msg)) // mismatch error message
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
	is.New(t).True(!job.do())
}

func TestTask_IsErrSkipped(t *testing.T) {
	var (
		oops = errors.New("oops")
		dt   = map[string]struct {
			in  *Task
			out bool
		}{
			"Default": {in: new(Task)},
			"Blank":   {in: &Task{err: oops, skipped: []error{}}},
			"OK":      {in: &Task{err: oops, skipped: []error{oops}}, out: true},
		}
	)
	t.Parallel()
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := tt.in.ErrorSkipped()
			is.New(t).Equal(tt.out, out)
		})
	}
}
