package main

import (
	"database/sql"
	"fmt"
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

func ss() {
	judgemap = judgemapStore()
	fmt.Println(judgemap)

}
