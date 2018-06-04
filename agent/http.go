package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/kexirong/monitor/agent/scriptplugin"
	"github.com/kexirong/monitor/common"
)

func startHTTPsrv() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/scriptplugin", func(w http.ResponseWriter, r *http.Request) {

		var req common.HttpReq
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Logger.Error.Println("fail to read requset data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		var ret = common.HttpResp{
			Code: 200,
			Msg:  "ok",
		}
		var sc common.ScriptConf
		req.Cause = &sc
		//不需要对Unmarshal 失败的错误信息进行处理
		json.Unmarshal(body, &req)

		if r.Method == "POST" {

			if sc.HostName != _hostname {
				ret.Code = 400
				ret.Msg = "hostname not is " + _hostname
			}

			switch req.Method {

			case "add":
				downloaurl := fmt.Sprintf("http://%s/downloadsscript/", conf.ServerHTTP)
				err := scriptplugin.CheckDownloads(downloaurl, path.Join(scriptPath, sc.FileName), false)
				if err != nil {
					Logger.Error.Println(err)
					break
				}
				tasker := scriptplugin.NewScripter(path.Join(scriptPath, sc.FileName),
					time.Duration(sc.Timeout)*time.Second)

				scriptScheduled.AddTask(time.Second*time.Duration(sc.Interval), tasker)

			case "delete":
				tasker := scriptplugin.NewScripter(path.Join(scriptPath, sc.FileName),
					time.Duration(sc.Timeout)*time.Second)

				scriptScheduled.DeleteTask(tasker.Name())

			case "getlist":

				ret.Result = scriptScheduled.EcheTaskList()
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				}
			}
		} else {
			ret.Code = 400
			ret.Msg = "bad request"
		}
		bret, _ := json.Marshal(ret)
		w.Write(bret)

	})

	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		var ret = common.HttpResp{
			Code: 200,
			Msg:  "ok",
		}

		if r.Method != "GET" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("BadRequest"))
			return
		}
		r.ParseForm()
		patterns := r.Form["pattern"]
		var pl common.ProcessList
		pl.Init()
		if err := pl.FilterCmdline(patterns); err != nil {
			ret.Code = 400
			ret.Msg = err.Error()
		}

		if ret.Code == 200 {
			pl.LoadsProcessInfo()
			ret.Result = pl
		}
		b, _ := json.Marshal(ret)
		w.Write(b)
	})

	log.Fatal(http.ListenAndServe(conf.HTTPListen, nil))
}
