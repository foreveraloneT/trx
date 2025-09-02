package op_test

import (
	"errors"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/foreveraloneT/trx"
	"github.com/foreveraloneT/trx/op"
)

var _ = Describe("Transformation Operations", func() {

	Describe("Map", func() {
		Context("when transforming values with a mapper function", func() {
			It("should transform each value according to the mapper", func() {
				source := op.Range(1, 4) // [1, 2, 3, 4]
				out := op.Map(source, func(value int, index int) (string, error) {
					return fmt.Sprintf("item-%d", value), nil
				})

				expectedValues := []string{"item-1", "item-2", "item-3", "item-4"}
				results := make([]string, 0)

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})

			It("should pass the correct index to the mapper", func() {
				source := op.Range(10, 3) // [10, 11, 12]
				indices := make([]int, 0)

				out := op.Map(source, func(value int, index int) (int, error) {
					indices = append(indices, index)
					return value * 2, nil
				})

				// Consume all results
				for range out {
				}

				Expect(indices).To(Equal([]int{0, 1, 2}))
			})

			It("should handle type transformations", func() {
				source := op.FormSlice([]int{1, 2, 3})
				out := op.Map(source, func(value int, index int) (float64, error) {
					return float64(value) * 1.5, nil
				})

				expectedValues := []float64{1.5, 3.0, 4.5}
				results := make([]float64, 0)

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})

			It("should handle empty sources", func() {
				source := op.Range(0, 0) // Empty range
				out := op.Map(source, func(value int, index int) (string, error) {
					return "should not be called", nil
				})

				count := 0
				for range out {
					count++
				}

				Expect(count).To(Equal(0))
			})
		})

		Context("when mapper function returns an error", func() {
			It("should propagate the error", func() {
				source := op.Range(0, 5)
				testError := errors.New("mapper error")

				out := op.Map(source, func(value int, index int) (string, error) {
					if value == 3 {
						return "", testError
					}
					return fmt.Sprintf("value-%d", value), nil
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
				source := op.Range(1, 5)
				out := op.Map(source, func(value int, index int) (int, error) {
					return value * value, nil // Square the values
				}, op.WithPoolSize(3))

				results := make([]int, 0)
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				// Results might not be in order due to concurrent processing
				expectedSet := []int{1, 4, 9, 16, 25}
				Expect(results).To(ConsistOf(expectedSet))
			})

			It("should maintain order with serialized processing", func() {
				source := op.Range(1, 4)
				out := op.Map(source, func(value int, index int) (int, error) {
					return value * 10, nil
				}, op.WithPoolSize(3), op.WithSerialize())

				expectedValues := []int{10, 20, 30, 40}
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

	Describe("BufferWithCount", func() {
		Context("when buffering values by count", func() {
			It("should group values into batches of specified size", func() {
				source := op.Range(0, 7) // [0, 1, 2, 3, 4, 5, 6]
				out := op.BufferWithCount(source, 3)

				expectedBatches := [][]int{
					{0, 1, 2},
					{3, 4, 5},
					{6}, // Last batch with remaining items
				}

				results := make([][]int, 0)
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					batch, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, batch)
				}

				Expect(results).To(Equal(expectedBatches))
			})

			It("should handle exact multiples of buffer size", func() {
				source := op.Range(0, 6) // [0, 1, 2, 3, 4, 5]
				out := op.BufferWithCount(source, 2)

				expectedBatches := [][]int{
					{0, 1},
					{2, 3},
					{4, 5},
				}

				results := make([][]int, 0)
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					batch, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, batch)
				}

				Expect(results).To(Equal(expectedBatches))
			})

			It("should handle buffer size of 1", func() {
				source := op.Range(0, 3)
				out := op.BufferWithCount(source, 1)

				expectedBatches := [][]int{
					{0}, {1}, {2},
				}

				results := make([][]int, 0)
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					batch, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, batch)
				}

				Expect(results).To(Equal(expectedBatches))
			})

			It("should handle empty sources", func() {
				source := op.Range(0, 0)
				out := op.BufferWithCount(source, 3)

				count := 0
				for range out {
					count++
				}

				Expect(count).To(Equal(0))
			})

			It("should work with different data types", func() {
				source := op.FormSlice([]string{"a", "b", "c", "d", "e"})
				out := op.BufferWithCount(source, 2)

				expectedBatches := [][]string{
					{"a", "b"},
					{"c", "d"},
					{"e"},
				}

				results := make([][]string, 0)
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					batch, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, batch)
				}

				Expect(results).To(Equal(expectedBatches))
			})
		})
	})

	Describe("BufferWithTime", func() {
		Context("when buffering values by time", func() {
			It("should emit batches after timeout", func() {
				source := make(chan trx.Result[int], 5)

				// Send values with delays
				go func() {
					defer close(source)
					source <- trx.Ok(1)
					source <- trx.Ok(2)
					source <- trx.Ok(3)
					time.Sleep(60 * time.Millisecond) // Force timeout
					source <- trx.Ok(4)
					source <- trx.Ok(5)
				}()

				out := op.BufferWithTime(op.FormChannel(source), 50*time.Millisecond, 0)

				batches := make([][]trx.Result[int], 0)
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					batch, err := result.Get()
					Expect(err).To(BeNil())
					batches = append(batches, batch)
				}

				// Should have at least 2 batches due to timeout
				Expect(len(batches)).To(BeNumerically(">=", 2))

				// Check first batch content
				if len(batches) > 0 {
					firstBatch := make([]int, 0)
					for _, res := range batches[0] {
						if res.IsOk() {
							val, _ := res.Get()
							firstBatch = append(firstBatch, val)
						}
					}
					Expect(firstBatch).To(Equal([]int{1, 2, 3}))
				}
			})

			It("should respect max size limit", func() {
				source := make(chan trx.Result[int], 10)

				// Send many values quickly
				go func() {
					defer close(source)
					for i := 1; i <= 8; i++ {
						source <- trx.Ok(i)
					}
				}()

				out := op.BufferWithTime(op.FormChannel(source), 100*time.Millisecond, 3)

				batches := make([][]trx.Result[int], 0)
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					batch, err := result.Get()
					Expect(err).To(BeNil())
					batches = append(batches, batch)
				}

				// Most batches should be limited by max size
				for i, batch := range batches[:len(batches)-1] { // Exclude last batch which might be partial
					Expect(len(batch)).To(BeNumerically("<=", 3), fmt.Sprintf("Batch %d exceeded max size", i))
				}
			})
		})
	})

	Describe("BufferWithTimeOrCount", func() {
		Context("when buffering values by time or count", func() {
			It("should emit when count is reached", func() {
				source := op.Range(0, 8)
				out := op.BufferWithTimeOrCount(source, 100*time.Millisecond, 3)

				batches := make([][]int, 0)
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					batch, err := result.Get()
					Expect(err).To(BeNil())
					batches = append(batches, batch)
				}

				// Most batches should have exactly 3 items (count limit)
				for i, batch := range batches[:len(batches)-1] {
					Expect(len(batch)).To(Equal(3), fmt.Sprintf("Batch %d should have 3 items", i))
				}

				// Last batch should have remaining items
				if len(batches) > 0 {
					lastBatch := batches[len(batches)-1]
					Expect(len(lastBatch)).To(BeNumerically("<=", 3))
				}
			})

			It("should emit when time expires", func() {
				source := make(chan trx.Result[int], 10)

				// Send some values then wait
				go func() {
					defer close(source)
					source <- trx.Ok(1)
					source <- trx.Ok(2)
					time.Sleep(60 * time.Millisecond) // Force timeout
					source <- trx.Ok(3)
					source <- trx.Ok(4)
				}()

				out := op.BufferWithTimeOrCount(op.FormChannel(source), 50*time.Millisecond, 5)

				batches := make([][]trx.Result[int], 0)
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					batch, err := result.Get()
					Expect(err).To(BeNil())
					batches = append(batches, batch)
				}

				// Should have multiple batches due to timeout
				Expect(len(batches)).To(BeNumerically(">=", 2))
			})
		})
	})

	Describe("Combined transformation operations", func() {
		Context("when chaining multiple transformations", func() {
			It("should apply transformations in sequence", func() {
				source := op.Range(1, 10)

				// Map to double values, then buffer by count
				doubled := op.Map(source, func(value int, index int) (int, error) {
					return value * 2, nil
				})
				out := op.BufferWithCount(doubled, 3)

				expectedBatches := [][]int{
					{2, 4, 6},
					{8, 10, 12},
					{14, 16, 18},
					{20},
				}

				results := make([][]int, 0)
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					batch, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, batch)
				}

				Expect(results).To(Equal(expectedBatches))
			})
		})
	})
})
