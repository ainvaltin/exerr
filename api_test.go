package exerr

import (
	"errors"
	"fmt"
	"testing"
)

func Test_Errorf(t *testing.T) {
	t.Parallel()

	t.Run("no fields attached", func(t *testing.T) {
		err := Errorf("some error")
		if err == nil {
			t.Fatal("unexpectedli Errorf returned nil error")
		}
		if n := len(Fields(err)); n != 0 {
			t.Errorf("unexpectedly error has %d fields attached", n)
		}
		if n := errChainLen(err); n != 2 {
			t.Errorf("expected chain lenght 2, got %d", n)
		}
	})

	t.Run("one fields attached", func(t *testing.T) {
		err := Errorf("some error").AddField("field1", 1)
		flds := Fields(err)
		if n := len(flds); n != 1 {
			t.Errorf("expecte that error has %d fields attached, got %d", 1, n)
		} else {
			containsField(t, flds, "field1", 1)
		}

		expectFieldValue(t, err, "field1", 1)

		if n := errChainLen(err); n != 2 {
			t.Errorf("expected chain lenght 2, got %d", n)
		}
	})

	t.Run("two fields attached", func(t *testing.T) {
		err := Errorf("some error").AddField("field1", 11).AddField("field2", 12)

		flds := Fields(err)
		if n := len(flds); n != 2 {
			t.Errorf("expecte that error has %d fields attached, got %d", 2, n)
		} else {
			containsField(t, flds, "field1", 11)
			containsField(t, flds, "field2", 12)
		}

		expectFieldValue(t, err, "field1", 11)
		expectFieldValue(t, err, "field2", 12)

		if n := errChainLen(err); n != 2 {
			t.Errorf("expected chain lenght 2, got %d", n)
		}
	})

	t.Run("Errorf wraps Errorf", func(t *testing.T) {
		err1 := Errorf("some error").AddField("field1", 11)
		err := Errorf("wrap error: %w", err1).AddField("field2", 12)

		flds := Fields(err)
		if n := len(flds); n != 2 {
			t.Errorf("expecte that error has %d fields attached, got %d", 2, n)
		} else {
			containsField(t, flds, "field1", 11)
			containsField(t, flds, "field2", 12)
		}

		expectFieldValue(t, err, "field1", 11)
		expectFieldValue(t, err, "field2", 12)

		if n := errChainLen(err); n != 4 {
			t.Errorf("expected chain lenght 4, got %d", n)
		}
	})
}

func Test_AddField(t *testing.T) {
	t.Parallel()

	t.Run("nil error as imput", func(t *testing.T) {
		err := AddField(nil, "fieldA", 1)

		flds := Fields(err)
		if n := len(flds); n != 1 {
			t.Errorf("expecte that error has %d fields attached, got %d", 1, n)
		} else {
			containsField(t, flds, "fieldA", 1)
		}

		expectFieldValue(t, err, "fieldA", 1)

		if n := errChainLen(err); n != 1 {
			t.Errorf("expected chain lenght 1, got %d", n)
		}
	})

	t.Run("stdlib error as input", func(t *testing.T) {
		origErr := fmt.Errorf("some error")
		err := AddField(origErr, "fieldA", 1)

		flds := Fields(err)
		if n := len(flds); n != 1 {
			t.Errorf("expecte that error has %d fields attached, got %d", 0, n)
		} else {
			containsField(t, flds, "fieldA", 1)
		}

		expectFieldValue(t, err, "fieldA", 1)

		if n := errChainLen(err); n != 2 {
			t.Errorf("expected chain lenght 2, got %d", n)
		} else {
			if !errors.Is(err, origErr) {
				t.Error("errors.Is doesn't detect the wrapped error")
			}
		}
	})

	t.Run("exerr without fields", func(t *testing.T) {
		origErr := Errorf("some error")
		err := AddField(origErr, "fieldA", 1)

		flds := Fields(err)
		if n := len(flds); n != 1 {
			t.Errorf("expecte that error has %d fields attached, got %d", 0, n)
		} else {
			containsField(t, flds, "fieldA", 1)
		}

		expectFieldValue(t, err, "fieldA", 1)

		if n := errChainLen(err); n != 2 {
			t.Errorf("expected chain lenght 2, got %d", n)
		} else {
			if !errors.Is(err, origErr) {
				t.Error("errors.Is doesn't detect the wrapped error")
			}
		}
	})

	t.Run("exerr with fields", func(t *testing.T) {
		origErr := Errorf("some error").AddField("fieldB", 2)
		err := AddField(origErr, "fieldA", 1)

		flds := Fields(err)
		if n := len(flds); n != 2 {
			t.Errorf("expecte that error has %d fields attached, got %d", 2, n)
		} else {
			containsField(t, flds, "fieldA", 1)
			containsField(t, flds, "fieldB", 2)
		}

		expectFieldValue(t, err, "fieldA", 1)
		expectFieldValue(t, err, "fieldB", 2)

		if n := errChainLen(err); n != 2 {
			t.Errorf("expected chain lenght 2, got %d", n)
		} else {
			if !errors.Is(err, origErr) {
				t.Error("errors.Is doesn't detect the wrapped error")
			}
		}
	})

	t.Run("AddFeld called twice", func(t *testing.T) {
		origErr := fmt.Errorf("some error")
		exErr := AddField(origErr, "fieldA", 1)
		err := AddField(exErr, "fieldB", 2)

		flds := Fields(err)
		if n := len(flds); n != 2 {
			t.Errorf("expecte that error has %d fields attached, got %d", 2, n)
		} else {
			containsField(t, flds, "fieldA", 1)
			containsField(t, flds, "fieldB", 2)
		}

		expectFieldValue(t, err, "fieldA", 1)
		expectFieldValue(t, err, "fieldB", 2)

		if n := errChainLen(err); n != 2 {
			t.Errorf("expected chain lenght 2, got %d", n)
		} else {
			if !errors.Is(err, origErr) {
				t.Error("errors.Is doesn't detect the wrapped error")
			}
			if !errors.Is(err, exErr) {
				t.Error("errors.Is doesn't detect the wrapped error")
			}
		}
	})

	t.Run("AddFeld chained", func(t *testing.T) {
		origErr := Errorf("some error").AddField("fieldB", 2)
		err := AddField(origErr, "fieldA", 1).AddField("fieldC", 3)

		flds := Fields(err)
		if n := len(flds); n != 3 {
			t.Errorf("expecte that error has %d fields attached, got %d", 3, n)
		} else {
			containsField(t, flds, "fieldA", 1)
			containsField(t, flds, "fieldB", 2)
			containsField(t, flds, "fieldC", 3)
		}

		expectFieldValue(t, err, "fieldA", 1)
		expectFieldValue(t, err, "fieldB", 2)
		expectFieldValue(t, err, "fieldC", 3)

		if n := errChainLen(err); n != 2 {
			t.Errorf("expected chain lenght 2, got %d", n)
		} else {
			if !errors.Is(err, origErr) {
				t.Error("errors.Is doesn't detect the wrapped error")
			}
		}
	})
}

func expectFieldValue(t *testing.T, err error, name string, value any) {
	t.Helper()

	v, ok := FieldValue(err, name)
	if !ok {
		t.Errorf(`field named %q was not found`, name)
	} else if v != value {
		t.Errorf("expected value to be %v, got %v", value, v)
	}
}

func containsField(t *testing.T, fields map[string]any, name string, value any) {
	t.Helper()

	v, ok := fields[name]
	if !ok {
		t.Errorf(`field named %q was not found`, name)
	} else if v != value {
		t.Errorf("expected value to be %v, got %v", value, v)
	}
}

func errChainLen(err error) int {
	len := 0
	for err != nil {
		len++
		err = errors.Unwrap(err)
	}
	return len
}
