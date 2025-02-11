package value

// IterList iterates over a list stored in the Value, calling the callback function
// for each element with its index and value.
//
// The Value must hold a slice of type []V, otherwise the function does nothing.
// If the Value is nil, or the underlying value is not a slice of the specified type,
// the callback will not be executed.
//
// Example:
//
//	IterList[int](myValue, func(i int, v int) {
//		fmt.Printf("Index: %d, Value: %d\n", i, v)
//	})
func IterList[V any](val *Value, callback func(int, V)) {
	if val.original == nil {
		return
	}
	list, ok := val.original.([]V)
	if !ok {
		return
	}
	for i, v := range list {
		callback(i, v)
	}
}

// IterMap iterates over a map stored in the Value, calling the callback function
// for each key-value pair.
//
// The Value must hold a map of type map[K]V, where K is comparable, otherwise
// the function does nothing.  If the Value is nil, or the underlying value is not a map
// of the specified type, the callback will not be executed.
//
// Example:
//
//	IterMap[string, int](myValue, func(k string, v int) {
//		fmt.Printf("Key: %s, Value: %d\n", k, v)
//	})
func IterMap[K comparable, V any](val *Value, callback func(K, V)) {
	if val.original == nil {
		return
	}
	mMap, ok := val.original.(map[K]V)
	if !ok {
		return
	}
	for k, v := range mMap {
		callback(k, v)
	}
}
