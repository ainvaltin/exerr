package exerr

import (
	"errors"
	"fmt"
)

/*
ErrorWithFields is extended error type which makes it easy to add fields to errors
returned by [New], [Errorf] and [AddField] using method chaining (fluent interface
style).
*/
type ErrorWithFields interface {
	error
	AddField(name string, value any) ErrorWithFields
}

/*
New is like [errors.New] but it returns [ErrorWithFields] which makes it easy to chain AddField
calls to the error. It also captures the location in the source code where the error was created.

Do not use this func to create sentinel errors - for that [errors.New] should be used!
*/
func New(text string) ErrorWithFields {
	return newExErr(errors.New(text))
}

/*
Errorf is like [fmt.Errorf] but it returns [ErrorWithFields] which makes it easy to
chain AddField calls to the error.

Errorf also captures the location in the source code where the error was created, to access
that information [Stack] func can be used or the error could be checked for implementing

	func Stack() []uintptr

which returns return program counters of function invocations on the place the error was created.

Do not use this func to create sentinel errors - for that [errors.New] should be used.
*/
func Errorf(format string, a ...any) ErrorWithFields {
	return newExErr(fmt.Errorf(format, a...))
}

/*
AddField allows to attach fields to a error without adding any new message to the error.

It is like [Errorf] except when the "err" parameter implements ErrorWithFields the field
is added to the "err" instead of creating new wrapper error.

In case nil is sent as "err" parameter AddField returns non nil error with field attached,
the error message in that case will be "<nil>". It is recommended not to use nil for the
err parameter.
*/
func AddField(err error, name string, value any) ErrorWithFields {
	if af, ok := err.(ErrorWithFields); ok {
		return af.AddField(name, value)
	}
	return newExErr(err).AddField(name, value)
}
