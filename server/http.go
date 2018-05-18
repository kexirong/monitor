package main

import (
	"net/http"
	"strings"

	"github.com/kexirong/monitor/common"
)

func init() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/config/plugin", func(w http.ResponseWriter, r *http.Request) {
		var ret = common.HttpResp{
			Code: 200,
			Msg:  "ok",
		}
		if r.Method == "GET" {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			conf := pluginconfGet(ip)
			ret.Result = conf
		} else {
			ret.Code = 400
			ret.Msg = "bad request"
		}
		bret, _ := json.Marshal(ret)
		w.Write(bret)
	})

	http.Handle("/getscript/", http.StripPrefix("/getscript/", http.FileServer(http.Dir("./scriptrepo/"))))

}
