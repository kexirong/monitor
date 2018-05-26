package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kexirong/monitor/common"
)

func init() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/console", func(w http.ResponseWriter, r *http.Request) {
		var ret common.HttpResp
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("BadRequest"))
			return
		}
		var console common.Console
		req, err := ioutil.ReadAll(r.Body)
		//defer r.Body.Close()
		if err != nil {
			Logger.Error.Println("fail to read requset data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//不需要对Unmarshal 失败的错误信息进行处理
		json.Unmarshal(req, &console)

		switch console.Category {
		case "scriptplugin":
			console.Events.UniqueID = common.NewUniqueID(10)
			if console.Events.Method == "add" {
				err := sp.CheckDownloads(fmt.Sprintf("http://%s/getscript/", conf.ServerHTTP), console.Events.Target, true)
				if err != nil {
					ret.Code = 500
					ret.Msg = err.Error()
					break
				}
			}
			nv := sp.AddEventAndWaitResult(console.Events)
			//console.Events.Result = nv.Result
			if console.Events.UniqueID != nv.UniqueID {
				ret.Code = 500
				ret.Msg = "please retry"
			}
			if nv.Result == "ok" {
				ret.Code = 200
				ret.Msg = "ok"
			}
			if err := json.Unmarshal([]byte(nv.Result), &ret.Result); err != nil {
				ret.Code = 400
				ret.Msg = nv.Result
				ret.Result = nil
			} else {
				ret.Code = 200
				ret.Msg = "ok"
			}

		default:
			ret.Code = 400
			ret.Msg = "unkown Category"
		}
		bt, _ := json.MarshalIndent(ret, "", "   ")
		w.Write(bt)
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
		fmt.Println(len(pl))
		b, _ := json.Marshal(ret)
		w.Write(b)
	})
}
