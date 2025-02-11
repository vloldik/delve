package delve_test

import (
	"testing"

	"github.com/vloldik/delve/v2"
	"github.com/vloldik/delve/v2/internal/quals"
)

func TestSetFunction(t *testing.T) {
	t.Run("Set value in new nested map", func(t *testing.T) {
		m := make(map[string]any)
		nav := delve.New(m)
		ok := nav.QualSet(quals.Q("a.b.c"), 10)
		if !ok {
			t.Fatal("QualSet failed")
		}

		if val := nav.Get("a.b.c").Int(); val != 10 {
			t.Errorf("Expected 10, got %v", val)
		}
	})

	t.Run("Overwrite existing value in map", func(t *testing.T) {
		m := map[string]any{"x": 5}
		nav := delve.New(m)
		ok := nav.QualSet(quals.CQ("x"), 10)
		if !ok {
			t.Fatal("QualSet failed")
		}
		if m["x"] != 10 {
			t.Errorf("Expected 10, got %v", m["x"])
		}
	})

	t.Run("Set in existing list", func(t *testing.T) {
		list := []any{0, 1, 2}
		nav := delve.New(list)
		ok := nav.QualSet(quals.Q("1"), 99)
		if !ok {
			t.Fatal("QualSet failed")
		}
		if val := nav.Get("1").Int(); val != 99 {
			t.Errorf("Expected 99, got %v", val)
		}
	})

	t.Run("Set in list with '+' index", func(t *testing.T) {
		list := []any{"a"}
		nav := delve.New(list)
		ok := nav.QualSet(quals.CQ("+"), "b")
		if !ok {
			t.Fatal("QualSet failed")
		}
		if val, ok := nav.QualGet(quals.Q("1")); !ok {
			t.Errorf("Expected [a b], got %v", nav)
		} else if val.(string) != "b" {
			t.Errorf("Expected b, got %v", val)
		}
	})

	t.Run("Nested list within map within list", func(t *testing.T) {
		nested := []any{
			map[string]any{
				"a": []any{10, 20},
			},
		}
		nav := delve.New(nested)
		ok := nav.QualSet(quals.CQ("0.a.1"), 30)
		if !ok {
			t.Fatal("QualSet failed")
		}

		elem0 := nested[0].(map[string]any)
		list := elem0["a"].([]any)
		if list[1] != 30 {
			t.Errorf("Expected 30, got %v", list[1])
		}
	})

	t.Run("Set with invalid list index", func(t *testing.T) {
		list := []any{1, 2, 3}
		nav := delve.New(list)
		ok := nav.QualSet(quals.CQ("abc"), 5)
		if ok {
			t.Error("Expected QualSet to fail with invalid index")
		}
	})

	t.Run("Set path through non-map node", func(t *testing.T) {
		m := map[string]any{"a": 5}
		nav := delve.New(m)
		nav.QualSet(quals.CQ("a.b"), 10)
		if val := nav.Get("a.b").Int(); val != 10 {
			t.Errorf("Expected value 10, got %d", val)
		}
	})

	t.Run("Set with negative list index", func(t *testing.T) {
		list := []any{1, 2, 3}
		nav := delve.New(list)
		ok := nav.QualSet(quals.CQ("-1"), 99)
		if !ok {
			t.Fatal("QualSet failed")
		}
		if list[2] != 99 {
			t.Errorf("Expected last element 99, got %v", list[2])
		}
	})

	t.Run("Set list index out of bounds", func(t *testing.T) {
		list := []any{0, 1, 2}
		nav := delve.New(list)
		ok := nav.QualSet(quals.CQ("3"), 99)
		if ok {
			t.Error("Expected failure due to out of bounds index")
		}
	})

	t.Run("Set in list element's map key", func(t *testing.T) {
		list := []any{map[string]any{}}
		nav := delve.New(list)
		ok := nav.QualSet(quals.CQ("0.key"), "value")
		if !ok {
			t.Fatal("QualSet failed")
		}

		elem0, ok := list[0].(map[string]any)
		if !ok {
			t.Fatal("element 0 is not a map")
		}
		if elem0["key"] != "value" {
			t.Errorf("Expected 'value', got %v", elem0["key"])
		}
	})
}
