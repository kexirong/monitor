package activeplugin

import (
	"testing"
)

func Test_httpandhttps(t *testing.T) {
	err := get("https://www.baidu.com/")
	t.Log(err)
}
