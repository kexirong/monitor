package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kexirong/monitor/server/models"
)

type channels struct {
	emails []string
	wchats []string
}

func sendAlarm(aq *models.AlarmQueue) {

	var chans channels
	aq.Stat = 1
	err := aq.Update(monitorDB)
	if err != nil {
		Logger.Error.Println(err)
	}

	al, err := models.AlarmLinkByAlarmName(monitorDB, aq.AlarmName+"["+aq.HostName+"]")

	if err != nil {
		if err != sql.ErrNoRows {
			Logger.Error.Println(err)
		}
		return
	}
	if al.Channel == 0 {

		return
	}
	switch al.Type.String() {
	case "team":
		teams := strings.Split(al.List, ",")
		if len(teams) < 1 {
			Logger.Error.Printf("sendAlarm: invalid list field in alarm_link table  alarmname=%s", al.AlarmName)
			return
		}
		s := strings.Join(teams, "','")
		//Logger.Info.Println(s)
		rows, err := monitorDB.Query(`SELECT b.email,b.wechat FROM opsmgt.monitor_staff_group a 
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
				chans.emails = append(chans.emails, e)
				chans.wchats = append(chans.wchats, w)
			} else {
				Logger.Error.Println(err)
			}
		}

	case "staff":
		staffs := strings.Split(al.List, ";")
		if len(staffs) < 1 {
			Logger.Error.Printf("sendAlarm: invalid list field in alarm_link table  alarmname=%s", al.AlarmName)
			return
		}
		s := strings.Join(staffs, "','")
		rows, err := monitorDB.Query("SELECT b.email,b.wechat FROM  opsmgt.monitor_staff b WHERE b.staffid in (?)", s)
		if err != nil {
			Logger.Error.Println(err)
			return
		}
		for rows.Next() {
			var e, w string
			err := rows.Scan(&e, &w)
			if err == nil {
				chans.emails = append(chans.emails, e)
				chans.wchats = append(chans.wchats, w)
			} else {
				Logger.Error.Println(err)
			}
		}
	default:
		Logger.Warning.Printf("sendAlarm: invalid type field in alarm_link table alarmname=%s", al.AlarmName)
		return
	}
	if al.Channel&1 == 1 && len(chans.emails) > 0 {
		data := fmt.Sprintf("to=%s&subject=MonitorAlarm&content=%s", strings.Join(chans.emails, ","), aq.String())
		Logger.Info.Println(data)
		ret, err := http.Post(conf.EmailURL, "application/x-www-form-urlencoded", strings.NewReader(data))

		if err != nil {
			Logger.Error.Println(err)
		} else {

			Logger.Info.Println(ret)
			aq.Stat = 1
			err := aq.Update(monitorDB)
			if err != nil {
				Logger.Error.Println(err)
			}
			ret.Body.Close()
		}
	}
	if al.Channel&2 == 2 && len(chans.wchats) > 0 {
		data := fmt.Sprintf("to=%s&content=%s", strings.Join(chans.wchats, "|"), aq.String())
		Logger.Info.Println(data)
		ret, err := http.Post(conf.WchatURL, "application/x-www-form-urlencoded", strings.NewReader(data))
		if err != nil {
			Logger.Error.Println(err)
		} else {
			Logger.Info.Println(ret)
			aq.Stat = 1
			err := aq.Update(monitorDB)
			if err != nil {
				Logger.Error.Println(err)
			}
			ret.Body.Close()
		}
	}

}

func alarmdo() {
	for range time.Tick(time.Second * 5) {
		avs, err := models.AlarmQueueByStat(monitorDB, 0)
		if err != nil {
			Logger.Error.Println(err)
		}
		if len(avs) < 1 {
			continue
		}
		for _, av := range avs {
			go sendAlarm(av)
		}
	}
}
