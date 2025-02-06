package delve

import (
	"strconv"
)

func New(source IGetter) *Navigator {
	return &Navigator{source: source}
}

func FromMap(source map[string]any) *Navigator {
	return New(getterMap(source))
}

func FromList(source []any) *Navigator {
	return New(getterList(source))
}

type Navigator struct {
	source IGetter
}

type getterMap map[string]any
type getterList []any

// IGetter defines an interface for getting values by key
type IGetter interface {
	Get(string) (any, bool)
}

// Get retrieves a value from delve by key
func (fm getterMap) Get(key string) (any, bool) {
	value, ok := fm[key]
	return value, ok
}

// Get retrieves a value from FlexList by index (passed as string)
func (fl getterList) Get(uncasted string) (any, bool) {
	key, err := strconv.Atoi(uncasted)
	if err != nil {
		return nil, false
	}
	if len(fl) < key {
		return nil, false
	}
	return fl[key], true
}

// GetByQual retrieves a nested value using a compiled qualifier
func (fm Navigator) GetByQual(qual CompiledQual) (any, bool) {
	var currentGetter IGetter = fm.source
	if currentGetter == nil {
		return nil, false
	}
	for i, part := range qual {
		if i == len(qual)-1 {
			return currentGetter.Get(part)
		}
		if value, ok := getInnerGetter(part, currentGetter); ok {
			currentGetter = value
		} else {
			return nil, false
		}
	}
	return nil, false
}

// Returns value by qual or panics
func (fm Navigator) MustGetByQual(qual CompiledQual) any {
	if val, ok := fm.GetByQual(qual); ok {
		return val
	}
	panic("could not get by qual " + qual.String())
}

// getInnerGetter retrieves nested delve or FlexList values for further access
func getInnerGetter(key string, from IGetter) (IGetter, bool) {
	result, ok := from.Get(key)
	if !ok {
		return nil, false
	}
	switch typed := result.(type) {
	case []any:
		return getterList(typed), true
	case map[string]any:
		return getterMap(typed), true
	case IGetter:
		return typed, true
	default:
		return nil, false
	}
}

func (fm *Navigator) SetMapSource(source map[string]any) {
	fm.SetSource(getterMap(source))
}

func (fm *Navigator) SetListSource(source []any) {
	fm.SetSource(getterList(source))
}

func (fm *Navigator) SetSource(source IGetter) {
	fm.source = source
}
