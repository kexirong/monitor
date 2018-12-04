package main

import "runtime"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go startHTTPsrv()
	go heartdeamo()
	go activePluginScheduler()

	Logger.Info.Println("runing, listen:,", conf.Service)
	startTCPsrv()
}
