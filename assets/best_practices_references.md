# Go Best Practices Reference

This document tracks key best practices from the [Google Go Style Guide - Best Practices](https://google.github.io/styleguide/go/best-practices) for reference during code reviews and refactoring.

## Table of Contents
1. [Naming](#naming)
2. [Error Handling](#error-handling)
3. [Function Arguments](#function-arguments)
4. [Documentation](#documentation)
5. [Variable Declarations](#variable-declarations)
6. [String Concatenation](#string-concatenation)
7. [Global State](#global-state)
8. [Tests](#tests)

---

## Naming

### Function and Method Names

**Avoid repetition:**
- Don't repeat package name in function names
- Don't repeat receiver type in method names
- Don't repeat parameter names in function names
- Don't repeat return value names/types

**Naming conventions:**
- Functions that return something: noun-like names (avoid `Get` prefix)
- Functions that do something: verb-like names
- Identical functions differing only by types: include type name at end

**Example:**
```go
// Bad:
func (c *Config) GetJobName(key string) (value string, ok bool)

// Good:
func (c *Config) JobName(key string) (value string, ok bool)
```

### Shadowing

- Use `:=` for initialization with non-zero values
- Use `var` for zero-value declarations
- Be careful with shadowing in new scopes
- Avoid shadowing standard package names

---

## Error Handling

### Error Structure

- Use structured errors (sentinel values, custom types) for programmatic inspection
- Don't distinguish errors based on string matching
- Use `errors.Is` and `errors.As` for error checking

### Adding Information to Errors

- Use `%v` for simple annotation or when creating new errors
- Use `%w` when you want callers to programmatically inspect the error chain
- Place `%w` at the end of error strings: `fmt.Errorf("context: %w", err)`
- Avoid redundant information that the underlying error already provides
- Don't add annotations that don't add new information

### Error Wrapping

```go
// Good: Adding context while preserving error chain
return fmt.Errorf("couldn't find remote file: %w", err)

// Good: Simple annotation without preserving chain
return fmt.Errorf("launch codes unavailable: %v", err)
```

### Logging Errors

- Don't log errors that you return (let caller decide)
- Use appropriate log levels (avoid excessive ERROR level)
- Be careful with PII in logs

---

## Function Arguments

### Option Structures

Use option structs when:
- All callers need to specify one or more options
- Large number of callers need to provide many options
- Options are shared between multiple functions

**Note:** Contexts are never included in option structs.

### Variadic Options

Use variadic options when:
- Most callers won't need to specify any options
- Most options are used infrequently
- There are a large number of options
- Options require arguments
- Options could fail or be set incorrectly

---

## Documentation

### Conventions

- Not every parameter must be enumerated
- Document error-prone or non-obvious fields/parameters
- Document why parameters are interesting, not just what they are
- Context cancellation is implied - don't restate it
- Document concurrency safety when it's unclear or non-standard
- Document cleanup requirements explicitly

### Error Documentation

- Document significant error sentinel values or error types
- Note whether errors are pointer receivers or not
- Document overall error conventions in package documentation

---

## Variable Declarations

### Initialization

- Prefer `:=` over `var` when initializing with non-zero value
- Use `var` for zero-value declarations
- Use composite literals when you know initial elements
- Use size hints for slices/maps when size is known

### Zero Values

```go
// Good:
var coords Point
if err := json.Unmarshal(data, &coords); err != nil {
    // ...
}

// Good for pointer types:
msg := new(pb.Bar)
```

---

## String Concatenation

- Prefer `+` for simple cases (few strings)
- Prefer `fmt.Sprintf` when formatting
- Prefer `strings.Builder` for constructing strings piecemeal
- Use backticks for constant multi-line string literals

---

## Global State

**Avoid global state when:**
- Multiple functions interact via global state despite being independent
- Independent test cases interact through global state
- Users are tempted to swap/replace global state for testing
- Users must consider special ordering requirements

**Safe global state:**
- Logically constant
- Stateless (caller can't distinguish cache hits from misses)
- Doesn't bleed into external systems
- No expectation of predictable behavior

**Guidelines:**
- Allow clients to create isolated instances
- Package-level APIs should be thin proxies
- Package-level API only for binary targets, not libraries
- Document and enforce invariants

---

## Tests

### Test Helpers vs Assertion Helpers

- **Test helpers**: Do setup/cleanup, call `t.Helper()`
- **Assertion helpers**: Not idiomatic in Go - prefer inline logic in `Test` function

### Error Handling in Test Helpers

- Test helpers should call `t.Fatal` when setup fails
- Include description of what happened in failure message
- Use `t.Cleanup` for cleanup functions

### Don't Call `t.Fatal` from Separate Goroutines

- Only call `t.Fatal` from the goroutine running the `Test` function
- Use `t.Error` and return from other goroutines

### Test Organization

- Keep setup code scoped to specific tests
- Use `sync.Once` for expensive setup that doesn't require teardown
- Use custom `TestMain` only when all tests require common setup with teardown

---

## Additional Best Practices

### Package Size

- Packages should be focused on a cohesive idea
- Related types that interact should be in the same package
- Split conceptually distinct functionality into separate packages

### Import Ordering

- Follow standard import grouping conventions
- Protocol buffer imports use `pb` and `grpc` suffixes

### Program Checks and Panics

- Prefer error return values over panics
- Use `log.Fatal` for consistency checks on invariants
- Don't recover panics to avoid crashes
- Panics are acceptable for API misuse (like standard library)

---

## References

- [Google Go Style Guide - Best Practices](https://google.github.io/styleguide/go/best-practices)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Blog on Errors](https://go.dev/blog/error-handling-and-go)
