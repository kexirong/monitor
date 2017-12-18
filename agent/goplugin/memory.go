package goplugin

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kexirong/monitor/common/packetparse"
)

var (
	meminfoFile = "/proc/meminfo"
)

//MEM exproted method has Init GetTarget
type MEM struct {
	plugin
}

//Gather scheduler use
func (m *MEM) Gather() ([]packetparse.Packet, error) {
	var hostname, _ = os.Hostname()
	var ret []packetparse.Packet
	var subret = packetparse.Packet{
		Plugin:    "net",
		HostName:  hostname,
		TimeStamp: packetparse.Nsecond2Unix(time.Now().UnixNano()),
		Type:      "derive",
		VlTags:    m.vltags,
	}
	valuemem, err := m.collect()
	if err != nil {
		return nil, err
	}
	subret.Value = valuemem
	ret = append(ret, subret)
	return ret, nil
}

func (m *MEM) init() error {
	var err error
	m.valueMap = map[string]int{
		"MemTotal":     0,
		"MemFree":      0,
		"MemAvailable": 0,
		"Buffers":      0,
		"Cached":       0,
		"SwapCached":   0,
		"SwapTotal":    0,
		"SwapFree":     0,
	}
	if !m.Config("vltags", "MemTotal|MemFree|SwapTotal|SwapFree") {
		return errors.New("MEM plugin： init set vltags error")
	}
	if !m.Config("step", 1) {
		return errors.New("MEM plugin： init set step error")
	}

	return err
}

func (m *MEM) collect() ([]float64, error) {
	var ret []float64
	lines, err := readFileToStrings(meminfoFile, 0, -1)
	if err != nil {
		return nil, err
	}
	memkv := parseLineMEM(lines)
	for _, v := range m.valueC {
		if s, ok := memkv[v]; ok {
			ss := strings.Fields(s)
			if len(ss) == 2 && ss[1] == "kB" {
				ssf, err := strconv.ParseFloat(ss[0], 64)
				if err != nil {
					return nil, errors.New("MEM plugin error: parse /proc/stat strconv.ParseInt error")
				}
				ret = append(ret, ssf*1024)
			}
		} else {
			return nil, fmt.Errorf("MEM plugin error: get %s value failed", v)
		}
	}
	if len(ret) != len(m.valueC) {
		return nil, errors.New("MEM plugin error: len(ret)!= len(m.valueC)")
	}
	return ret, nil
}

func parseLineMEM(lines []string) map[string]string {
	var memkv = map[string]string{}
	for _, line := range lines {
		sline := strings.Split(line, ":")
		memkv[sline[0]] = sline[1]
	}
	return memkv
}

func init() {
	mem := new(MEM)
	if err := mem.init(); err == nil {
		register("memory", mem)
	} else {
		fmt.Println(err)
	}

}
