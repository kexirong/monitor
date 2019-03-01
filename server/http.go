package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"monitor/common"
	"monitor/common/scheduler"
	"monitor/server/activeplugin"
	"monitor/server/models"
)

func startHTTPsrv() {
	//	models.XOLog = func(str string, param ...interface{}) { Logger.Info.Println(str, param) }
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/plugin", func(w http.ResponseWriter, r *http.Request) {

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

		var p = &models.Plugin{}
		req.Cause = p
		//不需要对Unmarshal 失败的错误信息进行处理
		json.Unmarshal(body, &req)

		if r.Method == "POST" {

			np, err := models.PluginByPluginName(monitorDB, p.PluginName)
			if err != nil && !(req.Method == "add" || req.Method == "getlist") {
				ret.Code = 400
				ret.Msg = err.Error()
				goto end
			}
			switch req.Method {

			case "get":
				ret.Result = np
			case "add", "update":
				if np == nil {
					np = p
				} else {
					np.FileName = p.FileName
					np.PluginType = p.PluginType
					np.PluginName = p.PluginName
					np.Comment = p.Comment
				}
				err = np.Save(monitorDB)
				if err != nil {
					Logger.Error.Println(err.Error())
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
				ret.Result, err = models.PluginsAll(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
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

	http.HandleFunc("/plugin_config", func(w http.ResponseWriter, r *http.Request) {

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

		var pc = &models.PluginConfig{}
		req.Cause = pc
		//不需要对Unmarshal 失败的错误信息进行处理
		json.Unmarshal(body, &req)

		if r.Method == "POST" {
			npc, err := models.PluginConfigByID(monitorDB, pc.ID)
			if err != nil && !(req.Method == "add" || req.Method == "getlist") {
				ret.Code = 400
				ret.Msg = err.Error()
				goto end
			}

			switch req.Method {

			case "get":
				ret.Result = npc

			case "add", "update":
				if npc == nil {
					npc = pc
				} else {
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
					break
				}
				clt := activeplugin.NewHTTPClient(5 * time.Second)
				err := npc.LoadPlugin(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = "query associate plugin failed"
					break
				}
				req := common.HttpReq{
					Method: "add",
					Cause:  npc,
				}
				bts, _ := json.Marshal(req)
				body := bytes.NewReader(bts)
				rsp, err := clt.Post(
					fmt.Sprintf("%s://%s:%d/scriptplugin",
						conf.Agent.Scheme,
						npc.HostIP,
						conf.Agent.HTTPPort),
					"application/json",
					body)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}
				defer rsp.Body.Close()
				rbts, _ := ioutil.ReadAll(rsp.Body)

				json.Unmarshal(rbts, &ret)
				//taskScheduled

			case "delete":
				err = npc.Delete(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}
				clt := activeplugin.NewHTTPClient(5 * time.Second)
				err := npc.LoadPlugin(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = "query associate plugin failed"
					break
				}
				req := common.HttpReq{
					Method: "delete",
					Cause:  npc,
				}
				bts, _ := json.Marshal(req)
				body := bytes.NewReader(bts)
				rsp, err := clt.Post(
					fmt.Sprintf("%s://%s:%d/scriptplugin",
						conf.Agent.Scheme,
						npc.HostIP,
						conf.Agent.HTTPPort),
					"application/json",
					body)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}
				defer rsp.Body.Close()
				rbts, _ := ioutil.ReadAll(rsp.Body)

				json.Unmarshal(rbts, &ret)

			case "getlist":
				/*	ret.Result, err = models.PluginConfigsAll(monitorDB)
						if err != nil {
							ret.Code = 400
							ret.Msg = err.Error()
						}

					case "getown":
				*/
				var err error
				var pcs []*models.PluginConfig
				if pc.HostName != "" {
					pcs, err = models.PluginConfigsByHostName(monitorDB, pc.HostName)
				} else {
					pcs, err = models.PluginConfigsAll(monitorDB)
				}
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}
				for i := range pcs {
					err = pcs[i].LoadPlugin(monitorDB)
					if err != nil {
						break
					}
				}
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				} else {
					ret.Result = pcs
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

	http.HandleFunc("/active_probe", func(w http.ResponseWriter, r *http.Request) {

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
		var ap = &models.ActiveProbe{}

		req.Cause = ap
		//不需要对Unmarshal 失败的错误信息进行处理
		json.Unmarshal(body, &req)

		if r.Method == "POST" {
			nap, err := models.ActiveProbeByID(monitorDB, ap.ID)
			if err != nil && !(req.Method == "add" || req.Method == "getlist" || req.Method == "getruninglist") {
				Logger.Error.Println(err.Error())
				ret.Code = 400
				ret.Msg = err.Error()
				goto end
			}

			switch req.Method {
			case "get":
				ret.Result = nap
			case "getlist":

				ret.Result, err = models.ActiveProbesAll(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				}

			case "add", "update":

				if nap == nil {
					nap = ap
				} else {
					nap.HostName = ap.HostName
					nap.Interval = ap.Interval
					nap.HostIP = ap.HostIP
					nap.PluginName = ap.PluginName
				}

				err = nap.Save(monitorDB)
				if err != nil {
					Logger.Error.Println(err.Error())
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}
				var tasker scheduler.Tasker
				switch ap.PluginName {
				case "http_probe":

					tasker = activeplugin.NewHTTPProbe(ap.HostName)

				case "process_probe":
					tasker = activeplugin.NewProcessProbe(ap.HostName, ap.HostIP)

				}

				taskScheduled.AddTask(time.Second*time.Duration(ap.Interval), tasker)

			case "delete":
				if nap == nil {
					ret.Code = 400
					ret.Msg = "not exist"
					break
				}
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
					tasker = activeplugin.NewProcessProbe(ap.HostName, ap.HostIP)

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

	http.HandleFunc("/active_probe_config", func(w http.ResponseWriter, r *http.Request) {

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

		var apc = &models.ActiveProbeConfig{}
		req.Cause = apc

		//不需要对Unmarshal 失败的错误信息进行处理
		json.Unmarshal(body, &req)

		if r.Method == "POST" {

			napc, err := models.ActiveProbeConfigByID(monitorDB, apc.ID)

			if err != nil && !(req.Method == "add" || req.Method == "getlist") {
				ret.Code = 400
				ret.Msg = err.Error()
				goto end
			}

			switch req.Method {

			case "get":
				ret.Result = napc
			case "add", "update":
				if napc == nil {
					napc = apc
				} else {
					napc.ActiveProbeID = apc.ActiveProbeID
					napc.Arg1 = apc.Arg1
					napc.Arg2 = apc.Arg2
					napc.Target = apc.Target
				}
				err = napc.Save(monitorDB)
				if err != nil {
					Logger.Error.Println(err.Error())
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}
				ap, err := models.ActiveProbeByID(monitorDB, apc.ActiveProbeID)
				if err != nil {
					Logger.Error.Println(err.Error())
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}

				var tasker scheduler.Tasker
				switch ap.PluginName {
				case "http_probe":
					tasker = activeplugin.NewHTTPProbe(ap.HostName)

				case "process_probe":
					tasker = activeplugin.NewProcessProbe(ap.HostName, ap.HostIP)
				}

				err = taskScheduled.AddJob(tasker.Name(), napc.Target, napc.Arg1, napc.Arg2)
				if err != nil {
					Logger.Error.Println(err.Error())
					ret.Code = 400
					ret.Msg = err.Error()
				}

			case "delete":
				if napc == nil {
					ret.Code = 400
					ret.Msg = "not exist"
					break
				}
				err = napc.Delete(monitorDB)
				if err != nil {
					Logger.Error.Println(err.Error())
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}

				ap, err := models.ActiveProbeByID(monitorDB, apc.ActiveProbeID)
				if err != nil {
					Logger.Error.Println(err.Error())
					ret.Code = 400
					ret.Msg = err.Error()
					break
				}

				var tasker scheduler.Tasker
				switch ap.PluginName {
				case "http_probe":
					tasker = activeplugin.NewHTTPProbe(ap.HostName)

				case "process_probe":
					tasker = activeplugin.NewProcessProbe(ap.HostName, ap.HostIP)

				}
				err = taskScheduled.DeleteJob(tasker.Name(), apc.Target)
				if err != nil {
					Logger.Error.Println(err.Error())
					ret.Code = 400
					ret.Msg = err.Error()
				}

			case "getlist":
				if apc.ActiveProbeID == 0 {
					ret.Result, err = models.ActiveProbeConfigsAll(monitorDB)
				} else {
					ret.Result, err = models.ActiveProbeConfigsByActiveProbeID(monitorDB, apc.ActiveProbeID)
				}

				if err != nil {
					Logger.Error.Println(err.Error())
					ret.Code = 400
					ret.Msg = err.Error()
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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		bret, _ := json.Marshal(ret)

		w.Write(bret)

	})

	http.Handle("/scriptdownloads/", http.StripPrefix("/scriptdownloads/", http.FileServer(http.Dir("./scriptrepo/"))))

	http.HandleFunc("/alarm_send", func(w http.ResponseWriter, r *http.Request) {

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
		var al = &models.AlarmSend{}

		req.Cause = al
		//不需要对Unmarshal 失败的错误信息进行处理
		json.Unmarshal(body, &req)

		if r.Method == "POST" {

			nal, err := models.AlarmSendByAnchorPoint(monitorDB, al.AnchorPoint)
			if err != nil && !(req.Method == "add" || req.Method == "getlist") {
				ret.Code = 400
				ret.Msg = err.Error()
				goto end
			}

			switch req.Method {
			case "get":
				ret.Result = nal
			case "getlist":

				ret.Result, err = models.AlarmSendsAll(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				}

			case "add", "update":
				if nal == nil {
					nal = al
				} else {
					nal.AnchorPoint = al.AnchorPoint
					nal.Channel = al.Channel
					nal.List = al.List
					nal.Type = al.Type
					nal.SendTick = al.SendTick

				}
				err = nal.Save(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()

				}

			case "delete":
				if nal == nil {
					ret.Code = 400
					ret.Msg = "not exist"
					break
				}
				err = nal.Delete(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()

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

	http.HandleFunc("/alarm_judge", func(w http.ResponseWriter, r *http.Request) {

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
		var aj = &models.AlarmJudge{}

		req.Cause = aj
		//不需要对Unmarshal 失败的错误信息进行处理
		err = json.Unmarshal(body, &req)
		if err != nil {
			req.Method = "other"
			ret.Code = 400
			ret.Msg = err.Error()
		}

		if r.Method == "POST" {

			naj, err := models.AlarmJudgeByID(monitorDB, aj.ID)
			if err != nil && !(req.Method == "add" || req.Method == "getlist") {
				ret.Code = 400
				ret.Msg = err.Error()
				goto end
			}

			switch req.Method {
			case "get":
				ret.Result = naj
			case "getlist":

				ret.Result, err = models.AlarmJudgesAll(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				} else {
					ret.Code = 200
					ret.Msg = "ok"
				}

			case "add", "update":
				if naj == nil {
					naj = aj
				} else {
					naj.AnchorPoint = aj.AnchorPoint
					naj.Express = aj.Express
					naj.Level = aj.Level
					naj.Title = aj.Title
				}
				err = naj.Save(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
				}

			case "delete":
				if naj == nil {
					ret.Code = 400
					ret.Msg = "not exist"
					break
				}
				err = naj.Delete(monitorDB)
				if err != nil {
					ret.Code = 400
					ret.Msg = err.Error()
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
	log.Fatal(http.ListenAndServe(":5001", nil))
}
