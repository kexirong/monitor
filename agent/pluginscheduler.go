package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kexirong/monitor/agent/goplugin"
	"github.com/kexirong/monitor/agent/pyplugin"
	"github.com/kexirong/monitor/common/packetparse"
	"github.com/kexirong/monitor/common/queue"
)

func pyPluginScheduler(qe *queue.BytesQueue) {
	pp, err := pyplugin.Initialize("./pyplugin")
	if err != nil {
		panic(err)
	}
	err = pp.InsertEntry("cpu", 5)
	if err != nil {
		fmt.Println(err)
	}
	err = pp.InsertEntry("cpus1", 2)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		events := []string{"add:cpus2|1", "add:cpus|1", "add:cpus3|1", "delete:cpus2|1"}
		for _, v := range events {
			pp.Event <- v
		}
	}()
	for {
		errs := pp.WaitAndEventDeal()
		if len(errs) > 0 {
			for _, e := range errs {
				Logger.Error.Println(e)
			}
		}
		ret, err := pp.Scheduler()
		if err != nil {
			Logger.Error.Println(err)
			continue
		}
		go func(ret string) {
			var tps []packetparse.TargetPacket
			err = json.Unmarshal([]byte(ret), &tps)
			if err != nil {
				Logger.Error.Println(err.Error())
			}
			for _, pk := range tps {
				fmt.Println(pk)
				gatherbs, err := packetparse.TargetPackage(pk)
				if err == nil {
					if err := qe.PutWait(gatherbs, 1000); err != nil {
						Logger.Error.Println("gopluginScheduler error: " + err.Error())
					}
				} else {
					Logger.Error.Println("gopluginScheduler error: " + err.Error())
				}
			}
		}(ret)

	}
}

//gopluginScheduler
func goPluginScheduler(qe *queue.BytesQueue) {
	for name, plugin := range goplugin.GopluginMap {
		go func(name string, plugin goplugin.PLUGIN) {
			for range time.Tick(time.Duration(plugin.GetInterval())) {
				gather, err := plugin.Gather()
				if err != nil {
					Logger.Error.Printf("gopluginScheduler error:%s, %s ", name, err.Error())
					return
				}
				for _, pk := range gather {
					gatherbs, err := packetparse.TargetPackage(pk)
					if err == nil {
						if err := qe.PutWait(gatherbs, 1000); err != nil {
							Logger.Error.Println("gopluginScheduler error: " + err.Error())
						}
					} else {
						Logger.Error.Println("gopluginScheduler error: " + err.Error())
					}
				}
			}
		}(name, plugin)
	}
}
