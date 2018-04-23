package main

import (
	"fmt"
	"time"

	"github.com/kexirong/monitor/agent/goplugin"
	"github.com/kexirong/monitor/common/packetparse"
	"github.com/kexirong/monitor/common/queue"
)

func pyPluginScheduler(qe *queue.BytesQueue) {
	err := pp.InsertEntry("cpu", 5)
	if err != nil {
		fmt.Println(err)
	}
	err = pp.InsertEntry("cpus1", 2)
	if err != nil {
		fmt.Println(err)
	}

	for {
		pp.WaitAndEventDeal()
		ret, err := pp.Scheduler()
		if err != nil {
			Logger.Error.Println(err)
			continue
		}
		go func(ret string) {
			//fmt.Println(ret)
			var tps []packetparse.TargetPacket
			err = json.Unmarshal([]byte(ret), &tps)
			if err != nil {
				Logger.Error.Println(err.Error())
			}
			for _, pk := range tps {
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
