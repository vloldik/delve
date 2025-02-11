package value

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
