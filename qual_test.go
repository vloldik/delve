package flexmap_test

import (
	"slices"
	"testing"

	"github.com/vloldik/flexmap"
)

func TestQualCompile(t *testing.T) {
	testString := `a\.b.c\.c\\.d`
	qual := flexmap.CompileQual(testString)
	compiled := flexmap.CompiledQual{"a.b", "c.c\\", "d"}

	if !slices.Equal(qual, compiled) {
		t.Fatalf("Slice %v is not equal %v", qual, compiled)
	}
}
