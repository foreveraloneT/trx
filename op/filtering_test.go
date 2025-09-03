package op_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/foreveralonet/trx"
	"github.com/foreveralonet/trx/op"
)

var _ = Describe("Filtering Operations", func() {

	Describe("Filter", func() {
		Context("when filtering with a predicate function", func() {
			It("should keep only elements that match the predicate", func() {
				source := op.Range(0, 10)
				out := op.Filter(source, func(value int, index int) (bool, error) {
					return value%2 == 0, nil // Keep even numbers
				})

				expectedValues := []int{0, 2, 4, 6, 8}
				results := make([]int, 0)

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})

			It("should pass the correct index to the predicate", func() {
				source := op.Range(10, 5) // [10, 11, 12, 13, 14]
				indices := make([]int, 0)

				out := op.Filter(source, func(value int, index int) (bool, error) {
					indices = append(indices, index)
					return true, nil // Keep all elements
				})

				// Consume all results
				for range out {
				}

				Expect(indices).To(Equal([]int{0, 1, 2, 3, 4}))
			})

			It("should handle empty sources", func() {
				source := op.Range(0, 0) // Empty range
				out := op.Filter(source, func(value int, index int) (bool, error) {
					return true, nil
				})

				count := 0
				for range out {
					count++
				}

				Expect(count).To(Equal(0))
			})

			It("should work with different data types", func() {
				source := op.FormSlice([]string{"apple", "banana", "cherry", "date"})
				out := op.Filter(source, func(value string, index int) (bool, error) {
					return len(value) > 5, nil // Keep strings longer than 5 chars
				})

				expectedValues := []string{"banana", "cherry"}
				results := make([]string, 0)

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})
		})

		Context("when predicate function returns an error", func() {
			It("should propagate the error", func() {
				source := op.Range(0, 5)
				testError := errors.New("predicate error")

				out := op.Filter(source, func(value int, index int) (bool, error) {
					if value == 3 {
						return false, testError
					}
					return true, nil
				})

				errorFound := false
				for result := range out {
					if result.IsErr() {
						errorFound = true
						_, err := result.Get()
						Expect(err).To(Equal(testError))
					}
				}

				Expect(errorFound).To(BeTrue())
			})
		})

		Context("when using with pool options", func() {
			It("should work with concurrent processing", func() {
				source := op.Range(0, 20)
				out := op.Filter(source, func(value int, index int) (bool, error) {
					return value%3 == 0, nil // Keep multiples of 3
				}, op.WithPoolSize(3))

				results := make([]int, 0)
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				// Results might not be in order due to concurrent processing
				expectedSet := []int{0, 3, 6, 9, 12, 15, 18}
				Expect(results).To(ConsistOf(expectedSet))
			})

			It("should maintain order with serialized processing", func() {
				source := op.Range(0, 10)
				out := op.Filter(source, func(value int, index int) (bool, error) {
					return value%2 == 1, nil // Keep odd numbers
				}, op.WithPoolSize(3), op.WithSerialize())

				expectedValues := []int{1, 3, 5, 7, 9}
				results := make([]int, 0)

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})
		})
	})

	Describe("Take", func() {
		Context("when taking a specific number of elements", func() {
			It("should emit exactly n elements from the source", func() {
				source := op.Range(0, 100)
				n := 5
				out := op.Take(source, n)

				expectedValues := []int{0, 1, 2, 3, 4}
				results := make([]int, 0)

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})

			It("should handle taking zero elements", func() {
				source := op.Range(0, 10)
				out := op.Take(source, 0)

				count := 0
				for range out {
					count++
				}

				Expect(count).To(Equal(0))
			})

			It("should handle taking more elements than available", func() {
				source := op.Range(5, 3) // Only 3 elements: [5, 6, 7]
				out := op.Take(source, 10)

				expectedValues := []int{5, 6, 7}
				results := make([]int, 0)

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})

			It("should work with different data types", func() {
				source := op.FormSlice([]string{"a", "b", "c", "d", "e"})
				out := op.Take(source, 3)

				expectedValues := []string{"a", "b", "c"}
				results := make([]string, 0)

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})
		})

		Context("when source contains errors", func() {
			It("should propagate errors from the source", func() {
				// Create a source where the first element is an error
				source := make(chan trx.Result[int], 1)
				source <- trx.Err[int](errors.New("source error"))
				close(source)

				out := op.Take[int](source, 5)

				errorFound := false
				resultCount := 0

				for result := range out {
					resultCount++
					if result.IsErr() {
						errorFound = true
						_, err := result.Get()
						Expect(err.Error()).To(Equal("source error"))
					}
				}

				// Debug: check if we got any results at all
				Expect(resultCount).To(BeNumerically(">", 0), "Should have received at least one result")
				Expect(errorFound).To(BeTrue())
			})
		})

		Context("with buffering options", func() {
			It("should work with buffered channels", func() {
				source := op.Range(0, 10, op.WithBufferSize(5))
				out := op.Take[int](source, 3, op.WithBufferSize(2))

				expectedValues := []int{0, 1, 2}
				results := make([]int, 0)

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})
		})
	})

	Describe("Combined filtering operations", func() {
		Context("when chaining Filter and Take", func() {
			It("should apply operations in sequence", func() {
				source := op.Range(0, 20)

				// First filter for even numbers, then take first 3
				filtered := op.Filter(source, func(value int, index int) (bool, error) {
					return value%2 == 0, nil
				})
				out := op.Take[int](filtered, 3)

				expectedValues := []int{0, 2, 4}
				results := make([]int, 0)

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})
		})
	})
})
