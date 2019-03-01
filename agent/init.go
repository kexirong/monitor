package main

import (
	"io/ioutil"
	"os"

	"monitor/common/scheduler"
)

//var sp *scriptplugin.ScriptPlugin

var scriptScheduled = scheduler.New()
var scriptPath string

var _hostname, _ = os.Hostname()
var conf = struct {
	Servers    []string
	ScriptPath string
	ServerHTTP string
	HTTPListen string
}{}

func init() {

	dat, err := ioutil.ReadFile("./agentconf.json")
	checkErr(err)
	err = json.Unmarshal(dat, &conf)
	checkErr(err)

}

func init() {
	//var err error
	//sp, err = scriptplugin.Initialize(conf.ScriptPath)
	//if err != nil {
	//	panic(err)
	//}

	scriptPath = conf.ScriptPath
}
