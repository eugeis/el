package core

import (
	"testing"
)

func TestSearch(t *testing.T) {
	el := El{}
	defer el.Close()
	err := el.Export()
	if err != nil {
		println(err.Error())
	}

}
