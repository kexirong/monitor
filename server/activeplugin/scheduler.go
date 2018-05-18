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

type Plugin interface {
	Gather() ([]packetparse.TargetPacket, error)
	AddJob(args ...interface{}) error
	Name() string
	DeleteJob(target string) error
}

//pluginEntry  is  a ring
type pluginEntry struct {
	runing       bool
	nextTime     time.Time
	interval     time.Duration
	pNext, pPrev *pluginEntry
	plugin       Plugin
}

func (p *pluginEntry) insert(e *pluginEntry) {
	//	n := p.pNext
	//	p.pNext = e
	//	e.pNext = n
	//	e.pPrev = p
	//	n.pPrev = e
	e.pNext = p.pNext
	p.pNext.pPrev = e
	p.pNext = e
	e.pPrev = p
}

//pop  while len==0  will return  nil
func (p *pluginEntry) pop() *pluginEntry {
	next := p
	if p.len() > 1 {
		p.pPrev.pNext = p.pNext
		p.pNext.pPrev = p.pPrev
	} else {
		p = nil
	}
	return next
}

func (p *pluginEntry) shift(n int) {
	cur := p.pop()

	if cur.nextTime.Before(next.nextTime) {
		return
	}
	if r.len < 3 {
		r.next()
		return
	}
	cur = r.pop()
	for i := 0; i < r.len-1; i++ {
		r.next()
		if cur.nextTime.Before(r.curEntry.nextTime) {
			r.prev()
			break
		}
	}
	r.insert(cur)
	r.curEntry = next
	return
}

func (p *pluginEntry) isRuning() bool {
	return p.runing
}

func (p *pluginEntry) Interval() time.Duration {
	return p.interval
}

func (p *pluginEntry) next() *pluginEntry {
	return p.pNext
}

func (p *pluginEntry) perv() *pluginEntry {
	return p.pPrev
}

/*
func (p *pluginEntry) init() *pluginEntry {
	p.pNext = p
	p.pPrev = p
	return p
}
*/
func (p *pluginEntry) len() int {
	n := 0
	if p != nil {
		n = 1
		for r := p.pNext; r != p; p = p.pNext {
			n++
		}
	}
	return n
}

//ActivePlugin   内嵌环形链表，不支持并发操作环形链表
type ActivePlugin struct {
	event      chan common.Event //"method:pluginnam[|interval]"
	result     chan common.Event
	pluginRing *pluginEntry
	mutex      *sync.Mutex
}

//Initialize  初始化, 等同New
func Initialize() (*ActivePlugin, error) {
	ring := new(ActivePlugin)
	ring.result = make(chan common.Event, 1)
	ring.event = make(chan common.Event, 1)
	ring.mutex = new(sync.Mutex)
	return ring, nil
}

//AddEventAndWaitResult add a event and wait eventDeal return result
func (r *ActivePlugin) AddEventAndWaitResult(event common.Event) common.Event {
	r.event <- event
	return <-r.result
}

func (r *ActivePlugin) eventDeal(event common.Event) {
	/*
		select {
		case <-r.result:
		default:
		}
	*/
	event.Result = "ok"
	switch event.Method {
	case "delete":
		if err := r.DeleteTask(event.Target); err != nil {
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

		if err := r.AddTask(event.Target, invl, timeout); err != nil {
			event.Result = err.Error()
		}
	case "getlist":
		res := r.foreche()
		event.Result = res
	default:
		event.Result = "unknown operation type"
	}
	r.result <- event
	return

}

//WaitAndEventDeal 等待阻塞结束和时间
func (r *ActivePlugin) WaitAndEventDeal() {
	for {
		var wait = time.After(3 * time.Second)
		len := r.Len()
		if len != 0 {
			now := time.Now()
			if r.pluginRing.nextTime.Before(now) {
				return
			}
			wait = time.After(r.pluginRing.nextTime.Sub(now))
		}

		for {
			select {
			case <-wait:
				if len == 0 {
					break
				}
				return
			case e := <-r.event:
				r.eventDeal(e)
			}
		}
	}
}

//Scheduler must be  after  initialize true
func (r *ActivePlugin) Scheduler() ([]byte, error) {
	/*
		if r.len == 0 {
			return nil, errors.New("pluginEntry is empty")
		}
	*/
	r.mutex.Lock()
	defer r.mutex.Unlock()
	pe := r.pluginRing

	pe.nextTime = pe.nextTime.Add(pe.interval)

	if pe.runing {
		return nil, fmt.Errorf("%sis runing, may interval Too brief", pe.plugin.Name())
	}
	pe.runing = true

	r.fixOrder()

	return pe.plugin.Gather()
}

//DeleteTask arg PythonModuleName ,as stop the module
func (r *ActivePlugin) DeleteTask(name string) error {
	cur := r.pluginRing
	len := cur.Len()
	for i := 0; i < len; i++ {
		if name == r.pluginRing.plugin.Name() {
			r.pop()
			if i != 0 {
				r.pluginRing = cur
			}
			return nil
		}
		r.next()
	}
	r.pluginRing = cur
	return errors.New("not exist")
}

/*
func (r *ActivePlugin) pop() *pluginEntry {
	cur := r.curEntry
	r.next()
	if r.len > 1 {
		cur.pPrev.pNext = cur.pNext
		cur.pNext.pPrev = cur.pPrev
	}
	r.len--
	return cur
}
*/

//Len return r.len
func (r *ActivePlugin) Len() int {
	return r.pluginRing.len()
}

//InsertEntry  为了不重复，插入前都尝试一次删除 ,interval单位s ,timeout 必须大于0，单位s
func (r *ActivePlugin) AddTask(name string, interval int, timeout int) error {
	r.DeleteTask(name)
	node, err := r.genEntry(name, interval, timeout)
	if err != nil {
		return err
	}

	if r.Len() > 1 {
		r.prev()
	}
	r.insert(node)
	r.next()
	r.fixOrder()
	return nil
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
func (r *ActivePlugin) fixOrder() {
	select {
	case <-r.plug:
	default:
	}
	defer func() {
		r.plug <- struct{}{}
	}()
	cur := r.curEntry
	next := cur.pNext

	if cur.nextTime.Before(next.nextTime) {
		return
	}
	if r.len < 3 {
		r.next()
		return
	}
	cur = r.pop()
	for i := 0; i < r.len-1; i++ {
		r.next()
		if cur.nextTime.Before(r.curEntry.nextTime) {
			r.prev()
			break
		}
	}
	r.insert(cur)
	r.curEntry = next
	return
}

func (r *ActivePlugin) next() {
	r.curEntry = r.curEntry.pNext

}
func (r *ActivePlugin) prev() {
	r.curEntry = r.curEntry.pPrev

}

func (r *ActivePlugin) genEntry(name string, interval int, timeout int) (*pluginEntry, error) {
	var pe pluginEntry
	if interval < 1 || timeout < 1 {
		return nil, errors.New("arg false: interval or timeout le 0")
	}
	pe.name = r.scriptPath + name
	if !common.CheckFileIsExist(pe.name) {
		return nil, errors.New("the script is not exist")
	}
	pe.interval = time.Second * time.Duration(interval)
	pe.timeout = time.Second * time.Duration(timeout)
	pe.nextTime = time.Now().Add(pe.interval)

	return &pe, nil
}
func (r *ActivePlugin) foreche() string {
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

func (r *ActivePlugin) CheckDownloads(url, filename string, check bool) error {
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
