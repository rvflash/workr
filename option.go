// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package workr

// Option represents a Task option.
type Option func(t *Task)

// ID defines the Task identifier.
func ID(id interface{}) Option {
	return func(t *Task) {
		t.ID = id
	}
}

// Metadata associates these metadata to the task.
func Metadata(list []interface{}) Option {
	return func(t *Task) {
		t.Metadata = append(t.Metadata, list)
	}
}

// SkipError adds an error to skip.
func SkipError(err error) Option {
	return func(t *Task) {
		t.skipped = append(t.skipped, err)
	}
}
