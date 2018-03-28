package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	client "github.com/influxdata/influxdb/client/v2"
)

var conf = struct {
	MysqlConnetString string
	InfluxDB          string
	InfluxUser        string
	InfluxPasswd      string
	InfluxHost        string
	WchatURL          string
	EmailURL          string
}{}

func init() {
	/*
		conf.MysqlConnetString = "monitor:monitor@tcp(10.1.1.222:3306)/monitor?charset=utf8"
		conf.InfluxDB = "monitor"
		conf.InfluxUser = "monitor"
		conf.InfluxPasswd = "monitor"
		conf.InfluxHost = "http://10.1.1.201:8086"
		conf.WchatURL = "http://10.8.12.152:4000/sender/wechat"
		conf.EmailURL = "http://10.8.12.152:4000/sender/mail"
	*/
	dat, err := ioutil.ReadFile("./conf.json")
	checkErr(err)
	err = json.Unmarshal(dat, &conf)
	checkErr(err)
	//logging init
	Logger.Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	//mysql init
	mysql, err = sql.Open("mysql", conf.MysqlConnetString)
	checkErr(err)
	mysql.SetMaxOpenConns(100)
	mysql.SetMaxIdleConns(20)

	//judge init 需要在 mysql init 后面
	judgemap = judgemapGet()

	//influxdb
	clt, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     conf.InfluxHost,
		Username: conf.InfluxUser,
		Password: conf.InfluxPasswd,
	})
	checkErr(err)

}
