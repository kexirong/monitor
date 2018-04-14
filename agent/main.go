package main

import (
	"flag"
	"os"
	"runtime"
	"strings"

	"github.com/kexirong/monitor/common/queue"
)

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
	/*
		go UnixTCPsrv(btq)
		go func() {
			path := getCurrentPath()
			cmd := exec.Command("/usr/bin/python", fmt.Sprintf("%s/pysched/Scheduler.py", path))

			if err := cmd.Start(); err != nil {
				Logger.Error.Fatalln("run python error:", err.Error())
			}

			if err := cmd.Wait(); err != nil {
				Logger.Error.Println(err.Error(), cmd.Args)
			}
		}()*/
	go goPluginScheduler(btq)

	select {}
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
