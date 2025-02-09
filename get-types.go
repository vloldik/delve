package delve

import (
	"reflect"
)

// Numeric is a type constraint that includes all the common numeric types in Go.
// The '~' before each type means that it includes any type whose *underlying* type is that type.
// For example, `~int` includes `int`, and also any named types defined as `type MyInt int`.
type Numeric interface {
	~float64 | ~float32 | ~int | ~int64 | ~int32 | ~int16 | ~int8 | ~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}

// compareNumerics compares two numeric values of different types.
// It takes two numeric values 'a' and 'b' of potentially different numeric types A and B,
// constrained by the Numeric interface.
// It returns the value of 'b' and a boolean indicating whether 'a' and 'b' represent the same numerical value.
// The comparison is done by converting 'b' back to the type of 'a' and checking for equality.
// This approach handles potential loss of precision during type conversion.
func compareNumerics[A, B Numeric](a A, b B) (B, bool) {
	aBack := A(b)    // Convert b to the type of a. (e.g., if `a` is int32, `b` float64, then `aBack` := int32(b))
	ok := a == aBack // Check if the converted value is equal to the original a. This verifies if the conversion was lossless.
	return b, ok     // Return the original 'b' value and the boolean indicating the success of the conversion.
}

// AnyToNumeric attempts to convert an `any` value (which could be any type) to a specific numeric type T.
// T is constrained by the Numeric interface, meaning it must be one of the numeric types specified in that interface.
// The function returns the converted value of type T and a boolean indicating success.
func AnyToNumeric[T Numeric](num any) (val T, ok bool) {
	ok = true
	switch casted := num.(type) {
	case T:
		return casted, true
	case float64:
		return compareNumerics(casted, T(casted))
	case float32:
		return compareNumerics(casted, T(casted))
	case int:
		return compareNumerics(casted, T(casted))
	case int64:
		return compareNumerics(casted, T(casted))
	case int32:
		return compareNumerics(casted, T(casted))
	case int16:
		return compareNumerics(casted, T(casted))
	case int8:
		return compareNumerics(casted, T(casted))
	case uint:
		return compareNumerics(casted, T(casted))
	case uint64:
		return compareNumerics(casted, T(casted))
	case uint32:
		return compareNumerics(casted, T(casted))
	case uint16:
		return compareNumerics(casted, T(casted))
	case uint8:
		return compareNumerics(casted, T(casted))

	default: // If 'num' is not any of the supported numeric types...
		ok = false // Set ok to false, because we cannot convert it.
	}
	return // Return the zero value of type T and ok=false.  This is the zero value of T and `false`.
}

func getTyped[T any](fm *Navigator, qual IQual, _default ...T) T {
	if val, ok := fm.QualGet(qual); ok {
		if casted, ok := val.(T); ok {
			return casted
		}
	}
	return defaultEmptyVal(_default)
}

func getNumeric[T Numeric](fm *Navigator, qual IQual, _default ...T) T {
	if val, ok := fm.QualGet(qual); ok {
		if casted, ok := AnyToNumeric[T](val); ok {
			return casted
		}
	}
	return defaultEmptyVal(_default)
}

// Len attempts to get the length of a value associated with a given qualifier.
//
// It retrieves the value associated with `qual`. If found, and the value is one of
// the following types: Chan, Map, Array, Slice, or String, the function returns
// the length of the value. Otherwise (the qualifier is not found, or the value
// has a type that doesn't have a defined length concept in Go), it returns -1.
func (fm *Navigator) Len(qual IQual) int {
	val, ok := fm.QualGet(qual)
	if !ok {
		return -1
	}
	refVal := reflect.ValueOf(val)
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
func (fm *Navigator) SafeInterface(qual IQual, defaultVal any) any {
	val, ok := fm.QualGet(qual)
	if !ok {
		return defaultVal
	}
	if defaultVal == nil {
		return val
	}
	if reflect.TypeOf(val).AssignableTo(reflect.TypeOf(defaultVal)) {
		return val
	}
	return defaultVal
}

// Get interface or default.
func (fm *Navigator) Interface(qual IQual, _default ...any) any {
	return getTyped(fm, qual, _default...)
}

// Get string or default
func (fm *Navigator) String(qual IQual, _default ...string) string {
	return getTyped(fm, qual, _default...)
}

// Get boolean or default
func (fm *Navigator) Bool(qual IQual, _default ...bool) bool {
	return getTyped(fm, qual, _default...)
}

// Get complex64 or default
func (fm *Navigator) Complex64(qual IQual, _default ...complex64) complex64 {
	return getTyped(fm, qual, _default...)
}

// Get complex128 or default
func (fm *Navigator) Complex128(qual IQual, _default ...complex128) complex128 {
	return getTyped(fm, qual, _default...)
}

// Get int or default
func (fm *Navigator) Int(qual IQual, _default ...int) int {
	return getNumeric(fm, qual, _default...)
}

// Get int64 or default
func (fm *Navigator) Int64(qual IQual, _default ...int64) int64 {
	return getNumeric(fm, qual, _default...)
}

// Get int32 or default
func (fm *Navigator) Int32(qual IQual, _default ...int32) int32 {
	return getNumeric(fm, qual, _default...)
}

// Get uint or default
func (fm *Navigator) Uint(qual IQual, _default ...uint) uint {
	return getNumeric(fm, qual, _default...)
}

// Get uint64 or default
func (fm *Navigator) Uint64(qual IQual, _default ...uint64) uint64 {
	return getNumeric(fm, qual, _default...)
}

// Get uint32 or default
func (fm *Navigator) Uint32(qual IQual, _default ...uint32) uint32 {
	return getNumeric(fm, qual, _default...)
}

// Get uint16 or default
func (fm *Navigator) Uint16(qual IQual, _default ...uint16) uint16 {
	return getNumeric(fm, qual, _default...)
}

// Get uint8 or default
func (fm *Navigator) Uint8(qual IQual, _default ...uint8) uint8 {
	return getNumeric(fm, qual, _default...)
}

// Get int16 or default
func (fm *Navigator) Int16(qual IQual, _default ...int16) int16 {
	return getNumeric(fm, qual, _default...)
}

// Get int8 or default
func (fm *Navigator) Int8(qual IQual, _default ...int8) int8 {
	return getNumeric(fm, qual, _default...)
}

// Get float64 or default
func (fm *Navigator) Float64(qual IQual, _default ...float64) float64 {
	return getNumeric(fm, qual, _default...)
}

// Get float32 or default
func (fm *Navigator) Float32(qual IQual, _default ...float32) float32 {
	return getNumeric(fm, qual, _default...)
}

func (fm *Navigator) Navigator(qual IQual) *Navigator {
	if val, ok := fm.QualGet(qual); ok {
		switch casted := val.(type) {
		case map[string]any:
			return FromMap(casted)
		case []any:
			return FromList(casted)
		case ISource:
			return New(casted)
		}
	}
	return nil
}
