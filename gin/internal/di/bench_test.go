package di

import (
	"reflect"
	"testing"
)

// BenchType and constructor used for microbenchmarks.
type BenchType struct{
	N int
}

func newBench() *BenchType { return &BenchType{N: 1} }

// Factory wrapper type to simulate registering a direct factory in the container.
type factoryFunc func() *BenchType

func BenchmarkDirectCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = newBench()
	}
}

func BenchmarkReflectCall_NoArgs(b *testing.B) {
	ctor := reflect.ValueOf(func() *BenchType { return newBench() })
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		results := ctor.Call(nil)
		_ = results[0].Interface().(*BenchType)
	}
}

func BenchmarkFactoryClosureCall(b *testing.B) {
	f := factoryFunc(newBench)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = f()
	}
}

// A tiny benchmark that simulates the container invoking a reflect-based constructor
// with one resolved argument. This mirrors typical DI where constructors accept deps.
func BenchmarkReflectCall_OneArg(b *testing.B) {
	ctor := reflect.ValueOf(func(x *BenchType) *BenchType { return &BenchType{N: x.N + 1} })
	// prepare arg (simulate resolved dependency)
	arg := reflect.ValueOf(newBench())
	args := []reflect.Value{arg}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		results := ctor.Call(args)
		_ = results[0].Interface().(*BenchType)
	}
}

// --- container-level benchmarks ---
// These measure end-to-end resolution using the Container: reflect vs direct factory.
func BenchmarkContainer_Resolve_Reflect(b *testing.B) {
	c := New()
	// register reflect-based constructor
	Register[*BenchType](c, func() *BenchType { return &BenchType{N: 1} }, Transient)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Resolve[*BenchType](c)
	}
}

func BenchmarkContainer_Resolve_Factory(b *testing.B) {
	c := New()
	// register direct factory
	RegisterFactory[*BenchType](c, func() *BenchType { return &BenchType{N: 1} }, Transient)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Resolve[*BenchType](c)
	}
}
