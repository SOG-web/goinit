package di

import (
	"testing"
)

type A struct{ Val int }

type B struct {
	A *A
}

var aCounter int

func newA() *A { aCounter++; return &A{Val: aCounter} }

func newB(a *A) (*B, error) { return &B{A: a}, nil }

// Impl provides a small implementation used by tests. Methods must be
// declared at package level (not inside functions), otherwise the Go
// compiler reports a syntax error.
type Impl struct{}

func (Impl) Do() string { return "ok" }

// Types for circular dependency test. Declared at package scope so they can
// reference each other.
type C struct{ D *D }
type D struct{ C *C }

func TestSingletonResolution(t *testing.T) {
	c := New()
	if err := Register[*A](c, newA, Singleton); err != nil {
		t.Fatal(err)
	}

	a1, _ := Resolve[*A](c)
	a2, _ := Resolve[*A](c)
	if a1 != a2 {
		t.Error("expected same singleton instance")
	}
}

func TestTransientResolution(t *testing.T) {
	c := New()
	if err := Register[*A](c, newA, Transient); err != nil {
		t.Fatal(err)
	}
	a1, _ := Resolve[*A](c)
	a2, _ := Resolve[*A](c)
	if a1 == a2 {
		t.Error("expected different transient instances")
	}
}

func TestConstructorWithError(t *testing.T) {
	c := New()
	if err := Register[*A](c, newA, Singleton); err != nil {
		t.Fatal(err)
	}
	if err := Register[*B](c, newB, Singleton); err != nil {
		t.Fatal(err)
	}
	b, err := Resolve[*B](c)
	if err != nil {
		t.Fatal(err)
	}
	if b.A == nil {
		t.Error("expected resolved A inside B")
	}
}

func TestSliceResolutionSingletons(t *testing.T) {
	c := New()
	if err := Register[*A](c, newA, Singleton, "one"); err != nil {
		t.Fatal(err)
	}
	if err := Register[*A](c, newA, Singleton, "two"); err != nil {
		t.Fatal(err)
	}

	list, err := Resolve[[]*A](c)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 elements, got %d", len(list))
	}
	// both elements should be distinct singletons per tag
	if list[0] == list[1] {
		t.Error("expected different singletons for different tags")
	}
	// subsequent slice resolution should reuse the same singletons
	list2, _ := Resolve[[]*A](c)
	if list[0] != list2[0] || list[1] != list2[1] {
		t.Error("expected slice elements to be stable singletons")
	}
}

func TestCircularDependency(t *testing.T) {
	c := New()
	Register[*C](c, func(d *D) *C { return &C{D: d} }, Singleton)
	Register[*D](c, func(cc *C) *D { return &D{C: cc} }, Singleton)

	_, err := Resolve[*C](c)
	if err == nil {
		t.Error("expected circular dependency error")
	}
}

func TestAmbiguousRegistrationError(t *testing.T) {
	c := New()
	Register[*A](c, newA, Singleton, "x")
	Register[*A](c, newA, Singleton, "y")
	_, err := Resolve[*A](c)
	if err == nil {
		t.Error("expected ambiguity error when multiple tags exist but resolving single")
	}
}

func TestProvideInterface(t *testing.T) {
	type I interface{ Do() string }

	c := New()
	if err := Provide[I](c, Impl{}); err != nil { t.Fatal(err) }

	v, err := Resolve[I](c)
	if err != nil { t.Fatal(err) }
	if v.Do() != "ok" { t.Error("expected provided implementation to be resolved") }
}

func TestPrimitiveParamValidation(t *testing.T) {
	type X struct{ S string }
	c := New()
	err := Register[*X](c, func(s string) *X { return &X{S: s} }, Singleton)
	if err == nil {
		t.Error("expected error due to primitive string dependency with no registration")
	}
}

func TestClearAndClose(t *testing.T) {
	c := New()
	Register[*A](c, newA, Singleton)
	c.Clear()
	if IsRegistered[*A](c) {
		t.Error("expected registrations cleared")
	}

	c = New()
	Register[*A](c, newA, Singleton)
	c.Close()
	_, err := Resolve[*A](c)
	if err == nil {
		t.Error("expected error resolving from closed container")
	}
}

// helper to verify Provide rejects non-interfaces
func TestProvideNonInterface(t *testing.T) {
	c := New()
	err := Provide[*A](c, &A{})
	if err == nil {
		t.Error("expected error when using Provide with non-interface type")
	}
}
