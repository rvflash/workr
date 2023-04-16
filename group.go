// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package workr provides synchronization, error propagation, context cancellation and execution details
// for groups of tasks running in parallel with a limited number of goroutines.
package workr

import (
	"context"
	"runtime"
	"sync"
)

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs first.
func WithContext(parent context.Context, opts ...Setting) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(parent)
	return New(append([]Setting{SetCancel(cancel)}, opts...)...), ctx
}

// New returns a new instance of Group based on the given settings.
func New(opts ...Setting) *Group {
	g := new(Group)
	for _, opt := range append([]Setting{
		setDefaultPoolSize(),
		setDefaultQueueSize(),
	}, opts...) {
		opt(g)
	}
	return g
}

// Group is a collection of tasks using a pool of workers to call their associated function in parallel.
// A zero Group is valid and does not cancel on error.
// If a cancellation method is associated to the Group,
// the first call to return a non-nil error cancels the group.
// Its error will be returned by Wait and WaitAndReturn.
type Group struct {
	poolSize  int
	queueSize int
	cancel    context.CancelFunc

	rs Result
	ch chan *Task
	wg sync.WaitGroup

	errAll  bool
	errOnce sync.Once
	err     error
}

const once = 1

// Go creates with the given function and associated options a Task.
// This Task will be running opts one of workers pool goroutines.
func (g *Group) Go(f func() error, opts ...Option) {
	if !g.ready() {
		return
	}
	t := newTask(f, opts...)
	g.rs = append(g.rs, t)
	g.wg.Add(once)
	g.ch <- t
}

// Wait blocks until all tasks function calls have returned,
// then returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	if !g.ready() {
		return ErrPanic
	}
	close(g.ch)

	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	if g.errAll {
		return g.rs.Error()
	}
	return g.err
}

// WaitAndReturn does the same job as Wait but also returns details about performing tasks.
func (g *Group) WaitAndReturn() (Result, error) {
	err := g.Wait()
	if g == nil {
		return nil, err
	}
	return g.rs, err
}

func (g *Group) init() {
	for i := 0; i < g.poolSize; i++ {
		go func(ch <-chan *Task) {
			for t := range ch {
				func(t *Task) {
					defer g.wg.Done()
					if !t.do() {
						g.errOnce.Do(func() {
							g.err = t.Error()
							if g.cancel != nil {
								g.cancel()
							}
						})
					}
				}(t)
			}
		}(g.ch)
	}
}

func (g *Group) ready() bool {
	if g == nil {
		// Avoid panic.
		return false
	}
	if g.poolSize < 1 {
		setDefaultPoolSize()(g)
	}
	if g.queueSize < 1 {
		setDefaultQueueSize()(g)
	}
	if g.ch == nil {
		g.rs = make(Result, 0, g.queueSize)
		g.ch = make(chan *Task, g.queueSize)
		g.init()
	}
	return true
}

func setDefaultPoolSize() Setting {
	return func(g *Group) {
		g.poolSize = runtime.NumCPU()
	}
}

func setDefaultQueueSize() Setting {
	return func(g *Group) {
		g.queueSize = runtime.NumCPU()
	}
}
