package test

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func AssertEqual(t *testing.T, got, want any) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v; want %v", got, want)
	}
}

func AssertNotEqual(t *testing.T, got, notwant any) {
	t.Helper()

	if reflect.DeepEqual(got, notwant) {
		t.Fatalf("got (but don't want) %v", got)
	}
}

func AssertAtLeast(t *testing.T, got, want int) {
	t.Helper()

	if got < want {
		t.Fatalf("got %v; want at least %v", got, want)
	}
}

func AssertStringContains(t *testing.T, got, want string) {
	t.Helper()

	if !strings.Contains(got, want) {
		t.Fatalf("got %q; want to contain: %q", got, want)
	}
}

func AssertSliceContains[T comparable](t *testing.T, got []T, want T) {
	t.Helper()

	if !slices.Contains(got, want) {
		t.Fatalf("got %v; want to contain: %v", got, want)
	}
}

func AssertSliceDoesNotContain[T comparable](t *testing.T, got []T, want T) {
	t.Helper()

	if slices.Contains(got, want) {
		t.Fatalf("got %v; should not contain: %v", got, want)
	}
}

func AssertNilError(t *testing.T, got error) {
	t.Helper()

	if got != nil {
		t.Fatalf("got: %v; want: nil", got)
	}
}

func AssertErrorIs(t *testing.T, got error, want error) {
	t.Helper()

	if got == nil {
		t.Fatalf("got: nil; want: %q", want)
	}

	if !errors.Is(got, want) {
		t.Fatalf("got %q; want: %q", got, want)
	}
}

func AssertErrorAs(t *testing.T, got error, want any) {
	t.Helper()

	if got == nil {
		t.Fatalf("got: nil; want: %T", want)
	}

	if !errors.As(got, want) {
		t.Fatalf("got %q; want: %T", got, want)
	}
}

func AssertErrorContains(t *testing.T, got error, want string) {
	t.Helper()

	if got == nil {
		t.Fatalf("got: nil; want: error to contain: %q", want)
	}

	if !strings.Contains(got.Error(), want) {
		t.Fatalf("got %q; want to contain: %q", got.Error(), want)
	}
}
