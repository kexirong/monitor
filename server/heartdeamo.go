package main

import (
	"fmt"
	"time"

	"github.com/kexirong/monitor/server/activeplugin"
)

var hostHeartRecorde = make(map[string]int64) //time.Now().Unix()
var hostIPMap = make(map[string]string)

func scanAssetdb() int {
	var host struct {
		name string
		ip   string
	}
	rows, err := mysql.Query("SELECT HostName,ip FROM opsmgt.cmdb_asset where is_active =1")
	checkErr(err)
	cur := time.Now().Unix()
	for rows.Next() {
		if err := rows.Scan(&host.name, &host.ip); err != nil {
			Logger.Error.Println(err)
			continue
		}

		hostIPMap[host.name] = host.ip
		hostHeartRecorde[host.name] = cur
	}

	return len(hostIPMap)
}

func hostIPMapAdd(host, ip string) bool {
	if v, ok := hostIPMap[host]; ok {
		if v == ip {
			return false
		}
	}
	hostIPMap[host] = ip
	hostHeartRecorde[host] = time.Now().Unix()
	return true
}
func heartdeamo() {
	scanAssetdb()
	var av alarmValue
	for range time.Tick(time.Second * 10) {
		now := time.Now().Unix()
		for k, v := range hostHeartRecorde {
			if now-v > 30 {
				av.Plugin = "heartbeat"
				av.Instance = "heartbeat.timeout"
				av.HostName = k
				av.Time = time.Unix(now, 0).Format("2006-01-02 15:04:05")
				av.Level = "level3"
				av.Value = float64(now - v)
				go func(ip string, av alarmValue) {
					out, err := activeplugin.HostPinger(4000, ip)
					if err == nil {
						av.Message = fmt.Sprintf("heartbeat lost %g；ping ok", av.Value)
					} else {
						av.Message = fmt.Sprintf("heartbeat lost %g；ping %s", av.Value, out)
					}
					alarmInsert(av)
				}(hostIPMap[k], av)
			}

		}

	}

}
