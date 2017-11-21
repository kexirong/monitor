package main

import (
	"log"
	"os"
)

type logs struct {
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

//Logger logging
var Logger = new(logs)

func init() {
	Logger.Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
