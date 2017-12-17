package main

import (
	"fmt"
	"time"

	"github.com/kexirong/monitor/agent/goplugin"
)

func unixTimeComepare(t1, t2 int64) int64 {
	if t1 == 0 {
		return t2
	}
	if t1 < t2 {
		return t2
	}
	return t2
}

//gopluginScheduler
func gopluginScheduler() {
	var timeNow, nextTime, tickTime int64
	for {
		nextTime = 0
		for _, plugin := range goplugin.GopluginMap {
			timeNow = time.Now().UnixNano()
			fmt.Println("timenow: ", timeNow, "plugin nextTime: ", plugin.NextTime)
			if timeNow >= plugin.NextTime {
				plugin.NextTime += plugin.Instance.GetStep()
				go func() {
					gather, err := plugin.Instance.Gather()
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println(gather)
				}()
				continue
			}
			nextTime = unixTimeComepare(nextTime, plugin.NextTime)
		}
		tickTime = nextTime - time.Now().UnixNano()
		if tickTime > 0 {
			<-time.After(time.Duration(tickTime))
		}
	}
}
