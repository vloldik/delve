package delve

import (
	"strconv"
)

// Qdelimiter is used to separate nested keys in qualified paths
const DefaultDelimiter = '.'

func New(source IGetter, _delimiter ...rune) *Navigator {
	delimiter := DefaultDelimiter
	if len(_delimiter) > 0 {
		delimiter = _delimiter[0]
	}
	return &Navigator{source: source, delimiter: delimiter}
}

func FromMap(source map[string]any, _delimiter ...rune) *Navigator {
	return New(getterMap(source), _delimiter...)
}

func FromList(source []any, _delimiter ...rune) *Navigator {
	return New(getterList(source), _delimiter...)
}

type Navigator struct {
	source    IGetter
	delimiter rune
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

func (fm *Navigator) SetDelimiter(delimiter rune) {
	fm.delimiter = delimiter
}
