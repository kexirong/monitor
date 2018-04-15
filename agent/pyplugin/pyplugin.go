package pyplugin

import (
	"errors"
	"strconv"
	"strings"
	"time"

	python "github.com/sbinet/go-python"
)

type pluginEntry struct {
	runing   bool
	nextTime int64
	interval int64
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
		return "", errors.New("is runing, may interval Too brief")
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
	plug       chan struct{}
	len        int
	event      chan string //"method:pluginnam[|interval]"
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
	ring.plug = make(chan struct{}, 1)
	ring.event = make(chan string, 1)
	ring.ready()
	ring.initialize = true
	return ring, nil
}

func (r *PythonPlugin) eventDeal(event string) error {
	ev := strings.SplitN(event, ":", 2)
	if len(ev) != 2 {
		return errors.New(event + "the arg format false")
	}
	switch ev[0] {
	case "delete":
		if err := r.DeleteEntry(ev[1]); err != nil {
			return errors.New(event + err.Error())
		}

	case "add":
		evs := strings.Split(ev[1], "|")
		if len(evs) != 2 {
			return errors.New(event + "the arg format false")
		}
		invl, err := strconv.Atoi(evs[1])
		if err != nil {
			return err
		}
		return r.InsertEntry(ev[1], invl)
	default:
		return errors.New(event + "unknown operation type")
	}

	return nil
}
func (r *PythonPlugin) ready() {
	r.plug <- struct{}{}

}

//WaitAndEventDeal 等待阻塞结束和时间
func (r *PythonPlugin) WaitAndEventDeal() []error {
	var err []error
	<-r.plug
	r.timer = time.Now().Unix()
	n := r.curEntry.nextTime - r.timer

	if n <= 0 {
		return nil
	}
	wait := time.After(time.Duration(n) * time.Second)
	for {
		select {
		case <-wait:
			return err
		case e := <-r.event:
			err1 := r.eventDeal(e)
			if err1 != nil {
				err = append(err, err1)
			}
		}
	}
}

//Scheduler must be  after  initialize true
func (r *PythonPlugin) Scheduler() (string, error) {
	pn := r.curEntry
	pn.nextTime += pn.interval
	r.fixOrder()
	r.ready()
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

//InsertEntry  m
func (r *PythonPlugin) InsertEntry(name string, interval int) error {
	node, err := r.genEntry(name, interval)
	if err != nil {
		return err
	}
	if r.len > 1 {
		r.prev()
	}
	r.insert(node)
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
func (r *PythonPlugin) fixOrder() {
	cur := r.curEntry
	next := cur.pNext

	if cur.nextTime < next.nextTime {
		return
	}
	if r.len < 3 {
		r.next()
		return
	}
	cur = r.pop()
	for i := 0; i < r.len-1; i++ {
		if cur.nextTime >= r.curEntry.nextTime {
			break
		}
		r.next()
	}
	r.insert(cur)
	r.curEntry = next
	return
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
	pe.interval = int64(interval)
	pe.nextTime = time.Now().Unix() + pe.interval

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
