package goplugin

import (
	"testing"
	"time"
)

func Test_cpuplugin(t *testing.T) {
	var cpu CPU
	//	t.Log(cpu.init(), gopluginMap)
	time.Sleep(time.Second * 1)
	ret, err := cpu.Gather()
	t.Log(err)
	t.Log(ret)
}
