package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kexirong/monitor/packetparse"
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

type judgeMap map[string]map[string]judge //key1 is plugin, key2 is instance

type judge struct {
	ajtype string
	level1 sql.NullFloat64
	level2 sql.NullFloat64
	level3 sql.NullFloat64
}

var judgemap judgeMap

func doJudge(av alarmValue) string {
	iv, ok := judgemap[av.Plugin]
	var cmp func(x, y float64) bool
	if !ok {
		return ""
	}
	jv, ok := iv[av.Instance]
	if !ok {
		return ""
	}
	switch jv.ajtype {
	case "le":
		cmp = func(x, y float64) bool {
			return x <= y
		}
	case "ge":
		cmp = func(x, y float64) bool {
			return x >= y
		}
	case "ne":
		cmp = func(x, y float64) bool {
			return x != y
		}
	default:
		return ""

	}

	switch true {
	case jv.level3.Valid:
		if cmp(av.Value, jv.level3.Float64) {
			return "level3"
		}
	case jv.level2.Valid:
		if cmp(av.Value, jv.level2.Float64) {
			return "level2"
		}
	case jv.level1.Valid:
		if cmp(av.Value, jv.level1.Float64) {
			return "level1"
		}

	}
	return ""

}

func alarmJudge(pk packetparse.Packet) error {
	var alarmvalue alarmValue
	alarmvalue.HostName = pk.HostName
	alarmvalue.Plugin = pk.Plugin
	alarmvalue.Time = time.Unix(int64(pk.TimeStamp), 0).Format("2006-01-02 15:04:05")
	alarmvalue.Message = pk.Message

	fmt.Println(judgemap)
	leng := len(pk.Value)
	if leng <= 0 {
		return fmt.Errorf("value error: %v", pk.Value)
	}

	if leng == 1 {
		alarmvalue.Value = pk.Value[0]
		alarmvalue.Level = doJudge(alarmvalue)
		if alarmvalue.Level == "" {
			return nil
		}

		return alarmInsert(alarmvalue)

	}

	if leng > 1 {

		if pk.VlTags == "" {
			return fmt.Errorf("VlTags error: value gt 0 but vltags is '' ")
		}

		sl := strings.Split(pk.VlTags, "|")

		if len(sl) < len(pk.Value) {
			return fmt.Errorf("VlTags error:  vltags is not enough ")
		}

		for idx, value := range pk.Value {

			alarmvalue.Value = value
			alarmvalue.Instance = pk.Instance + "_" + sl[idx]

			alarmvalue.Level = doJudge(alarmvalue)
			if alarmvalue.Level == "" {
				continue
			}

			if err := alarmInsert(alarmvalue); err != nil {
				return err
			}
		}

	}

	return nil

}
