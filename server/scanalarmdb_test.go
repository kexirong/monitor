package main

import (
	"fmt"
	"testing"
)

func Test_scanalarmdb(t *testing.T) {
	avs := Scanalarmdb()
	for _, av := range avs {
		fmt.Println(av.String())
	}
	//t.Log(err)
}
