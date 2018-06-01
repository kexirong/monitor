package activeplugin

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_httpandhttps(t *testing.T) {
	ret, err := http.Get("http://10.8.12.152:4000/health")
	fmt.Printf("%#v\n ERR:%s", ret, err)

}
