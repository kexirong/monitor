package queue

import (
	"errors"
	"runtime"
	"sync/atomic"
	"time"
)

//Queue interface
type Queue interface {
	Len() int
	Put(interface{}) (bool, error)
	Get() (interface{}, bool, error)
	PutWait(interface{}, ...time.Duration) error
	GetWait(...time.Duration) (interface{}, error)
}

var ErrEmpty = errors.New("queue is empty")
var ErrFull = errors.New("queue is full")
var ErrTimeout = errors.New("queue op. timeout")
var ErrLost = errors.New("queue op. get token failed")

type data struct {
	//isFull bool
	stat  uint32 //isFull(读写指针相同时有bug) 改为 stat , 0(可写), 1(写入中), 2(可读)，3(读取中)
	value interface{}
}

// BytesQueue .queque is a []byte slice
type BytesQueue struct {
	cap    uint32 //队列容量
	len    uint32 //队列长度
	ptrStd uint32 //ptr基准(cap-1)
	putPtr uint32 // queue[putPtr].stat must < 2
	getPtr uint32 // queue[getPtr].stat may < 2
	queue  []data //队列
}

//NewBtsQueue cap转换为的2的n次幂-1数,建议直接传入2^n数
//如出现 error0: the get pointer excess roll 说明需要加大队列容量(其它的非nil error 说明有bug)
//非 阻塞型的 put get方法 竞争失败时会立即返回
func NewBtsQueue(cap uint32) *BytesQueue {
	bq := new(BytesQueue)
	bq.ptrStd = minCap(cap)
	bq.cap = bq.ptrStd + 1
	bq.queue = make([]data, bq.cap)
	return bq
}

func minCap(u uint32) uint32 { //溢出环形计算需要，得出一个2的n次幂减1数（具体可百度kfifo）
	u-- //兼容0, as min as ,128->127 !255
	u |= u >> 1
	u |= u >> 2
	u |= u >> 4
	u |= u >> 8
	u |= u >> 16
	return u
}

//Len method
func (bq *BytesQueue) Len() uint32 {
	return atomic.LoadUint32(&bq.len)
}

//Empty method
func (bq *BytesQueue) Empty() bool {
	if bq.Len() > 0 {
		return false
	}
	return true
}

//Put return value
func (bq *BytesQueue) Put(bs interface{}) (bool, error) {
	var putPtr, stat uint32
	var dt *data

	putPtr = atomic.LoadUint32(&bq.putPtr)

	if bq.Len() >= bq.ptrStd {
		return false, ErrFull
	}

	if !atomic.CompareAndSwapUint32(&bq.putPtr, putPtr, putPtr+1) {
		return false, ErrLost
	}
	atomic.AddUint32(&bq.len, 1)
	dt = &bq.queue[putPtr&bq.ptrStd]

	for {
		stat = atomic.LoadUint32(&dt.stat) & 3
		if stat == 0 {
			//可写

			atomic.AddUint32(&dt.stat, 1)
			dt.value = bs
			atomic.AddUint32(&dt.stat, 1)
			return true, nil
		}
		runtime.Gosched()
	}
}

//Get method
func (bq *BytesQueue) Get() (interface{}, bool, error) {
	var getPtr, stat uint32
	var dt *data
	var bs interface{} //中间变量，保障数据完整性

	getPtr = atomic.LoadUint32(&bq.getPtr)

	if bq.Len() < 1 {
		return nil, false, ErrEmpty
	}

	if !atomic.CompareAndSwapUint32(&bq.getPtr, getPtr, getPtr+1) {
		return nil, false, ErrLost
	}
	atomic.AddUint32(&bq.len, 4294967295) //bq.len - 1
	dt = &bq.queue[getPtr&bq.ptrStd]

	for {
		stat = atomic.LoadUint32(&dt.stat)
		if stat == 2 {
			//可读
			atomic.AddUint32(&dt.stat, 1) // change stat to 读取中
			bs = dt.value
			dt.value = nil
			atomic.StoreUint32(&dt.stat, 0) //重置stat为0
			return bs, true, nil
		}
		runtime.Gosched()
	}
}

// PutWait 阻塞型put,ms 最大等待豪秒数,默认 1000
func (bq *BytesQueue) PutWait(bs interface{}, ms ...time.Duration) error {
	var timeout <-chan time.Time
	var msto time.Duration
	if len(ms) > 0 {
		msto = time.Millisecond * ms[0]
	} else {
		msto = time.Millisecond * 1000
	}
	timeout = time.After(msto)
	//wait := time.NewTicker(100 * time.Millisecond)
	//defer wait.Stop()
	for {
		select {
		case <-timeout:
			return ErrTimeout
		default:
			if ok, err := bq.Put(bs); ok {
				return nil
			} else if err == ErrFull {
				//<-wait.C
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// GetWait 阻塞型get, ms为 等待毫秒 默认1000
func (bq *BytesQueue) GetWait(ms ...time.Duration) (interface{}, error) {
	var timeout <-chan time.Time
	var msto time.Duration

	if len(ms) > 0 {
		msto = time.Millisecond * ms[0]
	} else {
		msto = time.Millisecond * 1000
	}
	timeout = time.After(msto)

	//wait := time.NewTicker(100 * time.Millisecond)
	//defer wait.Stop()
	for {
		select {
		case <-timeout:
			return nil, ErrTimeout

		default:
			if value, ok, err := bq.Get(); ok {
				return value, nil
			} else if err == ErrEmpty {
				//<-wait.C
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
