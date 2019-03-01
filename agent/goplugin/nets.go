package goplugin

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"monitor/common/packetparse"
)

var (
	netStatFile = "/proc/net/dev"
)

//NET exproted method has Init GetTarget
type NET struct {
	plugin
}

//Gather scheduler use
func (n *NET) Gather() (packetparse.TargetPackets, error) {
	var hostname, _ = os.Hostname()
	var ret packetparse.TargetPackets

	valuenet, err := n.collect()
	if err != nil {
		return nil, err
	}

	for k, v := range valuenet {
		var subret = &packetparse.TargetPacket{
			Plugin:    "net",
			HostName:  hostname,
			TimeStamp: packetparse.Nsecond2Unix(time.Now().UnixNano()),
			Type:      "derive",
			VlTags:    n.vltags,
		}
		d, err := fsliced(n.lastValue[k], v)
		if err != nil {
			return nil, err
		}
		perc, err := n.getValueC(d)
		if err != nil {
			return nil, err
		}
		subret.Value = perc
		subret.Instance = k
		ret = append(ret, subret)
	}
	n.lastValue = valuenet
	return ret, nil
}

// func NewNets()(Goplugin,error) 有时间再改
func (n *NET) init() error {
	var err error
	n.valueMap = map[string]int{
		"rx_bytes":      0,
		"rx_packets":    1,
		"rx_errs":       2,
		"rx_drop":       3,
		"rx_fifo":       4,
		"rx_frame":      5,
		"rx_compressed": 6,
		"rx_multicast":  7,
		"tx_bytes":      8,
		"tx_packets":    9,
		"tx_errs":       10,
		"tx_drop":       11,
		"tx_fifo":       12,
		"tx_colls":      13,
		"tx_carrier":    14,
		"tx_compressed": 15,
	}
	if !n.Config("vltags", "rx_bytes|tx_bytes") {
		return errors.New("NET plugin： init set vltags error")
	}
	if !n.Config("interval", 1) {
		return errors.New("NET plugin： init set interval error")
	}

	n.lastValue, err = n.collect()
	return err
}

func (n *NET) getValueC(value []float64) ([]float64, error) {
	var ret []float64
	for _, v := range n.valueC {
		r := value[n.valueMap[v]]
		if r < 0 {
			return nil, fmt.Errorf("net plugin calculate error: %s  value lt 0 ", v)
		}
		ret = append(ret, r)
	}
	return ret, nil
}

func (n *NET) collect() (procvalue, error) {
	lines, err := readFileToStrings(netStatFile, 2, -1)
	if err != nil {
		return nil, err
	}
	return parseLineNET(lines)
}

func parseLineNET(lines []string) (procvalue, error) {

	var ret = make(procvalue)
	for _, line := range lines {
		var vl = make([]float64, 0, 16)
		sline := strings.Fields(line)
		if !strings.HasSuffix(sline[0], ":") {
			return nil, errors.New("net plugin error: parse /proc/net/dev error")
		}
		for _, v := range sline[1:] {
			n, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, errors.New("net plugin error: parse /proc/net/dev strconv.ParseInt error")
			}
			vl = append(vl, n)
		}
		key := sline[0][:len(sline[0])-1]
		ret[key] = vl
	}
	return ret, nil
}

func init() {
	nets := new(NET)
	if err := nets.init(); err == nil {
		register("nets", nets)
	} else {
		fmt.Println(err)
	}
}
