// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Deep equality test via reflection

package util

import (
	`unsafe`
	`reflect`
)

// visit keeps track of checks in progress. The comparison algorithm assumes
// that all checks in progress are true when it reencounters them. Visited
// comparisons are stored in a map indexed by visit.
type visit struct {
	addr1		unsafe.Pointer
	addr2		unsafe.Pointer
	typ		reflect.Type
}

// Current keeps track of the object and property under examination.
type Current struct {
	Object		string
	Property	string
}

// set changes the current object and property.
func (this *Current) set(obj, prop string) {
	this.Object, this.Property = obj, prop
}

// Delta is a difference in the value of a property in two objects.
type Delta struct {
	Object		string
	Property	string
	Value1		interface{}
	Value2		interface{}
}

// Deltas is an ordered collection of deltas.
type Deltas []*Delta

// add appends a delta to a collection of deltas.
func (this *Deltas) add(obj, prop string, val1, val2 interface{}) {
	*this = append(*this, &Delta{obj, prop, val1, val2})
}

// Tests for deep equality using reflected types. The map argument tracks
// comparisons that have already been seen, which allows short circuiting on
// recursive types.
func deepDiff(v1, v2 reflect.Value, visited map[visit]bool, deltas *Deltas) {

	// Reduce number of entries in visited. For any reference cycle that
	// might be encountered, hard(t) needs to return true for at least one
	// of the types in the cycle.

	hard := func(k reflect.Kind) bool {

		switch k {
		case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
			return true
		}
		return false
	}

	if v1.CanAddr() && v2.CanAddr() && hard(v1.Kind()) {

		addr1 := unsafe.Pointer(v1.UnsafeAddr())
		addr2 := unsafe.Pointer(v2.UnsafeAddr())
		typ := v1.Type()

		// Reduce number of entries in visited by sorting addresses.
		if uintptr(addr1) > uintptr(addr2) {
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.

		v := visit{addr1, addr2, typ}

		if visited[v] {
			return
		}

		visited[v] = true
	}

	switch v1.Kind() {
	case reflect.Array:
		for i := 0; i < v1.Len(); i++ {
			if !deepDiff(v1.Index(i), v2.Index(i), visited, depth+1) {
				return false
			}
		}
		return true
	case reflect.Slice:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for i := 0; i < v1.Len(); i++ {
			if !deepDiff(v1.Index(i), v2.Index(i), visited, depth+1) {
				return false
			}
		}
		return true
	case Interface:
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() == v2.IsNil()
		}
		return deepDiff(v1.Elem(), v2.Elem(), visited, depth+1)
	case Ptr:
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		return deepDiff(v1.Elem(), v2.Elem(), visited, depth+1)
	case Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			if !deepDiff(v1.Field(i), v2.Field(i), visited, depth+1) {
				return false
			}
		}
		return true
	case Map:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for _, k := range v1.MapKeys() {
			val1 := v1.MapIndex(k)
			val2 := v2.MapIndex(k)
			if !val1.IsValid() || !val2.IsValid() || !deepDiff(v1.MapIndex(k), v2.MapIndex(k), visited, depth+1) {
				return false
			}
		}
		return true
	case Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		// Can't do better than this:
		return false
	default:
		// Normal equality suffices
		return valueInterface(v1, false) == valueInterface(v2, false)
		// TODO v1.Interface() == v2.Interface() ... but that uses safe=true
	}
}

// DeepEqual reports whether x and y are ``deeply equal,'' defined as follows.
// Two values of identical type are deeply equal if one of the following cases applies.
// Values of distinct types are never deeply equal.
//
// Array values are deeply equal when their corresponding elements are deeply equal.
//
// Struct values are deeply equal if their corresponding fields,
// both exported and unexported, are deeply equal.
//
// Func values are deeply equal if both are nil; otherwise they are not deeply equal.
//
// Interface values are deeply equal if they hold deeply equal concrete values.
//
// Map values are deeply equal when all of the following are true:
// they are both nil or both non-nil, they have the same length,
// and either they are the same map object or their corresponding keys
// (matched using Go equality) map to deeply equal values.
//
// Pointer values are deeply equal if they are equal using Go's == operator
// or if they point to deeply equal values.
//
// Slice values are deeply equal when all of the following are true:
// they are both nil or both non-nil, they have the same length,
// and either they point to the same initial entry of the same underlying array
// (that is, &x[0] == &y[0]) or their corresponding elements (up to length) are deeply equal.
// Note that a non-nil empty slice and a nil slice (for example, []byte{} and []byte(nil))
// are not deeply equal.
//
// Other values - numbers, bools, strings, and channels - are deeply equal
// if they are equal using Go's == operator.
//
// In general DeepEqual is a recursive relaxation of Go's == operator.
// However, this idea is impossible to implement without some inconsistency.
// Specifically, it is possible for a value to be unequal to itself,
// either because it is of func type (uncomparable in general)
// or because it is a floating-point NaN value (not equal to itself in floating-point comparison),
// or because it is an array, struct, or interface containing
// such a value.
// On the other hand, pointer values are always equal to themselves,
// even if they point at or contain such problematic values,
// because they compare equal using Go's == operator, and that
// is a sufficient condition to be deeply equal, regardless of content.
// DeepEqual has been defined so that the same short-cut applies
// to slices and maps: if x and y are the same slice or the same map,
// they are deeply equal regardless of content.
//
// As DeepEqual traverses the data values it may find a cycle. The
// second and subsequent times that DeepEqual compares two pointer
// values that have been compared before, it treats the values as
// equal rather than examining the values to which they point.
// This ensures that DeepEqual terminates.
func DeepDiff(i1, i2 interface{}) (Deltas, error) {

	deltas := new(Deltas)

	if reflect.DeepEqual(i1, i2) {
		return *deltas, nil
	} else if i1 == nil || i2 == nil {
		return *deltas, fmt.Errorf(`cannot compare nil values`)
	}

	v1, v2 := ValueOf(i1), ValueOf(i2)

	if v1.Type() != v2.Type() {
		return *deltas, fmt.Errorf(`cannot compare different types`)
	} else if !v1.IsValid() || !v2.IsValid() {
		return *deltas, fmt.Errorf(`cannot compare invalid values`)
	}

	current := &Current{Object: `None`, Property: `None`}
	deepDiff(v1, v2, make(map[visit]bool), current, deltas)

	return *deltas, nil
}
