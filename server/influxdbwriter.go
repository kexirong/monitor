package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/kexirong/monitor/common/packetparse"
)

var (
	dbinflux     = "monitor"
	userinflux   = "monitor"
	passwdinflux = "monitor"
	hostinflux   = "http://10.1.1.201:8086"
)

/*
type  TargetPacket struct {
    HostName  string        `json:"hostname"`
    TimeStamp float64       `json:"timestamp"`
    Plugin    string        `json:"plugin"`
    Instance  string        `json:"instance"`
    Type      string        `json:"type"`
    Value     []float64     `json:"value"`
    VlTags    string        `json:"vltags"`
    Message   string       	`json:"message"`
}
*/
//由于float64的精度问题所以保留到毫秒，舍弃纳秒部分
func timestamp2Time(ts float64) time.Time {
	if ts < 0 {
		return time.Now().Round(time.Millisecond)
	}
	deno := float64(time.Second)
	return time.Unix(0, int64(ts*deno)).Round(time.Millisecond)
}

func writeToInfluxdb(pk packetparse.TargetPacket) error {
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
		"type":     pk.Type,
	}
	if pk.Instance != "" {
		tags["instance"] = pk.Instance
	}
	fields := make(map[string]interface{})
	if len(pk.Value) <= 0 {
		return fmt.Errorf("value error: %v", pk.Value)
	}

	/*if len(pk.Value) == 1 {
		fields["value"] = pk.Value[0]
		pt, err := client.NewPoint(pk.Plugin, tags, fields, time.Unix(int64(pk.TimeStamp), 0))
		if err != nil {
			return err
		}
		bp.AddPoint(pt)
	} else {  */

	if pk.VlTags == "" {
		return fmt.Errorf("value gt 0 but vltags is '' ")
	}

	sl := strings.Split(pk.VlTags, "|")

	if len(sl) < len(pk.Value) {
		return fmt.Errorf("value  and  vltags is not equals ")
	}

	for idx, value := range pk.Value {
		fields[sl[idx]] = value
		//	tags["metric"] = sl[idx]
	}
	pt, err := client.NewPoint(pk.Plugin, tags, fields, timestamp2Time(pk.TimeStamp))
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	fmt.Println("writing...", bp)
	fmt.Println(pk.Plugin, tags, fields, timestamp2Time(pk.TimeStamp))

	err = clt.Write(bp)
	clt.Close()
	return err

}
