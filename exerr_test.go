package exerr

import (
	"fmt"
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
		if len(e.pcs) != len(e.Stack()) {
			t.Fatal("unexpectedly e.Stack() returns slice with different length")
		}

		if e.Error() != "<nil>" {
			t.Fatalf(`expected "<nil>" got %q`, e.Error())
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
		if e.err != e.Unwrap() {
			t.Fatal("Unwrap returns unexpected value")
		}

		if len(e.pcs) == 0 {
			t.Fatal("unexpectedly e.pcs is not assigned")
		}
		if len(e.pcs) != len(e.Stack()) {
			t.Fatal("unexpectedly e.Stack() returns slice with different length")
		}

		if e.Error() != expErr.Error() {
			t.Fatalf(`expected %q got %q`, expErr.Error(), e.Error())
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
