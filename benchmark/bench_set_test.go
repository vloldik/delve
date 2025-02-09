package delve_test

import (
	"testing"

	"github.com/vloldik/delve/v2"
)

func BenchmarkSetValueInMap(b *testing.B) {
	m := make(map[string]any)
	nav := delve.FromMap(m)
	q := delve.CQ("key")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nav.QualSet(q, i)
	}
}

func BenchmarkOverwriteValueInMap(b *testing.B) {
	m := map[string]any{"key": 0}
	nav := delve.FromMap(m)
	q := delve.CQ("key")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nav.QualSet(q, i)
	}
}

func BenchmarkSetValueInList(b *testing.B) {
	list := make([]any, 100)
	for i := range list {
		list[i] = 0
	}
	nav := delve.FromList(list)
	q := delve.CQ("50")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nav.QualSet(q, i)
	}
}

func BenchmarkSetNestedValueInMap(b *testing.B) {
	m := make(map[string]any)
	nav := delve.FromMap(m)
	q := delve.CQ("a.b.c")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nav.QualSet(q, i)
		nav.SetMapSource(map[string]any{})
	}
}

func BenchmarkAppendToList(b *testing.B) {
	list := []any{}
	nav := delve.FromList(list)
	q := delve.CQ("+")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nav.QualSet(q, i)
	}
}

func BenchmarkSetNegativeIndexInList(b *testing.B) {
	list := make([]any, 100)
	for i := range list {
		list[i] = 0
	}
	nav := delve.FromList(list)
	q := delve.CQ("-1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nav.QualSet(q, i)
	}
}
