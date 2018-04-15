package pyplugin

import (
	"fmt"
	"testing"
	"time"
)

func Test_pyplugin(t *testing.T) {
	pp, err := Initialize("")
	if err != nil {

		panic(err)
	}
	err = pp.InsertEntry("cpu", 1)
	if err != nil {
		fmt.Println(err)
	}

	err = pp.InsertEntry("cpus", 2)
	if err != nil {
		fmt.Println(err)
	}
	err = pp.InsertEntry("cpus1", 5)
	if err != nil {
		fmt.Println(err)
	}
	err = pp.InsertEntry("cpus2", 3)
	if err != nil {
		fmt.Println(err)
	}
	err = pp.InsertEntry("cpus3", 1)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 5; i++ {
		err := pp.WaitAndEventDeal()
		fmt.Println("########################################################################")
		for i := 0; i < pp.len; i++ {
			cur := pp.curEntry
			fmt.Printf("nextTime:%v, cur.name:%s, cur.interval:%v \n", cur.nextTime, cur.name, cur.interval)
			cur = cur.pNext
		}
		fmt.Println("########################################################################")
		if len(err) > 0 {
			fmt.Println(err)
		}

		go func() {
			fmt.Println(pp.Scheduler())
		}()
	}
	<-time.After(time.Second * 10)
}
