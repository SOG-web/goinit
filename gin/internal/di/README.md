# Dependency Injection Container

This package provides a robust, production-ready typed dependency injection (DI) container for Go applications.

## Overview

The DI container supports:

- Generic type-safe registrations
- Automatic dependency resolution
- Singleton and transient scopes
- Circular dependency detection
- Tagged/qualified injections
- Thread-safe operations
- Interface-to-implementation mapping

## Core Concepts

### Container

The `Container` is the main DI container that holds all registrations and handles dependency resolution.

### Scopes

- **Singleton**: Single shared instance throughout the application lifecycle
- **Transient**: New instance created each time it's resolved

## API Reference

### Register[T]

**Purpose**: Registers a constructor function that creates instances of type `T`

**When to use**:

- When you want the container to create instances for you
- When the type has dependencies that need automatic injection
- For concrete types that require initialization logic
- Supports singleton/transient scopes

**Signature**:

```go
func Register[T any](c *Container, constructor interface{}, scope Scope, tag ...string) error
```

**Examples**:

```go
// Register a service with dependencies
Register[*user.UserService](c, func(repo repo.UserRepository, email email.EmailServiceInterface) *user.UserService {
    return user.NewUserService(repo, email)
}, Singleton)

// Register a repository with DB dependency
Register[repo.UserRepository](c, func(db *gorm.DB) repo.UserRepository {
    return dataRepo.NewGormUserRepository(db)
}, Singleton)

// Register with tag for multiple implementations
Register[Logger](c, func() Logger { return &ConsoleLogger{} }, Singleton, "console")
```

### Provide[T]

**Purpose**: Provides an existing instance for interface type `T`

**When to use**:

- When you have an already created instance (like external services)
- For interface-to-implementation mapping
- When the instance is created outside the container (e.g., from config)
- Typically used for singleton instances that don't need recreation

**Signature**:

```go
func Provide[T any](c *Container, impl T, tag ...string) error
```

**Examples**:

```go
// Provide existing email service instance for the interface
Provide[email.EmailServiceInterface](c, emailSvc)

// Provide existing DB instance
Provide[*gorm.DB](c, gdb)

// Provide with tag
Provide[Logger](c, &FileLogger{}, "file")
```

### Resolve[T]

**Purpose**: Resolves and returns an instance of type `T`

**Signature**:

```go
func Resolve[T any](c *Container) (T, error)
```

**Examples**:

```go
// Resolve a service
userSvc, err := Resolve[*user.UserService](c)

// Resolve with tag
logger, err := ResolveWithTag[Logger](c, "console")

// Resolve slice of implementations
loggers, err := Resolve[[]Logger](c)
```

### MustResolve[T]

**Purpose**: Resolves T and panics on error

**Signature**:

```go
func MustResolve[T any](c *Container) T
```

## Key Differences: Register vs Provide

| Aspect           | Register                           | Provide                          |
| ---------------- | ---------------------------------- | -------------------------------- |
| **Input**        | Constructor function               | Existing instance                |
| **Dependencies** | Auto-injected                      | None (instance already has them) |
| **Scope**        | Configurable (Singleton/Transient) | Always Singleton                 |
| **Use Case**     | Creating new instances             | Registering existing instances   |
| **Type**         | Concrete types with factories      | Interfaces with implementations  |

## Usage Patterns

### Basic Setup

```go
c := New()

// Register dependencies
Register[*gorm.DB](c, func() *gorm.DB { return gdb }, Singleton)
Register[repo.UserRepository](c, func(db *gorm.DB) repo.UserRepository {
    return dataRepo.NewGormUserRepository(db)
}, Singleton)

// Resolve
repo := MustResolve[repo.UserRepository](c)
```

### Interface Mapping

```go
// Register implementation
Register[*ConsoleLogger](c, func() *ConsoleLogger { return &ConsoleLogger{} }, Singleton)

// Map to interface
Provide[Logger](c, MustResolve[*ConsoleLogger](c))

// Resolve interface
logger := MustResolve[Logger](c)
```

### Tagged Dependencies

```go
// Register multiple implementations with tags
Register[PaymentProcessor](c, func() PaymentProcessor { return &StripeProcessor{} }, Singleton, "stripe")
Register[PaymentProcessor](c, func() PaymentProcessor { return &PayPalProcessor{} }, Singleton, "paypal")

// Resolve specific implementation
stripe := MustResolveWithTag[PaymentProcessor](c, "stripe")
```

### Circular Dependency Detection

The container automatically detects circular dependencies:

```go
// This would panic with "circular dependency detected"
Register[A](c, func(b B) A { return A{b} }, Singleton)
Register[B](c, func(a A) B { return B{a} }, Singleton)
```

## Error Handling

All resolution operations return errors for:

- Missing registrations
- Circular dependencies
- Type mismatches

Use `MustResolve` for cases where you're certain the dependency exists.

## Thread Safety

The container is thread-safe and can be used concurrently from multiple goroutines.

## Integration Example

```go
// In main.go
c := di.New()

// Register infrastructure
di.Register[*gorm.DB](c, func() *gorm.DB { return gdb }, di.Singleton)
di.Provide[email.EmailServiceInterface](c, emailSvc)

// Register repositories
di.Register[repo.UserRepository](c, func(db *gorm.DB) repo.UserRepository {
    return dataRepo.NewGormUserRepository(db)
}, di.Singleton)

// Register services
di.Register[*user.UserService](c, func(repo repo.UserRepository, emailSvc email.EmailServiceInterface) *user.UserService {
    return user.NewUserService(repo, emailSvc)
}, di.Singleton)

// In handlers
userSvc := di.MustResolve[*user.UserService](c)
```
