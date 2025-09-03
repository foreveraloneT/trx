# trx

[![Go Reference](https://pkg.go.dev/badge/github.com/foreveraloneT/trx.svg)](https://pkg.go.dev/github.com/foreveraloneT/trx)
[![Go Report Card](https://goreportcard.com/badge/github.com/foreveraloneT/trx)](https://goreportcard.com/report/github.com/foreveraloneT/trx)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go library that provides utilities for handling channels with Rust-like Result types and reactive programming patterns. It combines type-safe error handling with powerful channel operations for building robust concurrent applications.

## Features

- **🦀 Rust-like Result Type**: Type-safe error handling without exceptions
- **⚡ Reactive Channel Operations**: Creation, transformation, and filtering operators
- **🔄 Concurrent Processing**: Built-in worker pools with configurable concurrency
- **🎯 Type Safety**: Full generic type support for compile-time safety

## Installation

```bash
go get github.com/foreveraloneT/trx
```

## Quick Start

### Basic Result Usage

```go
package main

import (
    "fmt"
    "github.com/foreveraloneT/trx"
)

func divide(a, b int) trx.Result[int] {
    if b == 0 {
        return trx.Err[int](fmt.Errorf("division by zero"))
    }
    return trx.Ok(a / b)
}

func main() {
    result := divide(10, 2)
    
    if result.IsOk() {
        value := result.Unwrap()
        fmt.Printf("Result: %d\n", value) // Result: 5
    }
    
    // Safe unwrapping with default value
    safeResult := divide(10, 0)
    value := safeResult.UnwrapOr(0)
    fmt.Printf("Safe result: %d\n", value) // Safe result: 0
}
```

### Channel Operations

```go
package main

import (
    "fmt"
    "time"
    "github.com/foreveraloneT/trx/op"
)

func main() {
    // Create a range of numbers
    source := op.Range(0, 10)
    
    // Transform each number
    doubled := op.Map(source, func(v int, index int) (int, error) {
        return v * 2, nil
    })
    
    // Filter even numbers
    filtered := op.Filter(doubled, func(v int, index int) (bool, error) {
        return v%4 == 0, nil
    })
    
    // Take only first 3 results
    result := op.Take(filtered, 3)
    
    // Process results
    for r := range result {
        if r.IsOk() {
            fmt.Println(r.Unwrap()) // Output: 0, 4, 8
        }
    }
}
```

## Core Concepts

### Result Type

The `Result[T]` type represents a value that can either be successful (`Ok`) or contain an error (`Err`), similar to Rust's Result enum.

```go
// Create successful result
result := trx.Ok(42)

// Create error result
result := trx.Err[int](errors.New("something went wrong"))

// Safe value extraction
value := result.UnwrapOr(0)          // Returns 0 if error
value := result.UnwrapOrElse(func(err error) int {
    log.Printf("Error: %v", err)
    return -1
})

// Check result state
if result.IsOk() {
    // Handle success
} else if result.IsErr() {
    // Handle error
}
```

### Channel Operators

The `op` package provides reactive programming operators for working with channels:

#### Creation Operators

```go
// Timer - emit after delay
timer := op.Timer(2 * time.Second)

// Interval - emit at regular intervals
interval := op.Interval(500 * time.Millisecond)

// Range - emit sequence of numbers
numbers := op.Range(1, 5) // [1, 2, 3, 4, 5]

// From slice
data := op.FormSlice([]string{"a", "b", "c"})

// From existing channel
existing := make(chan int, 3)
wrapped := op.FormChannel(existing)
```

#### Transformation Operators

```go
// Map - transform each value
source := op.Range(1, 5)
squared := op.Map(source, func(v int, index int) (int, error) {
    return v * v, nil
})

// Buffer by count
buffered := op.BufferWithCount(source, 3) // Groups into slices of 3

// Buffer by time
timeBuf := op.BufferWithTime(source, time.Second, 10)

// Buffer by time or count (whichever comes first)
flexBuf := op.BufferWithTimeOrCount(source, time.Second, 5)
```

#### Filtering Operators

```go
// Filter - keep only matching values
source := op.Range(1, 10)
evens := op.Filter(source, func(v int, index int) (bool, error) {
    return v%2 == 0, nil
})

// Take - limit number of items
limited := op.Take(source, 5)
```

## Advanced Features

### Concurrent Processing

Use worker pools for CPU-intensive operations:

```go
source := op.Range(1, 100)

// Process with 4 worker goroutines
result := op.Map(source, func(v int, index int) (int, error) {
    // Expensive computation
    time.Sleep(100 * time.Millisecond)
    return v * v, nil
}, op.WithPoolSize(4))

// Maintain order with serialized output
ordered := op.Map(source, func(v int, index int) (int, error) {
    return v * 2, nil
}, op.WithPoolSize(4), op.WithSerialize())
```

### Configuration Options

All operators support functional options:

```go
source := op.Range(1, 10,
    op.WithBufferSize(5),           // Channel buffer size
    op.WithContext(ctx),            // Cancellation context
)

transformed := op.Map(source, mapper,
    op.WithPoolSize(3),             // Worker pool size
    op.WithSerialize(),             // Maintain order
    op.WithBufferSize(10),          // Output buffer size
)
```

### Error Handling

Errors propagate through the operator chain:

```go
source := op.Range(1, 5)

result := op.Map(source, func(v int, index int) (string, error) {
    if v == 3 {
        return "", fmt.Errorf("error at value 3")
    }
    return fmt.Sprintf("value-%d", v), nil
})

for r := range result {
    if r.IsErr() {
        fmt.Printf("Error: %v\n", r.Err())
        continue
    }
    fmt.Printf("Success: %s\n", r.Unwrap())
}
```

## Real-World Examples

### Data Processing Pipeline

```go
func processLogFiles(filenames []string) <-chan trx.Result[ProcessedLog] {
    // Create source from filenames
    source := op.FormSlice(filenames)
    
    // Read files concurrently
    contents := op.Map(source, func(filename string, index int) ([]byte, error) {
        return os.ReadFile(filename)
    }, op.WithPoolSize(10))
    
    // Parse logs
    parsed := op.Map(contents, func(data []byte, index int) (RawLog, error) {
        return parseLog(data)
    })
    
    // Filter valid logs
    valid := op.Filter(parsed, func(log RawLog, index int) (bool, error) {
        return log.IsValid(), nil
    })
    
    // Transform to processed format
    return op.Map(valid, func(log RawLog, index int) (ProcessedLog, error) {
        return processLog(log)
    })
}
```

### Real-time Data Processing

```go
func processMetrics() {
    // Collect metrics every second
    metrics := op.Interval(time.Second)
    
    // Transform to metric data
    data := op.Map(metrics, func(tick int, index int) (Metric, error) {
        return collectSystemMetrics()
    })
    
    // Buffer for batch processing
    batched := op.BufferWithTimeOrCount(data, 10*time.Second, 100)
    
    // Process batches
    for batch := range batched {
        if batch.IsOk() {
            metrics := batch.Unwrap()
            processBatch(metrics)
        }
    }
}
```

## Performance Considerations

- **Buffer Sizes**: Use appropriate buffer sizes to prevent blocking
- **Worker Pools**: Match pool size to available CPU cores for CPU-bound tasks
- **Memory Usage**: Be mindful of buffering large amounts of data
- **Context Cancellation**: Always use contexts for long-running operations

## Testing

The library includes comprehensive tests using Ginkgo and Gomega:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE.txt](LICENSE.txt) file for details.

## Inspiration

This library is inspired by:
- [Rust's Result type](https://doc.rust-lang.org/std/result/)
- [RxJS](https://rxjs.dev/) reactive programming patterns
- [RxGo](https://github.com/ReactiveX/RxGo) reactive programming patterns for the Go Language
- [conc](https://github.com/sourcegraph/conc) structured concurrency for Go Language
