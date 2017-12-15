package goplugin

import (
	"errors"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/kexirong/monitor/common/packetparse"
)

var (
	cpuStatFile = "/proc/stat"
	cpuNum      = runtime.NumCPU()
)

//timesCPU ”/proc/stat“ times unitis are 10ms，so that‘s it

//CPU exproted method has Init GetTarget
type CPU struct {
	commonStruct
}

//Gather scheduler use
func (c *CPU) Gather() ([]packetparse.Packet, error) {
	var hostname, _ = os.Hostname()
	var ret []packetparse.Packet
	var subret = packetparse.Packet{
		Plugin:    "cpu",
		HostName:  hostname,
		TimeStamp: float64(time.Now().Unix()),
		Type:      "percent",
		VlTags:    c.vltags,
	}
	timescpu, err := c.collect()
	if err != nil {
		return nil, err
	}

	for k, v := range timescpu {
		d, err := fsliced(c.preValue[k], v)
		if err != nil {
			return nil, err
		}
		perc, err := c.timesPercent(d)
		if err != nil {
			return nil, err
		}
		subret.Value = perc
		subret.Instance = k
		ret = append(ret, subret)
	}
	c.preValue = timescpu
	return ret, nil
}

func (c *CPU) timesPercent(times []float64) ([]float64, error) {
	var timestot float64
	var ret []float64

	for _, v := range times {
		timestot += v
	}
	for _, v := range c.valueC {
		r := times[c.valueMap[v]] / timestot * 100
		if r < 0 {
			return nil, errors.New("cpu plugin calculate error:  precent lt 0 ")
		}
		ret = append(ret, r)
	}
	return ret, nil
}

func (c *CPU) collect() (procValue, error) {
	lines, err := readFileToStrings(cpuStatFile, 0, cpuNum+1) // 读取cpu跟各个核心的状态行 故cpuNum+1
	if err != nil {
		return nil, err
	}
	return parseLineCPU(lines)
}

//Config conf method
func (c *CPU) Config(key string, value string) bool {
	var cvalue []string
	switch key {
	case "vltags":
		tc := strings.Split(value, "|")
		for _, v := range tc {
			_, ok := c.valueMap[v]
			if ok {
				cvalue = append(cvalue, v)
			}
		}
		if len(cvalue) < 1 {
			return false
		}
		c.valueC = cvalue
		c.vltags = strings.Join(c.valueC, "|")
		return true

	default:
		return false
	}

}

func (c *CPU) init() error {
	var err error

	c.valueMap = map[string]int{
		"user":       0,
		"nice":       1,
		"system":     2,
		"idle":       3,
		"iowait":     4,
		"irq":        5,
		"softirq":    6,
		"steal":      7,
		"guest":      8,
		"guest_nice": 9,
	}
	if !c.Config("vltags", "user|nice|system|idle") {
		return errors.New("CPU plugin： init error")
	}
	c.preValue, err = c.collect()
	return err
}

func parseLineCPU(lines []string) (procValue, error) {
	var vl = make([]float64, 0, 10)
	var ret = make(procValue)
	for _, line := range lines {
		if !strings.HasPrefix(line, "cpu") {
			return nil, errors.New("cpu plugin error: parse /proc/stat error")
		}
		sline := strings.Fields(line)
		for _, v := range sline[1:] {
			n, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, errors.New("cpu plugin error: parse /proc/stat strconv.ParseInt error")
			}
			vl = append(vl, n)
			//fmt.Println("vl is cap：", cap(vl))
		}
		ret[sline[0]] = vl
		vl = make([]float64, 0, 10)

	}
	return ret, nil
}

func init() {
	cpu := new(CPU)
	if err := cpu.init(); err == nil {
		register("cpu", cpu)
	}

}
