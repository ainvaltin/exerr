package exerr

import (
	"runtime"
)

func newExErr(err error) *exErr {
	pcs := make([]uintptr, 16)
	n := runtime.Callers(3, pcs)
	return &exErr{err: err, pcs: pcs[:n:n]}
}

// error with location info (stack trace) and optional metadata (fields).
type exErr struct {
	err    error
	pcs    []uintptr
	fields map[string]any
}

func (e *exErr) Unwrap() error { return e.err }

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

func (e *exErr) Fields() map[string]any { return e.fields }

/*
Stack returns return program counters of function invocations on the the place error was created.
*/
func (e *exErr) Stack() []uintptr { return e.pcs }
