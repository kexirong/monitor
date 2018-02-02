package main

import (
	"fmt"
)

type alarmValue struct {
	HostName string
	Time     string
	Plugin   string
	Instance string
	Stat     int64
	Value    float64
	Level    string
	Message  string
}

func (a *alarmValue) String() string {
	return fmt.Sprintf("HostName: %s, Time: %s, Plugin: %s, Instance: %s, Value: %f, Level: %s, Message: %s",
		a.HostName, a.Time, a.Plugin, a.Instance, a.Value, a.Level, a.Message)
}
