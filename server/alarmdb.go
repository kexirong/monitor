package main

import (
	"fmt"
)

type alarmValue struct {
	ID       int64
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
	return fmt.Sprintf("[%s] seq: %d, HostName: %s, Time: %s, Plugin: %s, Instance: %s, Value: %g, Message: %s",
		a.Level, a.ID, a.HostName, a.Time, a.Plugin, a.Instance, a.Value, a.Message)
}
