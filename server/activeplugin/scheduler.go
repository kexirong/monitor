package activeplugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kexirong/monitor/common"
	"github.com/kexirong/monitor/common/packetparse"
)

type Tasker interface {
	Do() ([]packetparse.TargetPacket, error)
	AddJob(args ...interface{}) error
	Name() string
	DeleteJob(target string) error
}

//taskList  实现双向环形链表
type taskList struct {
	runing       bool
	nextTime     time.Time
	invl         time.Duration
	pNext, pPrev *taskList
	task         Tasker
}

/*
func (t *taskList) init() *taskList {
	t.pNext = t
	t.pPrev = t
	return t
}
*/

//insert 在t节点的前面插入
func (t *taskList) insert(node *taskList) *taskList {
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
	} else {
		t = nil
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
		next = t.pNext
		n++
	}
	return
}

//shift 使用findpos返回的值,当前节点移动n个位置插入taskList中
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
		for r := t.pNext; r != t; t = t.pNext {
			n++
		}
	}
	return n
}

//TaskScheduled   pTask 为双向环形链表，mutex保证安全操作链表
type TaskScheduled struct {
	event  chan common.Event //"method:pluginnam[|interval]"
	result chan common.Event
	pTask  *taskList
	mutex  *sync.Mutex
}

//New  retrun *TaskScheduled
func New() *TaskScheduled {
	ts := new(TaskScheduled)
	ts.result = make(chan common.Event, 1)
	ts.event = make(chan common.Event, 1)
	ts.mutex = new(sync.Mutex)
	return ts
}

//AddTask  为了不重复，添加前都尝试一次删除 ,interval单位s ,timeout 必须大于0，单位s
func (t *TaskScheduled) AddTask(tl *taskList) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.DeleteTask(tl.task.Name())
	t.pTask.insert(tl)
	t.pTask = t.nextTask()
}

//DeleteTask .
func (t *TaskScheduled) DeleteTask(name string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	len := t.Len()
	if t.Len() == 0 {
		return
	}
	cur := t.pTask
	for i := 0; i < len; i++ {
		if name != cur.name() {
			cur = cur.next()
			continue
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
func (t *TaskScheduled) nextTask() *taskList {
	n := t.pTask.findPos()
	next := t.pTask.next()
	if n > 0 {
		t.pTask.shift(n)
		return next
	}
	return t.pTask
}

//AddEventAndWaitResult add a event and wait eventDeal return result
func (t *TaskScheduled) AddEventAndWaitResult(event common.Event) common.Event {
	t.event <- event
	return <-t.result
}

func (t *TaskScheduled) eventDeal(event common.Event) {
	/*
		select {
		case <-r.result:
		default:
		}
	*/
	event.Result = "ok"
	switch event.Method {
	case "delete":
		if err := t.DeleteTask(event.Target); err != nil {
			event.Result = err.Error()
		}

	case "add":
		invl, err := strconv.Atoi(event.Args["interval"])
		if err != nil {
			event.Result = "Arg:" + err.Error()
			break
		}
		var timeout = 3
		if v, ok := event.Args["timeout"]; ok {
			timeout, err = strconv.Atoi(v)
			if err != nil {
				event.Result = "Arg:" + err.Error()
				break
			}
		}

		if err := t.AddTask(event.Target, invl, timeout); err != nil {
			event.Result = err.Error()
		}
	case "getlist":
		res := t.foreche()
		event.Result = res
	default:
		event.Result = "unknown operation type"
	}
	t.result <- event
	return

}

//WaitAndEventDeal 等待阻塞结束和时间
func (t *TaskScheduled) WaitAndEventDeal() {
	for {
		var wait = time.After(3 * time.Second)
		len := t.Len()
		if len != 0 {
			now := time.Now()
			if t.pTask.nextTime.Before(now) {
				return
			}
			wait = time.After(t.pTask.nextTime.Sub(now))
		}

		for {
			select {
			case <-wait:
				if len == 0 {
					break
				}
				return
			case e := <-t.event:
				t.eventDeal(e)
			}
		}
	}
}

//Scheduler must be  after  initialize true
func (t *TaskScheduled) Scheduler() ([]byte, error) {
	/*
		if t.Len() == 0 {
			return nil, errors.New("pluginEntry is empty")
		}
	*/
	t.mutex.Lock()
	defer t.mutex.Unlock()
	pe := t.pTask

	pe.nextTime = pe.nextTime.Add(pe.interval)

	if pe.runing {
		return nil, fmt.Errorf("%sis runing, may interval Too brief", pe.task.Name())
	}
	pe.runing = true

	t.fixOrder()

	return pe.task.Gather()
}

//Len return t.pTask.len()
func (t *TaskScheduled) Len() int {
	return t.pTask.len()
}

/*
func (r *ActivePlugin) insert(e *pluginEntry) {

	if r.len == 0 {
		r.curEntry = e
		r.curEntry.pNext = e
		r.curEntry.pPrev = e
	} else {
		p := r.curEntry.pNext
		e.pNext = p
		e.pPrev = r.curEntry
		r.curEntry.pNext = e
		p.pPrev = e
	}
	r.len++
}
*/
//调度后进行位置修正
func (t *TaskScheduled) fixOrder() {

	cur := t.pTask
	next := cur.pNext

	if cur.nextTime.Before(next.nextTime) {
		return
	}
	if t.pTask.len() < 3 {
		t.next()
		return
	}
	cur = t.pop()
	for i := 0; i < t.Len()-1; i++ {
		t.next()
		if cur.nextTime.Before(t.pTask.nextTime) {
			t.prev()
			break
		}
	}
	t.insert(cur)
	t.pTask = next
	return
}

func (t *TaskScheduled) next() {
	t.pTask = t.pTask.pNext

}
func (t *TaskScheduled) prev() {
	t.pTask = t.pTask.pPrev

}

func (t *TaskScheduled) genEntry(name string, interval int, timeout int) (*taskList, error) {
	var e taskList
	if interval < 1 || timeout < 1 {
		return nil, errors.New("arg false: interval or timeout le 0")
	}
	e.name = t.scriptPath + name
	if !common.CheckFileIsExist(pe.name) {
		return nil, errors.New("the script is not exist")
	}
	e.interval = time.Second * time.Duration(interval)
	e.timeout = time.Second * time.Duration(timeout)
	e.nextTime = time.Now().Add(pe.interval)

	return &e, nil
}
func (t *TaskScheduled) foreche() string {
	cur := r.curEntry
	var ret = `{"name":"%s","interval":"%v","nextime":"%s"}`
	var plugins []string
	for i := 0; i < r.len; i++ {
		plugins = append(plugins, fmt.Sprintf(ret,
			cur.name,
			cur.interval,
			cur.nextTime.Format("2006-01-02 15:04:05.000000")))
		cur = cur.pNext
	}
	return fmt.Sprintf("[%s]", strings.Join(plugins, ","))
}

func (t *TaskScheduled) CheckAndDownloads(url, filename string, check bool) error {
	if check && common.CheckFileIsExist(r.scriptPath+filename) {
		return nil
	}
	res, err := http.Get(url + filename)
	if err != nil {
		return err
	}
	robots, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}
	file, err := os.OpenFile(r.scriptPath+filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(robots)
	file.Sync()
	return nil
}
