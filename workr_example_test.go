// Copyright (c) 2023 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package workr_test

import (
	"errors"
	"fmt"
	"log"

	"github.com/rvflash/workr"
)

func ExampleGroup_Wait() {
	var (
		f = func() error {
			return nil
		}
		g = new(workr.Group)
	)
	g.Go(f)
	fmt.Println(g.Wait())
	// Output: <nil>
}

func ExampleReturnAllErrors() {
	g := workr.New(workr.ReturnAllErrors())
	g.Go(func() error {
		return errors.New("oops")
	})
	g.Go(func() error {
		return errors.New("oops")
	})
	fmt.Println(g.Wait())
	// Output: oops
	// oops
}

func ExampleResult_IDList() {
	g := workr.New()
	g.Go(
		func() error {
			return nil
		},
		workr.ID(1),
	)
	g.Go(
		func() error {
			return nil
		},
		workr.ID(2),
	)
	res, err := g.WaitAndReturn()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(len(res.IDList()))
	// Output: 2
}
