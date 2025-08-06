// Package trx provides utilities for handling Go channel
package trx

// Result represents a value that can either be successful (Ok) or contain an error (Err).
// It is a generic type similar to Rust's Result enum, providing safe error handling
// without using exceptions. The zero value is not useful; use Ok() or Err() constructors.
type Result[T any] struct {
	v   T     // The success value
	err error // The error, nil if the result is Ok
}

// Get returns both the value and error from the Result.
// This is useful when you want to handle both cases explicitly.
func (r *Result[T]) Get() (T, error) {
	return r.v, r.err
}

// IsOk returns true if the Result contains a successful value (no error).
func (r *Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the Result contains an error.
func (r *Result[T]) IsErr() bool {
	return r.err != nil
}

// Unwrap returns the success value if the Result is Ok, otherwise panics.
// Use this only when you are certain the Result is Ok, or when you want
// the program to panic on error. For safer alternatives, use UnwrapOr or Value.
func (r *Result[T]) Unwrap() T {
	if r.err != nil {
		panic("attempted to unwrap a Result with an error: " + r.err.Error())
	}

	return r.v
}

// UnwrapOr returns the success value if Ok, otherwise returns the provided default value.
// This is a safe way to extract a value from a Result without panicking.
func (r *Result[T]) UnwrapOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}

	return r.v
}

// UnwrapOrElse returns the success value if Ok, otherwise calls the provided function
// with the error and returns its result. This allows for custom error handling logic.
func (r *Result[T]) UnwrapOrElse(f func(error) T) T {
	if r.err != nil {
		return f(r.err)
	}

	return r.v
}

// Err returns the error from the Result, or nil if the Result is Ok.
func (r *Result[T]) Err() error {
	return r.err
}

// Ok creates a successful Result containing the given value.
// Example: result := Ok(42) creates a Result[int] with value 42.
func Ok[T any](v T) Result[T] {
	return Result[T]{v: v}
}

// Err creates an error Result containing the given error.
// Example: result := Err[int](errors.New("failed")) creates a Result[int] with an error.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// Map applies a function to the success value if Ok, returning a new Result.
func Map[T, U any](r Result[T], mapper func(T) (U, error)) Result[U] {
	if r.err != nil {
		return Err[U](r.err)
	}

	mapped, err := mapper(r.v)
	if err != nil {
		return Err[U](err)
	}

	return Ok(mapped)
}
