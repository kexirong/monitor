package main

import (
	"fmt"
	"time"

	"monitor/server/activeplugin"
	"monitor/server/models"
)

//暂时先这样临时实现，待改成收到数据就更新心跳数据，而不是单独的心跳包，将cmdb数据进行字段扩展，quic
var ipHeartRecorde = make(map[string]int64) //time.Now().Unix()
var ipHostnameMap = make(map[string]string)

func scanAssetdb() int {
	var host struct {
		name string
		ip   string
	}
	rows, err := monitorDB.Query("SELECT HostName,ip FROM dashboard.cmdb_asset where is_active =1")
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

	for range time.Tick(time.Second * 10) {
		now := time.Now().Unix()
		for k, v := range ipHeartRecorde {
			if now-v > 30 {
				ae := &models.AlarmEvent{}
				ae.Stat = 1
				ae.AnchorPoint = "heartbeat.timeout"
				ae.HostName = ipHostnameMap[k]
				ae.CreatedAt = time.Unix(now, 0)
				ae.Level = models.LevelWarning
				ae.Value = float64(now - v)

				go func(ip string, ae *models.AlarmEvent) {
					out, err := activeplugin.HostPinger(4000, ip)
					if err == nil {
						ae.Message = fmt.Sprintf("heartbeat lost %g；ping ok", ae.Value)
					} else {
						ae.Message = fmt.Sprintf("heartbeat lost %g；ping %s", ae.Value, out)
					}
					//ae.Insert(monitorDB)
					sendAlarm(ae)
				}(k, ae)
			}
		}
	}
}
