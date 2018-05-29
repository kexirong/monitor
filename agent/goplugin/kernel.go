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
	kernelfile = "/proc/sys/fs/file-nr"
)

//KERNEL exproted method has Init GetTarget
type KERNEL struct {
	plugin
}

//Gather scheduler use
func (k *KERNEL) Gather() ([]*packetparse.TargetPacket, error) {
	var hostname, _ = os.Hostname()
	var ret []*packetparse.TargetPacket
	var subret = &packetparse.TargetPacket{
		Plugin:    "kernel",
		HostName:  hostname,
		TimeStamp: packetparse.Nsecond2Unix(time.Now().UnixNano()),
		Type:      "gauge",
		VlTags:    k.vltags,
	}
	valueker, err := k.collect()
	if err != nil {
		return nil, err
	}
	subret.Value = valueker
	ret = append(ret, subret)
	return ret, nil
}

func (k *KERNEL) init() error {
	var err error
	k.valueMap = map[string]int{
		"fdused":   0,
		"fdunused": 1,
		"fdmax":    2,
	}
	if !k.Config("vltags", "fdused|fdunused|fdmax") {
		return errors.New("KERNEL plugin： init set vltags error")
	}
	if !k.Config("interval", 30) {
		return errors.New("KERNEL plugin： init set interval error")
	}
	return err
}

func (k *KERNEL) collect() ([]float64, error) {
	var ret []float64
	//var value float64
	line, err := readSingleLine(kernelfile)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return nil, errors.New("kernelfile fields ne 3")
	}
	for _, c := range k.valueC {
		value, err := strconv.ParseFloat(fields[k.valueMap[c]], 64)
		if err != nil {
			return nil, errors.New("KERNEL plugin error: ParseFloat " + err.Error())
		}
		ret = append(ret, value)
	}
	return ret, nil
}

func init() {
	ker := new(KERNEL)
	if err := ker.init(); err == nil {
		register("kernel", ker)
	} else {
		fmt.Println(err)
	}
}
