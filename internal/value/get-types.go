package value

import (
	"reflect"

	"github.com/vloldik/delve/v3/internal/defaultval"
)

type Value struct {
	original any
}

func New(original any) *Value {
	return &Value{original: original}
}

func getTyped[T any](v *Value, _default ...T) T {
	if val, ok := v.original.(T); ok {
		return val
	}
	return defaultval.WithDefaultEmpty(_default)
}

func getNumeric[T Numeric](v *Value, _default ...T) T {
	if casted, ok := AnyToNumeric[T](v.original); ok {
		return casted
	}
	return defaultval.WithDefaultEmpty(_default)
}

// Returns true if value source is nil
func (val *Value) IsNil() bool {
	return val.original == nil
}

// Len attempts to get the length of a value associated with a given qualifier.
//
// It retrieves the value associated with `qual`. If found, and the value is one of
// the following types: Chan, Map, Array, Slice, or String, the function returns
// the length of the value. Otherwise (the qualifier is not found, or the value
// has a type that doesn't have a defined length concept in Go), it returns -1.
func (val *Value) Len() int {
	refVal := reflect.ValueOf(val.original)
	switch refVal.Kind() {
	case reflect.Chan, reflect.Map, reflect.Array, reflect.Slice, reflect.String:
		return refVal.Len()
	}
	return -1
}

// **Do not works with interface as default value**
//
// SafeInterface retrieves a value associated with the given qualifier `qual`.
// It returns the value if found and if it's type-assignable to the type of `defaultVal`.
// If the qualifier is not found or the type is not assignable, it returns `defaultVal`.
func (val *Value) SafeInterface(defaultVal any) any {
	if defaultVal == nil {
		return val.original
	}
	if reflect.TypeOf(val.original).AssignableTo(reflect.TypeOf(defaultVal)) {
		return val.original
	}
	return defaultVal
}

// Get interface or default.
func (val *Value) Interface(_default ...any) any {
	return getTyped(val, _default...)
}

// Get string or default
func (val *Value) String(_default ...string) string {
	return getTyped(val, _default...)
}

// Get boolean or default
func (val *Value) Bool(_default ...bool) bool {
	return getTyped(val, _default...)
}

// Get complex64 or default
func (val *Value) Complex64(_default ...complex64) complex64 {
	return getTyped(val, _default...)
}

// Get complex128 or default
func (val *Value) Complex128(_default ...complex128) complex128 {
	return getTyped(val, _default...)
}

// Get int or default
func (val *Value) Int(_default ...int) int {
	return getNumeric(val, _default...)
}

// Get int64 or default
func (val *Value) Int64(_default ...int64) int64 {
	return getNumeric(val, _default...)
}

// Get int32 or default
func (val *Value) Int32(_default ...int32) int32 {
	return getNumeric(val, _default...)
}

// Get uint or default
func (val *Value) Uint(_default ...uint) uint {
	return getNumeric(val, _default...)
}

// Get uint64 or default
func (val *Value) Uint64(_default ...uint64) uint64 {
	return getNumeric(val, _default...)
}

// Get uint32 or default
func (val *Value) Uint32(_default ...uint32) uint32 {
	return getNumeric(val, _default...)
}

// Get uint16 or default
func (val *Value) Uint16(_default ...uint16) uint16 {
	return getNumeric(val, _default...)
}

// Get uint8 or default
func (val *Value) Uint8(_default ...uint8) uint8 {
	return getNumeric(val, _default...)
}

// Get int16 or default
func (val *Value) Int16(_default ...int16) int16 {
	return getNumeric(val, _default...)
}

// Get int8 or default
func (val *Value) Int8(_default ...int8) int8 {
	return getNumeric(val, _default...)
}

// Get float64 or default
func (val *Value) Float64(_default ...float64) float64 {
	return getNumeric(val, _default...)
}

// Get float32 or default
func (val *Value) Float32(_default ...float32) float32 {
	return getNumeric(val, _default...)
}

// Get string slice or default
func (val *Value) StringSlice(_default ...[]string) []string {
	return getTyped(val, _default...)
}

// Get bool slice or default
func (val *Value) BoolSlice(_default ...[]bool) []bool {
	return getTyped(val, _default...)
}

// Get int slice or default
func (val *Value) IntSlice(_default ...[]int) []int {
	return getTyped(val, _default...)
}

// Get int64 slice or default
func (val *Value) Int64Slice(_default ...[]int64) []int64 {
	return getTyped(val, _default...)
}

// Get float64 slice or default
func (val *Value) Float64Slice(_default ...[]float64) []float64 {
	return getTyped(val, _default...)
}

// Get string map or default
func (val *Value) StringMap(_default ...map[string]string) map[string]string {
	return getTyped(val, _default...)
}

// Get interface map or default
func (val *Value) InterfaceMap(_default ...map[string]any) map[string]any {
	return getTyped(val, _default...)
}

// Get bool map or default
func (val *Value) BoolMap(_default ...map[string]bool) map[string]bool {
	return getTyped(val, _default...)
}

// Get int map or default
func (val *Value) IntMap(_default ...map[string]int) map[string]int {
	return getTyped(val, _default...)
}

// Get int64 map or default
func (val *Value) Int64Map(_default ...map[string]int64) map[string]int64 {
	return getTyped(val, _default...)
}

// Get float64 map or default
func (val *Value) Float64Map(_default ...map[string]float64) map[string]float64 {
	return getTyped(val, _default...)
}
