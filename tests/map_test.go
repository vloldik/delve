package delve_test

import (
	"encoding/json"
	"testing"

	"github.com/vloldik/delve"
)

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
	nav := delve.FromMap(mMap)
	if nav.Int(delve.CQ("a.b.0.c")) != 3 {
		m.FailNow()
	}
	if !nav.Bool(delve.Q("a.b.-1.last")) {
		m.Fatal("Last index get fail")
	}
}

func TestScreening(m *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	nav := delve.FromMap(mMap)

	if a := nav.Int(delve.CQ("b.c.a\\.b")); a != 321 {
		m.Fatalf("%d is not eq 321", a)
	}
}

func TestCustomDelemiter(m *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	nav := delve.FromMap(mMap)
	if value, ok := nav.QualGet(delve.CQ("a/b/0/c", '/')); ok {
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
	nav := delve.FromMap(mMap)
	inner := nav.Navigator(delve.CQ("a.b"))
	if inner == nil {
		t.Fatal("Inner map is nil")
		return
	}
	if inner.Int16(delve.CQ("0.f")) != 1 {
		t.Fatal("Inner int is not equal 1")
	}
}

func TestTypeGets(t *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	nav := delve.FromMap(mMap)

	if nav.String(delve.CQ("a.b.0.e")) != "hello" {
		t.Errorf("String not equal")
	}

	if nav.Bool(delve.CQ("a.b.0.d")) != true {
		t.Errorf("Bool not equal")
	}

	byteSlice := []byte(nav.String(delve.CQ("b.c.bytes")))
	if string(byteSlice) != "AQID" {
		t.Errorf("ByteSlice not equal")
	}

	if nav.Float64(delve.CQ("a.b.0.c")) != 3.14 {
		t.Errorf("Float64 not equal")
	}

	if nav.Float32(delve.CQ("a.b.0.g")) != 1.1 {
		t.Errorf("Float32 not equal")
	}

	if nav.Int(delve.CQ("a.b.0.f")) != 1 {
		t.Errorf("Int not equal")
	}

	if nav.Int64(delve.CQ("a.b.0.f")) != 1 {
		t.Errorf("Int64 not equal")
	}

	if nav.Int32(delve.CQ("a.b.0.f")) != 1 {
		t.Errorf("Int32 not equal")
	}

	if nav.Int16(delve.CQ("a.b.0.f")) != 1 {
		t.Errorf("Int16 not equal")
	}

	if nav.Int8(delve.CQ("a.b.0.f")) != 1 {
		t.Errorf("Int8 not equal")
	}

	if nav.Uint(delve.CQ("a.b.0.f")) != 1 {
		t.Errorf("Uint not equal")
	}

	if nav.Uint64(delve.CQ("a.b.0.f")) != 1 {
		t.Errorf("Uint64 not equal")
	}

	if nav.Uint32(delve.CQ("a.b.0.f")) != 1 {
		t.Errorf("Uint32 not equal")
	}

	if nav.Uint16(delve.CQ("a.b.0.f")) != 1 {
		t.Errorf("Uint16 not equal")
	}

	if nav.Uint8(delve.CQ("a.b.0.f")) != 1 {
		t.Errorf("Uint8 not equal")
	}

	if nav.Navigator(delve.CQ("a.b.0.j")) == nil {
		t.Errorf("Navigator not equal")
	}
}

func TestTypeDefaults(t *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	nav := delve.FromMap(mMap)

	if nav.String(delve.CQ("notexist"), "default") != "default" {
		t.Errorf("String default not equal")
	}

	if nav.Bool(delve.CQ("notexist"), true) != true {
		t.Errorf("Bool default not equal")
	}

	if string(nav.ByteSlice(delve.CQ("notexist"), []byte("default"))) != "default" {
		t.Errorf("ByteSlice default not equal")
	}

	if nav.Float64(delve.CQ("notexist"), 3.14) != 3.14 {
		t.Errorf("Float64 default not equal")
	}

	if nav.Float32(delve.CQ("notexist"), 1.1) != 1.1 {
		t.Errorf("Float32 default not equal")
	}

	if nav.Int(delve.CQ("notexist"), 1) != 1 {
		t.Errorf("Int default not equal")
	}

	if nav.Int64(delve.CQ("notexist"), 1) != 1 {
		t.Errorf("Int64 default not equal")
	}

	if nav.Int32(delve.CQ("notexist"), 1) != 1 {
		t.Errorf("Int32 default not equal")
	}

	if nav.Int16(delve.CQ("notexist"), 1) != 1 {
		t.Errorf("Int16 default not equal")
	}

	if nav.Int8(delve.CQ("notexist"), 1) != 1 {
		t.Errorf("Int8 default not equal")
	}

	if nav.Uint(delve.CQ("notexist"), 1) != 1 {
		t.Errorf("Uint default not equal")
	}

	if nav.Uint64(delve.CQ("notexist"), 1) != 1 {
		t.Errorf("Uint64 default not equal")
	}

	if nav.Uint32(delve.CQ("notexist"), 1) != 1 {
		t.Errorf("Uint32 default not equal")
	}

	if nav.Uint16(delve.CQ("notexist"), 1) != 1 {
		t.Errorf("Uint16 default not equal")
	}

	if nav.Uint8(delve.CQ("notexist"), 1) != 1 {
		t.Errorf("Uint8 default not equal")
	}
}
