package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/kexirong/monitor/common"

	"github.com/kexirong/monitor/common/queue"
)

var testJSON = `{
    "category": "scriptplugin",
    "events": 
        {
            "method": "getlist"
        }
    
}`

/*
[
        {
            "method": "add",
            "target": "cpus1",
            "arg": {"interval":"1"}
        },
        {
            "method": "add",
            "target": "cpus2",
            "arg": {"interval":"3"}
        },
        {
            "method": "add",
            "target": "cpus3",
            "arg": {"interval":"2"}
        },
        {
            "method": "delete",
            "target": "cpus1"
        },
*/
func Test_scriptpluginConsole(t *testing.T) {
	go func() {
		log.Fatal(http.ListenAndServe(":5101", nil))
	}()
	time.Sleep(time.Second)
	go func() {
		btq := queue.NewBtsQueue(4096)

		scriptPluginScheduler(btq)
	}()
	time.Sleep(time.Second * 3)
	client := &http.Client{}
	j := strings.NewReader(testJSON)
	req, err := http.NewRequest("POST", "http://127.0.0.1:5101/console", j)
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	bt, err := ioutil.ReadAll(resp.Body)

	log.Println(err, string(bt))

}

func Test_getpluginconfig(t *testing.T) {
	res, err := http.Get(fmt.Sprintf("http://%s/config/plugin", conf.ServerHTTP))
	if err != nil {
		t.Error(err)
		return
	}
	robots, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	t.Log(string(robots))
}

func Test_temp(t *testing.T) {
	/*
		s := `
		{"code":200,"msg":"ok","result":[{"name":"cpu","filename":"cpu.py","hostname":"kk-debian","interval":1,"timeout":3}]}`

		var resp common.HttpResp
		var conf []common.ScriptConf
		resp.Result = &conf
		json.Unmarshal([]byte(s), &struct {
			*common.HttpResp
		}{&resp})

		fmt.Println((*((resp.Result).(*[]common.ScriptConf)))[0].Name)
		fmt.Println(conf)

		err := sp.CheckDownloads("http://127.0.0.1:5001/getscript/", "cpus.py", false)
		if err != nil {
			t.Error(err)
		}*/
	var pl common.ProcessList
	pl.Init()
	pl.LoadsProcessInfo()
	b, _ := json.MarshalIndent(pl, "", "   ")
	_ = b
	fmt.Println(string(b))
}

func Test_ProcessInfo(t *testing.T) {
	go func() {
		//	log.Fatal(http.ListenAndServe(":5101", nil))
	}()
	time.Sleep(1 * time.Second)
	v := url.Values{}
	//	v.Add("pattern", ".*code.*")
	v.Add("pattern", ".*chrome.*")

	resp, err := http.Get("http://127.0.0.1:5101/process?" + v.Encode())
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	var ret common.HttpResp
	var pl common.ProcessList
	ret.Result = &pl

	json.Unmarshal(bt, &ret)
	fmt.Printf("%#v", pl)
	for _, p := range pl {
		fmt.Printf("%d  %d  %#v\n", len([]byte(p.CmdLine)), len([]rune(p.CmdLine)), []byte(p.CmdLine))
	}
}
