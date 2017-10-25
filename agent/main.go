package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/kexirong/monitor/queue"
)

func main() {

	var addr = flag.String("s", "ip:port", "server addrs multi delimit  with ',' ")
	flag.Parse()

	if *addr == "ip:port" {
		Logger.Error.Fatalln("no servers;-h;exit")
	}

	servers := strings.Split(*addr, ",")

	btq := queue.NewBtsQueue(4096)
	var waitGroup = new(sync.WaitGroup)
	runtime.GOMAXPROCS(runtime.NumCPU())
	waitGroup.Add(1)
	go UnixTCPsrv(btq)
	go sendStart(servers, btq)
	go func() {
		path := getCurrentPath()
		cmd := exec.Command("/usr/bin/python", fmt.Sprintf("%s/pysched/Scheduler.py", path))

		if err := cmd.Start(); err != nil {

			Logger.Error.Fatalln("xxxxxxxxx:", err.Error())
		}

		if err := cmd.Wait(); err != nil {

			Logger.Error.Println(err.Error(), cmd.Args)
		}
	}()

	waitGroup.Wait()

}

func getCurrentPath() string {

	path, err := os.Getwd()
	checkErr(err)
	/*
		i := strings.LastIndex(path, "/")
		if i < 0 {
			Logger.Error.Fatalln("get the path error")
		}
		path = string(path[0:i])
	*/
	return path
}

func checkErr(err error) {
	if err != nil {
		Logger.Error.Panicf("error: %s, exit!", err.Error())
		os.Exit(1)
	}
}
