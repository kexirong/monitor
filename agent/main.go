package main

import (
	"os"
	"runtime"

	jsoniter "github.com/json-iterator/go"
	"monitor/common/queue"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {

	btq := queue.NewBtsQueue(4096)
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)

	go sendStart(conf.Servers, btq)

	go goPluginScheduler(btq)
	go scriptPluginScheduler(btq)

	startHTTPsrv()
}

func getCurrentPath() string {
	path, err := os.Getwd()
	checkErr(err)
	return path
}

func checkErr(err error) {
	if err != nil {
		Logger.Error.Panicf("error: %s, exit!", err.Error())
		os.Exit(1)
	}
}

/*
	cpuprofile := "./agent.prof"
	if isExist(cpuprofile) {
		err := os.Remove(cpuprofile)

		if err != nil {
			os.Exit(1)
		}
	}
	f, err := os.Create(cpuprofile)
	if err != nil {
		fmt.Println(err)
	}
	pprof.StartCPUProfile(f)
	go func() {
		<-time.After(time.Second * 600)
		fmt.Println("StopCPUProfile")
		pprof.StopCPUProfile()
	}()
*/
