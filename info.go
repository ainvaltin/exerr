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
	for err != nil {
		fv, ok := err.(interface{ FieldValue(name string) (any, bool) })
		if ok {
			if v, ok := fv.FieldValue(name); ok {
				return v, true
			}
		}

		err = errors.Unwrap(err)
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
	for err != nil {
		if fv, ok := err.(interface{ Fields() map[string]any }); ok {
			if f == nil {
				f = make(map[string]any)
			}
			for k, v := range fv.Fields() {
				f[k] = v
			}
		}
		err = errors.Unwrap(err)
	}

	return f
}

type stacked interface {
	Stack() []uintptr
}

func Stack(err error) []string {
	var se stacked
	for err != nil {
		if s, ok := err.(stacked); ok {
			se = s
		}
		err = errors.Unwrap(err)
	}

	if se != nil {
		return formatStack(se.Stack())
	}
	return nil
}

func formatStack(pcs []uintptr) (r []string) {
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		r = append(r, fmt.Sprintf("%s (%s : %d)", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return r
}
