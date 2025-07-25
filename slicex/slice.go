//nolint:gosec,gocritic
package slicex

import (
	"math/rand"
	"reflect"
	"sort"
	"strings"
)

// All iterate over all slice elements and stop when func returns false.
//
// If was iterated all elements - returns true
func All[T any](source []T, fn func(T) bool) bool {
	for _, element := range source {
		if !fn(element) {
			return false
		}
	}

	return true
}

// Any find element by provided condition func and if found - returns true
func Any[T any](source []T, fn func(T) bool) bool {
	for _, element := range source {
		if fn(element) {
			return true
		}
	}

	return false
}

// Each iterate over all slice elements and run provided function
func Each[T any](source []T, fn func(int, T)) {
	for index, element := range source {
		fn(index, element)
	}
}

// EachErr iterate over all slice elements and run provided function.
//
// Stop iterating when provided function returns error
func EachErr[T any](source []T, fn func(int, T) error) error {
	for index, element := range source {
		if err := fn(index, element); err != nil {
			return err
		}
	}

	return nil
}

// Filter slice with provided condition func.
//
// Element appends to new slice if condition func returns true.
//
// Important: function returns new slice
func Filter[T any](source []T, fn func(T) bool) []T {
	dst := make([]T, 0, len(source))
	for _, element := range source {
		if fn(element) {
			dst = append(dst, element)
		}
	}
	return dst
}

// FilterNot slice with provided condition func.
//
// Element appends to new slice if condition func returns false.
//
// Important: function returns new slice
func FilterNot[T any](source []T, fn func(T) bool) []T {
	dst := make([]T, 0, len(source))
	for _, element := range source {
		if !fn(element) {
			dst = append(dst, element)
		}
	}
	return dst
}

// Single returns element and found boolean by provided condition func.
//
// If element found - returns true
func Single[T any](source []T, fn func(T) bool) (T, bool) {
	for _, element := range source {
		if fn(element) {
			return element, true
		}
	}

	var empty T
	return empty, false
}

// Exist check if element exist by provided condition func.
//
// The check performs on first matched element
func Exist[T any](source []T, fn func(T) bool) bool {
	_, ok := Single(source, func(t T) bool { return fn(t) })
	return ok
}

// First returns element and found boolean by provided condition func.
//
// If element found - returns true.
//
// Element start matching from the start of the slice
func First[T any](source []T, fn func(T) bool) (T, bool) {
	return Single(source, fn)
}

// Last returns element and found boolean by provided condition func.
//
// If element found - returns true.
//
// Element start matching from the end of the slice
func Last[T any](source []T, fn func(T) bool) (T, bool) {
	for i := len(source) - 1; i >= 0; i-- {
		if fn(source[i]) {
			return source[i], true
		}
	}

	var empty T
	return empty, false
}

// Contains check if element exist in slice.
//
// Could be provided custom comparing function, by default compares by using reflect.DeepEqual.
//
// The check performs on first matched element
func Contains[T any](source []T, value T, fn ...func(T, T) bool) bool {
	var compareFunc func(T, T) bool
	if len(fn) > 0 {
		compareFunc = func(a, b T) bool {
			return fn[0](a, b)
		}
	} else {
		compareFunc = func(a, b T) bool {
			return reflect.DeepEqual(a, b)
		}
	}

	_, ok := Single(source, func(element T) bool {
		return compareFunc(element, value)
	})
	return ok
}

// Get return element by index.
//
// If index is out of slice range - returns empty value of slice type
func Get[T any](source []T, index int) T {
	if index < 0 || index > len(source) {
		var empty T
		return empty
	}

	return source[index]
}

// Map convert slice to another type of slice by provided converting function
func Map[T any, U any](source []T, fn func(T) U) []U {
	newSlice := make([]U, len(source))
	for i, element := range source {
		newSlice[i] = fn(element)
	}
	return newSlice
}

// MapErr convert slice to another type of slice by provided converting function.
//
// Converting function can return error
func MapErr[T any, U any](source []T, fn func(T) (U, error)) ([]U, error) {
	newSlice := make([]U, len(source))
	for i, element := range source {
		mapped, err := fn(element)
		if err != nil {
			return nil, err
		}

		newSlice[i] = mapped
	}
	return newSlice, nil
}

// Reverse slice
//
// Important: function returns new slice
func Reverse[T any](source []T) []T {
	length := len(source)
	out := make([]T, length)
	copy(out, source)
	for i := 0; i < length/2; i++ {
		out[i], out[length-i-1] = out[length-i-1], out[i]
	}
	return out
}

// Shuffle set elements in slice by random indexes.
//
// Could be provided custom rand.Source implementation.
//
// Important: function returns new slice
func Shuffle[T any](source []T, r ...rand.Source) []T {
	out := make([]T, len(source))
	copy(out, source)
	if r == nil || len(r) == 0 {
		rand.Shuffle(len(out), func(i, j int) {
			out[i], out[j] = out[j], out[i]
		})

		return out
	}

	rnd := rand.New(r[0])
	rnd.Shuffle(len(out), func(i, j int) {
		out[i], out[j] = out[j], out[i]
	})
	return out
}

// Sort slice by provided compare function.
//
// Important: function returns new slice
func Sort[T any](source []T, less func(a, b T) bool) []T {
	if len(source) <= 1 {
		return source
	}

	out := make([]T, len(source))
	copy(out, source)
	sort.Slice(out, func(i, j int) bool {
		return less(out[i], out[j])
	})
	return out
}

// Add appends new elements to slice.
//
// Important: function returns new slice
func Add[T any](source []T, elements ...T) []T {
	return Set(source, len(source), elements...)
}

// Join unions provided slices into one
func Join[T any](joins ...[]T) []T {
	var capacity int
	for _, join := range joins {
		capacity += len(join)
	}

	out := make([]T, 0, capacity)
	for _, join := range joins {
		out = append(out, join...)
	}
	return out
}

// AddLeft append new elements to the start of slice.
//
// Important: function returns new slice
func AddLeft[T any](source []T, elements ...T) []T {
	return Set(source, 0, elements...)
}

// Set append new elements to slice on provided index.
//
// Important: function returns new slice
func Set[T any](source []T, index int, elements ...T) []T {
	if index < 0 {
		index = 0
	}

	if index >= len(source) {
		return append(source, elements...)
	}

	return append(source[:index], append(elements, source[index:]...)...)
}

// Remove delete elements from slice by provided indexes.
//
// Important: function returns new slice
func Remove[T any](source []T, index ...int) []T {
	if len(index) == 0 {
		return source
	}

	if len(index) == 1 {
		i := index[0]

		if i < 0 || i >= len(source) {
			return source
		}

		return append(source[:i], source[i+1:]...)
	}

	sort.Ints(index)

	dst := make([]T, 0, len(source))

	prev := 0
	for _, i := range index {
		if i < 0 || i >= len(source) {
			continue
		}

		dst = append(dst, source[prev:i]...)
		prev = i + 1
	}

	return append(dst, source[prev:]...)
}

// IndexOf return index of found element by provided condition func.
//
// If element not found - returns -1
func IndexOf[T any](source []T, fn func(T) bool) int {
	for index, element := range source {
		if fn(element) {
			return index
		}
	}

	return -1
}

// RemoveWhere delete element from slice by provided condition func.
//
// Important: function returns new slice
func RemoveWhere[T any](source []T, fn func(T) bool) []T {
	index := IndexOf(source, fn)
	if index == -1 {
		return source
	}

	return Remove(source, index)
}

// SliceAny convert slice to "any" type elements slice
func SliceAny[T any](source []T, fn ...func(T) any) []any {
	var mapFn func(T) any
	if len(fn) > 0 {
		mapFn = fn[0]
	}

	sliceAny := make([]any, len(source))
	for i, element := range source {
		if mapFn != nil {
			sliceAny[i] = mapFn(element)
		} else {
			sliceAny[i] = element
		}
	}
	return sliceAny
}

// Sub return "sub slice" by provided start & end indexes.
//
// Important: function returns new slice
func Sub[T any](source []T, start, end int) []T {
	sub := make([]T, 0)
	if start < 0 || end < 0 {
		return sub
	}

	if start >= end {
		return sub
	}

	length := len(source)
	if start < length {
		if end <= length {
			sub = source[start:end]
		} else {
			zeroArray := make([]T, end-length)
			sub = append(source[start:length], zeroArray[:]...)
		}
	} else {
		zeroArray := make([]T, end-start)
		sub = zeroArray[:]
	}

	return sub
}

// Limit return "limited slice" by provided limit size.
//
// # If provided limit is 0 or slice is empty returns empty slice
//
// Important: function returns new slice
func Limit[T any](source []T, limit int) []T {
	if limit == 0 || source == nil || len(source) == 0 {
		return []T{}
	}

	if limit > len(source) {
		return source
	}

	return Sub(source, 0, limit)
}

// JoinString build string from slice elements.
//
// Every element string builds from provided func.
//
// Could be provided custom separator between element strings
func JoinString[T any](source []T, joiner func(T) string, sep ...string) string {
	result := strings.Builder{}
	for index, element := range source {
		result.WriteString(joiner(element))

		if index < len(source)-1 {
			separator := ","
			if len(sep) > 0 {
				separator = sep[0]
			}
			result.WriteString(separator)
		}
	}
	return result.String()
}

// Unique make slice unique by provided condition func.
//
// Important: function returns new slice
func Unique[T any](source []T, fn func(a, b T) bool) []T {
	if source == nil || len(source) == 0 {
		return source
	}

	uniqueSource := make([]T, 0, len(source))

	for _, element := range source {
		isUnique := true
		for _, uItem := range uniqueSource {
			if fn(element, uItem) {
				isUnique = false
				break
			}
		}
		if isUnique {
			uniqueSource = append(uniqueSource, element)
		}
	}

	return uniqueSource
}

// UniqueComparable make slice unique.
//
// Important: function returns new slice
func UniqueComparable[T comparable](source []T) []T {
	if source == nil || len(source) == 0 {
		return source
	}

	uniqueSource := make([]T, 0, len(source))

	for _, element := range source {
		isUnique := true
		for _, uItem := range uniqueSource {
			if element == uItem {
				isUnique = false
				break
			}
		}
		if isUnique {
			uniqueSource = append(uniqueSource, element)
		}
	}

	return uniqueSource
}

// AreUnique compares slice for unique elements by provided func
func AreUnique[T any](source []T, fn func(a, b T) bool) bool {
	return len(Unique(source, fn)) == len(source)
}

// AreUniqueComparable compares slice of comparable types for unique elements.
//
// Uses AreUnique function with default provided func
func AreUniqueComparable[T comparable](source []T) bool {
	return AreUnique(source, func(a, b T) bool {
		return a == b
	})
}

// AreEqual compares two slices by using provided func
func AreEqual[T any](source []T, against []T, fn func(T, T) bool) bool {
	if len(source) != len(against) {
		return false
	}

	for idx := range source {
		if !fn(source[idx], against[idx]) {
			return false
		}
	}

	return true
}

// AreEqualComparable compares two slices of comparable types.
//
// Uses AreEqual function with default provided func
func AreEqualComparable[T comparable](source []T, against []T) bool {
	return AreEqual(source, against, func(a, b T) bool {
		return a == b
	})
}

// Chunk divide slice for sub-slices by provided chunk size.
//
// Example:
//
//	texts := []string{"text #1", "text #2", "text #3", "text #4", "text #5"}
//	chunks := list.Chunk(texts, 2)
//	fmt.Println(chunks) // [[text #1 text #2] [text #3 text #4] [text #5]]
func Chunk[T ~[]E, E any](source T, size int) []T {
	if size <= 0 {
		size = len(source)
	}

	chunks := make([]T, 0, len(source)/size+1)

	for i := 0; i < len(source); i += size {
		end := i + size
		if end > len(source) {
			end = len(source)
		}

		chunks = append(chunks, source[i:end])
	}

	return chunks
}
