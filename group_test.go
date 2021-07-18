// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package workr_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/rvflash/workr"
)

func TestGroup_Go(t *testing.T) {
	var g *workr.Group
	g.Go(func() error {
		return errors.New("oops")
	})
	is.New(t).True(errors.Is(g.Wait(), workr.ErrPanic))
}

func TestGroup_Wait(t *testing.T) {
	var g *workr.Group
	is.New(t).True(errors.Is(g.Wait(), workr.ErrPanic))
}

func TestGroup_WaitAndReturn(t *testing.T) {
	var (
		oops = errors.New("oops")
		dt   = map[string]struct {
			// inputs
			conf  []workr.Setting
			tasks []func() error
			opts  [][]workr.Option
			// outputs
			err error
			msg string
			okN int
		}{
			"Default": {},
			"No function": {
				tasks: []func() error{nil},
				opts: [][]workr.Option{
					{workr.SetID(1)},
				},
				err: workr.ErrPanic,
				msg: "worker pool: panic: task ID(1): runtime error:",
			},
			"Return in error": {
				tasks: []func() error{
					func() error {
						return nil
					},
					func() error {
						return oops
					},
					func() error {
						return nil
					},
				},
				opts: [][]workr.Option{
					{workr.SetID(1)},
					{workr.SetID(2)},
					{workr.SetID(3)},
				},
				okN: 2,
				err: oops,
				msg: "worker pool: task ID(2):",
			},
			"Error skipped": {
				tasks: []func() error{
					func() error {
						return nil
					},
					func() error {
						return oops
					},
					func() error {
						return nil
					},
				},
				opts: [][]workr.Option{
					{workr.SetID(1)},
					{workr.SetID(2), workr.AddErrToSkip(oops)},
					{workr.SetID(3)},
				},
				okN: 2,
			},
			"OK": {
				tasks: []func() error{
					func() error {
						return nil
					},
					func() error {
						return nil
					},
					func() error {
						return nil
					},
				},
				opts: [][]workr.Option{
					{workr.SetID(1)},
					{workr.SetID(2)},
					{workr.SetID(3)},
				},
				okN: 3,
			},
		}
	)
	t.Parallel()
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := workr.New(tt.conf...)
			for k, v := range tt.tasks {
				g.Go(v, tt.opts[k]...)
			}
			var (
				res, err = g.WaitAndReturn()
				are      = is.New(t)
			)
			are.True(errors.Is(err, tt.err))                               // mismatch error
			are.True(err == nil || strings.HasPrefix(err.Error(), tt.msg)) // mismatch error message
			are.Equal(tt.okN, len(workr.SuccessfulTaskIDs(res)))           // mismatch successful task number
		})
	}
}

func TestGroup_WaitAndReturn2(t *testing.T) {
	var (
		g   *workr.Group
		are = is.New(t)
	)
	res, err := g.WaitAndReturn()
	are.True(errors.Is(err, workr.ErrPanic)) // mismatch error
	are.Equal(nil, res)                      // unexpected content
}

func TestWithContext(t *testing.T) {
	are := is.New(t)
	g, ctx := workr.WithContext(context.Background())
	are.True(ctx != nil) // context required
	are.NoErr(g.Wait())  // unexpected error
}

func TestWithContext2(t *testing.T) {
	parent, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	g, ctx := workr.WithContext(parent)
	g.Go(func() error {
		return nil
	})
	g.Go(func() error {
		for {
			<-ctx.Done()
			return ctx.Err()
		}
	})
	err := g.Wait()
	is.New(t).True(errors.Is(err, context.DeadlineExceeded))
}
