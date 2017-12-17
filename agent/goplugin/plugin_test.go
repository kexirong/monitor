package goplugin

import (
	"testing"
	"time"
)

func Test_goplugin(t *testing.T) {
	time.Sleep(time.Second * 1)
	for _, i := range GopluginMap {
		t.Log(i.Gather())

	}
}

func Test_cpuplugin(t *testing.T) {
	var cpu = new(CPU)
	cpu.init()
	time.Sleep(time.Second * 1)
	t.Log(cpu.GetStep())
}
