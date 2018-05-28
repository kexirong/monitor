package main

import "net/http"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	go http.ListenAndServe(":5001", nil)
	go heartdeamo()
	go activePluginScheduler()
	go alarmdo()

	Logger.Info.Println("runing,listen:,", conf.Service)
	startTCPsrv()
}
