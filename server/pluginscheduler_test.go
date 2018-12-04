package main

import (
	"testing"
)

func Test_activePluginScheduler(t *testing.T) {
	go activePluginScheduler()
	select {}
}
