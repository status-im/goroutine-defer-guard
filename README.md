# Goroutine Defer Guard

A static analysis tool that ensures all goroutines have proper panic recovery via a configurable panic handler (defaults to `HandlePanic`).

Example use case: ensure that every goroutine defers to your Sentry reporting wrapper, because Sentry needs panics recovered inside each goroutine.

## What it does

The linter analyzes Go code to find goroutines and ensures that:

1. **Anonymous goroutines** have `defer <target>()` as their first statement
2. **Function call goroutines** call functions that have `defer <target>()` as their first statement  
3. **Method call goroutines** call methods that have `defer <target>()` as their first statement

## Installation

```bash
go install github.com/status-im/goroutine-defer-guard/cmd/goroutine-defer-guard@latest
```

## Usage

```bash
# Run on current directory (defaults to target HandlePanic in the same package)
goroutine-defer-guard ./...

# Skip certain directories
goroutine-defer-guard -skip=./vendor ./...

# Specify a fully-qualified panic handler (import path + function)
goroutine-defer-guard -target=github.com/your/module/common.HandlePanic ./...

# Example: Sentry reporting handler
# Point the linter at your wrapper that reports panics to Sentry:
goroutine-defer-guard -target=github.com/yourorg/observability/panicutil.ReportToSentry ./...
```

## Examples

### ✅ Good - Anonymous goroutine with defer

```go
go func() {
    defer common.HandlePanic()
    // ... rest of function
}()
```

### ❌ Bad - Missing defer statement

```go
go func() {
    // Missing defer common.HandlePanic()
    // ... rest of function
}()
```

### ✅ Good - Function call with proper defer

```go
func worker() {
    defer common.HandlePanic()
    // ... rest of function
}

// Usage
go worker()
```

## How it works

The linter uses:

1. **AST analysis** to find `go` statements (goroutines)
2. **Go type information** to resolve function/method definitions  
3. **Static analysis** to verify the first statement is `defer common.HandlePanic()`

## Configuration

- `-target` (default `HandlePanic`): fully-qualified panic handler in the form `import/path.Func`. 
If you omit the import path the linter accepts a function in the current package or a selector it can resolve to that name.

## Requirements

- Go 1.21+

## License

Mozilla Public License 2.0
