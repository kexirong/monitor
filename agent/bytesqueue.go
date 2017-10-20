package agent

import (
	"fmt"
	"sync/atomic"
	"time"
)

type data struct {
	isNull bool
	value  []byte
}

// BytesQueue is a []byte slice
type BytesQueue struct {
	cap    uint32 //队列容量
	ptrStd uint32 //ptr基准(cap-1)
	putPtr uint32 // queue[putPtr].isNull must true
	getPtr uint32 // queue[putPtr].isNull may true
	queue  []data //队列
}

//NewBtsQueue cap转换为的2的n次幂-1数,建议直接传入2^n数
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

//Len return  has be used cap
func (bq *BytesQueue) Len() uint32 {
	l, _, _ := bq.len()
	return l
}

//Empty return queue was empty
func (bq *BytesQueue) Empty() bool {
	if bq.Len() > 0 {
		return false
	}
	return true
}

//Put is put value in queue
func (bq *BytesQueue) Put(bs []byte) (bool, error) {
	var leng, putPtr uint32
	var dt *data
	leng, _, putPtr = bq.len()
	if leng >= bq.ptrStd {
		return false, nil
	}
	dt = &bq.queue[putPtr&bq.ptrStd]

	if !atomic.CompareAndSwapUint32(&bq.putPtr, putPtr, putPtr+1) {
		return false, nil
	}
	if dt.isNull {
		dt.isNull = false
		dt.value = bs
		return true, nil
	}
	return false, fmt.Errorf("has happened logic error")

}

//Get is get value for queue
func (bq *BytesQueue) Get() ([]byte, bool) {
	var leng, getPtr uint32
	var dt *data
	leng, getPtr, _ = bq.len()
	if leng < 1 {
		return nil, false
	}

	dt = &bq.queue[getPtr&bq.ptrStd]
	if !atomic.CompareAndSwapUint32(&bq.getPtr, getPtr, getPtr+1) {
		return nil, false
	}

	if !dt.isNull {
		dt.isNull = true
		return dt.value, true
	}

	return nil, false

}

// PutWait 阻塞型put,sec 最大等待秒数
func (bq *BytesQueue) PutWait(bs []byte, sec ...int) error {
	var ok bool
	var i = 30

	if len(sec) > 0 {
		i = sec[0] * 10
	}

	for i = i * 10; i > 0; i-- {
		ok, _ = bq.Put(bs)
		if ok {
			return nil
		}
		time.Sleep(time.Millisecond * 100)
	}
	return fmt.Errorf("time out")
}

// GetWait 阻塞型get, sec为 最大等待秒数
func (bq *BytesQueue) GetWait(sec ...int) ([]byte, error) {
	var i = 30

	if len(sec) > 0 {
		i = sec[0] * 10
	}

	for i = i * 10; i > 0; i-- {
		value, ok := bq.Get()
		if ok {
			return value, nil
		}
		time.Sleep(time.Millisecond * 100)
	}
	return nil, fmt.Errorf("time out")

}
