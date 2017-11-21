package main

import (
	"log"
)

type logs struct {
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

//Logger logging
var Logger = new(logs)
