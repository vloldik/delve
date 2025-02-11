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

// IterList iterates over a slice stored in a Value, providing type-safe access
// to elements. Requires the Value to contain exactly []V, otherwise does nothing.
//
// Parameters:
//   - val: Value containing the slice to iterate over
//   - callback: Function receiving (index, value) pairs
//
// Type constraints:
//   - V: Must match the slice element type
//
// Example iterating a list of integers:
//
//	delve.IterList[int](myValue, func(i int, v int) {
//	    fmt.Printf("Index %d: %d", i, v)
//	})
func IterList[V any](val *Value, callback func(int, V)) {
	value.IterList(val, callback)
}

// IterMap iterates over a map stored in a Value, providing type-safe access
// to key-value pairs. Requires the Value to contain exactly map[K]V where K is
// comparable, otherwise does nothing.
//
// Parameters:
//   - val: Value containing the map to iterate over
//   - callback: Function receiving (key, value) pairs
//
// Type constraints:
//   - K: Must be comparable and match map key type
//   - V: Must match map value type
//
// Example iterating a string->int map:
//
//	delve.IterMap[string, int](myValue, func(k string, v int) {
//	    fmt.Printf("Key %s: %d", k, v)
//	})
func IterMap[K comparable, V any](val *Value, callback func(K, V)) {
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
