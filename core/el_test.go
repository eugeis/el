package core

import (
	"testing"
)

func TestSearch(t *testing.T) {
	config := NewConfig("logs")
	config.Folder = "D:\\TC_CACHE\\logs\\export\\"
	config.Elastics[0].Queries = []*Query{
		&Query{
			Name: "dur_gte_50000",
			Source: `{
	"sort" : [
        { "dur" : {"order" : "desc"}}
    ],
  "query": {
    "range": {
      "dur": {
        "gte": 50000
      }
    }
  }
}`}}
	el := El{config: config}
	defer el.Close()
	err := el.Export()
	if err != nil {
		println(err.Error())
	}

}
