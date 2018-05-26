package main

import (
	"database/sql"

	"github.com/kexirong/monitor/common"

	_ "github.com/go-sql-driver/mysql"
)

//('10.1.1.107',3306,'monitor','monitor','monitor')

func pluginconfGet(ip string) []common.ScriptConf {
	var conf common.ScriptConf
	var confs []common.ScriptConf
	rows, err := monitorDB.Query("SELECT a.pluginname, b.filename,a.hostname,a.interval,a.timeout FROM plugin_config a JOIN plugin b on  a.pluginname=b.pluginname WHERE hostip = ?", ip)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&conf.Name, &conf.FileName, &conf.HostName, &conf.Interval, &conf.TimeOut)
		checkErr(err)
		confs = append(confs, conf)
	}
	return confs
}

func judgemapGet() judgeMap {
	judgemap := make(judgeMap)
	rows, err := monitorDB.Query("SELECT * FROM alarm_judge")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}
	for rows.Next() {
		var plugin string
		var instance string
		var ajtype string
		var l1, l2, l3 sql.NullFloat64
		err = rows.Scan(&plugin, &instance, &ajtype, &l1, &l2, &l3)
		checkErr(err)
		if _, ok := judgemap[plugin]; !ok {
			judgemap[plugin] = map[string]judge{
				instance: judge{
					ajtype: ajtype,
					level1: l1,
					level2: l2,
					level3: l3,
				},
			}
			continue
		}
		judgemap[plugin][instance] = judge{
			ajtype: ajtype,
			level1: l1,
			level2: l2,
			level3: l3,
		}
	}

	//rows.Close()
	return judgemap
}
