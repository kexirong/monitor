package main

import (
	"time"

	"github.com/kexirong/monitor/agent/goplugin"
	"github.com/kexirong/monitor/common/packetparse"
	"github.com/kexirong/monitor/common/queue"
)

func pyPluginScheduler(qe *queue.BytesQueue) {

}

//gopluginScheduler
func goPluginScheduler(qe *queue.BytesQueue) {
	for name, plugin := range goplugin.GopluginMap {
		go func(name string, plugin goplugin.PLUGIN) {
			for range time.Tick(time.Duration(plugin.GetStep())) {
				gather, err := plugin.Gather()
				if err != nil {
					Logger.Error.Printf("gopluginScheduler errror:%s, %s ", name, err.Error())
					return
				}
				for _, pk := range gather {
					gatherbs, err := packetparse.TargetPackage(pk)
					if err == nil {
						if err := qe.PutWait(gatherbs, 500); err != nil {
							Logger.Error.Println("gopluginScheduler errror: " + err.Error())
						}
					} else {
						Logger.Error.Println("gopluginScheduler errror: " + err.Error())
					}
				}
			}
		}(name, plugin)
	}
}
