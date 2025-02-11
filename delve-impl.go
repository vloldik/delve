package delve

import (
	"github.com/vloldik/delve/v2/internal/sources"
	"github.com/vloldik/delve/v2/pkg/idelve"
)

func (fm *navigator) qualGet(qual idelve.IQual) (any, bool) {
	defer qual.Reset()

	var currentGetter idelve.ISource = fm.source
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
		if inner := getInnerGetter(part, currentGetter); inner != nil {
			currentGetter = inner
		} else {
			return nil, false
		}
	}
	return nil, false
}

func (fm *navigator) qualSet(qual idelve.IQual, value any) bool {
	defer qual.Reset()

	var currentGetter = fm.source
	if fm.source == nil {
		return false
	}
	var hasNext bool = true
	var part string

	pathExist := true
	for hasNext {
		if !pathExist {
			newGetter := sources.MapSource{}
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
		if inner := getInnerGetter(part, currentGetter); inner != nil {
			currentGetter = inner
		} else {
			pathExist = false
		}
	}

	return currentGetter.Set(part, value)
}

// getInnerGetter retrieves nested ISource for further access. Returns nil if not successed
func getInnerGetter(key string, from idelve.ISource) idelve.ISource {
	result, ok := from.Get(key)
	if !ok {
		return nil
	}
	return sources.GetSource(result)
}
