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
			t.Fatal("expected that constructor will return non nil value")
		}

		if e.err != nil {
			t.Errorf("expected e.err to be nil, got %v", e.err)
		} else if e.err != e.Unwrap() {
			t.Errorf("Unwrap returns different value than e.err field:\n%v\nvs\n%v", e.Unwrap(), e.err)
		}

		if len(e.pcs) == 0 {
			t.Error("expected e.pcs not to be empty")
		} else if len(e.pcs) != len(e.PC()) {
			t.Errorf("expected lengths of e.PC() and e.pcs to be equal but got %d != %d", len(e.pcs), len(e.PC()))
		}

		if e.Error() != "<nil>" {
			t.Errorf(`expected error message to be "<nil>" got %q`, e.Error())
		}

		if e.fields != nil {
			t.Errorf("expected e.fields to be nil, got %v", e.fields)
		}
	})

	t.Run("non nil error as input", func(t *testing.T) {
		expErr := fmt.Errorf("some error")
		e := newExErr(expErr)
		if e == nil {
			t.Fatal("expected that constructor will return non nil value")
		}

		if e.err != expErr {
			t.Errorf("expected e.err to be the error passed to the constructor but got %v", e.err)
		}
		if ue := e.Unwrap(); ue != nil {
			t.Errorf("expected Unwrap to return nil got %v", ue)
		}

		if len(e.pcs) == 0 {
			t.Error("expected length of e.pcs not to be zero")
		} else if len(e.pcs) != len(e.PC()) {
			t.Errorf("expected e.PC() length (%d) to be equal to the length of e.pcs (%d)", len(e.PC()), len(e.pcs))
		}

		if e.Error() != expErr.Error() {
			t.Errorf(`expected error message to be %q got %q`, expErr.Error(), e.Error())
		}

		if e.fields != nil {
			t.Errorf("expected e.fields to be nil, got %v", e.fields)
		}
	})
}

func Test_exErr_Fields(t *testing.T) {
	t.Parallel()

	t.Run("no fields", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		if n := len(e.fields); n != 0 {
			t.Errorf("expected no fields, got %d", n)
		}
		if v, ok := e.FieldValue("foo"); ok {
			t.Errorf(`unexpectedly field named "foo" was found: %v`, v)
		}
	})

	expectField := func(t *testing.T, err *exErr, name string, value any) {
		t.Helper()

		if v, ok := err.FieldValue(name); !ok {
			t.Errorf(`field named %q was not found`, name)
		} else if v != value {
			t.Errorf("expected value of the field %q to be %v, got %v", name, value, v)
		}
	}

	t.Run("one field", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		e.AddField("name", 42)
		if n := len(e.fields); n != 1 {
			t.Errorf("expected 1 field, got %d", n)
		}

		expectField(t, e, "name", 42)
	})

	t.Run("two fields, fluid", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		e.AddField("first", 1).AddField("second", 2)
		if n := len(e.fields); n != 2 {
			t.Errorf("expected 2 fields, got %d", n)
		}

		expectField(t, e, "first", 1)
		expectField(t, e, "second", 2)
	})

	t.Run("two fields, two calls", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		e.AddField("first", 1)
		e.AddField("second", 2)
		if n := len(e.fields); n != 2 {
			t.Errorf("expected 2 fields, got %d", n)
		}

		expectField(t, e, "first", 1)
		expectField(t, e, "second", 2)
	})

	t.Run("adding the same field twice", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		e.AddField("foo", 1)
		e.AddField("foo", 2)
		if n := len(e.fields); n != 1 {
			t.Errorf("expected 1 field, got %d", n)
		}

		expectField(t, e, "foo", 2)
	})
}

func Test_exErr_FieldValue(t *testing.T) {
	t.Parallel()

	t.Run("field doesn't exist", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		v, ok := e.FieldValue("foo")
		if ok {
			t.Errorf(`unexpectedly field named "foo" was found: %v`, v)
		}
	})

	t.Run("field exists", func(t *testing.T) {
		e := newExErr(fmt.Errorf("some error"))
		e.AddField("foo", 42)
		if v, ok := e.FieldValue("foo"); !ok {
			t.Error(`unexpectedly field named "foo" was not found`)
		} else if v != 42 {
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
			t.Fatal("expected that constructor will return non nil value")
		}
		if !e.Is(nil) {
			t.Error("Is returns unexpected value")
		}
		if e.Is(e) {
			t.Error("Is returns unexpected value")
		}
	})

	t.Run("wrapping non-nil error", func(t *testing.T) {
		expErr := fmt.Errorf("some error")
		e := newExErr(expErr)
		if !e.Is(expErr) {
			t.Error("Is returns unexpected value")
		}
		if e.Is(e) {
			t.Error("Is returns unexpected value")
		}
	})
}

func Test_errors_Is(t *testing.T) {
	t.Parallel()
	// test that errors.Is returns expected results

	t.Run("nil error", func(t *testing.T) {
		var e *exErr
		if errors.Is(e, nil) {
			t.Error("expected errors.Is to return true for nil error")
		}
	})

	t.Run("errors.Is(self, self)", func(t *testing.T) {
		err := Errorf("foobar")
		if !errors.Is(err, err) {
			t.Error("expected errors.Is(self, self) to return true")
		}
	})

	t.Run("errors.Is detects error wrapped to exErr", func(t *testing.T) {
		ee := &net.OpError{}
		err := Errorf("foobar: %w", ee)
		if !errors.Is(err, ee) {
			t.Error("unexpectedly wrapped stdlib error is not detected")
		}
	})

	t.Run("errors.Is detects exErr wrapped to error", func(t *testing.T) {
		ee := Errorf("exErr")
		err := fmt.Errorf("stderr: %w", ee)
		if !errors.Is(err, ee) {
			t.Error("unexpectedly error is not detected when wrapped to stdlib error")
		}
	})

	t.Run("errors.Is detects exErr wrapped to exErr", func(t *testing.T) {
		ee := Errorf("exErr")
		err := Errorf("other error: %w", ee)
		if !errors.Is(err, ee) {
			t.Error("unexpectedly wrapped error is not detected")
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
			t.Error("unexpectedly error doesn't support stacked interface")
		}
	})

	t.Run("error doesn't implement the interface", func(t *testing.T) {
		err := Errorf("foobar")
		var se interface{ Foo() }
		if errors.As(err, &se) {
			t.Error("error shouldn't support interface with Foo method")
		}
	})

	t.Run("wrapped error implements the interface", func(t *testing.T) {
		err := fmt.Errorf("err 0: %w", Errorf("foobar"))
		var se stacked
		if !errors.As(err, &se) {
			t.Fatal("unexpectedly wrapped error isn't detected by As")
		}
		if len(se.PC()) == 0 {
			t.Error("unexpectedly wrapped error returns zero length PC slice")
		}
	})

	t.Run("error is of type", func(t *testing.T) {
		err := Errorf("foobar")
		ee := &exErr{}
		if !errors.As(err, &ee) {
			t.Error("unexpectedly As doesn't recognize error as exErr")
		}
	})
}
