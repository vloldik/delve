package sources

type MapSource map[string]any

// Get retrieves a value from delve by key
func (fm MapSource) Get(key string) (any, bool) {
	value, ok := fm[key]
	return value, ok
}

func (fm MapSource) Set(key string, val any) bool {
	fm[key] = val
	return true
}
