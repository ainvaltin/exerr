/*
Package exerr implements extended error type/API.

Extensions to stdlib are:
 1. option to attach additional information to error as fields (key - value pair);
 2. errors created by API capture the stack trace of the call stack;

To fully take advantage of these extensions other parts of the system must implement
support for them (ie logger must take care to log that additional information attached
to a error). This package just implements APIs to create these extended errors and
to extract information from errors.
*/
package exerr
