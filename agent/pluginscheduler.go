package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kexirong/monitor/common"

	"github.com/kexirong/monitor/agent/goplugin"
	"github.com/kexirong/monitor/common/packetparse"
	"github.com/kexirong/monitor/common/queue"
)

func scriptPluginScheduler(qe *queue.BytesQueue) {
	res, err := http.Get(fmt.Sprintf("http://%s/config/plugin", conf.ServerHTTP))
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var resp common.HttpResp
	var sconf []common.ScriptConf
	resp.Result = &conf
	json.Unmarshal(body, &resp)
	downloaurl := fmt.Sprintf("http://%s/getscript/", conf.ServerHTTP)
	for _, ret := range sconf {
		err := sp.CheckDownloads(downloaurl, ret.FileName, false)
		if err != nil {
			Logger.Error.Println(err)
		}
		if err := sp.InsertEntry(ret.FileName, ret.Interval, ret.TimeOut); err != nil {
			Logger.Error.Println(err.Error())
		}
	}

	for {
		sp.WaitAndEventDeal()
		ret, err := sp.Scheduler()
		if err != nil {
			Logger.Error.Println(err)
			continue
		}
		//Logger.Info.Println(ret)
		go func(ret []byte) {
			//fmt.Println(ret)
			var tps []packetparse.TargetPacket
			err = json.Unmarshal(ret, &tps)
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
