package main

import (
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
}

func run() {
	for range time.Tick(time.Second * 10) {
		rows, err := mysql.Query("SELECT hostname,url,contenttype,method,reqmsg ROM http_probe WHERE  now()  > time + cycle ")
		checkErr(err)
		for rows.Next() {
			var hp = new(httpprobe)
			err = rows.Scan(&hp.hostname, &hp.url, &hp.contenttype, &hp.method, &hp.reqmsg)
			checkErr(err)
			var ret string
			switch hp.method {
			case "get":
				ret, err = activeplugin.Get(hp.url)

			case "post":
				ret, err = activeplugin.Post(hp.url, hp.contenttype, hp.reqmsg)
			default:
				continue
			}
			hp.result = ret
			if err != nil {
				hp.result = err.Error()
			}
			_, err = mysql.Exec(
				"UPDATE http_probe SET repson=? where hostname=?,url=?,method=?",
				hp.result,
				hp.hostname,
				hp.url,
				hp.method,
			)

		}
	}

}
