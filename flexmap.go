package flexmap

import (
	"strconv"
	"strings"
)

// Delemiter of query: "."
const QDelemiter = "."

type FlexMap map[string]any
type FlexList []any
type IAnyGetter interface {
	Get(k string) (any, bool)
}

func (fm FlexMap) Get(key string) (any, bool) {
	value, ok := fm[key]
	return value, ok
}

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

func GetTyped[T any](key string, from IAnyGetter, allowNil_ ...bool) (val T, ok bool) {
	var allowNil bool
	if len(allowNil_) > 0 {
		allowNil = allowNil_[0]
	}
	untyped, ok := from.Get(key)
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
