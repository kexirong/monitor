package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/kexirong/monitor/common/queue"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {

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
	var addr = flag.String("s", "ip:port", "server addrs multi delimit  with ',' ")
	flag.Parse()

	if *addr == "ip:port" {
		Logger.Error.Fatalln("no servers;-h;exit")
	}

	servers := strings.Split(*addr, ",")
	btq := queue.NewBtsQueue(4096)
	runtime.GOMAXPROCS(runtime.NumCPU())

	go sendStart(servers, btq)

	go goPluginScheduler(btq)
	go pyPluginScheduler(btq)

	log.Fatal(http.ListenAndServe(":5101", nil))
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
