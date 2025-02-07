package delve

import (
	"fmt"
	"strconv"
)

// Interface represents qualifier to access fields of navigator
type IQual interface {
	// Function to access next part of qualifier
	Next() (string, bool)
	// Function to reset qualifier back to zero offset start state.
	Reset()
}

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

// IGetter defines an interface for navigator data source
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
	if key == -1 {
		key = max(0, len(fl)-1)
	}
	if len(fl) < key {
		return nil, false
	}
	return fl[key], true
}

// GetByQual retrieves a nested value using a compiled qualifier
func (fm Navigator) GetByQual(qual IQual) (any, bool) {
	defer qual.Reset()

	var currentGetter IGetter = fm.source
	if currentGetter == nil {
		return nil, false
	}
	var hasNext bool = true
	var part string

	for hasNext {
		part, hasNext = qual.Next()
		if !hasNext {
			return currentGetter.Get(part)
		}
		if inner, ok := getInnerGetter(part, currentGetter); ok {
			currentGetter = inner
		}
	}
	return nil, false
}

// Returns value by qual or panics
func (fm Navigator) MustGetByQual(qual IQual) any {
	if val, ok := fm.GetByQual(qual); ok {
		return val
	}
	panic(fmt.Sprintf("could not get by qual %v", qual))
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
