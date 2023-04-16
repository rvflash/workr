// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package workr_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/workr"
)

func TestErrWorkr_Error(t *testing.T) {
	is.New(t).Equal("workr: panic recovered", workr.ErrPanic.Error())
}
