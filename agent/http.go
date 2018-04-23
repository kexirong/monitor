package main

import (
	"io/ioutil"
	"net/http"

	"github.com/kexirong/monitor/common"
)

func init() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/console", func(w http.ResponseWriter, r *http.Request) {
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
		err = json.Unmarshal(req, &console)
		if err != nil {
			Logger.Error.Println(err)
			return
		}
		switch console.Category {
		case "pyplugin":
			for i := 0; i < len(console.Events); i++ {
				nv := pp.AddEventAndWaitResult(console.Events[i])
				console.Events[i].Result = nv.Result
				if console.Events[i].Target != nv.Target ||
					console.Events[i].Method != nv.Method ||
					console.Events[i].Arg != nv.Arg {

					console.Events[i].Result = "server internal error"
				}
			}
			b, _ := json.MarshalIndent(console, "", "    ")
			w.Write(b)
		default:
			w.Write([]byte("unkown Category"))
		}

	})

}
