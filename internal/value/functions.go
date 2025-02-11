package value

// See delve.IterList
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

// See delve.IterMap
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
