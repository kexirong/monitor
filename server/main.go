package main

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	go heartdeamo()
	go httpprobesched()
	go alarmdo()
	Logger.Info.Println("runing,listen:,", conf.Service)
	startTCPsrv()
}
