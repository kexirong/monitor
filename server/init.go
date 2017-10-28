package main

import (
	"database/sql"
	"log"
	"os"
)

func init() {
	//logging init

	Logger.Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	//mysql init
	var err error
	mysql, err = sql.Open("mysql", "monitor:monitor@tcp(10.1.1.107:3306)/monitor?charset=utf8")
	checkErr(err)
	mysql.SetMaxOpenConns(2000)
	mysql.SetMaxIdleConns(500)

	//judge init 需要在 mysql init 后面
	judgemap = judgemapStore()
}