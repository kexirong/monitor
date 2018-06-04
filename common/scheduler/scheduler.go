package scheduler

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Tasker interface {
	Do() ([]byte, error)
	AddJob(param ...interface{}) error
	Name() string
	DeleteJob(param ...interface{}) error
}

//taskList  实现双向环形链表
type taskList struct {
	runing       bool
	nextTime     time.Time
	invl         time.Duration
	pNext, pPrev *taskList
	task         Tasker
}

//insert 在t节点的前面插入
func (t *taskList) insert(node *taskList) *taskList {
	if node == nil {
		return t
	}
	if t == nil {
		node.pPrev = node
		node.pNext = node
	} else {
		node.pNext = t
		node.pPrev = t.pPrev
		t.pPrev.pNext = node
		t.pPrev = node
	}
	return node
}

//pop  将下一个节点从taskList中去除，并返回此节点，此方法不存在len为0的调用场景
func (t *taskList) popNext() *taskList {
	next := t.pNext
	if t.len() > 1 {
		t.pNext = next.pNext
		t.pNext.pPrev = t
	}
	next.pNext = nil
	next.pPrev = nil
	return next
}

//findPos 找到合适的相对位置,n >= 0
func (t *taskList) findPos() (n int) {
	next := t.pNext
	len := t.len()
	for i := 0; i < len-1; i++ {
		if t.nextTime.Before(next.nextTime) {
			break
		}
		next = next.pNext
		n++
	}
	return
}

//shift 使用findpos返回的值,当前节点移动n个位置插入taskList中,t不能==nil
func (t *taskList) shift(n int) {
	next := t.pNext
	cur := t.pPrev.popNext()
	for i := 0; i < n; i++ {
		next = next.pNext
	}
	next.insert(cur)
}

func (t *taskList) isRuning() bool {
	return t.runing
}
func (t *taskList) do() ([]byte, error) {
	t.runing = true
	defer func() {
		t.runing = false
	}()
	return t.task.Do()
}

func (t *taskList) interval() time.Duration {
	return t.invl
}

func (t *taskList) next() *taskList {
	return t.pNext
}

func (t *taskList) perv() *taskList {
	return t.pPrev
}

func (t *taskList) name() string {
	return t.task.Name()
}

func (t *taskList) len() int {
	n := 0
	if t != nil {
		n = 1
		for r := t.pNext; r != t; r = r.pNext {
			n++
		}
	}
	return n
}

//TaskScheduled   pTask 为双向环形链表，mutex保证安全操作链表
type TaskScheduled struct {
	pTask  *taskList
	mutex  *sync.Mutex
	waiter chan time.Time
}

//New  retrun *TaskScheduled
func New() *TaskScheduled {
	ts := new(TaskScheduled)
	ts.mutex = new(sync.Mutex)
	return ts
}

//AddTask  为了不重复，添加前都尝试一次删除 ,interval 要大于 1秒
func (t *TaskScheduled) AddTask(interval time.Duration, task Tasker) {
	t.DeleteTask(task.Name())
	if interval < time.Second {
		interval = time.Second * 10
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.pTask = t.pTask.insert(genNode(interval, task))
	t.nextTask()
}

func genNode(interval time.Duration, task Tasker) *taskList {
	var tl = new(taskList)
	tl.invl = interval
	tl.nextTime = time.Now().Add(interval)
	tl.task = task
	return tl
}

//DeleteTask .
func (t *TaskScheduled) DeleteTask(name string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.pTask == nil {
		return
	}
	len := t.Len()
	cur := t.pTask
	for i := 0; i < len; i++ {
		if name != cur.name() {
			cur = cur.next()
			continue
		}
		if len == 1 {
			t.pTask = nil
			return
		}

		if cur == t.pTask {
			t.pTask = t.pTask.next()
		}
		cur = cur.perv()
		cur.popNext()

		return
	}
}

//nextTask 不使用锁,应在AddTask,do中调用，否则非安全的
func (t *TaskScheduled) nextTask() {
	n := t.pTask.findPos()
	next := t.pTask.next()
	if n > 0 {
		t.pTask.shift(n)
		t.pTask = next
	}

}

func (t *TaskScheduled) AddJob(taskname string, param ...interface{}) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.pTask == nil {
		return errors.New("Task List is nil")
	}
	len := t.Len()
	cur := t.pTask
	for i := 0; i < len; i++ {
		if taskname == cur.name() {
			return cur.task.AddJob(param)
		}
		cur = cur.next()
	}
	return errors.New(taskname + " is not exist")
}

func (t *TaskScheduled) DeleteJob(taskname string, param ...interface{}) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.pTask == nil {
		return errors.New("Task List is nil")
	}
	len := t.Len()
	cur := t.pTask
	for i := 0; i < len; i++ {
		if taskname == cur.name() {
			return cur.task.DeleteJob(param)
		}
		cur = cur.next()
	}
	return errors.New(taskname + " is not exist")
}

func (t *TaskScheduled) scheduled() (*taskList, error) {
	t.mutex.Lock()
	pe := t.pTask
	if pe.isRuning() {
		t.mutex.Unlock()
		return nil, fmt.Errorf("%sis runing, may interval Too brief", pe.task.Name())
	}
	now := time.Now()
	wait := time.After(pe.nextTime.Sub(now))
	pe.nextTime = pe.nextTime.Add(pe.invl)
	t.nextTask()
	t.mutex.Unlock()
	<-wait
	return pe, nil
}

type callback func([]byte, error)

func (t *TaskScheduled) Star(callback callback) {
	for {
		if t.pTask == nil {
			time.Sleep(512 * time.Millisecond)
			continue
		}
		pe, err := t.scheduled()
		if err != nil {
			go callback(nil, err)
			continue
		}

		go callback(pe.do())
	}
}

//Len return t.pTask.len()
func (t *TaskScheduled) Len() int {
	return t.pTask.len()
}

func (t *TaskScheduled) next() {
	t.pTask = t.pTask.pNext

}
func (t *TaskScheduled) prev() {
	t.pTask = t.pTask.pPrev

}

func (t *TaskScheduled) EcheTaskList() string {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	cur := t.pTask
	var ret = `{"name":"%s","interval":"%v","nextime":"%s"}`
	var plugins []string
	len := t.Len()
	for i := 0; i < len; i++ {
		plugins = append(plugins, fmt.Sprintf(ret,
			cur.name(),
			cur.interval(),
			cur.nextTime.Format("2006-01-02 15:04:05.000000")))
		cur = cur.pNext
	}
	return fmt.Sprintf("[%s]", strings.Join(plugins, ","))
}
