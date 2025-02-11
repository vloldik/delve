package delve_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vloldik/delve/v2"
	"github.com/vloldik/delve/v2/internal/quals"
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
	fm := delve.New(map[string]any{"test": map[string]any{"test": 123}})
	qual := quals.CQ("test.test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fm.QGet(qual).Int()
	}
}

func BenchmarkFlexStringLen(b *testing.B) {
	baseStr := "1"
	for n := 0; n < 10; n++ {
		str := baseStr + strings.Repeat("1", n*11+1) // Ensure unique strings
		fm := delve.New(map[string]any{str: map[string]any{"test": 123}})
		name := str + ".test"
		qual := quals.CQ(name)
		strQ := quals.Q(name)

		b.Run(fmt.Sprintf("CompiledQFlexStrLen-%d", len(str)), func(b *testing.B) { // Name benchmarks by string length
			for i := 0; i < b.N; i++ {
				_ = fm.QGet(qual).Float64()
			}
		})
		b.Run(fmt.Sprintf("StringQFlexStrLen-%d", len(str)), func(b *testing.B) { // Name benchmarks by string length
			for i := 0; i < b.N; i++ {
				_ = fm.QGet(strQ).Float64()
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
		fm := delve.New(nestedMap)

		accessString := ""
		for i := depth - 1; i >= 1; i-- {
			accessString += "level" + fmt.Sprintf("%d", i) + "."
		}
		accessString += "test"
		qual := quals.CQ(accessString)
		sQual := quals.Q(accessString)

		b.Run(fmt.Sprintf("CompiledFlexStrDepth-%d", depth), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = fm.QGet(qual).Float64()
			}
		})
		b.Run(fmt.Sprintf("StringFlexStrDepth-%d", depth), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = fm.QGet(sQual).Float64()
			}
		})
	}
}

// 48 ns/op
func BenchmarkGetTyped(b *testing.B) {
	fm := delve.New(map[string]any{"lebel1": map[string]any{"test1": map[string]any{"inner": []any{0}, "string": "string"}}})
	qual := quals.CQ("lebel1.test1.inner")
	qualForString := quals.CQ("lebel1.test1.string")
	zeroQ := quals.CQ("0")
	b.ResetTimer()
	b.Run("Get an inner value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fm.QGetNavigator(qual).QGet(zeroQ).Int()
		}
	})
	b.Run("Get an array (unsafe)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fm.QGet(qual).Interface().([]any)[0].(int)
		}
	})
	b.Run("Get an inner array (safe)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fm.QGet(qual).SafeInterface([]any{}).([]any)[0].(int)
		}
	})
	b.Run("Get len of a string with len function", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fm.QGet(qualForString).Len()
		}
	})
	b.Run("Get len of a string by get it directly", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = len(fm.QGet(qualForString).String())
		}
	})
}

func BenchmarkQualCreationDifferentStringLen(b *testing.B) {
	baseStr := "1"
	for n := 0; n < 10; n++ {
		str := baseStr + strings.Repeat("1", n*11+1) // Ensure unique strings
		name := str + ".test"

		b.Run(fmt.Sprintf("FlexStrLen-%d", len(str)), func(b *testing.B) { // Name benchmarks by string length
			for i := 0; i < b.N; i++ {
				_ = quals.CQ(name)
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
				_ = quals.CQ(accessString)
			}
		})
	}
}
