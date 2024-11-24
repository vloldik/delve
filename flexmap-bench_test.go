package flexmap_test

import (
	"testing"

	"github.com/vloldik/flexmap"
)

// 5 ns/op
func BenchmarkDirect(b *testing.B) {
	m := map[string]map[string]int{"test": {"test": 123}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m["test"]["test"] // Прямой доступ
	}
}

// 55 ns/op
func BenchmarkFlexMap(b *testing.B) {
	fm := flexmap.FlexMap{"test": map[string]any{"test": 123}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fm.Int("test.test")
	}
}

// 54 ns/op
func BenchmarkFlexMapCast(b *testing.B) {
	fm := flexmap.FlexMap{"test": map[string]any{"test": 123}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fm.Float64("test.test")
	}
}

// 55 ns/op
func BenchmarkGetTyped(b *testing.B) {
	fm := flexmap.FlexMap{"test": map[string]any{"test": 123}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = flexmap.MustGetTypedByQual[int]("test.test", fm)
	}
}
