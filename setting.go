// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package workr

import "context"

// Setting represents a Group setting.
type Setting func(*Group)

// SetCancel defines the cancellation function to be called when one of the tasks of the group has failed.
func SetCancel(cancel context.CancelFunc) Setting {
	return func(g *Group) {
		g.cancel = cancel
	}
}

// SetPoolSize defines the pool size.
func SetPoolSize(size int) Setting {
	return func(g *Group) {
		if size > 0 {
			g.poolSize = size
		}
	}
}

// SetQueueSize defines the queue size.
func SetQueueSize(size int) Setting {
	return func(g *Group) {
		if size > 0 {
			g.queueSize = size
		}
	}
}

// ReturnAllErrors enables the error reporting to all.
func ReturnAllErrors() Setting {
	return func(g *Group) {
		g.errAll = true
	}
}
