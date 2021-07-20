// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package workr

import (
	"testing"

	"github.com/matryer/is"
)

func TestNumCPU(t *testing.T) {
	is.New(t).True(numCPU() > 0)
}

func TestNew(t *testing.T) {
	var (
		cpu = numCPU()
		def = &Group{poolSize: cpu, queueSize: cpu}
		dt  = map[string]struct {
			in  []Setting
			out *Group
		}{
			"Default":            {out: def},
			"Invalid pool size":  {in: []Setting{SetPoolSize(-1)}, out: def},
			"Invalid queue size": {in: []Setting{SetQueueSize(-1)}, out: def},
			"Pool size":          {in: []Setting{SetPoolSize(1)}, out: &Group{poolSize: 1, queueSize: cpu}},
			"Queue size":         {in: []Setting{SetQueueSize(1)}, out: &Group{poolSize: cpu, queueSize: 1}},
		}
	)
	t.Parallel()
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := New(tt.in...)
			is.New(t).Equal(tt.out, out)
		})
	}
}

func TestNew2(t *testing.T) {
	var (
		are = is.New(t)
		cpu = numCPU()
		g   Group
	)
	are.NoErr(g.Wait())
	are.Equal(cpu, g.poolSize)  // mismatch pool size
	are.Equal(cpu, g.queueSize) // mismatch queue size
}
