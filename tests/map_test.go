package delve_test

import (
	"encoding/json"
	"testing"

	"github.com/vloldik/delve/v2"
	"github.com/vloldik/delve/v2/internal/quals"
)

type mockSource struct{}

func (mockSource) Set(string, any) bool {
	return true
}

func (mockSource) Get(string) (any, bool) {
	return nil, true
}

const jsonTestStruct = `{
	"a": {
		"b": [
			{"c": 3.14, "d": true, "e": "hello", "f": 1, "g": 1.1, "h": 1111111111111111111, "i": ["a", "b"], "j": {"k": "l"}},
			{"last": true}
		]
	}, 
	"b": {
		"c": {
			"f": 123,
			"a.b": 321,
			"bytes": "AQID"
		}
	}
}`

func TestUsage(m *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	nav := delve.New(mMap)
	if nav.QGet(quals.Q("a.b.0.c")).Float64() != 3.14 {
		m.FailNow()
	}
	if !nav.QGet(quals.Q("a.b.-1.last")).Bool() {
		m.Fatal("Last index get fail")
	}
}

func TestScreening(m *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	nav := delve.New(mMap)

	if a := nav.QGet(quals.CQ("b.c.a\\.b")).Int(); a != 321 {
		m.Fatalf("%d is not eq 321", a)
	}
}

func TestCustomDelemiter(m *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	nav := delve.New(mMap)
	if value, ok := nav.QualGet(quals.CQ("a/b/0/c", '/')); ok {
		if value.(float64) != 3.14 {
			m.FailNow()
		}
	} else {
		m.FailNow()
	}
}

func TestInnerGet(t *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	nav := delve.New(mMap)
	inner := nav.QGetNavigator(quals.CQ("a.b"))
	if inner == nil {
		t.Fatal("Inner map is nil")
		return
	}
	if inner.QGet(quals.CQ("0.f")).Int16() != 1 {
		t.Fatal("Inner int is not equal 1")
	}
}

func TestTypeGets(t *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	nav := delve.New(mMap)

	if nav.QGet(quals.CQ("a.b.0.e")).String() != "hello" {
		t.Errorf("String not equal")
	}

	if nav.QGet(quals.CQ("a.b.0.d")).Bool() != true {
		t.Errorf("Bool not equal")
	}

	if nav.QGet(quals.CQ("a.b.0.c")).Float64() != 3.14 {
		t.Errorf("Float64 not equal")
	}

	if nav.QGet(quals.Q("a.b.0.h")).Int8(-1) != -1 {
		t.Errorf("Overflow should be handled")
	}

	valfloat32 := nav.QGet(quals.CQ("a.b.0.g")).Float32(-1)
	valfloat64 := nav.QGet(quals.CQ("a.b.0.g")).Float64(0)

	if valfloat64 == 0 || (valfloat32 != -1 && float64(valfloat32) != valfloat64) {
		t.Errorf("Float64 should not be convertible to float32")
	}

	if nav.QGet(quals.CQ("a.b.0.f")).Int() != 1 {
		t.Errorf("Int not equal")
	}

	if nav.QGet(quals.CQ("a.b.0.f")).Int64() != 1 {
		t.Errorf("Int64 not equal")
	}

	if nav.QGet(quals.CQ("a.b.0.f")).Int32() != 1 {
		t.Errorf("Int32 not equal")
	}

	if nav.QGet(quals.CQ("a.b.0.f")).Int16() != 1 {
		t.Errorf("Int16 not equal")
	}

	if nav.QGet(quals.CQ("a.b.0.f")).Int8() != 1 {
		t.Errorf("Int8 not equal")
	}

	if nav.QGet(quals.CQ("a.b.0.f")).Uint() != 1 {
		t.Errorf("Uint not equal")
	}

	if nav.QGet(quals.CQ("a.b.0.f")).Uint64() != 1 {
		t.Errorf("Uint64 not equal")
	}

	if nav.QGet(quals.CQ("a.b.0.f")).Uint32() != 1 {
		t.Errorf("Uint32 not equal")
	}

	if nav.QGet(quals.CQ("a.b.0.f")).Uint16() != 1 {
		t.Errorf("Uint16 not equal")
	}

	if nav.QGet(quals.CQ("a.b.0.f")).Uint8() != 1 {
		t.Errorf("Uint8 not equal")
	}

	if nav.QGetNavigator(quals.CQ("a.b.0.j")) == nil {
		t.Errorf("Navigator not equal")
	}

	if len := nav.QGet(quals.Q("a.b")).Len(); len != 2 {
		t.Errorf("a.b. len is two, but got %d", len)
	}

	if nav.QGet(quals.Q("a.b.0.f")).Interface().(float64) != 1 {
		t.Errorf("a.b.0.f should be 1.")
	}

	if nav.QGet(quals.Q("a.b.0.g")).SafeInterface(float64(1)).(float64) != 1.1 {
		t.Error("a.b.0.g should be 1.1")
	}

	if nav.QGet(quals.Q("a.b.0.g")).SafeInterface(any(4)) != any(4) {
		t.Error("safe interface should allow usinh any")
	}

	testMap := map[string]any{"test": mockSource{}}

	var defaultVal any
	if delve.New(testMap).QGet(quals.Q("test")).SafeInterface(defaultVal) == defaultVal {
		t.Error("should assign mocksource to mocksource")
	}
}

func TestTypeDefaults(t *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	nav := delve.New(mMap)

	if nav.QGet(quals.CQ("notexist")).String("default") != "default" {
		t.Errorf("String default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Bool(true) != true {
		t.Errorf("Bool default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Float64(3.14) != 3.14 {
		t.Errorf("Float64 default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Float32(1.1) != 1.1 {
		t.Errorf("Float32 default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Int(1) != 1 {
		t.Errorf("Int default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Int64(1) != 1 {
		t.Errorf("Int64 default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Int32(1) != 1 {
		t.Errorf("Int32 default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Int16(1) != 1 {
		t.Errorf("Int16 default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Int8(1) != 1 {
		t.Errorf("Int8 default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Uint(1) != 1 {
		t.Errorf("Uint default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Uint64(1) != 1 {
		t.Errorf("Uint64 default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Uint32(1) != 1 {
		t.Errorf("Uint32 default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Uint16(1) != 1 {
		t.Errorf("Uint16 default not equal")
	}

	if nav.QGet(quals.CQ("notexist")).Uint8(1) != 1 {
		t.Errorf("Uint8 default not equal")
	}

	if nav.QGet(quals.Q("a.b.c")).Len() != -1 {
		t.Errorf("Len of non countable types should be -1")
	}

	if nav.QGet(quals.Q("notexists")).Interface() != nil {
		t.Errorf("interface default not equal")
	}

	if nav.QGet(quals.Q("a.b")).SafeInterface(1).(int) != 1 {
		t.Errorf("safe interface default not equal")
	}
}
