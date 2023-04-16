# Workr, another Golang worker pool

[![GoDoc](https://godoc.org/github.com/rvflash/workr?status.svg)](https://godoc.org/github.com/rvflash/workr)
[![Build Status](https://github.com/rvflash/workr/workflows/build/badge.svg)](https://github.com/rvflash/workr/actions?workflow=build)
[![Code Coverage](https://codecov.io/gh/rvflash/workr/branch/main/graph/badge.svg)](https://codecov.io/gh/rvflash/workr)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/workr)](https://goreportcard.com/report/github.com/rvflash/workr)

`workr` provides synchronization, error propagation, context cancellation and execution details 
for groups of tasks running in parallel with a limited number of goroutines.
This number is by default fixed by the number of CPU.

It provides an interface similar to `sync/errgroup` to manage a group of subtasks.

## Install

go 1.20 is the minimum required version due to the usage of `errors.Join`.

```bash
$ go get -u github.com/rvflash/workr
```

## Example

Simple use case that returns the first error occurred.

```go
    g := new(workr.Group)
    g.Go(func() error {
        return nil
    })
    err := g.Wait()
````

A more advanced one that returns all errors that occurred, not only the first one.

```go
    var (
        g   = workr.New(workr.ReturnAllErrors())
        ctx = context.Background()
    )
    g.Go(func() error {
        return doSomething(ctx)
    })
    g.Go(func() error {
        return doSomethingElse(ctx)
    })
    err := g.Wait()
````

It also provides a method `WaitAndReturn` to get details on each task done and 
functions to list those that were successful or `SuccessfulResult` not `FailedResult`.

By creating the worker with a context, the first task on error will cancel it and so, 
all tasks using it are also cancelled.

```go
    oops := errors.New("oops")
	
    g, ctx := workr.WithContext(context.Background(), workr.SetPoolSize(2))
    g.Go(
		func() error {
            return oops
        }, 
		workr.ID(1), 
		workr.SkipError(oops),
	)
    g.Go(
		func() error {
            return doSomething(ctx)
        },
		workr.SetID(2),
	)
    res, err := g.WaitAndReturn()
    if err != nil {
        // No error expected
    }
    log.Println(workr.SuccessfulResult(res).IDList())
```
