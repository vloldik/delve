package delve

func defaultEmptyVal[T any](_default []T) T {
	var factoryDefault T
	return defaultVal(factoryDefault, _default)
}

func defaultVal[T any](factoryDefault T, _default []T) T {
	if len(_default) > 0 {
		return _default[0]
	}
	return factoryDefault
}
