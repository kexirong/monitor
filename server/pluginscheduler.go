package main

import (
	"fmt"
	"time"

	"github.com/kexirong/monitor/common/scheduler"
	"github.com/kexirong/monitor/server/activeplugin"
	"github.com/kexirong/monitor/server/models"
)

var taskScheduled = scheduler.New()

func activePluginScheduler() {
	aps, err := models.ActiveProbeAll(monitorDB)
	if err != nil {
		panic(err)
	}

	for i, ap := range aps {
		switch ap.PluginName {
		case "http_probe":
			fmt.Println("range aps:", i)
			var tasker = activeplugin.NewHTTPProbe(ap.HostName)
			apcs, err := models.ActiveProbeConfigsByActiveProbeID(monitorDB, ap.ID)
			if err != nil {
				panic(err)
			}
			fmt.Println("range aps:", i)
			for _, apc := range apcs {
				err := tasker.AddJob(apc.Target, apc.Arg1, apc.Arg2)
				if err != nil {
					Logger.Error.Println(err)
				}
			}
			fmt.Println("range aps:", i)
			taskScheduled.AddTask(time.Second*time.Duration(ap.Interval), tasker)

		case "process_probe":
		}
		fmt.Println("range aps:", i)
	}
	var callback = func(b []byte, err error) {
		if err != nil {
			fmt.Println(err)
			return
		}
		Logger.Info.Println(string(b))

	}
	fmt.Println("staring")
	taskScheduled.Star(callback)

}
