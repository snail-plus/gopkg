package stream

import (
	"fmt"
	"gitee.com/eve_3/gopkg/util/reflect"
	"github.com/spf13/cast"
	"math/rand"
	"sort"
)

// Stream represents a lazily evaluated stream of values.
type Stream[T any] struct {
	ch <-chan T
}

// NewStream creates a new Stream from a slice.
func NewStream[T any](data []T) *Stream[T] {
	ch := make(chan T)

	go func() {
		for _, v := range data {
			ch <- v
		}
		close(ch)
	}()

	return &Stream[T]{ch: ch}
}

func New[T any](data ...T) *Stream[T] {
	return NewStream(data)
}

func Generate[T any](f func() T) *Stream[T] {
	ch := make(chan T)

	go func() {
		for {
			ch <- f()
		}
	}()

	return &Stream[T]{ch: ch}
}

func Concat[T any](a Stream[T], b Stream[T]) *Stream[T] {
	ch := make(chan T)

	go func() {
		slice1 := a.ToSlice()
		for _, v1 := range slice1 {
			ch <- v1
		}

		slice2 := b.ToSlice()
		for _, v2 := range slice2 {
			ch <- v2
		}

	}()

	return &Stream[T]{ch: ch}
}

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 | uintptr |
		float32 | float64 | complex64 | complex128
}

type MapperFunc[T, R any] func(T) []R

/*
Returns a stream consisting of the results of replacing each element of

  - this stream with the contents of a mapped stream produced by applying

  - the provided mapping function to each element

    func doubleStream[T comparable](value T) [][T] {
    return []T{value, value * 2}
    }

    func main() {
    data := []int{1, 2, 3}
    stream := NewStream(data)

    flattenedStream := flatMap(stream, doubleStream[int])

    for val := range flattenedStream.ch {
    fmt.Println(val) // 输出：1 2 2 4 3 6
    }
    }
*/
func FlatMap[T, R any](stream *Stream[T], mapper MapperFunc[T, R]) *Stream[R] {
	outCh := make(chan R)

	go func() {
		defer close(outCh)

		for v := range stream.ch {
			mappedStream := NewStream(mapper(v))
			for mappedVal := range mappedStream.ch {
				outCh <- mappedVal
			}
		}
	}()

	return &Stream[R]{ch: outCh}
}

type ToIntFunc[T any] func(T) int

func (s *Stream[T]) MapToInt(f ToIntFunc[T]) *Stream[int] {
	ch := make(chan int)

	go func() {
		for v := range s.ch {
			ch <- f(v)
		}
		close(ch)
	}()

	return &Stream[int]{ch: ch}
}

type ToInt32Func[T any] func(T) int32

func (s *Stream[T]) MapToInt32(f ToInt32Func[T]) *Stream[int32] {
	ch := make(chan int32)

	go func() {
		for v := range s.ch {
			ch <- f(v)
		}
		close(ch)
	}()

	return &Stream[int32]{ch: ch}
}

type ToInt64Func[T any] func(T) int64

func (s *Stream[T]) MapToInt64(f ToInt64Func[T]) *Stream[int64] {
	ch := make(chan int64)

	go func() {
		for v := range s.ch {
			ch <- f(v)
		}
		close(ch)
	}()

	return &Stream[int64]{ch: ch}
}

type ToFloat32Func[T any] func(T) float32

func (s *Stream[T]) MapToFloat32(f ToFloat32Func[T]) *Stream[float32] {
	ch := make(chan float32)

	go func() {
		for v := range s.ch {
			ch <- f(v)
		}
		close(ch)
	}()

	return &Stream[float32]{ch: ch}
}

type ToFloat64Func[T any] func(T) float64

func (s *Stream[T]) MapToFloat64(f ToFloat64Func[T]) *Stream[float64] {
	ch := make(chan float64)

	go func() {
		for v := range s.ch {
			ch <- f(v)
		}
		close(ch)
	}()

	return &Stream[float64]{ch: ch}
}

// Map applies the given function to each element lazily.
func (s *Stream[T]) Map(f func(T) T) *Stream[T] {
	ch := make(chan T)

	go func() {
		for v := range s.ch {
			ch <- f(v)
		}
		close(ch)
	}()

	return &Stream[T]{ch: ch}
}

// ForEach applies the given function to each element lazily.
func (s *Stream[T]) ForEach(f func(T)) {
	for v := range s.ch {
		f(v)
	}
}

// Peek Returns a stream consisting of the elements of this stream, additionally performing the provided.
func (s *Stream[T]) Peek(f func(T)) *Stream[T] {
	ch := make(chan T)

	go func() {
		for v := range s.ch {
			f(v)
			ch <- v
		}
		close(ch)
	}()

	return &Stream[T]{ch: ch}
}

// TransformStream To change the type of the Stream, we need to introduce a new function, not a method.
// This function takes a Stream[T] and returns a Stream[R] by applying the provided function.
func TransformStream[T any, R any](s *Stream[T], f func(T) R) *Stream[R] {
	ch := make(chan R)

	go func() {
		for v := range s.ch {
			ch <- f(v)
		}
		close(ch)
	}()

	return &Stream[R]{ch: ch}
}

// Transform To change the type of the Stream, we need to introduce a new function, not a method.
// This function takes a Stream[T] and returns a Stream[R] by applying the provided function.
func Transform[T any, R any](data []T, f func(T) R) *Stream[R] {
	return TransformStream(NewStream(data), f)
}

// Filter keeps only the elements that satisfy the given predicate lazily.
func (s *Stream[T]) Filter(f func(T) bool) *Stream[T] {
	ch := make(chan T)

	go func() {
		for v := range s.ch {
			if f(v) {
				ch <- v
			}
		}
		close(ch)
	}()

	return &Stream[T]{ch: ch}
}

// AllMatch keeps only the elements that satisfy the given predicate lazily.
func (s *Stream[T]) AllMatch(f func(T) bool) bool {
	for v := range s.ch {
		if !f(v) {
			return false
		}
	}
	return true
}

// AnyMatch keeps only the elements that satisfy the given predicate lazily.
func (s *Stream[T]) AnyMatch(f func(T) bool) bool {
	for v := range s.ch {
		if f(v) {
			return true
		}
	}
	return false
}

// NoneMatch Returns whether no elements of this stream match the provided predicate.
func (s *Stream[T]) NoneMatch(f func(T) bool) bool {
	return !s.AnyMatch(f)
}

// Contains keeps only the elements that satisfy the given predicate lazily.
func (s *Stream[T]) Contains(in T) bool {
	slice := s.ToSlice()
	for _, item := range slice {

		if &item == &in {
			return true
		}

		if fmt.Sprintf("%v", item) == fmt.Sprintf("%v", in) {
			return true
		}
	}

	return false
}

// Limit returns a new Stream that contains the first n elements lazily.
func (s *Stream[T]) Limit(n int) *Stream[T] {
	ch := make(chan T)

	go func() {
		count := 0
		for v := range s.ch {
			if count >= n {
				break
			}
			ch <- v
			count++
		}
		close(ch)
	}()

	return &Stream[T]{ch: ch}
}

// Skip returns a new Stream that contains the first n elements lazily.
func (s *Stream[T]) Skip(n int) *Stream[T] {
	ch := make(chan T)

	go func() {
		count := 0
		for v := range s.ch {
			count++
			if count <= n {
				break
			}
			ch <- v
		}
		close(ch)
	}()

	return &Stream[T]{ch: ch}
}

// DistinctBy returns a new Stream that contains only distinct elements lazily,
// using the provided keyFunc to determine uniqueness.
func (s *Stream[T]) DistinctBy(keyFunc func(T) any) *Stream[T] {
	ch := make(chan T)
	seen := make(map[any]struct{})

	go func() {
		for v := range s.ch {
			key := keyFunc(v)
			if _, ok := seen[key]; !ok {
				seen[key] = struct{}{}
				ch <- v
			}
		}
		close(ch)
	}()

	return &Stream[T]{ch: ch}
}

// Reduce applies the given binary function to the elements of the Stream,
// starting with the provided initial value, to reduce the Stream to a single value.
func (s *Stream[T]) Reduce(initialValue T, f func(T, T) T) T {
	value := initialValue
	for v := range s.ch {
		value = f(value, v)
	}
	return value
}

// GroupBy collects elements of the Stream into a map keyed by the result of
// applying the given function to each element. The value of each map entry is
// a slice containing the elements that mapped to that key.
func (s *Stream[T]) GroupBy(keyFunc func(T) any) map[any][]T {
	groups := make(map[any][]T)
	for elem := range s.ch {
		key := keyFunc(elem)
		groups[key] = append(groups[key], elem)
	}
	return groups
}

// ToSlice collects all elements from the stream into a slice.
func (s *Stream[T]) ToSlice() []T {
	var slice []T
	for v := range s.ch {
		slice = append(slice, v)
	}
	return slice
}

// Sort collects all elements from the stream, sorts them using the provided less function,
// and returns a new stream with the sorted elements.
func (s *Stream[T]) Sort(less func(a, b T) bool) *Stream[T] {
	slice := s.ToSlice()
	sort.Slice(slice, func(i, j int) bool {
		return less(slice[i], slice[j])
	})
	return NewStream(slice)
}

// Reverse reverses the order of elements in the stream.
// Simple implementation: collects all elements first, then reverses the slice.
func (s *Stream[T]) Reverse() *Stream[T] {
	data := s.Collect()
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return NewStream(data)
}

// Random selects a random subset of elements from the stream.
// If the stream size is less than the requested size, it returns all elements.
func (s *Stream[T]) Random(size ...uint) *Stream[T] {
	var n uint

	var data = s.ToSlice()
	if len(size) > 0 {
		n = size[0]
	} else {
		n = uint(len(data)) // If no size is provided, return all elements (not really random, but for illustration)
	}

	if n > uint(len(data)) {
		n = uint(len(data)) // Adjust size if it exceeds the number of elements in the stream
	}

	indices := rand.Perm(len(data))[:n] // Generate a random permutation of indices and slice to desired size

	resultCh := make(chan T)
	go func() {
		for _, idx := range indices {
			resultCh <- data[idx] // Send the randomly selected elements to the result channel
		}
		close(resultCh)
	}()

	return &Stream[T]{ch: resultCh} // Return a new stream with the randomly selected elements
}

// Count collects all elements from the stream, sorts them using the provided less function,
// and returns a new stream with the sorted elements.
func (s *Stream[T]) Count() int {
	count := 0
	for range s.ch {
		count++
	}
	return count
}

// Sum calculates the sum of all elements in the stream.
func (s *Stream[T]) Sum(key ...string) (sum float64) {
	slice := s.ToSlice()
	if len(key) == 0 {
		for _, v := range slice {
			sum += cast.ToFloat64(v)
		}
	} else {
		for _, item := range slice {
			sum += reflect.GetNumberField(item, key[0])
		}
	}
	return
}

// Max returns the maximum element in the stream.
// NOTE: This is not truly lazy as it needs to consume the entire stream.
func (s *Stream[T]) Max(key ...string) float64 {
	slice := s.ToSlice()
	if slice == nil {
		return 0
	}

	if len(key) == 0 {
		maxVal := cast.ToFloat64(slice[0])
		for _, k := range slice {
			item := cast.ToFloat64(k)
			if item > maxVal {
				maxVal = item
			}
		}

		return maxVal
	}

	maxVal := reflect.GetNumberField(slice[0], key[0])

	for _, v := range slice[1:] {
		item := reflect.GetNumberField(v, key[0])
		if item > maxVal {
			maxVal = item
		}
	}

	return maxVal
}

// Min returns the min element in the stream.
// NOTE: This is not truly lazy as it needs to consume the entire stream.
func (s *Stream[T]) Min(key ...string) float64 {
	slice := s.ToSlice()
	if slice == nil {
		return 0
	}

	if len(key) == 0 {
		minVal := cast.ToFloat64(slice[0])
		for _, k := range slice {
			item := cast.ToFloat64(k)
			if item < minVal {
				minVal = item
			}
		}

		return minVal
	}

	minVal := reflect.GetNumberField(slice[0], key[0])

	for _, v := range slice[1:] {
		item := reflect.GetNumberField(v, key[0])
		if item < minVal {
			minVal = item
		}
	}

	return minVal
}

// Range returns the max - min value of the stream.
func (s *Stream[T]) Range(key ...string) float64 {
	ts := s.Collect()
	maxVal := New(ts...).Max(key...)
	minVal := New(ts...).Min(key...)
	return maxVal - minVal
}

// Avg returns the maximum element in the stream.
// NOTE: This is not truly lazy as it needs to consume the entire stream.
func (s *Stream[T]) Avg(key ...string) float64 {
	slice := s.ToSlice()

	sum := float64(0)

	if len(key) == 0 {
		for _, v := range slice {
			item := cast.ToFloat64(v)
			sum += item
		}
	} else {
		for _, v := range slice {
			item := reflect.GetNumberField(v, key[0])
			sum += item
		}
	}

	return sum / float64(len(slice))
}

// First returns the first element of the stream, if available.
func (s *Stream[T]) First() *T {
	slice := s.ToSlice()
	if len(slice) == 0 {
		return nil
	}

	return &slice[0]
}

// Last returns the last element of the stream, if available.
// 注意：这个方法会消耗整个流，而且如果流是无限的，它会阻塞。
func (s *Stream[T]) Last() *T {
	slice := s.Reverse().ToSlice()
	if len(slice) == 0 {
		return nil
	}
	return &slice[0]
}

// Chunk splits the Stream into sub-slices of a specified size.
func (s *Stream[T]) Chunk(size int) [][]T {
	var chunks [][]T

	var chunk []T
	for v := range s.ch {
		chunk = append(chunk, v)
		if len(chunk) == size {
			chunks = append(chunks, chunk)
			chunk = make([]T, 0, size)
		}
	}

	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}

	return chunks
}

// Collect terminates the Stream by collecting all its values into a slice.
func (s *Stream[T]) Collect() []T {
	var result []T
	for v := range s.ch {
		result = append(result, v)
	}
	return result
}

func ToMap[T any, K comparable](stream *Stream[T], keyFunc func(T) K) map[K]T {
	result := make(map[K]T)
	for v := range stream.ch {
		result[keyFunc(v)] = v
	}
	return result
}
