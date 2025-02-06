package delve_test

import (
	"slices"
	"testing"

	"github.com/vloldik/delve"
)

func TestQualCompile(t *testing.T) {
	testString := `a\.b.c\.c\\.d`
	qual := delve.Qual(testString)
	compiled := delve.CompiledQual{"a.b", "c.c\\", "d"}

	if !slices.Equal(qual, compiled) {
		t.Fatalf("Slice %v is not equal %v", qual, compiled)
	}
}
