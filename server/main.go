package main

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	go startHTTPsrv()
	go heartdeamo()
	go activePluginScheduler()
	go alarmdo()

	Logger.Info.Println("runing,listen:,", conf.Service)
	startTCPsrv()
}
