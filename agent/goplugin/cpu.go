package goplugin

import (
	"errors"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/kexirong/monitor/common/packetparse"
)

var (
	//”/proc/stat“ times unitis are 10ms，so that‘s it
	cpuStatFile = "/proc/stat"
	cpuNum      = runtime.NumCPU()
)

//CPU exproted method has Init GetTarget
type CPU struct {
	plugin
}

//Gather scheduler use
func (c *CPU) Gather() ([]*packetparse.TargetPacket, error) {
	var hostname, _ = os.Hostname()
	var ret []*packetparse.TargetPacket

	timescpu, err := c.collect()
	if err != nil {
		return nil, err
	}
	for k, v := range timescpu {
		var subret = &packetparse.TargetPacket{
			Plugin:    "cpu",
			HostName:  hostname,
			TimeStamp: packetparse.Nsecond2Unix(time.Now().UnixNano()),
			Type:      "percent",
			VlTags:    c.vltags,
		}
		d, err := fsliced(c.lastValue[k], v)
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
	c.lastValue = timescpu
	return ret, nil
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
	if !c.Config("interval", 1) {
		return errors.New("CPU plugin： init set interval error")
	}

	if !c.Config("vltags", "user|nice|system|idle") {
		return errors.New("CPU plugin： init set vltags error")
	}

	c.lastValue, err = c.collect()
	return err
}

func (c *CPU) collect() (procvalue, error) {
	lines, err := readFileToStrings(cpuStatFile, 0, cpuNum+1) // 读取cpu跟各个核心的状态行 故cpuNum+1
	if err != nil {
		return nil, err
	}
	return parseLineCPU(lines)
}

func (c *CPU) timesPercent(times []float64) ([]float64, error) {
	var timestot float64
	var ret []float64

	for _, v := range times {
		timestot += v
	}
	for _, v := range c.valueC {
		r := times[c.valueMap[v]] / timestot * 100
		if r < 0 || math.IsNaN(r) {
			return nil, fmt.Errorf("cpu plugin calculate error:  %s precent lt 0 ,times:%v", v, times)
		}
		ret = append(ret, r)
	}
	return ret, nil
}

func parseLineCPU(lines []string) (procvalue, error) {
	var ret = make(procvalue)
	for _, line := range lines {
		var vl = make([]float64, 0, 10)
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
		}
		ret[sline[0]] = vl

	}
	return ret, nil
}

func init() {
	cpu := new(CPU)
	if err := cpu.init(); err == nil {
		register("cpu", cpu)
	} else {
		fmt.Println(err)
	}

}
