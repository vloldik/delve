package flexmap

type Numeric interface {
	float64 | float32 | int | int64 | int32 | int16 | int8 | uint | uint64 | uint32 | uint16 | uint8
}

func AnyToNumeric[T Numeric](num any) (val T, ok bool) {
	ok = true
	switch casted := num.(type) {
	case T:
		return val, true
	case float64:
		val = T(casted)
	case float32:
		val = T(casted)
	case int:
		val = T(casted)
	case int64:
		val = T(casted)
	case int32:
		val = T(casted)
	case int16:
		val = T(casted)
	case int8:
		val = T(casted)
	case uint:
		val = T(casted)
	case uint64:
		val = T(casted)
	case uint32:
		val = T(casted)
	case uint16:
		val = T(casted)
	case uint8:
		val = T(casted)
	default:
		ok = false
	}
	return
}

func getTyped[T any](fm *FlexMap, qual CompiledQual, _default ...T) T {
	var defaultVal T
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := val.(T); ok {
			return casted
		}
	}
	return defaultVal
}

// Get string or default
func (fm *FlexMap) String(qual CompiledQual, _default ...string) string {
	return getTyped(fm, qual, _default...)
}

// Get boolean or default
func (fm *FlexMap) Bool(qual CompiledQual, _default ...bool) bool {
	return getTyped(fm, qual, _default...)
}

// Get byte slice or default
func (fm *FlexMap) ByteSlice(qual CompiledQual, _default ...[]byte) []byte {
	return getTyped(fm, qual, _default...)
}

// Get rune slice or default
func (fm *FlexMap) RuneSlice(qual CompiledQual, _default ...[]rune) []rune {
	return getTyped(fm, qual, _default...)
}

// Get complex64 or default
func (fm *FlexMap) Complex64(qual CompiledQual, _default ...complex64) complex64 {
	return getTyped(fm, qual, _default...)
}

// Get complex128 or default
func (fm *FlexMap) Complex128(qual CompiledQual, _default ...complex128) complex128 {
	return getTyped(fm, qual, _default...)
}

// Get string slice or default
func (fm *FlexMap) StringSlice(qual CompiledQual, _default ...[]string) []string {
	return getTyped(fm, qual, _default...)
}

// Get map[string]string or default
func (fm *FlexMap) StringMap(qual CompiledQual, _default ...map[string]string) map[string]string {
	return getTyped(fm, qual, _default...)
}

// Get int or default
func (fm *FlexMap) Int(qual CompiledQual, _default ...int) int {
	var defaultVal int
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[int](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get int64 or default
func (fm *FlexMap) Int64(qual CompiledQual, _default ...int64) int64 {
	var defaultVal int64
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[int64](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get int32 or default
func (fm *FlexMap) Int32(qual CompiledQual, _default ...int32) int32 {
	var defaultVal int32
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[int32](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get uint or default
func (fm *FlexMap) Uint(qual CompiledQual, _default ...uint) uint {
	var defaultVal uint
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[uint](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get uint64 or default
func (fm *FlexMap) Uint64(qual CompiledQual, _default ...uint64) uint64 {
	var defaultVal uint64
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[uint64](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get uint32 or default
func (fm *FlexMap) Uint32(qual CompiledQual, _default ...uint32) uint32 {
	var defaultVal uint32
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[uint32](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get uint16 or default
func (fm *FlexMap) Uint16(qual CompiledQual, _default ...uint16) uint16 {
	var defaultVal uint16
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[uint16](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get uint8 or default
func (fm *FlexMap) Uint8(qual CompiledQual, _default ...uint8) uint8 {
	var defaultVal uint8
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[uint8](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get int16 or default
func (fm *FlexMap) Int16(qual CompiledQual, _default ...int16) int16 {
	var defaultVal int16
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[int16](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get int8 or default
func (fm *FlexMap) Int8(qual CompiledQual, _default ...int8) int8 {
	var defaultVal int8
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[int8](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get float64 or default
func (fm *FlexMap) Float64(qual CompiledQual, _default ...float64) float64 {
	var defaultVal float64
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[float64](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get float64 or default
func (fm *FlexMap) Float32(qual CompiledQual, _default ...float32) float32 {
	var defaultVal float32
	if len(_default) > 0 {
		defaultVal = _default[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[float32](val); ok {
			return casted
		}
	}
	return defaultVal
}

func (fm *FlexMap) FlexMap(qual CompiledQual) *FlexMap {
	var defaultVal FlexMap
	if val, ok := fm.GetByQual(qual); ok {
		switch casted := val.(type) {
		case map[string]any:
			return FromMap(casted)
		case []any:
			return FromList(casted)
		case IGetter:
			return New(casted)
		}
	}
	return &defaultVal
}
