package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kexirong/monitor/agent/goplugin"
	"github.com/kexirong/monitor/common"
	"github.com/kexirong/monitor/common/packetparse"
	"github.com/kexirong/monitor/common/queue"
)

func scriptPluginScheduler(qe *queue.BytesQueue) {
	res, err := http.Post(fmt.Sprintf("http://%s/plugin_config", conf.ServerHTTP),
		"application/json",
		strings.NewReader(`{"method":"getlist"}`),
	)
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
	if resp.Code != 200 {
		log.Fatal(errors.New(resp.Msg))
	}
	json.Unmarshal(body, &resp)
	if resp.Code != 200 {
		log.Fatal(errors.New(resp.Msg))
	}
	downloaurl := fmt.Sprintf("http://%s/downloadsscript/", conf.ServerHTTP)
	for _, ret := range sconf {
		if ret.HostName != _hostname {
			continue
		}
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
		go func(ret []byte) {
			//fmt.Println(ret)
			var tps []*packetparse.TargetPacket
			err = json.Unmarshal(ret, &tps)
			if err != nil {
				Logger.Error.Println(err.Error())
			}
			btps, err := packetparse.TargetPacketsMarshal(tps)
			if err == nil {
				if err := qe.PutWait(btps, 1000); err != nil {
					Logger.Error.Println("scriptPluginScheduler error: " + err.Error())
				}
			} else {
				Logger.Error.Println("scriptPluginScheduler error: " + err.Error())
			}
		}(ret)
	}
}

//gopluginScheduler
func goPluginScheduler(qe *queue.BytesQueue) {
	for name, plugin := range goplugin.GopluginMap {
		go func(name string, plugin goplugin.PLUGIN) {
			for range time.Tick(time.Duration(plugin.GetInterval())) {
				tps, err := plugin.Gather()
				if err != nil {
					Logger.Error.Printf("gopluginScheduler error:%s, %s ", name, err.Error())
					return
				}
				btps, err := packetparse.TargetPacketsMarshal(tps)
				if err == nil {
					if err := qe.PutWait(btps, 1000); err != nil {
						Logger.Error.Println("gopluginScheduler error: " + err.Error())
					}
				} else {
					Logger.Error.Println("gopluginScheduler error: " + err.Error())
				}
			}
		}(name, plugin)
	}
}
