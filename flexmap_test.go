package flexmap_test

import (
	"encoding/json"
	"testing"

	"github.com/vloldik/flexmap"
)

const jsonTestStruct = `{
	"a": {
		"b": [
			{"c": 3}
		]
	}, 
	"b": {
		"c": {
			"f": 123,
			"a.b": 321
		}
	}
}`

func TestUsage(m *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	flexMap := flexmap.FromMap(mMap)
	if flexMap.Int("a.b.0.c") != 3 {
		m.FailNow()
	}
}

func TestScreening(m *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	flexMap := flexmap.FromMap(mMap)
	if a := flexMap.Int("b.c.a\\.b"); a != 321 {
		m.Fatalf("%d is not eq 321", a)
	}
}

func TestCustomDelemiter(m *testing.T) {
	mMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonTestStruct), &mMap)
	if err != nil {
		panic(err)
	}
	flexMap := flexmap.FromMap(mMap, '/')
	if value, ok := flexMap.GetByQual("a/b/0/c"); ok {
		if value.(float64) != 3 {
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
	flexMap := flexmap.FromMap(mMap)
	inner := flexMap.FlexMap("a.b")
	if inner == nil {
		t.Fatal("Inner map is nil")
		return
	}
	if inner.Int16("0.c") != 3 {
		t.Fatal("Inner int is not equal 3")
	}
}
