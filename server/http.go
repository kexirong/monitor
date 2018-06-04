package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kexirong/monitor/common/scheduler"

	"github.com/kexirong/monitor/common"
	"github.com/kexirong/monitor/server/activeplugin"
	"github.com/kexirong/monitor/server/models"
)

func startHTTPsrv() {

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	///get_plugin_config 即将废弃
	/*
		http.HandleFunc("/get_plugin_config", func(w http.ResponseWriter, r *http.Request) {
			var ret = common.HttpResp{
				Code: 200,
				Msg:  "ok",
			}
			if r.Method == "GET" {
				ip := strings.Split(r.RemoteAddr, ":")[0]
				conf, err := models.GetPluginConfigsByHostIP(monitorDB, ip)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				} else {
					ret.Result = conf
				}

			} else {
				ret.Code = 400
				ret.Msg = "bad request"
			}
			bret, _ := json.Marshal(ret)
			w.Write(bret)
		})
	*/
	http.HandleFunc("/plugin", func(w http.ResponseWriter, r *http.Request) {

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

		var p = &models.Plugin{}
		req.Cause = p
		//不需要对Unmarshal 失败的错误信息进行处理
		json.Unmarshal(body, &req)

		if r.Method == "POST" {

			np, err := models.PluginByPluginName(monitorDB, p.PluginName)
			if err != nil {
				ret.Code = 400
				ret.Msg = err.Error()
			}
			switch req.Method {

			case "get":
				ret.Result = np
			case "add", "update":
				{
					np.FileName = p.FileName
					np.PluginType = p.PluginType
					np.PluginName = p.PluginName
					np.Comment = p.Comment
				}
				err = np.Save(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				}

			case "delete":

				err = np.Delete(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				}

			case "getlist":
				conf, err := models.PluginAll(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				} else {
					ret.Result = conf
					ret.Code = 200
					ret.Msg = "ok"
				}
			}
		} else {
			ret.Code = 400
			ret.Msg = "bad request"
		}
		bret, _ := json.Marshal(ret)
		w.Write(bret)
	})

	http.HandleFunc("/plugin_config", func(w http.ResponseWriter, r *http.Request) {

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

		var pc = &models.PluginConfig{}
		req.Cause = pc
		//不需要对Unmarshal 失败的错误信息进行处理
		json.Unmarshal(body, &req)

		if r.Method == "POST" {

			npc, err := models.PluginConfigByID(monitorDB, pc.ID)
			if err != nil {
				ret.Code = 400
				ret.Msg = err.Error()
			}

			switch req.Method {

			case "get":
				ret.Result = npc

			case "add", "update":
				{
					npc.HostIP = pc.HostIP
					npc.PluginName = pc.PluginName
					npc.Interval = pc.Interval
					npc.HostName = pc.HostName
					npc.Timeout = pc.Timeout
				}
				err = npc.Save(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				} else {

					//taskScheduled.
				}

			case "delete":
				err = npc.Delete(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				} else {

					//taskScheduled.
				}

			case "getlist":

				ret.Result, err = models.PluginConfigsAll(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				} else {
					ret.Code = 200
					ret.Msg = "ok"
				}
			}
		} else {
			ret.Code = 400
			ret.Msg = "bad request"
		}
		bret, _ := json.Marshal(ret)
		w.Write(bret)
	})

	http.HandleFunc("/active_probe", func(w http.ResponseWriter, r *http.Request) {

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
		var ap = &models.ActiveProbe{}

		req.Cause = ap
		//不需要对Unmarshal 失败的错误信息进行处理
		json.Unmarshal(body, &req)

		if r.Method == "POST" {

			nap, err := models.ActiveProbeByID(monitorDB, ap.ID)
			if err != nil {
				ret.Code = 400
				ret.Msg = err.Error()
			}

			switch req.Method {
			case "get":
				ap = nap
			case "getlist":

				ret.Result, err = models.ActiveProbeAll(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				} else {
					ret.Code = 200
					ret.Msg = "ok"
				}

			case "add", "update":
				{
					nap.HostName = ap.HostName
					nap.Interval = ap.Interval
					nap.IP = ap.IP
					nap.PluginName = ap.PluginName
				}
				err = nap.Save(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}
				var tasker scheduler.Tasker
				switch ap.PluginName {
				case "http_probe":

					tasker = activeplugin.NewHTTPProbe(ap.HostName)

				case "process_probe":
					tasker = activeplugin.NewProcessProbe(ap.HostName, ap.IP)

				}

				taskScheduled.AddTask(time.Second*time.Duration(ap.Interval), tasker)

			case "delete":

				err = nap.Delete(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}
				var tasker scheduler.Tasker
				switch ap.PluginName {
				case "http_probe":

					tasker = activeplugin.NewHTTPProbe(ap.HostName)

				case "process_probe":
					tasker = activeplugin.NewProcessProbe(ap.HostName, ap.IP)

				}
				taskScheduled.DeleteTask(tasker.Name())

			case "getruninglist":

				taskList := taskScheduled.EcheTaskList()

				if err := json.Unmarshal([]byte(taskList), &ret.Result); err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
					ret.Result = nil
				} else {
					ret.Code = 200
					ret.Msg = "ok"
				}

			}
		} else {
			ret.Code = 400
			ret.Msg = "bad request"
		}
		bret, _ := json.Marshal(ret)
		w.Write(bret)
	})

	http.HandleFunc("/active_probe_config", func(w http.ResponseWriter, r *http.Request) {

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
		var apc models.ActiveProbeConfig
		req.Cause = &apc
		//不需要对Unmarshal 失败的错误信息进行处理
		json.Unmarshal(body, &req)

		if r.Method == "POST" {

			napc, err := models.ActiveProbeConfigByID(monitorDB, apc.ID)

			if err != nil {
				ret.Code = 400
				ret.Msg = err.Error()
			}

			switch req.Method {

			case "get":
				ret.Result = napc
			case "add", "update":
				{
					napc.ActiveProbeID = apc.ActiveProbeID
					napc.Arg1 = apc.Arg1
					napc.Arg2 = apc.Arg2
					napc.Target = apc.Target
				}
				err = napc.Save(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}
				ap, err := models.ActiveProbeByID(monitorDB, apc.ActiveProbeID)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}

				var tasker scheduler.Tasker
				switch ap.PluginName {
				case "http_probe":
					tasker = activeplugin.NewHTTPProbe(ap.HostName)

				case "process_probe":
					tasker = activeplugin.NewProcessProbe(ap.HostName, ap.IP)

				}
				err = taskScheduled.AddJob(tasker.Name(), apc.Target, apc.Arg1, apc.Arg2)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				}

			case "delete":
				err = napc.Delete(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}

				ap, err := models.ActiveProbeByID(monitorDB, apc.ActiveProbeID)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}

				var tasker scheduler.Tasker
				switch ap.PluginName {
				case "http_probe":
					tasker = activeplugin.NewHTTPProbe(ap.HostName)

				case "process_probe":
					tasker = activeplugin.NewProcessProbe(ap.HostName, ap.IP)

				}
				err = taskScheduled.DeleteJob(tasker.Name(), apc.Target)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				}

			case "getlist":

				ret.Result, err = models.ActiveProbeConfigsByActiveProbeID(monitorDB, apc.ActiveProbeID)
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

	http.Handle("/downloadsscript/", http.StripPrefix("/downloadsscript/", http.FileServer(http.Dir("./scriptrepo/"))))

	log.Fatal(http.ListenAndServe(":5001", nil))
}
