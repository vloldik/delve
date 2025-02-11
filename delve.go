package delve

import (
	"fmt"

	"github.com/vloldik/delve/v2/internal/quals"
	"github.com/vloldik/delve/v2/internal/sources"
	"github.com/vloldik/delve/v2/internal/value"
	"github.com/vloldik/delve/v2/pkg/interfaces"
)

type Navigator = *navigator

type sourceType interface{ map[string]any | []any }

func New[T sourceType](source T) Navigator {
	return &navigator{source: sources.GetSource(source)}
}

func From(source interfaces.ISource) Navigator {
	return &navigator{source: source}
}

type navigator struct {
	source interfaces.ISource
}

func (fm *navigator) Source() interfaces.ISource {
	return fm.source
}

func (fm *navigator) QualGet(qual interfaces.IQual) (any, bool) {
	return fm.qualGet(qual)
}

func (fm *navigator) QualSet(qual interfaces.IQual, value any) bool {
	return fm.qualSet(qual, value)
}

func (fm *navigator) Get(qual string, _delimiter ...rune) *value.Value {
	return fm.QGet(quals.Q(qual, _delimiter...))
}

func (fm *navigator) QGet(qual interfaces.IQual) *value.Value {
	v, _ := fm.qualGet(qual)
	return value.New(v)
}

func (fm *navigator) QGetNavigator(qual interfaces.IQual) Navigator {
	v, ok := fm.qualGet(qual)
	if !ok {
		return nil
	}
	if source := sources.GetSource(v); source != nil {
		return &navigator{source: source}
	} else {
		return nil
	}
}

func (fm *navigator) GetNavigator(qual string, _delimiter ...rune) Navigator {
	return fm.QGetNavigator(quals.Q(qual, _delimiter...))
}

// Returns value by qual or panics
func (fm *navigator) MustGetByQual(qual interfaces.IQual) any {
	if val, ok := fm.QualGet(qual); ok {
		return val
	}
	panic(fmt.Sprintf("could not get by qual %v", qual))
}

func (fm *navigator) SetMapSource(source map[string]any) {
	fm.SetSource(sources.MapSource(source))
}

func (fm *navigator) SetListSource(source []any) {
	fm.SetSource(sources.NewList(source))
}

func (fm *navigator) SetSource(source interfaces.ISource) {
	fm.source = source
}
