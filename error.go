package workr

type errWorkr string

// Error implements the error interface.
func (e errWorkr) Error() string {
	return "workr: " + string(e)
}

// ErrPanic is returned when a task or the worker panics.
const ErrPanic = errWorkr("panic recovered")
