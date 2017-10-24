package main

import (
	"flag"
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

		cmd := exec.Command("/usr/bin/python", "/var/golang/src/opsAPI/agent/pysched/Scheduler.py")

		if err := cmd.Start(); err != nil {

			Logger.Error.Fatalln("xxxxxxxxx:", err.Error())
		}

		if err := cmd.Wait(); err != nil {

			Logger.Error.Println(err.Error())
		}
	}()

	waitGroup.Wait()

}
