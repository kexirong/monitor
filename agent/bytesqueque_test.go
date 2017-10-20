package agent

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func testQueuePut(bq *BytesQueue, times int, errors *uint32, wait *sync.WaitGroup) {

	for i := 0; i < times; i++ {
		ok, err := bq.Put([]byte(fmt.Sprintf("element %d", i)))
		if !ok {
			atomic.AddUint32(errors, 1)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		//r := rand.Intn(100)
		//	time.Sleep(time.Millisecond * time.Duration(r))
	}
	wait.Done()

}

func testQueueGet(bq *BytesQueue, times int, errors *uint32, wait *sync.WaitGroup) {

	for i := 0; i < times; i++ {
		vl, ok := bq.Get()
		if !ok {
			atomic.AddUint32(errors, 1)

		} else {
			fmt.Println(string(vl))
		}
		//r := rand.Intn(100)
		//time.Sleep(time.Millisecond * time.Duration(r))
	}
	wait.Done()

}
func Test_QueueGetAndPut(t *testing.T) {

	runtime.GOMAXPROCS(runtime.NumCPU())
	var waitGroup = new(sync.WaitGroup)
	var total int
	var perr, gerr uint32

	total = 1000
	bq := NewBtsQueue(10240)
	start := time.Now()
	for i := 0; i < runtime.NumCPU(); i++ {
		waitGroup.Add(1)
		go testQueuePut(bq, total, &perr, waitGroup)
		waitGroup.Add(1)
		go testQueueGet(bq, total, &gerr, waitGroup)

	}
	waitGroup.Wait()
	end := time.Now()
	use := end.Sub(start)
	total = total * runtime.NumCPU()
	op := use / time.Duration(total)
	fmt.Printf(" Grp: %3d, Times: %10d, perr:%6v,gerr:%6v, use: %12v, %8v/opn",
		runtime.NumCPU(), total, perr, gerr, use, op)

}
