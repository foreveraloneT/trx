package op_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/foreveralonet/trx/op"
)

var _ = Describe("Creation Operations", func() {

	Describe("Timer", func() {
		Context("when creating a timer with specific duration", func() {
			It("should emit a single value after the specified duration", func() {
				start := time.Now()
				duration := 50 * time.Millisecond

				out := op.Timer(duration)

				result := <-out
				elapsed := time.Since(start)

				Expect(result.IsOk()).To(BeTrue())
				Expect(result.IsErr()).To(BeFalse())

				value, err := result.Get()
				Expect(value).To(Equal(0))
				Expect(err).To(BeNil())

				Expect(elapsed).To(BeNumerically(">=", duration))
			})

			It("should close the channel after emission", func() {
				duration := 10 * time.Millisecond
				out := op.Timer(duration)

				// Wait for the value
				<-out

				// Channel should be closed
				_, ok := <-out
				Expect(ok).To(BeFalse())
			})
		})

		Context("with different durations", func() {
			It("should respect different timer durations", func() {
				shortDuration := 10 * time.Millisecond
				longDuration := 30 * time.Millisecond

				start := time.Now()
				shortTimer := op.Timer(shortDuration)
				longTimer := op.Timer(longDuration)

				// Short timer should fire first
				<-shortTimer
				shortElapsed := time.Since(start)

				<-longTimer
				longElapsed := time.Since(start)

				Expect(shortElapsed).To(BeNumerically("<", longElapsed))
			})
		})
	})

	Describe("Interval", func() {
		Context("when creating an interval with specific duration", func() {
			It("should emit incrementing values at regular intervals", func() {
				interval := 10 * time.Millisecond
				out := op.Take(op.Interval(interval), 3)

				expectedValues := []int{0, 1, 2}
				index := 0

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					Expect(value).To(Equal(expectedValues[index]))

					index++
				}

				Expect(index).To(Equal(3))
			})

			It("should maintain consistent timing between intervals", func() {
				interval := 20 * time.Millisecond
				out := op.Take(op.Interval(interval), 2)

				start := time.Now()

				// First value
				<-out

				// Second value
				<-out
				secondElapsed := time.Since(start)

				// Should be approximately 2x the interval
				expectedSecond := 2 * interval
				tolerance := 5 * time.Millisecond

				Expect(secondElapsed).To(BeNumerically("~", expectedSecond, tolerance))
			})
		})
	})

	Describe("FormSlice", func() {
		Context("when converting a slice to a channel", func() {
			It("should emit all slice elements in order", func() {
				input := []string{"hello", "world", "test"}
				out := op.FormSlice(input)

				results := make([]string, 0, len(input))
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(input))
			})

			It("should work with different types", func() {
				intInput := []int{1, 2, 3, 4, 5}
				out := op.FormSlice(intInput)

				results := make([]int, 0, len(intInput))
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(intInput))
			})

			It("should handle empty slices", func() {
				var emptySlice []string
				out := op.FormSlice(emptySlice)

				count := 0
				for range out {
					count++
				}

				Expect(count).To(Equal(0))
			})
		})
	})

	Describe("FormChannel", func() {
		Context("when converting a channel to a Result channel", func() {
			It("should emit all channel values as Ok results", func() {
				input := make(chan int, 3)
				input <- 10
				input <- 20
				input <- 30
				close(input)

				out := op.FormChannel(input)

				expectedValues := []int{10, 20, 30}
				results := make([]int, 0, len(expectedValues))

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})

			It("should handle unbuffered channels", func() {
				input := make(chan string)
				out := op.FormChannel(input)

				// Send values in a goroutine
				go func() {
					defer close(input)
					input <- "first"
					input <- "second"
				}()

				results := make([]string, 0, 2)
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal([]string{"first", "second"}))
			})

			It("should close output when input channel closes", func() {
				input := make(chan int)
				out := op.FormChannel(input)

				// Close input immediately
				close(input)

				count := 0
				for range out {
					count++
				}

				Expect(count).To(Equal(0))
			})
		})
	})

	Describe("Range", func() {
		Context("when creating a range of numbers", func() {
			It("should emit consecutive integers from start", func() {
				start := 5
				count := 4
				out := op.Range(start, count)

				expectedValues := []int{5, 6, 7, 8}
				results := make([]int, 0, count)

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})

			It("should handle zero count", func() {
				out := op.Range(10, 0)

				count := 0
				for range out {
					count++
				}

				Expect(count).To(Equal(0))
			})

			It("should handle negative start values", func() {
				start := -3
				count := 5
				out := op.Range(start, count)

				expectedValues := []int{-3, -2, -1, 0, 1}
				results := make([]int, 0, count)

				for result := range out {
					Expect(result.IsOk()).To(BeTrue())

					value, err := result.Get()
					Expect(err).To(BeNil())
					results = append(results, value)
				}

				Expect(results).To(Equal(expectedValues))
			})

			It("should handle single value range", func() {
				out := op.Range(42, 1)

				result := <-out
				Expect(result.IsOk()).To(BeTrue())

				value, err := result.Get()
				Expect(value).To(Equal(42))
				Expect(err).To(BeNil())

				// Channel should be closed
				_, ok := <-out
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("Integration with options", func() {
		Context("when using WithBufferSize option", func() {
			It("should create buffered channels", func() {
				out := op.Range(0, 3, op.WithBufferSize(2))

				// Should be able to emit at least buffer size without blocking
				count := 0
				for result := range out {
					Expect(result.IsOk()).To(BeTrue())
					count++
				}

				Expect(count).To(Equal(3))
			})
		})
	})
})
