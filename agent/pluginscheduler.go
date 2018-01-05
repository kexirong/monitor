package main

import (
	"fmt"
	"time"

	"github.com/kexirong/monitor/common/packetparse"

	"github.com/kexirong/monitor/common/queue"

	"github.com/kexirong/monitor/agent/goplugin"
)

func unixTimeComepare(t1, t2 int64) int64 {
	if t1 <= 0 {
		return t2
	}
	if t1 < t2 {
		return t1
	}
	return t2
}

//gopluginScheduler
func gopluginScheduler(qe *queue.BytesQueue) {
	var timeNow, nextTime, tickTime int64
	fmt.Println(goplugin.GopluginMap)
	for {
		nextTime = 0
		for name, plugin := range goplugin.GopluginMap {
			timeNow = time.Now().UnixNano()
			fmt.Println("name:", name, plugin.Instance.GetStep()/(1000*1000*1000), "timenow: ", timeNow/(1000*1000*1000), "pluginnextTime: ", plugin.NextTime/(1000*1000*1000))
			if timeNow >= plugin.NextTime {
				plugin.NextTime += plugin.Instance.GetStep()
				go func() {
					gather, err := plugin.Instance.Gather()
					if err != nil {
						fmt.Println(err)
						return
					}
					for _, pk := range gather {
						gatherbs, err := packetparse.TargetPackage(pk)
						if err == nil {
							if err := qe.PutWait(gatherbs); err != nil {
								fmt.Println("gopluginScheduler errror: " + err.Error())
							}
						}
					}
				}()
			}
			nextTime = unixTimeComepare(nextTime, plugin.NextTime)
		}
		tickTime = nextTime - time.Now().UnixNano()
		if tickTime > 0 {
			fmt.Println("sleeping....", tickTime/(1000*1000*1000))
			<-time.After(time.Duration(tickTime))
		}
	}
}

//gopluginScheduler
func gopluginScheduler2(qe *queue.BytesQueue) {

	fmt.Println(goplugin.GopluginMap)

	for name, plugin := range goplugin.GopluginMap {
		fmt.Println(name, ":")
		go func(plugin *goplugin.Goplugintype) {
			for {
				gather, err := plugin.Instance.Gather()
				if err != nil {
					fmt.Println(err)
					return
				}
				for _, pk := range gather {
					gatherbs, err := packetparse.TargetPackage(pk)
					if err == nil {
						if err := qe.PutWait(gatherbs); err != nil {
							fmt.Println("gopluginScheduler errror: " + err.Error())
						}
					}
				}
				<-time.After(time.Duration(plugin.Instance.GetStep()))
			}

		}(plugin)

	}

}
