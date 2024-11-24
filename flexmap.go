package flexmap

import (
	"strconv"
	"strings"
)

// QDelemiter is used to separate nested keys in qualified paths
var QDelemiter = "."

// FlexMap is a map type that can store values of any type
type FlexMap map[string]any

// FlexList is a slice type that can store values of any type
type FlexList []any

// IAnyGetter defines an interface for getting values by key
type IAnyGetter interface {
	Get(k string) (any, bool)
}

// Get retrieves a value from FlexMap by key
func (fm FlexMap) Get(key string) (any, bool) {
	value, ok := fm[key]
	return value, ok
}

// Get retrieves a value from FlexList by index (passed as string)
func (fl FlexList) Get(uncasted string) (any, bool) {
	key, err := strconv.Atoi(uncasted)
	if err != nil {
		return nil, false
	}
	if len(fl) < key {
		return nil, false
	}
	return fl[key], true
}

// GetByQual retrieves a nested value using a qualified path (e.g. "a.b.c")
func (fm FlexMap) GetByQual(qual string) (any, bool) {
	parts := strings.Split(qual, QDelemiter)
	var currentGetter IAnyGetter = fm
	for i, part := range parts {
		if i == len(parts)-1 {
			return currentGetter.Get(part)
		}
		if value, ok := GetInnerGetter(part, currentGetter); ok {
			currentGetter = value
		} else {
			return nil, false
		}
	}
	return nil, false
}

// Returns value by qual or panics
func (fm FlexMap) MustGetByQual(qual string) any {
	if val, ok := fm.GetByQual(qual); ok {
		return val
	}
	panic("could not get by qual " + qual)
}

// GetInnerGetter retrieves nested FlexMap or FlexList values for further access
func GetInnerGetter(key string, from IAnyGetter) (IAnyGetter, bool) {
	result, ok := from.Get(key)
	if !ok {
		return nil, false
	}
	switch typed := result.(type) {
	case []any:
		return FlexList(typed), true
	case map[string]any:
		return FlexMap(typed), true
	default:
		return nil, false
	}
}

// GetTypedByQual retrieves a nested value and attempts to cast it to the specified type T
// Optional allowNil parameter controls whether nil values are considered valid
func GetTypedByQual[T any](qual string, from FlexMap, allowNil_ ...bool) (val T, ok bool) {
	var allowNil bool
	if len(allowNil_) > 0 {
		allowNil = allowNil_[0]
	}
	untyped, ok := from.GetByQual(qual)
	if !ok {
		return
	}
	switch typed := untyped.(type) {
	case T:
		return typed, true
	case nil:
		ok = allowNil
		return
	default:
		return
	}
}

// Returns typed value by qual or panics
func MustGetTypedByQual[T any](qual string, from FlexMap, allowNil_ ...bool) T {
	if val, ok := GetTypedByQual[T](qual, from, allowNil_...); ok {
		return val
	}
	panic("could not get by qual: " + qual)
}
