package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/kexirong/monitor/common/packetparse"
)

/*
const (
	dbinflux     = "monitor"
	userinflux   = "monitor"
	passwdinflux = "monitor"
	hostinflux   = "http://10.1.1.201:8086"
)
*/

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

//Influxdb writeToInfluxdb
type Influxdb struct {
	batchSize int
	mu        *sync.Mutex

	clt client.Client
	bp  client.BatchPoints
}

func (db *Influxdb) Write(tp *packetparse.TargetPacket) error {

	// Make client+
	db.mu.Lock()
	defer db.mu.Unlock()
	// Create a new point batch
	if db.bp == nil {
		db.bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
			Database:  conf.Influx.Database,
			Precision: "s",
		})
	}

	// Create a point and add to batch
	tags := map[string]string{
		"hostname": tp.HostName,
		"type":     tp.Type,
	}
	if tp.Instance != "" {
		tags["instance"] = tp.Instance
	}
	fields := make(map[string]interface{})
	if len(tp.Value) == 0 {
		return fmt.Errorf("value error: %v", tp.Value)
	}

	sl := strings.Split(tp.VlTags, "|")

	for idx, value := range tp.Value {
		fields[sl[idx]] = value
		//	tags["metric"] = sl[idx]
	}
	pt, err := client.NewPoint(tp.Plugin, tags, fields, timestamp2Time(tp.TimeStamp))
	if err != nil {
		return err
	}
	db.bp.AddPoint(pt)
	//fmt.Println(pk.Plugin, tags, fields, timestamp2Time(pk.TimeStamp))
	if db.batchSize > len(db.bp.Points()) {
		return nil
	}
	err = db.clt.Write(db.bp)
	db.bp = nil
	return err
}
