package main

import (
	"fmt"
	"time"

	"github.com/kexirong/monitor/server/activeplugin"
	"github.com/kexirong/monitor/server/models"
)

var ipHeartRecorde = make(map[string]int64) //time.Now().Unix()
var ipHostnameMap = make(map[string]string)

func scanAssetdb() int {
	var host struct {
		name string
		ip   string
	}
	rows, err := monitorDB.Query("SELECT HostName,ip FROM opsmgt.cmdb_asset where is_active =1")
	checkErr(err)
	cur := time.Now().Unix()
	for rows.Next() {
		if err := rows.Scan(&host.name, &host.ip); err != nil {
			Logger.Error.Println(err)
			continue
		}

		ipHostnameMap[host.ip] = host.name
		ipHeartRecorde[host.ip] = cur
	}

	return len(ipHostnameMap)
}

func hostIPMapAdd(hostname, ip string) bool {
	if v, ok := ipHostnameMap[ip]; ok {
		if v == hostname {
			return false
		}
	}
	ipHostnameMap[ip] = hostname
	ipHeartRecorde[ip] = time.Now().Unix()
	return true
}

func heartdeamo() {
	scanAssetdb()
	var av models.AlarmQueue
	for range time.Tick(time.Second * 10) {
		now := time.Now().Unix()
		for k, v := range ipHeartRecorde {
			if now-v > 30 {
				av.AlarmName = "heartbeat"
				av.Alarmele = "heartbeat.timeout"
				av.HostName = ipHostnameMap[k]
				av.CreatedAt = time.Unix(now, 0)
				av.Level = models.LevelLevel2
				av.Value = float64(now - v)
				go func(ip string, av models.AlarmQueue) {
					out, err := activeplugin.HostPinger(4000, ip)
					if err == nil {
						av.Message = fmt.Sprintf("heartbeat lost %g；ping ok", av.Value)
					} else {
						av.Message = fmt.Sprintf("heartbeat lost %g；ping %s", av.Value, out)
					}
					av.Insert(monitorDB)
				}(k, av)
			}
		}
	}
}
