package main

import (
	"testing"

	"github.com/kexirong/monitor/server/models"
)

func Test_activePluginScheduler(t *testing.T) {
	go activePluginScheduler()
	select {}
}

func Test_alarmEventModels(t *testing.T) {
	aes, err := models.AlarmEventsByHostNameAnchorPointRule(monitorDB, "test", "test", "test")
	if err != nil {
		t.Error(err)
	}

	for _, ae := range aes {
		t.Log(ae)
		ae.Count++
		err := ae.Save(monitorDB)
		if err != nil {
			t.Error(err)
		}
	}
}
