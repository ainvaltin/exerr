# exerr

This package aims to solve two problems:

 1. need to attach additional data to a error (which is not part of the error message);
 2. provide stack trace of the location where the error was "raised";

while being as "stdlib like" as possible.

## Attaching additional data to a error

Main problem I have with logging / error handling is that I like to have logging only in the
"top level handler" but sometimes it would be useful to log some additional info with the error
(without having that info in the error message). Good example is when database query fails - it would
be nice to log the query and it's parameters with the error but not as part of the error message.

Obvious solution is to log the query in the method where database is accessed but this has two problems:
1. need to have logger available in the method (ie dependency must be passed down);
2. we then have two "detached" messages in the log while I'd prefer to have a single log entry for each failure;

So to solve this problem this package implements error type which has a option to attach "fields" to it
which are not visible in the error message but which can be logged by a logger.

```go
query := "select x, y, z from t where i=? and k=?"
qparam := []interface{}{42, "foo"}
rows, err := db.QueryContext(ctx, query, qparam...)
if err != nil {
	return exerr.Errorf("failed to open query: %w", err).AddField("sql_query", query).AddField("sql_param", qparam)
}
```

## Stack trace of the error

In my experience developers coming from languages which use exceptions complain about
Go errors not having stack trace attached to them. I personally am not bothered by this - I find
that when I have a trace of the log site and error handling is done properly (ie context is added
to the error every time it is handled) it is easy to start from the log site and move "down the call
chain" to the place error originates from. Doing it that way (rather than having the stack trace
of the first error) gives me much better understanding what happened and what went wrong...

Some third-party error libraries have methods like `err.WithStack()` to record that information but I
do not like that. So errors created with this library capture the stack by default and make it
available for loggers via helper methods.

## Possible improvements

 - integrate with [slog](https://pkg.go.dev/golang.org/x/exp/slog);
 - more flexible stack formatting;
 - unwrap behavior - currently `Errorf` results as two errors in the chain, should it behave like single error?
 - integration with popular log libraries;
