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

// Get string or default
func (fm FlexMap) String(qual string, default_ ...string) string {
	var defaultVal string
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := GetTypedByQual[string](qual, fm); ok {
		return val
	}
	return defaultVal
}

// Get int or default
func (fm FlexMap) Int(qual string, default_ ...int) int {
	var defaultVal int
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[int](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get int64 or default
func (fm FlexMap) Int64(qual string, default_ ...int64) int64 {
	var defaultVal int64
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[int64](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get int32 or default
func (fm FlexMap) Int32(qual string, default_ ...int32) int32 {
	var defaultVal int32
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[int32](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get uint or default
func (fm FlexMap) Uint(qual string, default_ ...uint) uint {
	var defaultVal uint
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[uint](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get uint64 or default
func (fm FlexMap) Uint64(qual string, default_ ...uint64) uint64 {
	var defaultVal uint64
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[uint64](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get uint32 or default
func (fm FlexMap) Uint32(qual string, default_ ...uint32) uint32 {
	var defaultVal uint32
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[uint32](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get uint16 or default
func (fm FlexMap) Uint16(qual string, default_ ...uint16) uint16 {
	var defaultVal uint16
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[uint16](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get uint8 or default
func (fm FlexMap) Uint8(qual string, default_ ...uint8) uint8 {
	var defaultVal uint8
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[uint8](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get int16 or default
func (fm FlexMap) Int16(qual string, default_ ...int16) int16 {
	var defaultVal int16
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[int16](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get int8 or default
func (fm FlexMap) Int8(qual string, default_ ...int8) int8 {
	var defaultVal int8
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[int8](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get float64 or default
func (fm FlexMap) Float64(qual string, default_ ...float64) float64 {
	var defaultVal float64
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[float64](val); ok {
			return casted
		}
	}
	return defaultVal
}

// Get float64 or default
func (fm FlexMap) Float32(qual string, default_ ...float32) float32 {
	var defaultVal float32
	if len(default_) > 0 {
		defaultVal = default_[0]
	}
	if val, ok := fm.GetByQual(qual); ok {
		if casted, ok := AnyToNumeric[float32](val); ok {
			return casted
		}
	}
	return defaultVal
}
