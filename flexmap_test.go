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
			"f": 123
		}
	}
}`

func TestUsage(m *testing.T) {
	flexMap := make(flexmap.FlexMap)
	err := json.Unmarshal([]byte(jsonTestStruct), &flexMap)
	if err != nil {
		panic(err)
	}
	if flexMap.Int("a.b.0.c") != 3 {
		m.FailNow()
	}
}

func TestCustomDelemiter(m *testing.T) {
	flexMap := make(flexmap.FlexMap)
	err := json.Unmarshal([]byte(jsonTestStruct), &flexMap)
	if err != nil {
		panic(err)
	}
	flexmap.QDelemiter = "/"
	if value, ok := flexMap.GetByQual("a/b/0/c"); ok {
		if value.(float64) != 3 {
			m.FailNow()
		}
	} else {
		m.FailNow()
	}
}
