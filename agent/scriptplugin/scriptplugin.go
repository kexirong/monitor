package scriptplugin

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kexirong/monitor/common"
)

type pluginEntry struct {
	runing   bool
	nextTime time.Time
	interval time.Duration
	timeout  time.Duration
	name     string
	pNext    *pluginEntry
	pPrev    *pluginEntry
}

func (p *pluginEntry) done() {
	p.runing = false
}

func (p *pluginEntry) run() ([]byte, error) {
	if p.runing {
		return nil, errors.New(p.name + "is runing, may interval Too brief")
	}
	p.runing = true
	defer p.done()
	return common.Command(p.name, p.timeout)

}

//ScriptPlugin   内嵌环形链表，不支持并发操作环形链表
type ScriptPlugin struct {
	initialize bool
	len        int
	timer      int64
	result     chan common.Event
	event      chan common.Event //"method:pluginnam[|interval]"
	curEntry   *pluginEntry
	scriptPath string
	plug       chan struct{}
}

//Initialize  初始化, 等同New
func Initialize(scriptPath string) (*ScriptPlugin, error) {

	ring := new(ScriptPlugin)
	ring.result = make(chan common.Event, 1)
	ring.event = make(chan common.Event, 1)
	if !common.CheckFileIsExist(scriptPath) {
		return nil, errors.New("scriptPath not IsExist")
	}
	ring.scriptPath = scriptPath
	if scriptPath[len(scriptPath)-1] != '/' {
		ring.scriptPath = scriptPath + "/"
	}
	ring.plug = make(chan struct{}, 1)
	ring.initialize = true

	return ring, nil
}

//AddEventAndWaitResult add a event and wait eventDeal return result
func (r *ScriptPlugin) AddEventAndWaitResult(event common.Event) common.Event {
	r.event <- event
	return <-r.result
}

func (r *ScriptPlugin) eventDeal(event common.Event) {
	/*
		select {
		case <-r.result:
		default:
		}
	*/
	event.Result = "ok"
	switch event.Method {
	case "delete":
		if err := r.DeleteEntry(event.Target); err != nil {
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

		if err := r.InsertEntry(event.Target, invl, timeout); err != nil {
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
func (r *ScriptPlugin) WaitAndEventDeal() {
	for {
		var wait = time.After(3 * time.Second)
		len := r.Len()
		if len != 0 {
			now := time.Now()
			if r.curEntry.nextTime.Before(now) {
				return
			}
			wait = time.After(r.curEntry.nextTime.Sub(now))
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
func (r *ScriptPlugin) Scheduler() ([]byte, error) {
	/*
		if r.len == 0 {
			return nil, errors.New("pluginEntry is empty")
		}
	*/
	<-r.plug
	pn := r.curEntry
	pn.nextTime = pn.nextTime.Add(pn.interval)
	r.fixOrder()

	return pn.run()
}

//DeleteEntry arg PythonModuleName ,as stop the module
func (r *ScriptPlugin) DeleteEntry(name string) error {
	cur := r.curEntry
	for i := 0; i < r.len; i++ {
		if name == r.curEntry.name {
			r.pop()
			if i != 0 {
				r.curEntry = cur
			}
			return nil
		}
		r.next()
	}
	r.curEntry = cur
	return errors.New("not exist")
}

func (r *ScriptPlugin) pop() *pluginEntry {
	cur := r.curEntry
	r.next()
	if r.len > 1 {
		cur.pPrev.pNext = cur.pNext
		cur.pNext.pPrev = cur.pPrev
	}
	r.len--
	return cur
}

//Len return r.len
func (r *ScriptPlugin) Len() int {
	return r.len
}

//InsertEntry  为了不重复，插入前都尝试一次删除 ,interval单位s ,timeout 必须大于0，单位s
func (r *ScriptPlugin) InsertEntry(name string, interval int, timeout int) error {
	r.DeleteEntry(name)
	node, err := r.genEntry(name, interval, timeout)
	if err != nil {
		return err
	}

	if r.len > 1 {
		r.prev()
	}
	r.insert(node)
	r.next()
	r.fixOrder()
	return nil
}

func (r *ScriptPlugin) insert(e *pluginEntry) {

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

//调度后进行位置修正
func (r *ScriptPlugin) fixOrder() {
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

func (r *ScriptPlugin) next() {
	r.curEntry = r.curEntry.pNext

}
func (r *ScriptPlugin) prev() {
	r.curEntry = r.curEntry.pPrev

}

func (r *ScriptPlugin) genEntry(name string, interval int, timeout int) (*pluginEntry, error) {
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
func (r *ScriptPlugin) foreche() string {
	cur := r.curEntry
	var ret = "["
	var plugins = `{"name":"%s","interval":%v,"nextime":"%s"},`
	for i := 0; i < r.len; i++ {
		ret += fmt.Sprintf(plugins, cur.name, cur.interval, cur.nextTime.Format("2006-01-02 15:04:05.000000"))
		cur = cur.pNext
	}

	ret = strings.TrimSuffix(ret, ",") + "]"
	return ret
}
