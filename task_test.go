// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package workr_test

/*
import (
	"errors"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/workr"
)

func TestErrors(t *testing.T) {
	var (
		are  = is.New(t)
		oops = errors.New("oops")
		dt   = map[string]struct {
			in  func() ([]*workr.Task, error)
			out error
		}{
			"Default": {
				in: func() ([]*workr.Task, error) {
					return nil, nil
				},
			},
			"Blank": {
				in: func() ([]*workr.Task, error) {
					return []*workr.Task{}, nil
				},
			},
			"Mixed": {
				in: func() ([]*workr.Task, error) {
					return []*workr.Task{
						{ID: 1},
						{ID: 2, err: oops},
						{ID: 3},
					}, oops
				},
				out: oops,
			},
			"OK": {
				in: func() ([]*workr.Task, error) {
					return []*workr.Task{
						{ID: 1, err: oops},
						{ID: 2},
						{ID: 3, err: oops},
					}, oops
				},
				out: errors.Join(oops, oops),
			},
		}
	)
	t.Parallel()
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := workr.Errors(tt.in())
			if tt.out != nil {
				are.Equal(out.Error(), tt.out.Error()) // mismatch error message
			} else {
				are.NoErr(out) // unexpected error
			}
		})
	}
}

func TestSuccessfulTaskIDs(t *testing.T) {
	var (
		oops = errors.New("oops")
		dt   = map[string]struct {
			in  []*workr.Task
			out []interface{}
		}{
			"Default": {},
			"Blank":   {in: []*workr.Task{}},
			"Mixed":   {in: []*workr.Task{{ID: 1}, {ID: 2, err: oops}, {ID: 3}}, out: []interface{}{1, 3}},
			"OK":      {in: []*workr.Task{{ID: 1}, {ID: 2}, {ID: 3}}, out: []interface{}{1, 2, 3}},
		}
	)
	t.Parallel()
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := workr.SuccessfulTaskIDs(tt.in)
			is.New(t).Equal(tt.out, out)
		})
	}
}

*/
