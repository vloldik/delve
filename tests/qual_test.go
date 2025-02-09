package delve_test

import (
	"slices"
	"testing"

	"github.com/vloldik/delve/v2"
)

func TestQualCompile(t *testing.T) {
	testString := `a\.b.c\.c\\.d`
	qual := delve.CQ(testString)
	if qual.String() != testString {
		t.Fatalf("String %s is not equals %s", qual.String(), testString)
	}
}

func IQualTest(t *testing.T, qual delve.IQual, expectedParts []string) {
	realParts := []string{}
	lastPart := ""
	var hasNext bool = true
	var part string
	var i int

	for i = 0; hasNext; i++ {
		part, hasNext = qual.Next()
		if !hasNext {
			lastPart = part
		}
		realParts = append(realParts, part)

		// Added check: Ensure number of iterations doesn't exceed the expected number of parts.
		if i > len(expectedParts) {
			t.Fatalf("Too many iterations of Next(), expected at most %d, got %d", len(expectedParts), i+1)
		}
	}

	// Check if Next() returned false too early.
	if len(realParts) < len(expectedParts) {
		t.Fatalf("Too few iterations of Next(), expected %d, got %d", len(expectedParts), len(realParts))
	}
	if lastPart != realParts[len(realParts)-1] {
		t.Fatalf("Last element is incorrent: %s in %#v", lastPart, realParts)
	}
	if !slices.Equal(realParts, expectedParts) {
		t.Fatalf("Real parts %#v not equal expected %#v", realParts, expectedParts)
	}

	// Add check:  Calling Next() again after the end should return an empty string and hasNext = false.
	part, hasNext = qual.Next()
	if part != "" {
		t.Fatalf("After iteration is complete, part should be \"\", but equals: %#v", part)
	}
	if hasNext {
		t.Fatalf("After iteration is complete, hasNext should be false, but equals: %#v", hasNext)
	}
}

func TestStringQual(t *testing.T) {
	qual := `a.b.c\.d.e\\.f\..`
	expected := []string{"a", "b", "c.d", "e\\", "f."}
	IQualTest(t, delve.Q(qual), expected)
}

func TestCompiledQual(t *testing.T) {
	qual := `a.b.c\.d.e\\.f\..`
	expected := []string{"a", "b", "c.d", "e\\", "f."}
	IQualTest(t, delve.CQ(qual), expected)
}

func TestStringQualEmpty(t *testing.T) {
	qual := ""
	expected := []string{""}
	IQualTest(t, delve.Q(qual), expected)
}

func TestCompiledQualEmpty(t *testing.T) {
	qual := ""
	expected := []string{""}
	IQualTest(t, delve.CQ(qual), expected)
}

func TestStringQualOnlyDelimiters(t *testing.T) {
	qual := "..."
	expected := []string{"", "", ""}
	IQualTest(t, delve.Q(qual), expected)
}

func TestCompiledQualOnlyDelimiters(t *testing.T) {
	qual := "..."
	expected := []string{"", "", ""}
	IQualTest(t, delve.CQ(qual), expected)
}

func TestStringQualEscapedDelimiterAtEnd(t *testing.T) {
	qual := "a\\."
	expected := []string{"a."}
	IQualTest(t, delve.Q(qual), expected)
}

func TestCompiledQualEscapedDelimiterAtEnd(t *testing.T) {
	qual := "a\\."
	expected := []string{"a."}
	IQualTest(t, delve.CQ(qual), expected)
}

func TestStringQualEscapedBackslash(t *testing.T) {
	qual := "a\\\\.b"
	expected := []string{"a\\", "b"}
	IQualTest(t, delve.Q(qual), expected)
}

func TestCompiledQualEscapedBackslash(t *testing.T) {
	qual := "a\\\\.b"
	expected := []string{"a\\", "b"}
	IQualTest(t, delve.CQ(qual), expected)
}

func TestStringQualMultipleEscapedDelimiters(t *testing.T) {
	qual := `a\.b\.c`
	expected := []string{"a.b.c"}
	IQualTest(t, delve.Q(qual), expected)
}

func TestCompiledQualMultipleEscapedDelimiters(t *testing.T) {
	qual := `a\.b\.c`
	expected := []string{"a.b.c"}
	IQualTest(t, delve.CQ(qual), expected)
}

func TestStringQualMixedEscapedAndUnescaped(t *testing.T) {
	qual := `a.b\.c.d\\.e`
	expected := []string{"a", "b.c", "d\\", "e"}
	IQualTest(t, delve.Q(qual), expected)
}

func TestCompiledQualMixedEscapedAndUnescaped(t *testing.T) {
	qual := `a.b\.c.d\\.e`
	expected := []string{"a", "b.c", "d\\", "e"}
	IQualTest(t, delve.CQ(qual), expected)
}

func TestStringQualLeadingAndTrailingDelimiters(t *testing.T) {
	qual := `.a.b.`
	expected := []string{"", "a", "b"}
	IQualTest(t, delve.Q(qual), expected)
}

func TestCompiledQualLeadingAndTrailingDelimiters(t *testing.T) {
	qual := `.a.b.`
	expected := []string{"", "a", "b"}
	IQualTest(t, delve.CQ(qual), expected)
}

func TestStringQualConsecutiveEscapedBackslashes(t *testing.T) {
	qual := "a\\\\\\\\.b"
	expected := []string{"a\\\\", "b"}
	IQualTest(t, delve.Q(qual), expected)
}

func TestCompiledQualConsecutiveEscapedBackslashes(t *testing.T) {
	qual := "a\\\\\\\\.b"
	expected := []string{"a\\\\", "b"}
	IQualTest(t, delve.CQ(qual), expected)
}
