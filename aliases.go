package delve

import (
	"github.com/vloldik/delve/v3/internal/quals"
	"github.com/vloldik/delve/v3/internal/value"
	"github.com/vloldik/delve/v3/pkg/idelve"
)

// Value represents a wrapped arbitrary value with type-safe access methods.
// This type alias exposes the internal value.Value type through the delve package.
// Use the Get/QGet methods to obtain Value instances from a Navigator.
type Value = value.Value

// IterList provides a type-safe way to iterate over a slice contained within a `Value`.
//
// It requires the `Value` to hold a slice of type `[]V`. If the `Value` contains
// any other type, the function does nothing.  The provided callback function
// is invoked for each element in the slice, receiving the index and the
// element value.
//
// Parameters:
//   - val: A pointer to a `Value` instance containing the slice.
//   - callback: A function that takes the index (int) and the element value (V) as arguments.
//     It should return a boolean. The loop breaks if this boolean is true.
//
// Type constraints:
//   - V: The type of elements within the slice.  This type must match the actual
//     slice element type contained in the `Value`.
//
// Example:
//
//		delve.IterList(myValue, func(i int, v int) bool {
//		    fmt.Printf("Index %d: %d\n", i, v)
//	     	return false // continue
//		})
func IterList[V any](val *Value, callback func(int, V) bool) {
	value.IterList(val, callback)
}

// IterMap provides a type-safe way to iterate over a map contained within a `Value`.
//
// It requires the Value to hold a map of type `map[K]V`. If the `Value`
// contains any other type, the function does nothing.  The callback
// function is invoked for each key-value pair in the map, receiving the key
//
//	and value.
//
// Parameters:
//   - val: A pointer to a `Value` instance containing the map.
//   - callback: A function to be called for each key-value pair.  It takes
//     the key (K) and value (V) as arguments. It should return a boolean.
//     The loop breaks if this boolean is true.
//
// Type constraints:
//   - K: The type of the map keys. This must be a comparable type and match
//     the actual map key type in the `Value`.
//   - V: The type of the map values. This must match the actual map value
//     type in the `Value`.
//
// Example:
//
//		delve.IterMap(myValue, func(k string, v int) bool {
//		    fmt.Printf("Key %s: %d\n", k, v)
//	     return false; // continue
//		})
func IterMap[K comparable, V any](val *Value, callback func(K, V) bool) {
	value.IterMap(val, callback)
}

// Q creates a new qualifier from a string path using optional delimiters.
// Default delimiter is '.'. Qualifiers can be reused for multiple operations.
//
// Use for one-off path operations or when delimiter customization is needed.
// For frequently used paths, consider using CQ with compiled qualifiers.
//
// Example:
//
//	qual := delve.Q("user.address.street")
//	street := navigator.QGet(qual)
func Q(qual string, _delimiter ...rune) idelve.IQual {
	return quals.Q(qual, _delimiter...)
}

// CQ creates a pre-compiled qualifier from a string path with optional delimiters.
// Compiled qualifiers offer better performance for frequently used paths.
//
// Prefer over Q when:
// - The same path is used multiple times
// - Working with performance-sensitive code
// - Delimiters need to be specified once
//
// Example:
//
//	var userStreetQual = delve.CQ("user.address.street")
//	street := navigator.QGet(userStreetQual)
func CQ(qual string, _delimiter ...rune) idelve.IQual {
	return quals.CQ(qual, _delimiter...)
}
