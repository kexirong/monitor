package activeplugin

import (
	"fmt"
	"net/url"
	"testing"
)

func Test_pingger(t *testing.T) {
	str, err := HostPinger(1, "www.kexirong.info")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("ok")
	}

	v := url.Values{}

	v.Add("pattern", ".*msg-sender.*")

	p := NewProcessProbe("server", "10.8.12.152")
	resp, err := p.client.Get("http://" + p.ip + ":5101/process?" + v.Encode())
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	fmt.Println("output", str)
}
