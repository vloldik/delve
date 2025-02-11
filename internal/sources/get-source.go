package sources

import "github.com/vloldik/delve/v2/pkg/interfaces"

func GetSource(unknown any) interfaces.ISource {
	switch typed := unknown.(type) {
	case []any:
		return NewList(typed)
	case map[string]any:
		return MapSource(typed)
	case interfaces.ISource:
		return typed
	default:
		return nil
	}
}
