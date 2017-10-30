package main

type alarmValue struct {
	HostName string
	Time     string
	Plugin   string
	Instance string
	//Stat     int64
	Value   float64
	Level   string
	Message string
}
