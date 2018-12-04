package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/kexirong/monitor/agent/goplugin"
	"github.com/kexirong/monitor/agent/scriptplugin"
	"github.com/kexirong/monitor/common"
	"github.com/kexirong/monitor/common/packetparse"
	"github.com/kexirong/monitor/common/queue"
)

func scriptPluginScheduler(qe *queue.BytesQueue) {
	res, err := http.Post(fmt.Sprintf("http://%s/plugin_config", conf.ServerHTTP),
		"application/json",
		strings.NewReader(`{"method":"getown"}`),
	)
	if err != nil {
		Logger.Warning.Println("The server may be down!")
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var resp common.HttpResp
	var scs []*common.ScriptConf
	resp.Result = &scs
	//fmt.Println(string(body))
	json.Unmarshal(body, &resp)
	if resp.Code != 200 {
		log.Fatal(errors.New(resp.Msg))
	}
	downloaurl := fmt.Sprintf("http://%s/downloadsscript/", conf.ServerHTTP)
	for _, sc := range scs {
		if sc.HostName != _hostname {
			continue
		}
		fmt.Println(sc.FileName)
		err := scriptplugin.CheckDownloads(downloaurl, filepath.Join(scriptPath, sc.FileName), false)

		if err != nil {
			Logger.Error.Println(err)
			continue
		}
		tasker := scriptplugin.NewScripter(filepath.Join(scriptPath, sc.FileName),
			time.Duration(sc.Timeout)*time.Second)

		scriptScheduled.AddTask(time.Duration(sc.Interval)*time.Second, tasker)
	}

	var callback = func(b []byte, err error) {
		if err != nil {
			Logger.Error.Println(err.Error())
			return
		}
		if b == nil {
			Logger.Error.Println("callback error arg b is nil")
			return
		}
		Logger.Info.Println(string(b))

		var tps packetparse.TargetPackets
		err = tps.UnmarshalJSON(b)
		if err != nil {
			Logger.Error.Println("callback json.Unmarshal TargetPacket error:", err.Error())
			return
		}
		//	btps, err := packetparse.TargetPacketsMarshal(tps)
		//	if err != nil {
		//		Logger.Error.Println("callback packetparse.TargetPacketsMarshal TargetPackets error:", err.Error())
		//		return
		//	}
		for _, tp := range tps {
			if err := qe.PutWait(tp, 1000); err != nil {
				Logger.Error.Println("scriptPluginScheduler error: " + err.Error())
			}
		}
	}

	Logger.Info.Println("activePluginScheduler staring")
	scriptScheduled.Star(callback)
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
				//btps, err := packetparse.TargetPacketsMarshal(tps)
				//	if err == nil {
				for _, tp := range tps {
					if err := qe.PutWait(tp, 1000); err != nil {
						Logger.Error.Println("gopluginScheduler error: " + err.Error())
					}
				}
			}
		}(name, plugin)
	}
}
