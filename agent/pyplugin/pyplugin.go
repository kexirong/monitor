package pyplugin

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kexirong/monitor/common"
	python "github.com/sbinet/go-python"
)

// 由于Python 的GLI机制CGO并发调度性能不如线性调度
// 由于并发问题多，Python脚本不采用这种调用
type pluginEntry struct {
	runing   bool
	nextTime time.Time
	interval time.Duration
	name     string
	plugin   *python.PyObject
	pNext    *pluginEntry
	pPrev    *pluginEntry
}

func (p *pluginEntry) done() {
	p.runing = false
}
func (p *pluginEntry) run() (string, error) {
	if p.runing {
		return "", errors.New(p.name + "is runing, may interval Too brief")
	}
	p.runing = true
	defer p.done()

	ret := p.plugin.CallFunction()
	if ret == nil {
		_, value, _ := python.PyErr_Fetch()
		return "", errors.New(python.PyString_AS_STRING(value.Str()))
	}

	if !python.PyString_Check(ret) {
		return "", errors.New("PyString_Check(ret)==false")
	}

	return python.PyString_AS_STRING(ret), nil
}

//PythonPlugin   内嵌环形链表，不支持并发操作环形链表
type PythonPlugin struct {
	initialize bool
	timer      int64
	result     chan common.Event
	len        int
	event      chan common.Event //"method:pluginnam[|interval]"
	curEntry   *pluginEntry
}

//Initialize  初始化, 等同New
func Initialize(pluginPath string) (*PythonPlugin, error) {
	err := python.Initialize()
	if err != nil {
		return nil, err
	}

	sysModule := python.PyImport_ImportModule("sys")
	path := sysModule.GetAttrString("path")

	python.PyList_Insert(path, 0, python.PyString_FromString(pluginPath))
	if python.PyErr_ExceptionMatches(python.PyExc_Exception) {
		_, value, _ := python.PyErr_Fetch()
		return nil, errors.New(python.PyString_AS_STRING(value.Str()))
	}

	ring := new(PythonPlugin)
	ring.result = make(chan common.Event, 1)
	ring.event = make(chan common.Event, 1)
	ring.initialize = true

	return ring, nil
}

//AddEventAndWaitResult add a event and wait eventDeal return result
func (r *PythonPlugin) AddEventAndWaitResult(event common.Event) common.Event {
	r.event <- event
	return <-r.result

}

func (r *PythonPlugin) eventDeal(event common.Event) {

	select {
	case <-r.result:
		fmt.Println("<-r.result ERROR!!!!")
	default:
	}

	event.Result = "ok"
	switch event.Method {
	case "delete":
		if err := r.DeleteEntry(event.Target); err != nil {
			event.Result = err.Error()
		}

	case "add":
		arg1 := event.Args["interval"]
		invl, err := strconv.Atoi(arg1)
		if err != nil {
			event.Result = "Arg:" + err.Error()
			break
		}
		if err := r.InsertEntry(event.Target, invl); err != nil {
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
func (r *PythonPlugin) WaitAndEventDeal() {

	wait := time.After(3 * time.Second)
	if r.len != 0 {
		now := time.Now()

		if r.curEntry.nextTime.Before(now) {
			return
		}
		wait = time.After(r.curEntry.nextTime.Sub(now))
	}
	for {
		select {
		case <-wait:
			return
		case e := <-r.event:
			r.eventDeal(e)

		}
	}
}

//Scheduler must be  after  initialize true
func (r *PythonPlugin) Scheduler() (string, error) {
	//defer r.ready()
	if r.len == 0 {
		return "", nil
	}
	pn := r.curEntry
	pn.nextTime = pn.nextTime.Add(pn.interval)
	r.fixOrder()

	return pn.run()
}

//DeleteEntry arg PythonModuleName ,as stop the module
func (r *PythonPlugin) DeleteEntry(name string) error {
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

func (r *PythonPlugin) pop() *pluginEntry {
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
func (r *PythonPlugin) Len() int {
	return r.len
}

//InsertEntry  为了不重复，插入前都尝试一次删除 ,interval单位s
func (r *PythonPlugin) InsertEntry(name string, interval int) error {
	r.DeleteEntry(name)
	node, err := r.genEntry(name, interval)
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

func (r *PythonPlugin) insert(e *pluginEntry) {

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
func (r *PythonPlugin) fixOrder() bool {
	cur := r.curEntry
	next := cur.pNext

	if cur.nextTime.Before(next.nextTime) {
		return false
	}
	if r.len < 3 {
		r.next()
		return true
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
	return true
}

func (r *PythonPlugin) next() {
	r.curEntry = r.curEntry.pNext

}
func (r *PythonPlugin) prev() {
	r.curEntry = r.curEntry.pPrev

}

func (PythonPlugin) genEntry(name string, interval int) (*pluginEntry, error) {
	var pe pluginEntry
	if interval < 1 {
		return nil, errors.New("the interval le 0")
	}
	pe.name = name
	pe.interval = time.Second * time.Duration(interval)

	pe.nextTime = time.Now().Add(pe.interval)
	module := python.PyImport_ImportModule(pe.name)
	if module == nil {
		_, value, _ := python.PyErr_Fetch()
		return nil, errors.New(python.PyString_AS_STRING(value.Str()))
	}
	if !python.PyModule_Check(module) {
		return nil, errors.New("PyModule_Check(module)==false")
	}

	getvalue := module.GetAttrString("getvalue")
	if getvalue == nil {
		_, value, _ := python.PyErr_Fetch()
		return nil, errors.New(python.PyString_AS_STRING(value.Str()))
	}
	if !python.PyFunction_Check(getvalue) {
		return nil, errors.New("PyFunction_Check(getvalue)==false")
	}

	ret := getvalue.CallFunction()
	if ret == nil {
		_, value, _ := python.PyErr_Fetch()
		return nil, errors.New(python.PyString_AS_STRING(value.Str()))
	}
	if !python.PyString_Check(ret) {
		return nil, errors.New("PyString_Check(ret)==false")
	}
	pe.plugin = getvalue

	return &pe, nil
}
func (r *PythonPlugin) foreche() string {
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
