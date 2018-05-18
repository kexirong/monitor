package common

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var pageSize = syscall.Getpagesize()

const userHZ = 100 //linux normal

type ProcessList map[string]Process

type Process struct {
	Pid         int     `json:"pid"`
	PPid        int     `json:"ppid"`
	Name        string  `json:"name"`
	CmdLine     string  `json:"cmdline"`
	MemoryUse   int     `json:"memroyused"`
	CPUPercent  float64 `json:"cpupercent"`
	ThreadCount int     `json:"threads"`
	Stat        string  `json:"stat"`
	CPUtimes    cpuInfo `json:"-"`
}

type cpuInfo struct {
	Time    time.Time
	Jiffies int
}

var t int
var trimFunc = func(r rune) bool {
	if r == '(' {
		t++
	}
	if r == ')' && t > 0 {
		t--
	}
	if r == ' ' && t == 0 {
		return true
	}
	return false
}

func isInt(s string) bool {
	for _, c := range s {
		if '0' > c && '9' < c {
			return false
		}
	}
	return true
}

func (pl *ProcessList) Init(pids ...string) {
	if len(pids) == 0 {
		s, _ := ioutil.ReadDir("/proc")
		for _, v := range s {
			if !v.IsDir() {
				continue
			}
			name := v.Name()
			if isInt(name) {
				pids = append(pids, name)
			}
		}
	}
	*pl = make(ProcessList, len(pids))
	for _, pid := range pids {
		cmdline, _ := ioutil.ReadFile(fmt.Sprintf("/proc/%s/cmdline", pid))
		(*pl)[pid] = Process{CmdLine: string(cmdline)}
	}
}
func (pl *ProcessList) FilterCmdline(patterns []string) error {
	for k, pro := range *pl {
		var i int
		for _, pattern := range patterns {
			matched, err := regexp.MatchString(pattern, pro.CmdLine)
			if err != nil {
				return err
			}
			if matched {
				i++
			}
		}
		if i == 0 {
			delete(*pl, k)
		}
	}
	return nil
}

func (pl *ProcessList) LoadsProcessInfo() {
	for pid, process := range *pl {
		now := time.Now()
		line, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s/stat", pid))
		if err != nil {
			delete(*pl, pid)
			continue
		}
		t = 0
		lines := strings.FieldsFunc(string(line), trimFunc)
		if len(lines) != 52 {
			delete(*pl, pid)
			continue
		}
		process.Pid, _ = strconv.Atoi(lines[0])
		process.Name = strings.Trim(lines[1], "()")
		process.Stat = lines[2]
		process.PPid, _ = strconv.Atoi(lines[3])

		mPage, _ := strconv.Atoi(lines[23])
		process.MemoryUse = mPage * pageSize
		process.ThreadCount, _ = strconv.Atoi(lines[19])
		t14, _ := strconv.Atoi(lines[13])
		t15, _ := strconv.Atoi(lines[14])
		t16, _ := strconv.Atoi(lines[15])
		t17, _ := strconv.Atoi(lines[16])
		process.CPUtimes.Jiffies = t14 + t15 + t16 + t17
		process.CPUtimes.Time = now
		(*pl)[pid] = process
	}
	time.Sleep(time.Millisecond * 100)
	for k, v := range *pl {
		now := time.Now()
		line, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s/stat", k))
		if err != nil {
			delete(*pl, k)
			continue
		}
		t = 0
		lines := strings.FieldsFunc(string(line), trimFunc)
		if len(lines) != 52 {
			delete(*pl, k)
			continue
		}
		t14, _ := strconv.Atoi(lines[13])
		t15, _ := strconv.Atoi(lines[14])
		t16, _ := strconv.Atoi(lines[15])
		t17, _ := strconv.Atoi(lines[16])
		//计算公式 Jiffies2-Jiffies1 / ((time2-time1)*hertz) * 100

		v.CPUPercent = float64(t14+t15+t16+t17-v.CPUtimes.Jiffies) / now.Sub(v.CPUtimes.Time).Seconds()
		(*pl)[k] = v
		if (*pl)[k].CPUPercent > 0 {
			fmt.Println(k, v.Name, v.CPUPercent)
		}
	}
}
