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
	if *addr == "ip:port" {
		Logger.Error.Fatalln("no servers;-h;exit")
	}

	servers := strings.Split(*addr, ",")

	btq := queue.NewBtsQueue(4096)
	var waitGroup = new(sync.WaitGroup)
	runtime.GOMAXPROCS(runtime.NumCPU())
	waitGroup.Add(2)
	go UnixTCPsrv(btq)
	go sendStart(servers, btq)

	cmd := exec.Command("python ./pysched/Scheduler.py > Scheduler.log ")
	cmd.Start()

	waitGroup.Wait()

}
