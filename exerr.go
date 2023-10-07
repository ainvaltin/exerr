package exerr

import (
	"errors"
	"runtime"
)

func newExErr(err error) *exErr {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(3, pcs) // 1=newExErr; 2=Errorf|AddField|New; 3=user code
	return &exErr{err: err, pcs: pcs[:n:n]}
}

// error with location info (stack trace) and optional metadata (fields).
type exErr struct {
	err    error
	pcs    []uintptr
	fields map[string]any
}

func (e *exErr) As(target any) bool { return errors.As(e.err, target) }

func (e *exErr) Is(target error) bool { return errors.Is(e.err, target) }

func (e *exErr) Unwrap() error { return errors.Unwrap(e.err) }

func (e *exErr) Error() string {
	if e.err == nil {
		return "<nil>"
	}
	return e.err.Error()
}

func (e *exErr) AddField(name string, value any) ErrorWithFields {
	if e.fields == nil {
		e.fields = make(map[string]any)
	}
	e.fields[name] = value

	return e
}

func (e *exErr) FieldValue(name string) (any, bool) {
	v, ok := e.fields[name]
	return v, ok
}

/*
Fields returns the name -> value map of the fields attached to the error.

Should be considered to be read-only, ie do not modify!
*/
func (e *exErr) Fields() map[string]any { return e.fields }

/*
PC returns return program counters of function invocations on the the place error was created.
*/
func (e *exErr) PC() []uintptr { return e.pcs }
