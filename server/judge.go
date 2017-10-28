package main

import (
	"database/sql"
	"fmt"
	"strings"

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

func alarmJudge(pk packetparse.Packet) error {
	fmt.Println(judgemap)

	if len(pk.Value) <= 0 {
		return fmt.Errorf("value error: %v", pk.Value)
	}

	if len(pk.Value) == 1 {

	} else {

		if pk.VlTags == "" {
			return fmt.Errorf("value gt 0 but vltags is '' ")
		}

		sl := strings.Split(pk.VlTags, "|")

		if len(sl) != len(pk.Value) {
			return fmt.Errorf("value  and  vltags is not equals ")
		}

		for idx, value := range pk.Value {

			tags["type"] = pk.Type + "_" + sl[idx]

			if err != nil {

				return err
			}

		}

	}

}
