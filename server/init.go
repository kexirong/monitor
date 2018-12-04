package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	client "github.com/influxdata/influxdb/client/v2"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var conf = struct {
	Service           string
	MysqlConnetString string
	Influx            struct {
		Database  string
		User      string
		Passwd    string
		Host      string
		Precision string
		BatchSize int
	}
	WchatURL string
	EmailURL string
}{}

var monitorDB *sql.DB
var influxdbwriter *Influxdb

func init() {
	dat, err := ioutil.ReadFile("./conf.json")
	checkErr(err)
	err = json.Unmarshal(dat, &conf)
	checkErr(err)
	//logging init
	Logger.Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	//mysql init
	monitorDB, err = sql.Open("mysql", conf.MysqlConnetString)
	checkErr(err)
	monitorDB.SetMaxOpenConns(5)
	monitorDB.SetMaxIdleConns(20)

	//judge init 需要在 mysql init 后面
	Judge = judgeInit()

	//influxdb

	clt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     conf.Influx.Host,
		Username: conf.Influx.User,
		Password: conf.Influx.Passwd,
	})

	checkErr(err)
	influxdbwriter = &Influxdb{
		clt:       clt,
		mu:        new(sync.Mutex),
		batchSize: conf.Influx.BatchSize,
	}

}
