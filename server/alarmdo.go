package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kexirong/monitor/server/activeplugin"
)

type channels struct {
	emails []string
	wchats []string
}
type alarmLink struct {
	alarmname string
	_type     string
	list      string
	channel   int32
	channels
}

func sendAlarm(av alarmValue) {
	var al alarmLink

	al.channels = channels{
		emails: make([]string, 0),
		wchats: make([]string, 0),
	}
	al.alarmname = fmt.Sprintf("%s[%s]", av.Plugin, av.HostName)
	err := mysql.QueryRow("select type,list,channel from alarm_link where alarmname=?",
		al.alarmname).Scan(&al._type, &al.list, &al.channel)
	if err != nil && err != sql.ErrNoRows {
		Logger.Error.Println(al.alarmname, err)
		return
	}

	if al.channel == 0 {
		err := alarmUpdate(av)
		if err != nil {
			Logger.Error.Println(err)
		}
		return
	}
	switch al._type {
	case "team":
		teams := strings.Split(al.list, ";")
		if len(teams) < 1 {
			Logger.Error.Printf("sendAlarm: invalid list field in alarm_link table  alarmname=%s", al.alarmname)
			return
		}
		s := strings.Join(teams, "','")
		//Logger.Info.Println(s)
		rows, err := mysql.Query(`SELECT b.email,b.wechat FROM opsmgt.monitor_staff_group a 
			join opsmgt.monitor_staff b on a.staff_id=b.staffid 
			join opsmgt.monitor_staffgroup c on c.id=a.staffgroup_id 
			WHERE c.groupname = (?) `, s)
		if err != nil {
			Logger.Error.Println(err)
			return
		}

		for rows.Next() {
			var e, w string
			err := rows.Scan(&e, &w)
			if err == nil {
				al.emails = append(al.emails, e)
				al.wchats = append(al.wchats, w)
			} else {
				Logger.Error.Println(err)
			}
		}

	case "staff":
		staffs := strings.Split(al.list, ";")
		if len(staffs) < 1 {
			Logger.Error.Printf("sendAlarm: invalid list field in alarm_link table  alarmname=%s", al.alarmname)
			return
		}
		s := strings.Join(staffs, "','")
		rows, err := mysql.Query("SELECT b.email,b.wechat FROM  opsmgt.monitor_staff b WHERE b.staffid in (?)", s)
		if err != nil {
			Logger.Error.Println(err)
			return
		}
		for rows.Next() {
			var e, w string
			err := rows.Scan(&e, &w)
			if err == nil {
				al.emails = append(al.emails, e)
				al.wchats = append(al.wchats, w)
			} else {
				Logger.Error.Println(err)
			}
		}
	default:
		Logger.Warning.Printf("sendAlarm: invalid type field in alarm_link table alarmname=%s", al.alarmname)
		return
	}
	if al.channel&1 == 1 && len(al.emails) > 0 {
		data := fmt.Sprintf("to=%s&subject=MonitorAlarm&content=%s", strings.Join(al.emails, ","), av.String())
		Logger.Info.Println(data)
		ret, err := activeplugin.Post(conf.EmailURL, "application/x-www-form-urlencoded", data)
		if err != nil {
			Logger.Error.Println(err)
		} else {
			Logger.Info.Println(ret)
			Logger.Error.Println(alarmUpdate(av))
		}
	}
	if al.channel&2 == 2 && len(al.wchats) > 0 {
		data := fmt.Sprintf("to=%s&content=%s", strings.Join(al.wchats, "|"), av.String())
		Logger.Info.Println(data)
		ret, err := activeplugin.Post(conf.WchatURL, "application/x-www-form-urlencoded", data)
		if err != nil {
			Logger.Error.Println(err)
		} else {
			Logger.Info.Println(ret)
			Logger.Error.Println(alarmUpdate(av))
		}
	}

}

func alarmdo() {
	for range time.Tick(time.Second * 3) {
		avs := scanalarmdb()
		if len(avs) < 1 {
			continue
		}
		for _, av := range avs {
			go sendAlarm(av)
		}
	}
}
