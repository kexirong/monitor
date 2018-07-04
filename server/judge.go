package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kexirong/monitor/common/packetparse"
	"github.com/kexirong/monitor/server/models"
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
    Message   string        `json:"message"`
}*/

type judgeMap map[string]map[string]judge //key1 is plugin, key2 is instance

type judge struct {
	ajtype string
	level1 sql.NullFloat64
	level2 sql.NullFloat64
	level3 sql.NullFloat64
}

var judgemap judgeMap

func doJudge(aq models.AlarmQueue, jv judge) string {

	var cmp func(x, y float64) bool

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
		if cmp(aq.Value, jv.level3.Float64) {
			return "level3"
		}
	case jv.level2.Valid:
		if cmp(aq.Value, jv.level2.Float64) {
			return "level2"
		}
	case jv.level1.Valid:
		if cmp(aq.Value, jv.level1.Float64) {
			return "level1"
		}

	}
	return ""

}

func alarmJudge(pk packetparse.TargetPacket) error {
	var aq models.AlarmQueue
	var jkey string
	iv, ok := judgemap[pk.Plugin+"["+pk.HostName+"]"]

	if !ok {
		return nil
	}
	aq.HostName = pk.HostName
	aq.AlarmName = pk.Plugin
	aq.CreatedAt = time.Unix(int64(pk.TimeStamp), 0)
	aq.Message = pk.Message

	leng := len(pk.Value)
	if leng <= 0 {
		return fmt.Errorf("value error: %v", pk.Value)
	}

	if leng == 1 {
		aq.Value = pk.Value[0]
		aq.Alarmele = pk.Instance + "." + pk.VlTags

		if pk.Instance == "" {
			aq.Alarmele = pk.VlTags
		}
		jkey = fmt.Sprintf("%s.%s", aq.Alarmele, pk.Type)
		jv, ok := iv[jkey]
		if !ok {
			return nil
		}

		err := aq.Level.UnmarshalText([]byte(doJudge(aq, jv)))
		if err != nil {
			return nil
		}

		return aq.Insert(monitorDB)

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

			aq.Value = value
			if pk.Instance != "" {
				aq.Alarmele = pk.Instance + "." + sl[idx]
			} else {
				aq.Alarmele = sl[idx]
			}
			//alarmvalue.Instance = pk.Instance + "." + sl[idx]
			jkey = fmt.Sprintf("%s.%s", aq.Alarmele, pk.Type)
			jv, ok := iv[jkey]
			if !ok {
				continue
			}
			err := aq.Level.UnmarshalText([]byte(doJudge(aq, jv)))
			if err != nil {
				continue
			}

			if err := aq.Insert(monitorDB); err != nil {
				return err
			}
		}

	}

	return nil

}
