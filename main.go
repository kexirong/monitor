package main

import (
	"fmt"

	"github.com/kexirong/monitor/influxdbwriter"
	. "github.com/kexirong/monitor/packetparse"
	// "encoding/binary"
)

func main() {

	var pp = Packet{
		HostName:  "hostname",
		TimeStamp: 1232131123,
		Plugin:    "test",
		Type:      "count",
		Instance:  "instance",
		Value:     []float64{12.1237, 128.123123},
		VlTags:    "aa|bb",
	}

	bb, err := Package(pp)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(bb) //binary.LittleEndian.Uint16(bb[0:2]),Network.BytesToUint16(bb[2:4]))

	st, err := Parse(bb)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(st)

	err = influxdbwriter.WriteToInfluxdb(st)

	if err != nil {

		fmt.Println(err.Error())
	}

}
