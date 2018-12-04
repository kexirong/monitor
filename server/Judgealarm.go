package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/kexirong/monitor/common/packetparse"
	"github.com/kexirong/monitor/server/judge"
	"github.com/kexirong/monitor/server/models"
)

var Judge *judge.Judge

func judgeInit() *judge.Judge {
	var Judge = judge.NewJudge()
	alarmJudges, err := models.AlarmJudgesAll(monitorDB)
	if err != nil {
		panic(err)
	}
	for _, aj := range alarmJudges {
		err := Judge.AddRule(aj)
		if err != nil {
			Logger.Error.Println(err)
		}
	}
	return Judge
}

func judgeAlarm(tp *packetparse.TargetPacket) {
	ret := Judge.DoJudge(tp)
	for _, v := range ret {
		sendAlarm(v)
	}
}

type channels struct {
	emails []string
	wchats []string
}

func sendAlarm(aq *models.AlarmEvent) {

	var chans channels

	al, err := models.AlarmSendByAnchorPoint(monitorDB, aq.AnchorPoint)

	if err != nil {
		if err != sql.ErrNoRows {
			Logger.Error.Println(err)
		}
		return
	}
	if al.Channel == 0 {
		return
	}
	if len(al.List) < 1 {
		Logger.Error.Printf("sendAlarm: invalid list field in alarm_send table  AnchorPoint=%s", al.AnchorPoint)
		return
	}
	switch al.Type {
	case models.TypeTeam:
		teams := strings.Split(al.List, ",")

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

	case models.TypeStaff:
		staffs := strings.Split(al.List, ";")

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
		Logger.Error.Println("sendAlarm: invalid type field in alarm_link table AnchorPoint=%s", al.AnchorPoint)
		return
	}
	if al.Channel&1 == 1 && len(chans.emails) > 0 {
		data := fmt.Sprintf("to=%s&subject=MonitorAlarm&content=%s", strings.Join(chans.emails, ","), aq.String())
		//Logger.Info.Println(data)
		ret, err := http.Post(conf.EmailURL, "application/x-www-form-urlencoded", strings.NewReader(data))
		if err != nil {
			Logger.Error.Println(err)
		}

		//Logger.Info.Println(ret)
		aq.Stat = 1
		err = aq.Save(monitorDB)
		if err != nil {
			Logger.Error.Println(err)
		}
		ret.Body.Close()

	}
	if al.Channel&2 == 2 && len(chans.wchats) > 0 {
		data := fmt.Sprintf("to=%s&content=%s", strings.Join(chans.wchats, "|"), aq.String())
		//Logger.Info.Println(data)
		ret, err := http.Post(conf.WchatURL, "application/x-www-form-urlencoded", strings.NewReader(data))
		if err != nil {
			Logger.Error.Println(err)
		}
		//Logger.Info.Println(ret)
		aq.Stat = 1
		err = aq.Save(monitorDB)
		if err != nil {
			Logger.Error.Println(err)
		}
		ret.Body.Close()

	}

}
