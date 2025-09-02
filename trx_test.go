package trx_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/foreveraloneT/trx"
)

var _ = Describe("Result", func() {

	Describe("Ok constructor", func() {
		Context("when creating a successful result", func() {
			It("should create an Ok result with the given value", func() {
				result := trx.Ok(42)

				Expect(result.IsOk()).To(BeTrue())
				Expect(result.IsErr()).To(BeFalse())

				value, err := result.Get()
				Expect(value).To(Equal(42))
				Expect(err).To(BeNil())
			})

			It("should work with different types", func() {
				stringResult := trx.Ok("hello")
				Expect(stringResult.IsOk()).To(BeTrue())

				value, err := stringResult.Get()
				Expect(value).To(Equal("hello"))
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Err constructor", func() {
		Context("when creating an error result", func() {
			It("should create an Err result with the given error", func() {
				testErr := errors.New("test error")
				result := trx.Err[int](testErr)

				Expect(result.IsOk()).To(BeFalse())
				Expect(result.IsErr()).To(BeTrue())

				_, err := result.Get()
				Expect(err).To(Equal(testErr))
				Expect(result.Err()).To(Equal(testErr))
			})
		})
	})

	Describe("Unwrap method", func() {
		Context("when the result is Ok", func() {
			It("should return the value", func() {
				result := trx.Ok(42)
				value := result.Unwrap()
				Expect(value).To(Equal(42))
			})
		})

		Context("when the result is Err", func() {
			It("should panic", func() {
				result := trx.Err[int](errors.New("test error"))
				Expect(func() { result.Unwrap() }).To(Panic())
			})
		})
	})

	Describe("UnwrapOr method", func() {
		Context("when the result is Ok", func() {
			It("should return the original value", func() {
				result := trx.Ok(42)
				value := result.UnwrapOr(10)
				Expect(value).To(Equal(42))
			})
		})

		Context("when the result is Err", func() {
			It("should return the default value", func() {
				result := trx.Err[int](errors.New("test error"))
				value := result.UnwrapOr(10)
				Expect(value).To(Equal(10))
			})
		})
	})

	Describe("UnwrapOrElse method", func() {
		Context("when the result is Ok", func() {
			It("should return the original value without calling the function", func() {
				result := trx.Ok(42)
				called := false
				value := result.UnwrapOrElse(func(err error) int {
					called = true
					return 10
				})

				Expect(value).To(Equal(42))
				Expect(called).To(BeFalse())
			})
		})

		Context("when the result is Err", func() {
			It("should call the function with the error and return its result", func() {
				testErr := errors.New("test error")
				result := trx.Err[int](testErr)

				value := result.UnwrapOrElse(func(err error) int {
					Expect(err).To(Equal(testErr))
					return 99
				})

				Expect(value).To(Equal(99))
			})
		})
	})

	Describe("Map function", func() {
		Context("when mapping an Ok result", func() {
			It("should apply the mapper function and return Ok result", func() {
				result := trx.Ok(42)
				mapped := trx.Map(result, func(v int) (string, error) {
					return "value: 42", nil
				})

				Expect(mapped.IsOk()).To(BeTrue())
				value, err := mapped.Get()
				Expect(value).To(Equal("value: 42"))
				Expect(err).To(BeNil())
			})

			It("should return Err if mapper function returns error", func() {
				result := trx.Ok(42)
				mapperErr := errors.New("mapper error")
				mapped := trx.Map(result, func(v int) (string, error) {
					return "", mapperErr
				})

				Expect(mapped.IsErr()).To(BeTrue())
				_, err := mapped.Get()
				Expect(err).To(Equal(mapperErr))
			})
		})

		Context("when mapping an Err result", func() {
			It("should propagate the error without calling mapper", func() {
				originalErr := errors.New("original error")
				result := trx.Err[int](originalErr)
				called := false

				mapped := trx.Map(result, func(v int) (string, error) {
					called = true
					return "should not be called", nil
				})

				Expect(mapped.IsErr()).To(BeTrue())
				_, err := mapped.Get()
				Expect(err).To(Equal(originalErr))
				Expect(called).To(BeFalse())
			})
		})
	})

	Describe("Edge cases", func() {
		Context("with nil values", func() {
			It("should handle nil pointers correctly", func() {
				var nilPtr *string
				result := trx.Ok(nilPtr)

				Expect(result.IsOk()).To(BeTrue())
				value, err := result.Get()
				Expect(value).To(BeNil())
				Expect(err).To(BeNil())
			})
		})

		Context("with zero values", func() {
			It("should handle zero values correctly", func() {
				result := trx.Ok(0)

				Expect(result.IsOk()).To(BeTrue())
				value, err := result.Get()
				Expect(value).To(Equal(0))
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Type safety", func() {
		Context("when using different types", func() {
			It("should maintain type safety across operations", func() {
				intResult := trx.Ok(42)
				stringMapped := trx.Map(intResult, func(v int) (string, error) {
					return "number: 42", nil
				})

				Expect(stringMapped.IsOk()).To(BeTrue())
				value, err := stringMapped.Get()
				Expect(value).To(BeAssignableToTypeOf(""))
				Expect(value).To(Equal("number: 42"))
				Expect(err).To(BeNil())
			})
		})
	})
})
