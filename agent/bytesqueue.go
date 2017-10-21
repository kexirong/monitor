package agent

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

type data struct {
	//isFull bool
	stat  uint32 //isFull(读写指针相同时有bug) 改为 stat ,通过 &3 运算得到 0(可写), 1(写入中), 2(可读)，3(读取中)
	value []byte
}

// BytesQueue .queque is a []byte slice
type BytesQueue struct {
	cap    uint32 //队列容量
	ptrStd uint32 //ptr基准(cap-1)
	putPtr uint32 // queue[putPtr].stat must < 2
	getPtr uint32 // queue[putPtr].stat may < 2
	queue  []data //队列
}

//NewBtsQueue cap转换为的2的n次幂-1数,建议直接传入2^n数
//
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

func (bq *BytesQueue) len() (leng, getPtr, putPtr uint32) {

	getPtr = atomic.LoadUint32(&bq.getPtr)
	putPtr = atomic.LoadUint32(&bq.putPtr)
	if putPtr >= getPtr {
		return putPtr - getPtr, getPtr, putPtr
	}
	return bq.cap + (putPtr - getPtr), getPtr, putPtr

}

//Len method
func (bq *BytesQueue) Len() uint32 {
	l, _, _ := bq.len()
	return l
}

//Empty method
func (bq *BytesQueue) Empty() bool {
	if bq.Len() > 0 {
		return false
	}
	return true
}

//Put method
func (bq *BytesQueue) Put(bs []byte) (bool, error) {
	var leng, putPtr, stat uint32
	var dt *data
	leng, _, putPtr = bq.len()
	if leng >= bq.ptrStd {
		return false, nil
	}
	dt = &bq.queue[putPtr&bq.ptrStd]

	if !atomic.CompareAndSwapUint32(&bq.putPtr, putPtr, putPtr+1) {
		return false, nil
	}
	for {
		stat = atomic.LoadUint32(&dt.stat) & 3
		switch stat {
		case 0: //可写
			atomic.AddUint32(&dt.stat, 1)
			dt.value = bs
			atomic.AddUint32(&dt.stat, 1)
			return true, nil
		case 3: //读取中
			runtime.Gosched() //出让cpu
		default:
			return false, fmt.Errorf("error%v: the put pointer excess roll  :%v, %v", stat, putPtr, dt.value)
		}
	}

}

//Get method
func (bq *BytesQueue) Get() ([]byte, bool, error) {
	var leng, getPtr, stat uint32
	var dt *data
	var bs []byte //中间变量，保障数据完整性
	leng, getPtr, _ = bq.len()
	if leng < 1 {
		return nil, false, nil
	}

	dt = &bq.queue[getPtr&bq.ptrStd]
	if !atomic.CompareAndSwapUint32(&bq.getPtr, getPtr, getPtr+1) {
		return nil, false, nil
	}

	for {
		stat = atomic.LoadUint32(&dt.stat) & 3
		switch stat {
		case 2: //可读
			atomic.AddUint32(&dt.stat, 1)
			bs = dt.value
			dt.value = nil
			atomic.AddUint32(&dt.stat, 1)
			return bs, true, nil
		case 1: //写入中
			runtime.Gosched()
		default:
			return nil, false, fmt.Errorf("error%v: the get pointer excess roll  :%v, %v", stat, getPtr, dt.value)
		}
	}

}

// PutWait 阻塞型put,ms 最大等待豪秒数
func (bq *BytesQueue) PutWait(bs []byte, ms ...int) error {
	var ok bool
	var i = 50

	if len(ms) > 0 {
		i = ms[0] / 10
	}

	for ; i > 0; i-- {
		ok, _ = bq.Put(bs)
		if ok {
			return nil
		}
		time.Sleep(time.Millisecond * 10)
	}
	return fmt.Errorf("time out")
}

// GetWait 阻塞型get, ms为 最大等待毫秒数
func (bq *BytesQueue) GetWait(ms ...int) ([]byte, error) {
	var i = 50

	if len(ms) > 0 {
		i = ms[0] / 10
	}

	for ; i > 0; i-- {
		value, ok, _ := bq.Get()
		if ok {
			return value, nil
		}

		time.Sleep(time.Millisecond * 10)

	}
	return nil, fmt.Errorf("time out")

}
