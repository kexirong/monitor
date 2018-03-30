package activeplugin

import (
	"testing"
)

func Test_httpandhttps(t *testing.T) {
	str, err := Get("http://10.1.1.201:8086/ping")
	t.Log(str, "|||", err)
}
