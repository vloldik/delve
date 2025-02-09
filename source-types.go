package delve

import "strconv"

type getterMap map[string]any
type getterList struct {
	list []any
}

// Get retrieves a value from delve by key
func (fm getterMap) Get(key string) (any, bool) {
	value, ok := fm[key]
	return value, ok
}

func (fm getterMap) Set(key string, val any) bool {
	fm[key] = val
	return true
}

func (fl *getterList) parseIndex(stringIndex string) (int, bool) {
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
func (fl *getterList) Get(uncasted string) (any, bool) {
	if index, ok := fl.parseIndex(uncasted); ok {
		return fl.list[index], true
	} else {
		return nil, false
	}
}

func (fl *getterList) Set(uncasted string, val any) bool {
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
