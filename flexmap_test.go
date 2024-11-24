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
	if value, ok := flexMap.GetByQual("a.b.0.c"); ok {
		if value.(float64) != 3 {
			m.FailNow()
		}
	} else {
		m.FailNow()
	}
}
