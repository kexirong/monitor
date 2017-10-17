package influxdbwriter

import (
	"fmt"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/kexirong/monitor/packetparse"
)

const (
	DB       = "monitor"
	username = "monitor"
	password = "monitor"
	host     = "http://10.1.1.201:8086"
)

/*
type  Packet struct {
    HostName  string        `json:"hostname"`
    TimeStamp float64       `json:"timestamp"`
    Plugin    string        `json:"plugin"`
    Instance  string        `json:"instance"`
    Type      string        `json:"type"`
    Value     []float64     `json:"value"`
    VlTags    string        `json:"vltags"`
    Message   string       ` json:"message"`
}*/

func WriteToInfluxdb(pk packetparse.Packet) error {
	// Make client
	clt, err := client.NewHTTPClient(client.HTTPConfig{Addr: host})
	if err != nil {
		panic(err.Error())
	}

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  DB,
		Precision: "s",
	})

	// Create a point and add to batch
	tags := map[string]string{
		"hostname": pk.HostName,
		"instance": pk.Instance,
	}

	fields := map[string]interface{}{}

	if len(pk.Value) == 1 {

		fields["value"] = pk.Value[0]
		tags["type"] = pk.Type

		pt, err := client.NewPoint(pk.Plugin, tags, fields, time.Unix(int64(pk.TimeStamp), 0))

		if err != nil {
			panic(err.Error())
		}
		bp.AddPoint(pt)

	} else if len(pk.Value) > 1 {

		if pk.VlTags == "" {
			return fmt.Errorf("value gt 0 but vltags is '' ")
		}

		sl := strings.Split(pk.VlTags, "|")

		if len(sl) != len(pk.Value) {
			return fmt.Errorf("value  and  vltags is not equals ")
		}

		for idx, value := range pk.Value {
			fields["value"] = value
			tags["type"] = pk.Type + "_" + sl[idx]

			pt, err := client.NewPoint(pk.Plugin, tags, fields, time.Unix(int64(pk.TimeStamp), 0))

			if err != nil {
				panic(err.Error())
			}

			bp.AddPoint(pt)

		}

	} else {
		return fmt.Errorf("value error: %v", pk.Value)
	}
	fmt.Println("writing...", bp)
	clt.Write(bp)
	fmt.Println("write done.")
	return nil

}
