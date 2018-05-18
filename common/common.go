package common

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"
)

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

type ScriptConf struct {
	Name     string `json:"name"`
	FileName string `json:"filename"`
	HostName string `json:"hostname"`
	Interval int    `json:"interval"`
	TimeOut  int    `json:"timeout"`
}

type HttpResp struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result,omitempty"`
}
