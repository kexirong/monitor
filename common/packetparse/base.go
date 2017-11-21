package packetparse

import (
	"reflect"
)

/*
   | 0 | 1 | 2 | 3 |
   +-------+-------+
   |  type |length |
   +-------+-------+
   | **** data***  |
   | ************  |
   +---------------+
*/

type Packet struct {
	HostName  string    `json:"hostname"`  //ops201
	TimeStamp float64   `json:"timestamp"` //23123131.123131
	Plugin    string    `json:"plugin"`    // cpu
	Instance  string    `json:"instance"`  // 0,1,2,3 (eth0,eth1)(sda,sdb)
	Type      string    `json:"type"`      //percent(百分比),counter(正数速率,主要是趋势),gauge(原值),derive(速率)
	Value     []float64 `json:"value"`
	VlTags    string    `json:"vltags"`  // "idle|user|system"(rx|tx)(read|write|use|free...)
	Message   string    `json:"message"` // description ,e.g: the disk is full please clean
}

var packMap = map[string]uint16{
	"hostname":  0x0000,
	"timestamp": 0x0001,
	"plugin":    0x0002,
	"instance":  0x0003,
	"type":      0x0004,
	"value":     0x0005,
	"vltags":    0x0006,
	"message":   0x0007,
}

var typesMap map[string]string //[name]type
var parseMap map[uint16]string //[id]name

func init() {
	parseMap = make(map[uint16]string)
	typesMap = make(map[string]string)

	for key, vl := range packMap {
		parseMap[vl] = key
	}

	var packet Packet
	t := reflect.TypeOf(packet)
	v := reflect.ValueOf(packet)

	for k := 0; k < t.NumField(); k++ {
		typesMap[t.Field(k).Tag.Get("json")] = v.Field(k).Kind().String()
	}
}
