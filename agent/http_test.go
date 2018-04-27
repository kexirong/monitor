package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kexirong/monitor/common/queue"
)

var testJSON = `{
    "category": "pyplugin",
    "events": [
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
        {
            "method": "getlist"
        }
    ]
}`

func Test_pypluginConsole(t *testing.T) {
	go func() {
		log.Fatal(http.ListenAndServe(":5101", nil))
	}()

	go func() {
		btq := queue.NewBtsQueue(4096)

		pyPluginScheduler(btq)
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
