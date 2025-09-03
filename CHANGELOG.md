# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.2] - 2025-09-03

### Changed
- **Documentation**: Updated README.md to reflect the new module name consistently
  - Updated installation instructions to use `github.com/foreveralonet/trx`
  - Updated import statements in code examples to use the new module name
  - Updated all URLs and references to use the correct module path

### Fixed
- **Import Consistency**: Fixed remaining import statements in documentation examples that were still using the old module name

## [0.1.1] - 2025-09-03

### Added
- N/A (No added feature)

### Changed
- **Module Name**: Updated module name from `github.com/foreveraloneT/trx` to `github.com/foreveralonet/trx`
  - **Breaking Change**: Users will need to update their import statements

## [0.1.0] - 2025-09-03

### Added

#### Core Features
- **Result Type**: Implemented generic `Result[T]` type for safe error handling
  - `Ok(value)` and `Err(error)` constructors
  - `Get()`, `IsOk()`, `IsErr()` methods for result inspection
  - `Unwrap()`, `UnwrapOr()`, `UnwrapOrElse()` methods for value extraction
  - `Map()` function for result transformation
  - Full type safety with Go generics

#### Channel Operators Package (`op`)
- **Creation Operators**:
  - `Timer(duration)` - Emit single value after delay
  - `Interval(duration)` - Emit incrementing values at regular intervals
  - `Range(start, count)` - Generate sequence of consecutive integers
  - `FormSlice(slice)` - Convert slice to channel stream
  - `FormChannel(channel)` - Wrap existing channel with Result types

- **Transformation Operators**:
  - `Map(source, mapper)` - Transform each value with concurrent processing
  - `BufferWithCount(source, count)` - Group values into fixed-size batches
  - `BufferWithTime(source, duration, maxSize)` - Group values by time with optional size limit
  - `BufferWithTimeOrCount(source, duration, count)` - Group by time OR count (whichever comes first)

- **Filtering Operators**:
  - `Filter(source, predicate)` - Keep only values matching predicate
  - `Take(source, n)` - Limit number of emitted values

#### Configuration System
- **Functional Options Pattern**:
  - `WithBufferSize(size)` - Configure channel buffer sizes
  - `WithPoolSize(size)` - Set number of worker goroutines
  - `WithSerialize()` - Enable ordered output in concurrent operations
  - `WithContext(ctx)` - Add cancellation support

#### Concurrency Features
- **Worker Pool Integration**:
  - Built on [Sourcegraph's conc](https://github.com/sourcegraph/conc) library
  - Configurable concurrent processing for Map and Filter operations
  - Optional result serialization to maintain order
  - Automatic resource cleanup and goroutine management

#### Documentation and Examples
- **Comprehensive Examples**:
  - Creation operators examples (`_example/creation/`)
  - Transformation operators examples (`_example/transformation/`)
  - Filtering operators examples (`_example/filtering/`)
  - Real-world usage patterns in documentation

- **Testing Suite**:
  - Full test coverage using Ginkgo and Gomega
  - Unit tests for all operators and Result type methods
  - Integration tests for operator chaining
  - Performance and concurrency testing

#### Developer Experience
- **MIT License** - Open source licensing
- **Go 1.23+ Support** - Modern Go features and generics
- **Zero Dependencies** - Core library has no external dependencies
- **Type Safety** - Full compile-time type checking with generics

### Technical Details

#### Dependencies
- **Core**:
  - `github.com/sourcegraph/conc v0.3.0`
- **Testing**: 
  - `github.com/onsi/ginkgo/v2 v2.25.2`
  - `github.com/onsi/gomega v1.38.2`
- **Concurrency**: 
  - `github.com/sourcegraph/conc v0.3.0`

#### Supported Go Versions
- Go 1.23.0 and later (requires generics support)

#### Key Features
- **Reactive Programming**: Channel-based operators inspired by [RxGo](https://github.com/ReactiveX/RxGo)
- **Functional Error Handling**: Monadic Result type for safe error propagation
- **High Performance**: Optimized concurrent processing with worker pools
- **Type Safety**: Full generic type support for compile-time safety
- **Flexible Configuration**: Functional options pattern for customization

### Breaking Changes
- N/A (Initial release)

### Migration Guide
- N/A (Initial release)

---

## How to Read This Changelog

- **Added**: New features and capabilities
- **Changed**: Changes to existing functionality
- **Deprecated**: Features that will be removed in future versions
- **Removed**: Features that have been removed
- **Fixed**: Bug fixes
- **Security**: Security-related changes

## Version Format

This project follows [Semantic Versioning](https://semver.org/):
- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality in a backwards compatible manner
- **PATCH**: Backwards compatible bug fixes

## Links

- [Repository](https://github.com/foreveraloneT/trx)
- [Go Package Documentation](https://pkg.go.dev/github.com/foreveraloneT/trx)
- [Issues](https://github.com/foreveraloneT/trx/issues)
- [Releases](https://github.com/foreveraloneT/trx/releases)
