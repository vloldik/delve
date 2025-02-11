package sources

import "strconv"

type ListSource struct {
	list []any
}

func NewList(list []any) *ListSource {
	return &ListSource{list: list}
}

func (fl *ListSource) parseIndex(stringIndex string) (int, bool) {
	key, err := strconv.Atoi(stringIndex)
	if err != nil {
		return -1, false
	}
	if key < 0 {
		key = len(fl.list) + key
	}
	if key >= len(fl.list) || key < 0 {
		return -1, false
	}
	return key, true
}

// Get retrieves a value from FlexList by index (passed as string)
func (fl *ListSource) Get(uncasted string) (any, bool) {
	if index, ok := fl.parseIndex(uncasted); ok {
		return fl.list[index], true
	} else {
		return nil, false
	}
}

func (fl *ListSource) Set(uncasted string, val any) bool {
	if uncasted == "+" {
		fl.list = append(fl.list, val)
		return true
	}
	if index, ok := fl.parseIndex(uncasted); ok {
		fl.list[index] = val
		return true
	} else {
		return false
	}
}
