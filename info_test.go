package exerr

import (
	"fmt"
	"testing"
)

func Test_FieldValue(t *testing.T) {
	t.Parallel()

	t.Run("stdlib error", func(t *testing.T) {
		err := fmt.Errorf("some error")

		v, ok := FieldValue(err, "foo")
		if ok {
			t.Error("unexpectedly field was found")
		}
		if v != nil {
			t.Errorf("expected nil value, got %v", v)
		}
	})

	t.Run("nil error", func(t *testing.T) {
		v, ok := FieldValue(nil, "foo")
		if ok {
			t.Error("unexpectedly field was found")
		}
		if v != nil {
			t.Errorf("expected nil value, got %v", v)
		}
	})

	t.Run("empty field name, stdlib error", func(t *testing.T) {
		v, ok := FieldValue(fmt.Errorf("some error"), "")
		if ok {
			t.Error("unexpectedly field was found")
		}
		if v != nil {
			t.Errorf("expected nil value, got %v", v)
		}
	})

	t.Run("exerr without any field", func(t *testing.T) {
		err := Errorf("some error")

		v, ok := FieldValue(err, "foo")
		if ok {
			t.Error("unexpectedly field was found")
		}
		if v != nil {
			t.Errorf("expected nil value, got %v", v)
		}
	})

	t.Run("field exists", func(t *testing.T) {
		err := Errorf("some error").AddField("fldName", "foo")
		// field should be found
		v, ok := FieldValue(err, "fldName")
		if !ok {
			t.Error("unexpectedly field was not found")
		} else if v != "foo" {
			t.Errorf("expected value to be %q, got %v", "foo", v)
		}
		// unknown field name shouldn't be found
		v, ok = FieldValue(err, "foo")
		if ok {
			t.Error("unexpectedly field was found")
		}
		if v != nil {
			t.Errorf("expected nil value, got %v", v)
		}
		// empty field name shouldn't be found
		v, ok = FieldValue(err, "")
		if ok {
			t.Error("unexpectedly field was found")
		}
		if v != nil {
			t.Errorf("expected nil value, got %v", v)
		}
	})

	t.Run("two fields", func(t *testing.T) {
		err := Errorf("some error").AddField("fldA", "foo").AddField("fldB", 42)
		// field should be found
		v, ok := FieldValue(err, "fldA")
		if !ok {
			t.Error("unexpectedly field was not found")
		} else if v != "foo" {
			t.Errorf("expected value to be %q, got %v", "foo", v)
		}

		v, ok = FieldValue(err, "fldB")
		if !ok {
			t.Error("unexpectedly field was not found")
		} else if v != 42 {
			t.Errorf("expected value to be %d, got %v", 42, v)
		}
		// unknown field name shouldn't be found
		v, ok = FieldValue(err, "foo")
		if ok {
			t.Error("unexpectedly field was found")
		}
		if v != nil {
			t.Errorf("expected nil value, got %v", v)
		}
		// empty field name shouldn't be found
		v, ok = FieldValue(err, "")
		if ok {
			t.Error("unexpectedly field was found")
		}
		if v != nil {
			t.Errorf("expected nil value, got %v", v)
		}
	})
}

func Test_Fields(t *testing.T) {
	t.Parallel()

	t.Run("stdlib error", func(t *testing.T) {
		err := fmt.Errorf("some error")

		flds := Fields(err)
		if n := len(flds); n != 0 {
			t.Errorf("expected that error has %d fields attached, got %d", 0, n)
		}
	})

	t.Run("nil error", func(t *testing.T) {
		flds := Fields(nil)
		if n := len(flds); n != 0 {
			t.Errorf("expected that error has %d fields attached, got %d", 0, n)
		}
	})

	t.Run("exerr without any field", func(t *testing.T) {
		err := Errorf("some error")

		flds := Fields(err)
		if n := len(flds); n != 0 {
			t.Errorf("expected that error has %d fields attached, got %d", 0, n)
		}
	})

	t.Run("field exists", func(t *testing.T) {
		err := Errorf("some error").AddField("field", 8)

		flds := Fields(err)
		if n := len(flds); n != 1 {
			t.Errorf("expected that error has %d fields attached, got %d", 1, n)
		}

		v, ok := flds["field"]
		if !ok {
			t.Error("unexpectedly field was not found")
		} else if v != 8 {
			t.Errorf("expected value to be 8, got %v", v)
		}
	})

	t.Run("two fields", func(t *testing.T) {
		err := Errorf("some error").AddField("A", 10).AddField("B", "field value")

		flds := Fields(err)
		if n := len(flds); n != 2 {
			t.Errorf("expected that error has %d fields attached, got %d", 2, n)
		}

		if v, ok := flds["A"]; !ok {
			t.Error("unexpectedly field A was not found")
		} else if v != 10 {
			t.Errorf("expected value to be 10, got %v", v)
		}

		if v, ok := flds["B"]; !ok {
			t.Error("unexpectedly field B was not found")
		} else if v != "field value" {
			t.Errorf("expected value to be 'field value', got %v", v)
		}
	})
}
