package delve_test

import (
	"testing"

	"github.com/vloldik/delve/v2"
)

func TestSetFunction(t *testing.T) {
	t.Run("Set value in new nested map", func(t *testing.T) {
		m := make(map[string]any)
		nav := delve.FromMap(m)
		ok := nav.QualSet(delve.Q("a.b.c"), 10)
		if !ok {
			t.Fatal("QualSet failed")
		}

		if val := nav.Int(delve.Q("a.b.c")); val != 10 {
			t.Errorf("Expected 10, got %v", val)
		}
	})

	t.Run("Overwrite existing value in map", func(t *testing.T) {
		m := map[string]any{"x": 5}
		nav := delve.FromMap(m)
		ok := nav.QualSet(delve.CQ("x"), 10)
		if !ok {
			t.Fatal("QualSet failed")
		}
		if m["x"] != 10 {
			t.Errorf("Expected 10, got %v", m["x"])
		}
	})

	t.Run("Set in existing list", func(t *testing.T) {
		list := []any{0, 1, 2}
		nav := delve.FromList(list)
		ok := nav.QualSet(delve.Q("1"), 99)
		if !ok {
			t.Fatal("QualSet failed")
		}
		if val := nav.Int(delve.Q("1")); val != 99 {
			t.Errorf("Expected 99, got %v", val)
		}
	})

	t.Run("Set in list with '+' index", func(t *testing.T) {
		list := []any{"a"}
		nav := delve.FromList(list)
		ok := nav.QualSet(delve.CQ("+"), "b")
		if !ok {
			t.Fatal("QualSet failed")
		}
		if val, ok := nav.QualGet(delve.Q("1")); !ok {
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
		nav := delve.FromList(nested)
		ok := nav.QualSet(delve.CQ("0.a.1"), 30)
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
		nav := delve.FromList(list)
		ok := nav.QualSet(delve.CQ("abc"), 5)
		if ok {
			t.Error("Expected QualSet to fail with invalid index")
		}
	})

	t.Run("Set path through non-map node", func(t *testing.T) {
		m := map[string]any{"a": 5}
		nav := delve.FromMap(m)
		nav.QualSet(delve.CQ("a.b"), 10)
		if val := nav.Int(delve.Q("a.b")); val != 10 {
			t.Errorf("Expected value 10, got %d", val)
		}
	})

	t.Run("Set with negative list index", func(t *testing.T) {
		list := []any{1, 2, 3}
		nav := delve.FromList(list)
		ok := nav.QualSet(delve.CQ("-1"), 99)
		if !ok {
			t.Fatal("QualSet failed")
		}
		if list[2] != 99 {
			t.Errorf("Expected last element 99, got %v", list[2])
		}
	})

	t.Run("Set list index out of bounds", func(t *testing.T) {
		list := []any{0, 1, 2}
		nav := delve.FromList(list)
		ok := nav.QualSet(delve.CQ("3"), 99)
		if ok {
			t.Error("Expected failure due to out of bounds index")
		}
	})

	t.Run("Set in list element's map key", func(t *testing.T) {
		list := []any{map[string]any{}}
		nav := delve.FromList(list)
		ok := nav.QualSet(delve.CQ("0.key"), "value")
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
