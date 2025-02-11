package defaultval

func WithDefaultEmpty[T any](_default []T) T {
	var factoryDefault T
	return WithDefaultVal(factoryDefault, _default)
}

func WithDefaultVal[T any](factoryDefault T, _default []T) T {
	if len(_default) > 0 {
		return _default[0]
	}
	return factoryDefault
}
