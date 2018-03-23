package activeplugin

import (
	"testing"
)

func Test_httpandhttps(t *testing.T) {
	str, err := Get("https://www.baidu.com/")
	t.Log(str, err)
}
