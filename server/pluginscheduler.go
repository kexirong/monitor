package main

import (
	"time"

	"monitor/common/packetparse"
	"monitor/common/scheduler"
	"monitor/server/activeplugin"
	"monitor/server/models"
)

var taskScheduled = scheduler.New()

func activePluginScheduler() {
	aps, err := models.ActiveProbesAll(monitorDB)
	if err != nil {
		panic(err)
	}

	for _, ap := range aps {
		switch ap.PluginName {
		case "http_probe":

			var tasker = activeplugin.NewHTTPProbe(ap.HostName)
			apcs, err := models.ActiveProbeConfigsByActiveProbeID(monitorDB, ap.ID)
			if err != nil {
				panic(err)
			}

			for _, apc := range apcs {
				err := tasker.AddJob(apc.Target, apc.Arg1, apc.Arg2)
				if err != nil {
					Logger.Error.Println(err)
				}
			}

			taskScheduled.AddTask(time.Second*time.Duration(ap.Interval), tasker)

		case "process_probe":
			var tasker = activeplugin.NewProcessProbe(ap.HostName, ap.HostIP)
			apcs, err := models.ActiveProbeConfigsByActiveProbeID(monitorDB, ap.ID)
			if err != nil {
				panic(err)
			}
			for _, apc := range apcs {
				err := tasker.AddJob(apc.Target)
				if err != nil {
					Logger.Error.Println(err)
				}
			}
			taskScheduled.AddTask(time.Second*time.Duration(ap.Interval), tasker)
		}

	}
	var callback = func(b []byte, err error) {
		if err != nil {
			Logger.Error.Println(err.Error())
			return
		}
		//	Logger.Info.Println("callback arg b is: ", string(b))
		var tps packetparse.TargetPackets
		_, err = tps.Unmarshal(b)
		if err != nil {
			Logger.Error.Println("callback Unmarshal TargetPacket error:", err.Error())
			return
		}
		for _, tp := range tps {
			go func(p *packetparse.TargetPacket) {

				err := influxdbwriter.Write(p)
				if err != nil {
					Logger.Error.Println("writeToInfluxdb error:", err.Error(), "\n", p.String())
				}

				judgeAlarm(p)
			}(tp)

		}
	}
	Logger.Info.Println("activePluginScheduler staring")
	taskScheduled.Start(callback)

}
