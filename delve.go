package delve

import (
	"fmt"
)

// Interface represents qualifier to access fields of navigator
type IQual interface {
	// Function to access next part of qualifier
	Next() (string, bool)
	// Function to reset qualifier back to zero offset start state.
	Reset()
	// Function to get an independent copy of current qual
	Copy() IQual
}

// ISource defines an interface for navigator data source
type ISource interface {
	Get(string) (any, bool)
	Set(string, any) bool
}

func New(source ISource) *Navigator {
	return &Navigator{source: source}
}

func FromMap(source map[string]any) *Navigator {
	return New(getterMap(source))
}

func FromList(source []any) *Navigator {
	return New(&getterList{source})
}

type Navigator struct {
	source ISource
}

func (fm *Navigator) Get(qual string) (any, bool) {
	return fm.source.Get(qual)
}

func (fm *Navigator) Set(qual string, value any) bool {
	return fm.source.Set(qual, value)
}

// QualGet retrieves a nested value using a compiled qualifier
func (fm *Navigator) QualGet(qual IQual) (any, bool) {
	defer qual.Reset()

	var currentGetter ISource = fm.source
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
		} else {
			return nil, false
		}
	}
	return nil, false
}

func (fm *Navigator) QualSet(qual IQual, value any) bool {
	var currentGetter = fm.source
	if fm.source == nil {
		return false
	}
	var hasNext bool = true
	var part string

	pathExist := true
	for hasNext {
		if !pathExist {
			newGetter := getterMap{}
			if !currentGetter.Set(part, newGetter) {
				return false
			}
			currentGetter = newGetter
			part, hasNext = qual.Next()
			continue
		}
		part, hasNext = qual.Next()
		if !hasNext {
			break
		}
		if inner, ok := getInnerGetter(part, currentGetter); ok {
			currentGetter = inner
		} else {
			pathExist = false
		}
	}

	return currentGetter.Set(part, value)
}

// Returns value by qual or panics
func (fm *Navigator) MustGetByQual(qual IQual) any {
	if val, ok := fm.QualGet(qual); ok {
		return val
	}
	panic(fmt.Sprintf("could not get by qual %v", qual))
}

// getInnerGetter retrieves nested delve or FlexList values for further access
func getInnerGetter(key string, from ISource) (ISource, bool) {
	result, ok := from.Get(key)
	if !ok {
		return nil, false
	}
	switch typed := result.(type) {
	case []any:
		return &getterList{list: typed}, true
	case map[string]any:
		return getterMap(typed), true
	case ISource:
		return typed, true
	default:
		return nil, false
	}
}

func (fm *Navigator) SetMapSource(source map[string]any) {
	fm.SetSource(getterMap(source))
}

func (fm *Navigator) SetListSource(source []any) {
	fm.SetSource(&getterList{source})
}

func (fm *Navigator) SetSource(source ISource) {
	fm.source = source
}
