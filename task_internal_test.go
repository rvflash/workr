// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package workr

import (
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
			"Default": {out: &Task{}},
			"Identifier": {
				opts: []Option{ID(math.MaxInt8)},
				out:  &Task{ID: math.MaxInt8},
			},
			"Metadata": {
				opts: []Option{Metadata([]interface{}{math.MaxInt8, math.MaxInt16})},
				out:  &Task{Metadata: []interface{}{math.MaxInt8, math.MaxInt16}},
			},
			"Error to skip": {
				opts: []Option{SkipError(oops)},
				out:  &Task{skipped: []error{oops}},
			},
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
			is.New(t).Equal("", cmp.Diff(tt.out, out, cmp.AllowUnexported(Task{}), cmpopts.EquateErrors()))
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

func TestTask_ErrorSkipped(t *testing.T) {
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

func TestSuccessfulResult(t *testing.T) {
	var (
		oops = errors.New("oops")
		dt   = map[string]struct {
			in   Result
			ids  []interface{}
			data []interface{}
		}{
			"Default": {},
			"Blank":   {in: Result{}},
			"Mixed": {
				in: Result{
					{ID: 1},
					{ID: 2, Metadata: []interface{}{1, 2}, err: oops},
					{ID: 3, Metadata: []interface{}{3}},
				},
				ids:  []interface{}{1, 3},
				data: []interface{}{3},
			},
			"OK": {
				in: Result{
					{ID: 1, Metadata: []interface{}{1, 2}},
					{ID: 2},
					{ID: 3, Metadata: []interface{}{3, 4}},
				},
				ids:  []interface{}{1, 2, 3},
				data: []interface{}{1, 2, 3, 4},
			},
		}
	)
	t.Parallel()
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var (
				are = is.New(t)
				out = SuccessfulResult(tt.in)
			)
			are.Equal(tt.ids, out.IDList())    // mismatch identifiers
			are.Equal(tt.data, out.Metadata()) // mismatch metadata
		})
	}
}

func TestFailedResult(t *testing.T) {
	var (
		oops = errors.New("oops")
		dt   = map[string]struct {
			in   Result
			ids  []interface{}
			data []interface{}
		}{
			"Default": {},
			"Blank":   {in: Result{}},
			"Mixed": {
				in: Result{
					{ID: 1, Metadata: []interface{}{1, 2}},
					{ID: 2, Metadata: []interface{}{3, 4}, err: oops},
					{ID: 3, Metadata: []interface{}{5, 6}},
				},
				ids:  []interface{}{2},
				data: []interface{}{3, 4},
			},
			"OK": {
				in: Result{
					{ID: 1, Metadata: []interface{}{1, 2}},
					{ID: 2, Metadata: []interface{}{3, 4}},
					{ID: 3, Metadata: []interface{}{5, 6}},
				},
			},
		}
	)
	t.Parallel()
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var (
				are = is.New(t)
				out = FailedResult(tt.in)
			)
			are.Equal(tt.ids, out.IDList())    // mismatch identifiers
			are.Equal(tt.data, out.Metadata()) // mismatch metadata
		})
	}
}
