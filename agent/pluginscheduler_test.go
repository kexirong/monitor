package main

import (
	"testing"

	"github.com/kexirong/monitor/common/queue"
)

func Test_scriptPluginScheduler(t *testing.T) {
	btq := queue.NewBtsQueue(4096)

	scriptPluginScheduler(btq)
}
