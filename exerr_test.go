package exerr

import (
	"errors"
	"fmt"
	"net"
	"testing"
)

func Test_newExErr(t *testing.T) {
	t.Parallel()

	t.Run("nil error as input", func(t *testing.T) {
		e := newExErr(nil)
		if e == nil {
			t.Fatal("unexpectedly nil was returned")
		}

		if e.err != nil {
			t.Fatal("unexpectedly e.err is not nil")
		}
		if e.err != e.Unwrap() {
			t.Fatal("Unwrap returns unexpected value")
		}

		if len(e.pcs) == 0 {
			t.Fatal("unexpectedly e.pcs is not assigned")
		}
		if len(e.pcs) != len(e.PC()) {
			t.Fatalf("unexpectedly e.PC() returns slice with different length (%d vs %d)", len(e.pcs), len(e.PC()))
		}

		if e.Error() != "<nil>" {
			t.Fatalf(`expected error message to be "<nil>" got %q`, e.Error())
		}

		if e.fields != nil {
			t.Fatal("unexpectedly e.fields is not nil")
		}
	})

	t.Run("non nil error as input", func(t *testing.T) {
		expErr := fmt.Errorf("some error")
		e := newExErr(expErr)
		if e == nil {
			t.Fatal("unexpectedly nil was returned")
		}

		if e.err != expErr {
			t.Fatal("unexpectedly e.err is not the error passed to the constructor")
		}
		if e.Unwrap() != nil {
			t.Fatal("Unwrap returns unexpected value")
		}

		if len(e.pcs) == 0 {
			t.Fatal("unexpectedly e.pcs is not assigned")
		}
		if len(e.pcs) != len(e.PC()) {
			t.Fatal("unexpectedly e.PC() returns slice with different length")
		}

		if e.Error() != expErr.Error() {
			t.Fatalf(`expected %q got %q`, expErr.Error(), e.Error())
		}

		if e.fields != nil {
			t.Fatal("unexpectedly e.fields is not nil")
		}
	})
}

func Test_exErr_Fields(t *testing.T) {
	t.Parallel()

	t.Run("no fields", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		if n := len(e.fields); n != 0 {
			t.Fatalf("expected no fields, got %d", n)
		}
		v, ok := e.FieldValue("foo")
		if ok {
			t.Fatalf(`unexpectedly field named "foo" was found: %v`, v)
		}
	})

	expectField := func(t *testing.T, err *exErr, name string, value any) {
		t.Helper()

		v, ok := err.FieldValue(name)
		if !ok {
			t.Fatalf(`field named %q was not found`, name)
		}
		if v != value {
			t.Fatalf("expected value to be %v, got %v", value, v)
		}
	}

	t.Run("one field", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		e.AddField("name", 42)
		if n := len(e.fields); n != 1 {
			t.Fatalf("expected 1 field, got %d", n)
		}

		expectField(t, e, "name", 42)
	})

	t.Run("two fields, fluid", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		e.AddField("first", 1).AddField("second", 2)
		if n := len(e.fields); n != 2 {
			t.Fatalf("expected 2 fields, got %d", n)
		}

		expectField(t, e, "first", 1)
		expectField(t, e, "second", 2)
	})

	t.Run("two fields, two calls", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		e.AddField("first", 1)
		e.AddField("second", 2)
		if n := len(e.fields); n != 2 {
			t.Fatalf("expected 2 fields, got %d", n)
		}

		expectField(t, e, "first", 1)
		expectField(t, e, "second", 2)
	})

	t.Run("adding the same field twice", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		e.AddField("foo", 1)
		e.AddField("foo", 2)
		if n := len(e.fields); n != 1 {
			t.Fatalf("expected 1 field, got %d", n)
		}

		expectField(t, e, "foo", 2)
	})

	t.Run("", func(t *testing.T) {})
}

func Test_exErr_FieldValue(t *testing.T) {
	t.Parallel()

	t.Run("field doesn't exist", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		v, ok := e.FieldValue("foo")
		if ok {
			t.Fatalf(`unexpectedly field named "foo" was found: %v`, v)
		}
	})

	t.Run("field exists", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		e.AddField("foo", 42)
		v, ok := e.FieldValue("foo")
		if !ok {
			t.Fatalf(`unexpectedly field named "foo" was not found`)
		}
		if v != 42 {
			t.Errorf("expected value to be 42, got %v", v)
		}
	})
}

func Test_exErr_Is(t *testing.T) {
	t.Parallel()
	// exerr acts as a container so Is, As and Unwrap should skip it and return the wrapped error

	t.Run("wrapping nil error", func(t *testing.T) {
		e := newExErr(nil)
		if e == nil {
			t.Fatal("unexpectedly nil was returned")
		}
		if !e.Is(nil) {
			t.Fatal("Is returns unexpected value")
		}
		if e.Is(e) {
			t.Fatal("Is returns unexpected value")
		}
	})

	t.Run("wrapping non-nil error", func(t *testing.T) {
		expErr := fmt.Errorf("some error")
		e := newExErr(expErr)
		if !e.Is(expErr) {
			t.Fatal("Is returns unexpected value")
		}
		if e.Is(e) {
			t.Fatal("Is returns unexpected value")
		}
	})
}

func Test_errors_Is(t *testing.T) {
	t.Parallel()
	// test that errors.Is returns expected results

	t.Run("nil error", func(t *testing.T) {
		var e *exErr
		if errors.Is(e, nil) {
			t.Fatal("errors.Is returns unexpected value")
		}
	})

	t.Run("errors.Is(self, self)", func(t *testing.T) {
		err := Errorf("foobar")
		if !errors.Is(err, err) {
			t.Error("unexpected")
		}
	})

	t.Run("errors.Is detects error wrapped to exErr", func(t *testing.T) {
		ee := &net.OpError{}
		err := Errorf("foobar: %w", ee)
		if !errors.Is(err, ee) {
			t.Error("unexpected")
		}
	})

	t.Run("errors.Is detects exErr wrapped to error", func(t *testing.T) {
		ee := Errorf("exErr")
		err := fmt.Errorf("stderr: %w", ee)
		if !errors.Is(err, ee) {
			t.Error("unexpected")
		}
	})

	t.Run("errors.Is detects exErr wrapped to exErr", func(t *testing.T) {
		ee := Errorf("exErr")
		err := Errorf("other error: %w", ee)
		if !errors.Is(err, ee) {
			t.Error("unexpected")
		}
	})
}

func Test_errors_As(t *testing.T) {
	t.Parallel()
	// test that errors.As returns expected results

	t.Run("error implements interface", func(t *testing.T) {
		err := Errorf("foobar")
		var se stacked
		if !errors.As(err, &se) {
			t.Error("unexpected")
		}
	})

	t.Run("error doesn't implement the interface", func(t *testing.T) {
		err := Errorf("foobar")
		var se interface{ Foo() }
		if errors.As(err, &se) {
			t.Error("unexpected")
		}
	})

	t.Run("wrapped error implements the interface", func(t *testing.T) {
		err := fmt.Errorf("err 0: %w", Errorf("foobar"))
		var se stacked
		if !errors.As(err, &se) {
			t.Error("unexpected")
		}
		if len(se.PC()) == 0 {
			t.Error("unexpected")
		}
	})

	t.Run("error is of type", func(t *testing.T) {
		err := Errorf("foobar")
		ee := &exErr{}
		if !errors.As(err, &ee) {
			t.Error("unexpected")
		}
	})
}
