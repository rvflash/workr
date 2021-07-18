// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package workr

// Option represents a Task option.
type Option func(t *Task)

func setFunc(f func() error) Option {
	return func(t *Task) {
		t.Func = f
	}
}

// SetID defines the Task identifier.
func SetID(id interface{}) Option {
	return func(t *Task) {
		t.ID = id
	}
}

// AddErrToSkip adds an error to skip.
func AddErrToSkip(err error) Option {
	return func(t *Task) {
		t.ErrToSkip = append(t.ErrToSkip, err)
	}
}
