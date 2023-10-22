# exerr

This package aims to solve two problems:

 1. need to attach additional data to a error (which is not part of the error message);
 2. provide stack trace of the location where the error was "raised";

while being as "stdlib like" as possible.

## Attaching additional data to a error

One problem with error handling and logging is that sometimes there is additional
information available at the place where the error happens which isn't suitable
to be included into error message but which would be useful when investigating
the error.

An example would be failing database query ‒ it's useful to have the SQL statement
and it's parameters available in addition to the error message but adding them to
the message is generally not acceptable.

Usually this means that the extra info is logged at the error site and error
returned to he caller, to be logged (again) at some point up in the call chain.
This means that when investigating the error one has somehow realize that these
two log records describe the same incident...

To solve the problem this package implements error type which has a option to
attach "fields" to it which are not visible in the error message but which can
be logged by logger ‒ so now there is no need to log at the error site and thus
every error gets logged only once.

```go
query := "select x, y, z from t where i=? and k=?"
param := []interface{}{42, "foo"}
rows, err := db.QueryContext(ctx, query, param...)
if err != nil {
	return exerr.Errorf("failed to open query: %w", err).AddField("sql_query", query).AddField("sql_param", param)
}
```

As a bonus the logger doesn't have to be available for the code which deals
with the database meaning there is one less dependency to pass down!


## Stack trace of the error

[Go proverb](https://go-proverbs.github.io/) says _"Don't just check errors,
handle them gracefully"_. Among other things this means that when error is
passed up to the caller context is added to it using `fmt.Errorf`. This
additional context makes it usually easy enough to follow the error path from
the know location where it was logged down to the error site. 
Still there are developers who would like to have stack trace of error.

Errors created by this library (ie `exerr.New` and `exerr.Errorf`) also record
the source code location so that stack trace can be logged when logging the
error. 

The stack trace is not included into error message, it has to be logged
separately (ie logger would have to support this feature).


## Possible improvements

 - use [slog.Attr](https://pkg.go.dev/log/slog) for fields;
 - more flexible stack formatting;
 - integration with popular log libraries;
