package goplugin

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kexirong/monitor/common/packetparse"
)

var (
	netStatFile = "/proc/net/dev"
)

//timesCPU ”/proc/stat“ times unitis are 10ms，so that‘s it

//NET exproted method has Init GetTarget
type NET struct {
	commonStruct
}

//Gather scheduler use
func (n *NET) Gather() ([]packetparse.Packet, error) {
	var hostname, _ = os.Hostname()
	var ret []packetparse.Packet
	var subret = packetparse.Packet{
		Plugin:    "net",
		HostName:  hostname,
		TimeStamp: float64(time.Now().Unix()),
		Type:      "derive",
		VlTags:    n.vltags,
	}
	valuenet, err := n.collect()
	if err != nil {
		return nil, err
	}

	for k, v := range valuenet {
		d, err := fsliced(n.preValue[k], v)
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
	return ret, nil
}

func (n *NET) getValueC(value []float64) ([]float64, error) {
	var ret []float64
	for _, v := range n.valueC {
		r := value[n.valueMap[v]]
		if r < 0 {
			return nil, errors.New("cpu plugin calculate error:  precent lt 0 ")
		}
		ret = append(ret, r)

	}
	return ret, nil

}

func (n *NET) collect() (procValue, error) {
	lines, err := readFileToStrings(netStatFile, 2, -1)
	if err != nil {
		return nil, err
	}
	return parseLineNET(lines)
}

//Init must use  befor of GetTarget method
func (n *NET) Init(VlTags string) error {
	var err error
	var tc []string
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
	tc = strings.Split(VlTags, "|")
	for _, v := range tc {
		_, ok := n.valueMap[v]
		if ok {
			n.valueC = append(n.valueC, v)
		}
	}
	if len(n.valueC) < 1 {
		return errors.New("net plugin init error: VlTags none hit")
	}
	n.vltags = strings.Join(n.valueC, "|")
	n.preValue, err = n.collect()
	return err
}

func parseLineNET(lines []string) (procValue, error) {
	var vl = make([]float64, 0, 16)
	var ret = make(procValue)
	for _, line := range lines {
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
			//fmt.Println("vl is cap：", cap(vl))
		}
		key := sline[0][:len(sline[0])-1]
		ret[key] = vl
		vl = make([]float64, 0, 10)

	}
	return ret, nil
}
