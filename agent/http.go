package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"monitor/agent/scriptplugin"
	"monitor/common"
	"monitor/server/models"
)

func startHTTPsrv() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/scriptplugin", func(w http.ResponseWriter, r *http.Request) {
		Logger.Info.Println("/scriptplugin")
		var req common.HttpReq
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Logger.Error.Println("fail to read requset data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//defer r.Body.Close()
		var ret = common.HttpResp{
			Code: 200,
			Msg:  "ok",
		}
		var pc models.PluginConfig
		req.Cause = &pc
		//不需要对Unmarshal 失败的错误信息进行处理
		json.Unmarshal(body, &req)
		Logger.Info.Println("/scriptplugin:", req.Cause)
		Logger.Info.Println("/scriptplugin:", pc.Plugin)
		if r.Method == "POST" && (pc.Plugin != nil || req.Method == "getlist") {

			if pc.HostName != _hostname {
				ret.Code = 400
				ret.Msg = "hostname not is " + _hostname
				goto end
			}

			switch req.Method {
			case "add":
				downloadurl := fmt.Sprintf("http://%s/scriptdownloads/", conf.ServerHTTP)
				err := scriptplugin.CheckDownloads(downloadurl, path.Join(scriptPath, pc.FileName), true)
				if err != nil {
					Logger.Error.Println(err)
					break
				}
				tasker := scriptplugin.NewScripter(path.Join(scriptPath, pc.FileName),
					time.Duration(pc.Timeout)*time.Second)

				scriptScheduled.AddTask(time.Second*time.Duration(pc.Interval), tasker)

			case "delete":
				tasker := scriptplugin.NewScripter(path.Join(scriptPath, pc.FileName),
					time.Duration(pc.Timeout)*time.Second)

				scriptScheduled.DeleteTask(tasker.Name())

			case "getlist":
				taskList := scriptScheduled.EcheTaskList()
				if err := json.Unmarshal([]byte(taskList), &ret.Result); err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
					ret.Result = nil
				}
			default:
				ret.Code = 400
				ret.Msg = "unkown method"
			}
		} else {
			ret.Code = 400
			ret.Msg = "bad request"
		}
	end:
		bret, _ := json.Marshal(ret)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
