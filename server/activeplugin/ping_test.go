package activeplugin

import (
	"fmt"
	"testing"
)

func Test_pingger(t *testing.T) {
	str, err := HostPinger(4000, "www.kexirong.info")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("ok")
	}

	fmt.Println("output", str)
}
