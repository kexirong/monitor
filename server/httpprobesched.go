package main

import (
	"fmt"
	"time"

	"github.com/kexirong/monitor/server/activeplugin"
)

type httpprobe struct {
	hostname    string
	url         string
	cycle       int32
	time        string
	result      string
	contenttype string
	method      string
	reqmsg      string
	counter     int32
}

func httpprobesched() {
	for range time.Tick(time.Second * 10) {
		rows, err := mysql.Query("SELECT hostname,url,contenttype,method,reqmsg,counter FROM http_probe WHERE  unix_timestamp(now())  > unix_timestamp(time) + cycle ")
		checkErr(err)
		for rows.Next() {
			var hp httpprobe
			err = rows.Scan(&hp.hostname, &hp.url, &hp.contenttype, &hp.method, &hp.reqmsg, &hp.counter)
			checkErr(err)
			var ret string
			switch hp.method {
			case "get":
				ret, err = activeplugin.Get(hp.url + hp.reqmsg) //strings.TrimSpace(hp.reqmsg)

			case "post":
				ret, err = activeplugin.Post(hp.url, hp.contenttype, hp.reqmsg)
			default:
				continue
			}
			if len(ret) > 1024 {
				ret = ret[:1024]
			}
			hp.result = ret
			hp.time = time.Now().Format("2006-01-02 15:04:05")
			if err != nil {
				hp.result = err.Error()
				hp.counter++
				av := alarmValue{
					HostName: hp.hostname,
					Time:     hp.time,
					Plugin:   "httpprobe",
					Instance: hp.url,

					Value:   float64(hp.counter),
					Level:   "level2",
					Message: fmt.Sprintf("[%s]%s", hp.method, hp.result),
				}
				if err := alarmInsert(av); err != nil {
					Logger.Error.Println(err)
				}

			} else {
				hp.counter = 0
			}
			_, err = mysql.Exec(
				"UPDATE http_probe SET respon=?,counter=?,time=? where hostname=? and url=? and method=?",
				hp.result,
				hp.counter,
				hp.time,
				hp.hostname,
				hp.url,
				hp.method,
			)
			if err != nil {
				Logger.Error.Println(err)
			}
		}
	}

}
