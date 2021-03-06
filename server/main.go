package main

import "runtime"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	go startHTTPsrv()
	go heartdeamo()
	go activePluginScheduler()

	Logger.Info.Println("runing, listen:,", conf.Service)
	Logger.Info.Printf("Judge:%v\n", Judge)

	startTCPsrv()
}
