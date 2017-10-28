package main

import "time"

type alarmValue struct {
	HostName  string
	TimeStamp time.Time
	Plugin    string
	Instance  string
	stat      int64
	Value     float64
	Message   string
}
