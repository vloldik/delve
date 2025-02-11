package sources

import "github.com/vloldik/delve/v3/pkg/idelve"

func GetSource(unknown any) idelve.ISource {
	switch typed := unknown.(type) {
	case []any:
		return NewList(typed)
	case map[string]any:
		return MapSource(typed)
	case idelve.ISource:
		return typed
	default:
		return nil
	}
}
