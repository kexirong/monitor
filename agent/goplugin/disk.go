package goplugin

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/kexirong/monitor/common/packetparse"
)

const (
	diskmtabFlie = "/proc/mounts"
	diskfsFile   = "/proc/filesystems"
)

//硬盘容量值使用了float64保存，单位字节。目前硬盘远小于1PB，理论上不会有精度问题

//DISK exproted method has Init GetTarget
type DISK struct {
	plugin
	devmtp map[string]string
}

//Gather scheduler use
func (d *DISK) Gather() (packetparse.TargetPackets, error) {
	var hostname, _ = os.Hostname()
	var ret packetparse.TargetPackets

	diskinfo, err := d.collect()
	if err != nil {
		return nil, err
	}
	for k, v := range diskinfo {
		var subret = &packetparse.TargetPacket{
			Plugin:    "disk",
			HostName:  hostname,
			TimeStamp: packetparse.Nsecond2Unix(time.Now().UnixNano()),
			Type:      "gauge",
			VlTags:    d.vltags,
		}
		value := make([]float64, len(d.valueC))
		for i, c := range d.valueC {
			value[i] = v[d.valueMap[c]]
		}
		if err != nil {
			return nil, err
		}
		subret.Value = value
		subret.Instance = k
		ret = append(ret, subret)
	}
	return ret, nil

}

func (d *DISK) init() error {
	var err error
	d.devmtp, err = getmtpmap()
	if err != nil {
		return err
	}
	d.valueMap = map[string]int{
		"all":   0,
		"used":  1,
		"avail": 2,
	}
	if !d.Config("vltags", "all|used|avail") {
		return errors.New("DISK plugin error： init set vltags error")
	}
	if !d.Config("interval", 60) {
		return errors.New("DISK plugin error: init set interval error")
	}
	return nil
}

func (d *DISK) collect() (procvalue, error) {
	var ret = procvalue{}
	var value []float64
	fs := syscall.Statfs_t{}
	for _, v := range d.devmtp {
		value = make([]float64, 3)
		err := syscall.Statfs(v, &fs)
		if err != nil {
			return nil, err
		}
		value[0] = float64(fs.Blocks * uint64(fs.Bsize))
		value[1] = float64((fs.Blocks - fs.Bfree) * uint64(fs.Bsize))
		value[2] = float64(fs.Bavail * uint64(fs.Bsize))
		ret[v] = value
	}
	return ret, nil

}

func getfilesystem() (map[string]bool, error) {
	var fs = map[string]bool{}
	lines, err := readFileToStrings(diskfsFile, 0, -1)
	if err != nil {
		return nil, err
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "nodev") {
			continue
		}
		sline := strings.TrimSpace(line)
		fs[sline] = true
	}
	return fs, nil
}

func getmtpmap() (map[string]string, error) {
	var mtp = map[string]string{}
	fs, err := getfilesystem()
	if err != nil {
		return nil, err
	}
	lines, err := readFileToStrings(diskmtabFlie, 0, -1)
	if err != nil {
		return nil, err
	}
	for _, line := range lines {
		if !strings.HasPrefix(line, "/dev") {
			continue
		}
		sline := strings.Fields(line)
		if len(sline) < 4 {
			continue
		}
		if _, ok := fs[sline[2]]; ok && strings.HasPrefix(sline[3], "rw") {
			mtp[sline[0]] = sline[1]
		}
	}
	if len(mtp) > 0 {
		return mtp, nil
	}
	return nil, errors.New("disk plugin error: get mount point failed")
}

func init() {
	disk := new(DISK)
	if err := disk.init(); err == nil {
		register("disk", disk)
	} else {
		fmt.Println(err)
	}
}
