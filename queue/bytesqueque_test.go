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

	runtime.GOMAXPROCS(2)

	//cnt := make([]map[string]int, 4)

	var waitGroup = new(sync.WaitGroup)
	var total int
	var perr, gerr uint32

	total = 10000
	bq := NewBtsQueue(40960)

	start := time.Now()

	waitGroup.Add(1)
	go testQueuePut(bq, total, &perr, waitGroup)
	waitGroup.Add(1)
	go testQueueGet(bq, total, &gerr, waitGroup)

	waitGroup.Wait()
	end := time.Now()
	use := end.Sub(start)

	op := use / time.Duration(total)
	fmt.Printf(" Grp: %3d, Times: %10d, perr:%6v,gerr:%6v, use: %12v, %8v/opn",
		runtime.NumCPU(), total, perr, gerr, use, op)

}

func Test_QueueWaitGetAndPut(t *testing.T) {
	num := runtime.NumCPU()
	runtime.GOMAXPROCS(num)

	cnt := make([]map[string]int, 4)
	var waitGroup = new(sync.WaitGroup)
	var total int
	var perr, gerr uint32

	total = 10000
	bq := NewBtsQueue(10240)
	start := time.Now()
	for i := 0; i < num; i++ {

		waitGroup.Add(1)
		go testQueuePutWait(bq, total, &perr, waitGroup)
		waitGroup.Add(1)
		go testQueueGetWait(bq, total, &gerr, waitGroup, cnt, i)

	}
	waitGroup.Wait()
	end := time.Now()
	use := end.Sub(start)
	total = total * num
	op := use / time.Duration(total*num)
	//fmt.Println(cnt)
	fmt.Printf(" Grp: %3d, Times: %10d, perr:%6v,gerr:%6v, use: %12v, %8v/opn",
		num, total, perr, gerr, use, op)
	conter := make(map[string]int)
	for _, v := range cnt {
		for k, l := range v {
			conter[k] += l

		}
	}
	ii := 0
	for k, v := range conter {
		ii += v
		//	fmt.Printf("%v----%v\n", k, v)
		if v != 4 {
			fmt.Println("errorsss:", k)
		}

	}
	fmt.Println(ii)
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
				fmt.Println(err.Error(), i)
			}
		} else if vl == nil {
			fmt.Print("fuck the value is nil")
		}

	}
	wait.Done()

}

func testQueuePutWait(bq *BytesQueue, times int, errors *uint32, wait *sync.WaitGroup) {

	for i := 0; i < times; i++ {
		err := bq.PutWait([]byte(fmt.Sprintf("element %d", i)), 500)
		if err != nil {
			atomic.AddUint32(errors, 1)

			fmt.Println(err.Error(), i)

		}

	}
	wait.Done()

}

func testQueueGetWait(bq *BytesQueue, times int, errors *uint32, wait *sync.WaitGroup, cnt []map[string]int, i int) {
	nt := make(map[string]int)
	for i := 0; i < times; i++ {
		vl, err := bq.GetWait(500)
		if err != nil {
			fmt.Println(err.Error(), i)
			atomic.AddUint32(errors, 1)

		} else {
			if vl == nil {
				fmt.Println("fuck value is nil")
				continue
			}
			nt[string(vl)]++
		}

	}
	cnt[i] = nt
	wait.Done()

}
