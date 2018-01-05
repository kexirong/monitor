package goplugin

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/kexirong/monitor/common/packetparse"
)

//Goplugin interface
type Goplugin interface {
	Gather() ([]packetparse.TargetPacket, error)
	Config(key string, value interface{}) bool
	GetStep() int64
}

//Goplugintype .
type Goplugintype struct {
	NextTime int64
	Instance Goplugin
}

// GopluginMap .
var GopluginMap = map[string]*Goplugintype{}

func register(name string, plugin Goplugin) error {
	if _, ok := GopluginMap[name]; ok {
		return fmt.Errorf("plugin regist error: %s is exist", name)
	}
	GopluginMap[name] = &Goplugintype{
		NextTime: time.Now().UnixNano() + plugin.GetStep(),
		Instance: plugin,
	}
	return nil
}

type plugin struct {
	vltags    string
	valueMap  map[string]int
	valueC    []string
	step      int64
	lastValue procvalue
}

func (p *plugin) Config(key string, value interface{}) bool {
	var cvalue []string
	switch key {
	case "vltags":
		if _, ok := value.(string); !ok {
			return false
		}
		tc := strings.Split(value.(string), "|")
		for _, v := range tc {
			_, ok := p.valueMap[v]
			if ok {
				cvalue = append(cvalue, v)
			}
		}
		if len(cvalue) < 1 {
			return false
		}
		p.valueC = cvalue
		p.vltags = strings.Join(p.valueC, "|")
		return true
	case "step":
		if v, ok := value.(int); ok && v > 0 {
			p.step = int64(v) * int64(time.Second)
			return true
		}
		return false
	default:
		return false
	}
}

func (p *plugin) GetStep() int64 {
	return p.step
}

type procvalue map[string][]float64

func fsliced(fs1 []float64, fs2 []float64) ([]float64, error) {
	var leng int
	leng = len(fs1)
	if leng < 1 || leng != len(fs2) {
		return nil, fmt.Errorf("fsliced error: args len notequal or 0")
	}
	ret := make([]float64, leng)
	for i := range fs1 {
		ret[i] = fs2[i] - fs1[i]
	}
	return ret, nil
}

func readFileToStrings(filepath string, offset uint, n int) ([]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var ret []string
	r := bufio.NewReader(f)
	for i := 0; i < n+int(offset) || n < 0; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if i < int(offset) {
			continue
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}
	return ret, nil
}

func readSingleLine(filepath string) (string, error) {
	lines, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	n := bytes.IndexByte(lines, '\n')
	if n > 0 {
		return string(lines[:n]), nil
	}
	if n < 0 && len(lines) > 0 {
		return string(lines), nil
	}
	return "", fmt.Errorf("readSingleLine error: %s may is none", filepath)
}

//
