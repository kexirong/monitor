package goplugin

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kexirong/monitor/common/packetparse"
)

var (
	// ”/proc/stat“ times unitis are 10ms，so that‘s it
	diskmtabFlie = "/proc/mounts"
	diskfsFile   = "/proc/filesystems"
)

//DISK exproted method has Init GetTarget
type DISK struct {
	plugin
	devmtp map[string]string
}

//Gather scheduler use
func (d *DISK) Gather() ([]packetparse.Packet, error) {
	var hostname, _ = os.Hostname()
	var ret []packetparse.Packet
	var subret = packetparse.Packet{
		Plugin:    "disk",
		HostName:  hostname,
		TimeStamp: packetparse.Nsecond2Unix(time.Now().UnixNano()),
		Type:      "percent",
		VlTags:    d.vltags,
	}

	return ret, nil
}

func (d *DISK) init() error {
	var err error
	mtp, err := getmtpmap()
	if err != nil {
		return err
	}
	d.devmtp = mtp
	d.valueMap = map[string]int{}
	for k, v := range mtp {
		d.valueMap[k] = 0
		d.valueC = append(d.valueC, v)
	}
	d.vltags = strings.Join(d.valueC, "|")
	if !d.Config("step", 60) {
		return errors.New("DISK plugin： init set step error")
	}
	return nil
}

func (d *DISK) collect() (procvalue, error) {

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
		if _, ok := fs[sline[2]]; ok {
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
