package activeplugin

import (
	"testing"
)

func Test_httpandhttps(t *testing.T) {
	err := get("http://www.werq.com/")
	t.Log(err)
}
