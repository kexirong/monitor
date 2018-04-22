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
	loadavgfile = "/proc/loadavg"
)

//LOADAVG exproted method has Init GetTarget
type LOADAVG struct {
	plugin
}

//Gather scheduler use
func (l *LOADAVG) Gather() ([]packetparse.TargetPacket, error) {
	var hostname, _ = os.Hostname()
	var ret []packetparse.TargetPacket
	var subret = packetparse.TargetPacket{
		Plugin:    "loadavg",
		HostName:  hostname,
		TimeStamp: packetparse.Nsecond2Unix(time.Now().UnixNano()),
		Type:      "gauge",
		VlTags:    l.vltags,
	}
	valueker, err := l.collect()
	if err != nil {
		return nil, err
	}
	subret.Value = valueker
	ret = append(ret, subret)
	return ret, nil
}

func (l *LOADAVG) init() error {
	var err error
	l.valueMap = map[string]int{
		"1min":  0,
		"5min":  1,
		"15min": 2,
	}
	if !l.Config("vltags", "1min|5min|15min") {
		return errors.New("NET plugin： init set vltags error")
	}
	if !l.Config("interval", 60) {
		return errors.New("NET plugin： init set interval error")
	}
	return err
}

func (l *LOADAVG) collect() ([]float64, error) {
	var ret []float64
	//var value float64
	line, err := readSingleLine(loadavgfile)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(line)
	if len(fields) != 5 {
		return nil, errors.New("loadavgfile fields ne 3")
	}
	for _, c := range l.valueC {
		value, err := strconv.ParseFloat(fields[l.valueMap[c]], 64)
		if err != nil {
			return nil, errors.New("KERNEL plugin error: ParseFloat " + err.Error())
		}
		ret = append(ret, value)
	}
	return ret, nil
}

func init() {
	load := new(LOADAVG)
	if err := load.init(); err == nil {
		register("loadavg", load)
	} else {
		fmt.Println(err)
	}
}
