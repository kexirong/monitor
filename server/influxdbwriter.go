package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/kexirong/monitor/packetparse"
)

var (
	dbinflux     = "monitor"
	userinflux   = "monitor"
	passwdinflux = "monitor"
	hostinflux   = "http://10.1.1.201:8086"
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

func writeToInfluxdb(pk packetparse.Packet) error {

	// Make client
	clt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     hostinflux,
		Username: userinflux,
		Password: passwdinflux,
	})
	if err != nil {
		panic(err.Error())
	}

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbinflux,
		Precision: "s",
	})

	// Create a point and add to batch
	tags := map[string]string{
		"hostname": pk.HostName,
		"instance": pk.Instance,
		"type":     pk.Type,
	}

	fields := make(map[string]interface{})

	if len(pk.Value) <= 0 {
		return fmt.Errorf("value error: %v", pk.Value)
	}

	if len(pk.Value) == 1 {

		fields["value"] = pk.Value[0]

		pt, err := client.NewPoint(pk.Plugin, tags, fields, time.Unix(int64(pk.TimeStamp), 0))

		if err != nil {
			return err
		}
		bp.AddPoint(pt)

	} else {

		if pk.VlTags == "" {
			return fmt.Errorf("value gt 0 but vltags is '' ")
		}

		sl := strings.Split(pk.VlTags, "|")

		if len(sl) != len(pk.Value) {
			return fmt.Errorf("value  and  vltags is not equals ")
		}

		for idx, value := range pk.Value {
			fields["value"] = value
			tags["instance"] = pk.Instance + "." + sl[idx]

			pt, err := client.NewPoint(pk.Plugin, tags, fields, time.Unix(int64(pk.TimeStamp), 0))

			if err != nil {

				return err
			}

			bp.AddPoint(pt)

		}

	}

	fmt.Println("writing...", bp)
	fmt.Println(pk.Plugin, tags, fields, time.Unix(int64(pk.TimeStamp), 0))

	err = clt.Write(bp)
	clt.Close()
	return err

}
