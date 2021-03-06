package goplugin

import (
	"testing"
	"time"
)

func Test_goplugin(t *testing.T) {
	time.Sleep(time.Second * 1)
	for n, Instance := range GopluginMap {
		t.Log(n)
		gather, err := Instance.Gather()
		if err != nil {
			t.Log(err)
			continue
		}
		t.Log(gather)
	}
}

func Test_cpuPlugin(t *testing.T) {
	var cpu = new(CPU)
	cpu.init()
	time.Sleep(time.Second * 1)
	t.Log(cpu.Gather())
}

func Test_diskPlugin(t *testing.T) {
	var disk = new(DISK)
	disk.init()
	time.Sleep(time.Second * 1)
	t.Log(disk.Gather())
}
