// Package comparator provides types and functions related to
// comparing values.
package comparator

// Comparator is a function type for comparing two values.
//
// It must return
//
//	-1 if a is less than b,
//	 0 if a equals b,
//	 1 if a is greater than b.
type Comparator[T any] func(a, b T) int

// Reverse returns a new comparator that reverses cmp
func Reverse[T any](cmp Comparator[T]) Comparator[T] {
	return func(a, b T) int {
		return cmp(b, a)
	}
}
