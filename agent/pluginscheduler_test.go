package main

import (
	"testing"

	"github.com/kexirong/monitor/common/queue"
)

func Test_pypluginscheduler(t *testing.T) {
	btq := queue.NewBtsQueue(4096)

	pyPluginScheduler(btq)
}
