package exerr

import (
	"errors"
	"fmt"
	"runtime"
)

/*
FieldValue returns the first value of the field named "name" in the error chain.
To check does only this particular err have the field check does it implement

	FieldValue(name string) (any, bool)

method and if it does call it.
*/
func FieldValue(err error, name string) (value any, ok bool) {
	for ; err != nil; err = errors.Unwrap(err) {
		if fv, ok := err.(interface{ FieldValue(name string) (any, bool) }); ok {
			if v, ok := fv.FieldValue(name); ok {
				return v, true
			}
		}
	}

	return nil, false
}

/*
Fields returns all the fields in the error chain.
If multiple errors do have a field with the same name random one ends up in the result.

To check which fields this particular err has check does it implement

	Fields() map[string]any

method and if it does call it.
*/
func Fields(err error) map[string]any {
	var f map[string]any
	for ; err != nil; err = errors.Unwrap(err) {
		if fv, ok := err.(interface{ Fields() map[string]any }); ok {
			if f == nil {
				f = make(map[string]any)
			}
			for k, v := range fv.Fields() {
				f[k] = v
			}
		}
	}

	return f
}

type stacked interface {
	PC() []uintptr
}

func Stack(err error) []string {
	var se stacked
	for ; err != nil; err = errors.Unwrap(err) {
		if s, ok := err.(stacked); ok {
			se = s
		}
	}

	if se != nil {
		return formatStack(se.PC())
	}
	return nil
}

func formatStack(pcs []uintptr) (r []string) {
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		r = append(r, fmt.Sprintf("%s (%s:%d)", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return r
}
