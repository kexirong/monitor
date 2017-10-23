package queue

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Test_QueueGetAndPut(t *testing.T) {

	runtime.GOMAXPROCS(runtime.NumCPU())
	var waitGroup = new(sync.WaitGroup)
	var total int
	var perr, gerr uint32

	total = 100000
	bq := NewBtsQueue(4096)
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
	op := use / time.Duration(total*runtime.NumCPU())
	fmt.Printf(" Grp: %3d, Times: %10d, perr:%6v,gerr:%6v, use: %12v, %8v/opn",
		runtime.NumCPU(), total, perr, gerr, use, op)

}

func Test_QueueWaitGetAndPut(t *testing.T) {

	runtime.GOMAXPROCS(runtime.NumCPU())
	var waitGroup = new(sync.WaitGroup)
	var total int
	var perr, gerr uint32

	total = 10000
	bq := NewBtsQueue(10240)
	start := time.Now()
	for i := 0; i < runtime.NumCPU(); i++ {
		waitGroup.Add(1)
		go testQueuePutWait(bq, total, &perr, waitGroup)
		waitGroup.Add(1)
		go testQueueGetWait(bq, total, &gerr, waitGroup)

	}
	waitGroup.Wait()
	end := time.Now()
	use := end.Sub(start)
	total = total * runtime.NumCPU()
	op := use / time.Duration(total*runtime.NumCPU())
	fmt.Printf(" Grp: %3d, Times: %10d, perr:%6v,gerr:%6v, use: %12v, %8v/opn",
		runtime.NumCPU(), total, perr, gerr, use, op)

}

func testQueuePut(bq *BytesQueue, times int, errors *uint32, wait *sync.WaitGroup) {

	for i := 0; i < times; i++ {
		ok, err := bq.Put([]byte(fmt.Sprintf("element %d", i)))
		if !ok {
			atomic.AddUint32(errors, 1)
			if err != nil {
				fmt.Println(err.Error())
			}
		}

	}
	wait.Done()

}

func testQueueGet(bq *BytesQueue, times int, errors *uint32, wait *sync.WaitGroup) {

	for i := 0; i < times; i++ {
		vl, ok, err := bq.Get() //Wait(10)
		if !ok {
			atomic.AddUint32(errors, 1)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else if vl == nil {
			fmt.Print("fuck the value is nil")
		}

	}
	wait.Done()

}

func testQueuePutWait(bq *BytesQueue, times int, errors *uint32, wait *sync.WaitGroup) {

	for i := 0; i < times; i++ {
		err := bq.PutWait([]byte(fmt.Sprintf("element %d", i)), 100)
		if err != nil {
			atomic.AddUint32(errors, 1)

			fmt.Println(err.Error())

		}

	}
	wait.Done()

}

func testQueueGetWait(bq *BytesQueue, times int, errors *uint32, wait *sync.WaitGroup) {

	for i := 0; i < times; i++ {
		vl, err := bq.GetWait(100)
		if err != nil {
			atomic.AddUint32(errors, 1)

			if vl == nil {
				fmt.Println("fuck value is nil")
			}
		} else {
			fmt.Println(string(vl))
		}

	}
	wait.Done()

}
