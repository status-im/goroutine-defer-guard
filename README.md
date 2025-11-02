# Goroutine Defer Guard

A static analysis tool that ensures all goroutines have proper panic recovery via `defer common.LogOnPanic()` calls.

This linter was migrated from [status-go](https://github.com/status-im/status-go) to be a standalone tool for checking that all goroutines include the necessary defer statement to handle panics properly.

## What it does

The linter analyzes Go code to find goroutines and ensures that:

1. **Anonymous goroutines** have `defer common.LogOnPanic()` as their first statement
2. **Function call goroutines** call functions that have `defer common.LogOnPanic()` as their first statement  
3. **Method call goroutines** call methods that have `defer common.LogOnPanic()` as their first statement

## Installation

```bash
go install github.com/status-im/goroutine-defer-guard/cmd/goroutine-defer-guard@latest
```

## Usage

```bash
# Run on current directory
goroutine-defer-guard ./...

# Skip certain directories
goroutine-defer-guard -skip=./vendor ./...
```

## Examples

### ✅ Good - Anonymous goroutine with defer

```go
go func() {
    defer common.LogOnPanic()
    // ... rest of function
}()
```

### ❌ Bad - Missing defer statement

```go
go func() {
    // Missing defer common.LogOnPanic()
    // ... rest of function
}()
```

### ✅ Good - Function call with proper defer

```go
func worker() {
    defer common.LogOnPanic()
    // ... rest of function
}

// Usage
go worker()
```

## How it works

The linter uses:

1. **AST analysis** to find `go` statements (goroutines)
2. **Go type information** to resolve function/method definitions  
3. **Static analysis** to verify the first statement is `defer common.LogOnPanic()`

## Requirements

- Go 1.21+

## License

Mozilla Public License 2.0
