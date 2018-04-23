package main

import "github.com/kexirong/monitor/agent/pyplugin"

var pp *pyplugin.PythonPlugin

func init() {
	var err error
	pp, err = pyplugin.Initialize("./pyplugin")
	if err != nil {
		panic(err)
	}
}
