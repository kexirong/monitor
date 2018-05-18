package activeplugin

import (
	"testing"
)

func Test_httpandhttps(t *testing.T) {
	str, err := Get("http://bing.com/search?q=heh")
	t.Log(str, "|||", err)
}
