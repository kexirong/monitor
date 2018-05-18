package main

import (
	"net/http"
	"testing"
)

func Test_httplisent(t *testing.T) {
	t.Fatal(http.ListenAndServe(":5001", nil))
}
