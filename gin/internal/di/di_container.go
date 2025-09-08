// Package di provides an improved dependency injection container
package di

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

type Scope int

const (
	Singleton Scope = iota
	Transient
)

type Container struct {
	mu            sync.RWMutex
	registrations map[typeKey]*registrationList
	// Singletons and their per-type creation mutexes (string key = typeKey.String())
	singletons      sync.Map
	singletonsMutex sync.Map
	// Circular dependency tracking by type+tag key
	resolving sync.Map
	closed    int32
}

type typeKey struct {
	typ reflect.Type
	tag string
}

func (tk typeKey) String() string {
	if tk.tag == "" {
		return tk.typ.String()
	}
	return tk.typ.String() + ":" + tk.tag
}

type registrationList struct {
	mu    sync.RWMutex
	items []*registration
}

type registration struct {
	constructor reflect.Value
	scope       Scope
	tag         string
	// optional direct factory that avoids reflect at resolve time
	directFactory func() (any, error)
}

func New() *Container {
	return &Container{registrations: make(map[typeKey]*registrationList)}
}

func (c *Container) isClosed() bool { return atomic.LoadInt32(&c.closed) != 0 }

// Close marks the container as closed and clears internal maps.
// It acquires c.mu so that Close and concurrent Register cannot interleave.
func (c *Container) Close() {
	atomic.StoreInt32(&c.closed, 1)
	c.mu.Lock()
	defer c.mu.Unlock()
	// Clear maps while holding the lock to avoid races with Register/Resolve
	c.registrations = make(map[typeKey]*registrationList)
	c.singletons.Range(func(key, _ any) bool { c.singletons.Delete(key); return true })
	c.singletonsMutex.Range(func(key, _ any) bool { c.singletonsMutex.Delete(key); return true })
	c.resolving.Range(func(key, _ any) bool { c.resolving.Delete(key); return true })
}

// Register registers a constructor for type T with the given scope and optional tag.
// Patches applied:
//   - Allow constructors that return (T) or (T, error).
//   - Basic validation on constructor input parameters: forbid primitives unless explicitly registered.
//   - Double-check closed state while holding c.mu to avoid Close/Register races.
func Register[T any](c *Container, constructor any, scope Scope, tag ...string) error {
	if c.isClosed() {
		return errors.New("container is closed")
	}

	t := reflect.TypeOf((*T)(nil)).Elem()
	ctorVal := reflect.ValueOf(constructor)
	ctorType := ctorVal.Type()

	if ctorType.Kind() != reflect.Func {
		return errors.New("constructor must be a function")
	}

	// Validate outputs: either 1 value (T) or 2 values (T, error)
	nOut := ctorType.NumOut()
	if nOut != 1 && nOut != 2 {
		return errors.New("constructor must return T or (T, error)")
	}
	if !ctorType.Out(0).AssignableTo(t) {
		return fmt.Errorf("constructor return type %v is not assignable to %v", ctorType.Out(0), t)
	}
	if nOut == 2 {
		if !isErrorType(ctorType.Out(1)) {
			return errors.New("if constructor returns two values, second must be error")
		}
	}

	// Validate inputs: disallow primitives (string, bool, numbers, etc.)
	// unless there's already a registration for that exact type.
	for i := 0; i < ctorType.NumIn(); i++ {
		inT := ctorType.In(i)
		if isPrimitiveKind(inT.Kind()) {
			c.mu.RLock()
			_, ok := c.registrations[typeKey{typ: inT, tag: ""}]
			c.mu.RUnlock()
			if !ok {
				return fmt.Errorf("constructor parameter %d is primitive %v and not registered", i, inT)
			}
		}
	}

	regTag := ""
	if len(tag) > 0 {
		regTag = tag[0]
	}

	key := typeKey{typ: t, tag: regTag}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed() {
		return errors.New("container is closed")
	}

	regList := c.registrations[key]
	if regList == nil {
		regList = &registrationList{}
		c.registrations[key] = regList
	}

	reg := &registration{constructor: ctorVal, scope: scope, tag: regTag}

	regList.mu.Lock()
	regList.items = append(regList.items, reg)
	regList.mu.Unlock()

	return nil
}

// RegisterFactory registers a typed factory function for T that will be called
// directly (no reflection) when resolving. This is useful for hot-path types.
// The factory must be a `func() T` (no error). If you need error-returning
// factories, use Register with a constructor that returns (T, error).
func RegisterFactory[T any](c *Container, factory func() T, scope Scope, tag ...string) error {
	if c.isClosed() {
		return errors.New("container is closed")
	}

	t := reflect.TypeOf((*T)(nil)).Elem()

	regTag := ""
	if len(tag) > 0 {
		regTag = tag[0]
	}

	key := typeKey{typ: t, tag: regTag}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed() {
		return errors.New("container is closed")
	}

	regList := c.registrations[key]
	if regList == nil {
		regList = &registrationList{}
		c.registrations[key] = regList
	}

	// wrap factory into a func() (any, error)
	wrapper := func() (any, error) {
		v := factory()
		return any(v), nil
	}

	reg := &registration{scope: scope, tag: regTag, directFactory: wrapper}

	regList.mu.Lock()
	regList.items = append(regList.items, reg)
	regList.mu.Unlock()

	return nil
}

// Provide binds a concrete implementation value to an interface type T as a singleton.
// This helper intentionally supports interface T only.
func Provide[T any](c *Container, impl T, tag ...string) error {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Interface {
		return errors.New("Provide is for interface types only")
	}
	return Register[T](c, func() T { return impl }, Singleton, tag...)
}

func Resolve[T any](c *Container) (T, error) {
	var zero T
	if c.isClosed() {
		return zero, errors.New("container is closed")
	}

	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() == reflect.Slice {
		res, err := c.resolveSlice(t.Elem())
		if err != nil {
			return zero, err
		}
		return res.(T), nil
	}

	res, err := c.resolveType(t, "")
	if err != nil {
		return zero, err
	}
	return res.(T), nil
}

func ResolveWithTag[T any](c *Container, tag string) (T, error) {
	var zero T
	if c.isClosed() {
		return zero, errors.New("container is closed")
	}
	t := reflect.TypeOf((*T)(nil)).Elem()
	res, err := c.resolveType(t, tag)
	if err != nil {
		return zero, err
	}
	return res.(T), nil
}

func MustResolve[T any](c *Container) T {
	v, err := Resolve[T](c)
	if err != nil {
		panic(err)
	}
	return v
}

func MustResolveWithTag[T any](c *Container, tag string) T {
	v, err := ResolveWithTag[T](c, tag)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Container) resolveType(t reflect.Type, tag string) (any, error) {
	key := typeKey{typ: t, tag: tag}
	keyStr := key.String()
	// resolving key: keyStr

	// Circular dependency detection
	if _, isResolving := c.resolving.LoadOrStore(keyStr, true); isResolving {
		return nil, fmt.Errorf("circular dependency detected for type %s", keyStr)
	}
	defer c.resolving.Delete(keyStr)

	// Read registrations
	c.mu.RLock()
	regList := c.registrations[key]
	c.mu.RUnlock()

	if regList == nil {
		return nil, fmt.Errorf("no registration found for type %s", keyStr)
	}

	regList.mu.RLock()
	cnt := len(regList.items)
	if cnt == 0 {
		regList.mu.RUnlock()
		return nil, fmt.Errorf("no registration found for type %s", keyStr)
	}
	// If more than one registration for this type/tag and we're resolving a single T,
	// this is ambiguous. Ask user to use tags or resolve []T.
	if cnt > 1 && t.Kind() != reflect.Slice {
		regList.mu.RUnlock()
		return nil, fmt.Errorf("multiple registrations for %s: specify a tag or resolve a slice ([]%s)", keyStr, t.String())
	}
	reg := regList.items[0]
	regList.mu.RUnlock()

	if reg.scope == Singleton {
		if instance, ok := c.singletons.Load(keyStr); ok {
			return instance, nil
		}
		// per-type creation mutex
		mutexI, _ := c.singletonsMutex.LoadOrStore(keyStr, &sync.Mutex{})
		m := mutexI.(*sync.Mutex)
		m.Lock()
		defer m.Unlock()

		// Double-check after locking
		if instance, ok := c.singletons.Load(keyStr); ok {
			return instance, nil
		}

		result, err := c.createInstance(reg)
		if err != nil {
			return nil, err
		}
		c.singletons.Store(keyStr, result)
		// Avoid leaking per-type mutexes after successful creation
		c.singletonsMutex.Delete(keyStr)
		return result, nil
	}

	return c.createInstance(reg)
}

// createInstance calls the constructor (supporting T or (T,error)) and returns the created value.
func (c *Container) createInstance(reg *registration) (any, error) {
	// If a directFactory is provided, use it (no reflection)
	if reg.directFactory != nil {
		return reg.directFactory()
	}

	ctorType := reg.constructor.Type()
	numIn := ctorType.NumIn()
	args := make([]reflect.Value, numIn)

	for i := 0; i < numIn; i++ {
		paramType := ctorType.In(i)
		arg, err := c.resolveType(paramType, "")
		if err != nil {
			return nil, fmt.Errorf("failed to resolve parameter %d (%v): %w", i, paramType, err)
		}
		args[i] = reflect.ValueOf(arg)
	}

	results := reg.constructor.Call(args)
	switch len(results) {
	case 1:
		return results[0].Interface(), nil
	case 2:
		// (T, error)
		if !results[1].IsNil() {
			return nil, results[1].Interface().(error)
		}
		return results[0].Interface(), nil
	default:
		return nil, errors.New("constructor must return T or (T, error)")
	}
}

// resolveSlice resolves all registrations whose type matches elemType across all tags.
// It uses resolveType(elemType, reg.tag) to preserve singleton semantics.
func (c *Container) resolveSlice(elemType reflect.Type) (any, error) {
	// Gather matching registrations by (type, tag)
	c.mu.RLock()
	var entries []typeKey
	for k := range c.registrations {
		if k.typ == elemType {
			entries = append(entries, k)
		}
	}
	c.mu.RUnlock()

	// Resolve each entry with its tag via resolveType to respect Singleton/Transient
	resolved := make([]reflect.Value, 0, len(entries))
	for _, k := range entries {
		inst, err := c.resolveType(elemType, k.tag)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve slice element for %s: %w", k.String(), err)
		}
		v := reflect.ValueOf(inst)
		if !v.Type().AssignableTo(elemType) {
			return nil, fmt.Errorf("resolved element type %v is not assignable to %v", v.Type(), elemType)
		}
		resolved = append(resolved, v)
	}

	slice := reflect.MakeSlice(reflect.SliceOf(elemType), len(resolved), len(resolved))
	for i, v := range resolved {
		slice.Index(i).Set(v)
	}
	return slice.Interface(), nil
}

// IsRegistered reports whether T (and optional tag) has at least one registration.
func IsRegistered[T any](c *Container, tag ...string) bool {
	t := reflect.TypeOf((*T)(nil)).Elem()
	regTag := ""
	if len(tag) > 0 {
		regTag = tag[0]
	}
	key := typeKey{typ: t, tag: regTag}
	c.mu.RLock()
	regList := c.registrations[key]
	c.mu.RUnlock()
	if regList == nil {
		return false
	}
	regList.mu.RLock()
	defer regList.mu.RUnlock()
	return len(regList.items) > 0
}

// Clear removes all registrations, singletons and resolving state.
func (c *Container) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.registrations = make(map[typeKey]*registrationList)
	c.singletons.Range(func(key, _ any) bool { c.singletons.Delete(key); return true })
	c.singletonsMutex.Range(func(key, _ any) bool { c.singletonsMutex.Delete(key); return true })
	c.resolving.Range(func(key, _ any) bool { c.resolving.Delete(key); return true })
}

// GetRegisteredTypes returns all registered types (with tags) for debugging.
func (c *Container) GetRegisteredTypes() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var types []string
	for key := range c.registrations {
		types = append(types, key.String())
	}
	return types
}

// --- helpers ---

func isErrorType(t reflect.Type) bool {
	var errType = reflect.TypeOf((*error)(nil)).Elem()
	return t.Implements(errType)
}

func isPrimitiveKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	default:
		return false
	}
}
