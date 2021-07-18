# Workr, another Golang worker pool

[![GoDoc](https://godoc.org/github.com/rvflash/workr?status.svg)](https://godoc.org/github.com/rvflash/workr)
[![Build Status](https://api.travis-ci.com/rvflash/workr.svg?branch=main)](https://travis-ci.com/rvflash/workr?branch=main)
[![Code Coverage](https://codecov.io/gh/rvflash/workr/branch/main/graph/badge.svg)](https://codecov.io/gh/rvflash/workr)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/workr?)](https://goreportcard.com/report/github.com/rvflash/workr)


`workr` provides synchronization, error propagation, context cancellation and execution details 
for groups of tasks running in parallel with a limited number of goroutines.
This number is by default fixed by the number of CPU.

It provides an interface similar to` sync / errgroup` to manage a group of subtasks.


## Index 

```go
type Group
    func New(opts ...Setting) *Group
    func WithContext(ctx context.Context) (*Group, context.Context)
    func (g *Group) Go(f func() error, opts ...Option)
    func (g *Group) Wait() error
    func (g *Group) WaitAndReturn() ([]*Task, error)

type Option
    func SetID(id interface{}) Option
    func AddErrToSkip(err error) Option

type Setting
    func SetCancel(cancel context.CancelFunc) Setting
    func SetPoolSize(size int) Setting
    func SetQueueSize(size int) Setting

type Task
    func SuccessfulTaskIDs([]*Task) []interface{}
```


## Example

Simple use case.

```go
    g := new(workr.Group)
    g.Go(func() error {
        return nil
    })
    err := g.Wait()
````

It also provides a method `WaitAndReturn` to get details on each task done and one function to list those that were successful.

```go
    oops := errors.New("oops")
    
    g, ctx := workr.WithContext(context.Background())
	g.Go(
		func() error {
		    return oops
	    },
		workr.SetID(1),
		workr.AddErrToSkip(oops), 
	)
    g.Go(
        func() error {
            return nil
        },
        workr.SetID(2),
    )
	res, err := g.WaitAndReturn()
	if err != nil {
        log.Fatalln(err)
    }
	log.Println(workr.SuccessfulTaskIDs(res))
}
```