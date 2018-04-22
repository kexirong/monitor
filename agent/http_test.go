package main

import (
	"log"
	"net/http"
	"testing"
)

func Test_pyplugin(t *testing.T) {
	t.Log("")
	httpInit()
	log.Fatal(http.ListenAndServe(":5101", nil))

}
