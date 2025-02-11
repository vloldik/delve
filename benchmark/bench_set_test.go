package delve_test

import (
	"testing"

	"github.com/vloldik/delve/v3"
	"github.com/vloldik/delve/v3/internal/quals"
)

func BenchmarkSetValueInMap(b *testing.B) {
	m := make(map[string]any)
	nav := delve.New(m)
	q := quals.CQ("key")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nav.QSet(q, i)
	}
}

func BenchmarkOverwriteValueInMap(b *testing.B) {
	m := map[string]any{"key": 0}
	nav := delve.New(m)
	q := quals.CQ("key")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nav.QSet(q, i)
	}
}

func BenchmarkSetValueInList(b *testing.B) {
	list := make([]any, 100)
	for i := range list {
		list[i] = 0
	}
	nav := delve.New(list)
	q := quals.CQ("50")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nav.QSet(q, i)
	}
}

func BenchmarkSetNestedValueInMap(b *testing.B) {
	m := make(map[string]any)
	nav := delve.New(m)
	q := quals.CQ("a.b.c")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nav.QSet(q, i)
		nav.SetMapSource(map[string]any{})
	}
}

func BenchmarkAppendToList(b *testing.B) {
	list := []any{}
	nav := delve.New(list)
	q := quals.CQ("+")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nav.QSet(q, i)
	}
}

func BenchmarkSetNegativeIndexInList(b *testing.B) {
	list := make([]any, 100)
	for i := range list {
		list[i] = 0
	}
	nav := delve.New(list)
	q := quals.CQ("-1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nav.QSet(q, i)
	}
}
