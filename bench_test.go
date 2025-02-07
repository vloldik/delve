package delve_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vloldik/delve"
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
func BenchmarkDelve(b *testing.B) {
	fm := delve.FromMap(map[string]any{"test": map[string]any{"test": 123}})
	qual := delve.CQ("test.test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fm.Int(qual)
	}
}

func BenchmarkFlexStringLen(b *testing.B) {
	baseStr := "1"
	for n := 0; n < 10; n++ {
		str := baseStr + strings.Repeat("1", n*11+1) // Ensure unique strings
		fm := delve.FromMap(map[string]any{str: map[string]any{"test": 123}})
		name := str + ".test"
		qual := delve.CQ(name)
		strQ := delve.Q(name)

		b.Run(fmt.Sprintf("CompiledQFlexStrLen-%d", len(str)), func(b *testing.B) { // Name benchmarks by string length
			for i := 0; i < b.N; i++ {
				_ = fm.Float64(qual)
			}
		})
		b.Run(fmt.Sprintf("StringQFlexStrLen-%d", len(str)), func(b *testing.B) { // Name benchmarks by string length
			for i := 0; i < b.N; i++ {
				_ = fm.Float64(strQ)
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
		fm := delve.FromMap(nestedMap)

		accessString := ""
		for i := depth - 1; i >= 1; i-- {
			accessString += "level" + fmt.Sprintf("%d", i) + "."
		}
		accessString += "test"
		qual := delve.CQ(accessString)
		sQual := delve.Q(accessString)

		b.Run(fmt.Sprintf("CompiledFlexStrDepth-%d", depth), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = fm.Float64(qual)
			}
		})
		b.Run(fmt.Sprintf("StringFlexStrDepth-%d", depth), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = fm.Float64(sQual)
			}
		})
	}
}

// 48 ns/op
func BenchmarkGetTyped(b *testing.B) {
	fm := delve.FromMap(map[string]any{"lebel1": map[string]any{"test1": map[string]any{"inner": []int{0}}}})
	qual := delve.CQ("lebel1.test1.inner")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fm.Navigator(qual).Int(delve.CQ("0"))
	}
}

func BenchmarkQualCreationDifferentStringLen(b *testing.B) {
	baseStr := "1"
	for n := 0; n < 10; n++ {
		str := baseStr + strings.Repeat("1", n*11+1) // Ensure unique strings
		name := str + ".test"

		b.Run(fmt.Sprintf("FlexStrLen-%d", len(str)), func(b *testing.B) { // Name benchmarks by string length
			for i := 0; i < b.N; i++ {
				_ = delve.CQ(name)
			}
		})
		baseStr = str
	}
}

func BenchmarkQualCreationStringDepth(b *testing.B) {
	for depth := 1; depth <= 10; depth++ {
		nestedMap := map[string]any{"test": 123}
		for i := 1; i < depth; i++ {
			nestedMap = map[string]any{"level" + fmt.Sprintf("%d", i): nestedMap}
		}

		accessString := ""
		for i := depth - 1; i >= 1; i-- {
			accessString += "level" + fmt.Sprintf("%d", i) + "."
		}
		accessString += "test"

		b.Run(fmt.Sprintf("FlexStrDepth-%d", depth), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = delve.CQ(accessString)
			}
		})
	}
}
