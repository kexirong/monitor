package pyplugin

import (
	"errors"

	python "github.com/sbinet/go-python"
)

type pluginNode struct {
	tally  uint8
	step   int
	name   string
	plugin *python.PyObject
	pNext  *pluginNode
}

func (p *pluginNode) run() (string, error) {
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
	timer      int
	nodeLen    int
	plug       chan struct{}
	event      chan string //"method:pluginname"
	curNode    *pluginNode
}

func initialize(pluginPath string) (*PythonPlugin, error) {

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
	ring.plug = make(chan struct{})
	ring.event = make(chan string)
	ring.initialize = true
	return ring, nil

}

func (PythonPlugin) genNode(name string) (*pluginNode, error) {
	var pn pluginNode

	pn.name = name

	module := python.PyImport_ImportModule(pn.name)

	if module == nil {
		_, value, _ := python.PyErr_Fetch()
		return nil, errors.New(python.PyString_AS_STRING(value.Str()))
	}
	if !python.PyModule_Check(module) {
		return nil, errors.New("PyModule_Check(module)==false")
	}
	/*
		step := module.GetAttrString("STEP")
		if step == nil {
			_, value, _ := python.PyErr_Fetch()
			return nil, errors.New(python.PyString_AS_STRING(value.Str()))
		}
		if !python.PyInt_Check(step) {
			return nil, errors.New("python.PyInt_Check(step)==false")
		}

		pn.step = python.PyInt_AS_LONG(step)
	*/
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
	pn.plugin = getvalue

	return &pn, nil
}

func (r *PythonPlugin) deleteNode(name string) bool {

	for i := 0; i < r.nodeLen; i++ {
		if name == r.curNode.pNext.name {
			if r.len() == 1 {
				r.nodeLen = 0
			} else {
				r.curNode.pNext = r.curNode.pNext.pNext
			}
			r.nodeLen--
			return true
		}
	}
	return false
}

func (r *PythonPlugin) next() *pluginNode {
	r.curNode = r.curNode.pNext
	return r.curNode
}

func (r *PythonPlugin) len() int {
	return r.nodeLen
}

// 随机插入
func (r *PythonPlugin) insertNode(node *pluginNode) bool {
	if r.nodeLen == 0 {
		r.curNode = node
		r.curNode.pNext = node
	} else {
		node.pNext = r.curNode.pNext
		r.curNode.pNext = node
	}
	r.nodeLen++
	return true
}
