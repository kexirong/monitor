package common

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"
)

/*easyjson:json
type ScriptConf struct {
	ID         int    `json:"id"`          // id
	HostIP     string `json:"host_ip"`     // host_ip
	HostName   string `json:"host_name"`   // host_name
	PluginName string `json:"plugin_name"` // plugin_name
	Interval   int    `json:"interval"`    // interval
	Timeout    int    `json:"timeout"`     // timeout
	FileName   string `json:"file_name"`
	PluginType string `json:"plugin_type"`
}
*/

type HttpResp struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result,omitempty"`
}

type HttpReq struct {
	Method string      `json:"method"`
	Cause  interface{} `json:"cause,omitempty"`
}

//CheckFileIsExist file or directory
func CheckFileIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

//NewUniqueID  gen a string
func NewUniqueID(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		//预防在rand.read 失败函数可用
		nano := time.Now().UnixNano()
		for i := 0; i < n; i++ {
			b[i] = byte(nano >> uint(i&63))
		}
	}
	return base64.StdEncoding.EncodeToString(b)[:n]
}
