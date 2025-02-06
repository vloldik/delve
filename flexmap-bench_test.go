package flexmap_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vloldik/flexmap"
)

// 4.2 ns/op
func BenchmarkDirect(b *testing.B) {
	m := map[string]map[string]int{"testdogpjsdiogndsiogfnsdiogsngiodg": {"test": 123}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m["testdogpjsdiogndsiogfnsdiogsngiodg"]["test"] // Прямой доступ
	}
}

// 45 ns/op
func BenchmarkFlexMap(b *testing.B) {
	fm := flexmap.FromMap(map[string]any{"test": map[string]any{"test": 123}})
	qual := flexmap.CompileQual("test.test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fm.Int(qual)
	}
}

func BenchmarkFlexStringLen(b *testing.B) {
	baseStr := "1"
	for n := 0; n < 10; n++ {
		str := baseStr + strings.Repeat("1", n*11+1) // Ensure unique strings
		fm := flexmap.FromMap(map[string]any{str: map[string]any{"test": 123}})
		name := str + ".test"
		qual := flexmap.CompileQual(name)

		b.Run(fmt.Sprintf("FlexStrLen-%d", len(str)), func(b *testing.B) { // Name benchmarks by string length
			for i := 0; i < b.N; i++ {
				_ = fm.Float64(qual)
			}
		})
		baseStr = str
	}
}

func BenchmarkFlexStringDepth(b *testing.B) {
	for depth := 1; depth <= 10; depth++ {
		nestedMap := map[string]any{"test": 123}
		for i := 1; i < depth; i++ {
			nestedMap = map[string]any{"level" + fmt.Sprintf("%d", i): nestedMap}
		}
		fm := flexmap.FromMap(nestedMap)

		accessString := ""
		for i := depth - 1; i >= 1; i-- {
			accessString += "level" + fmt.Sprintf("%d", i) + "."
		}
		accessString += "test"
		qual := flexmap.CompileQual(accessString)

		b.Run(fmt.Sprintf("FlexStrDepth-%d", depth), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = fm.Float64(qual)
			}
		})
	}
}

// 48 ns/op
func BenchmarkGetTyped(b *testing.B) {
	fm := flexmap.FromMap(map[string]any{"lebel1": map[string]any{"test1": map[string]any{"inner": []string{"test"}}}})
	qual := flexmap.CompileQual("lebel1.test1.inner.test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fm.StringSlice(qual)
	}
}
