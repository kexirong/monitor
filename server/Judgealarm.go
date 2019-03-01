package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"monitor/common/packetparse"
	"monitor/server/judge"
	"monitor/server/models"
)

// 发送周期
const sendTick = 1

//Judge 全局线程安全
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
	if Judge == nil {
		Logger.Error.Println("the Judge is null")
	}
	ret := Judge.DoJudge(tp)

	for _, v := range ret {
		sendAlarm(v)
	}
}

type channels struct {
	emails []string
	wchats []string
}

func sendAlarm(ae *models.AlarmEvent) {

	var chans channels
	aes, err := models.AlarmEventsByHostNameAnchorPointRule(monitorDB, ae.HostName, ae.AnchorPoint, ae.Rule)
	if err != nil {
		Logger.Error.Println(err)
		return
	}
	if len(aes) > 0 {
		aes[0].Count++
		aes[0].Value = ae.Value
		aes[0].Stat = ae.Stat
		aes[0].Level = ae.Level
		aes[0].Message = ae.Message
		ae = aes[0]
	}
	ae.Save(monitorDB)
	al, err := models.AlarmSendByAnchorPoint(monitorDB, ae.AnchorPoint)

	if err != nil {
		if err != sql.ErrNoRows {
			Logger.Error.Println(err)
		}
		return
	}
	if al.SendTick > 1 {
		if ae.Count%al.SendTick != 0 {
			Logger.Info.Printf("#SendTick check# id:%d, sendtick:%d, count:%d ->ignore\n", ae.ID, al.SendTick, ae.Count)
			return
		}
	}
	ae.Stat = 2
	err = ae.Save(monitorDB)
	if err != nil {
		Logger.Error.Println(err)
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
		teams := strings.Split(al.List, ";")

		s := strings.Join(teams, "','")
		//Logger.Info.Println(s)
		rows, err := monitorDB.Query(`SELECT a.email,a.wxwork FROM dashboard.user a 
			join dashboard.user_groups b on a.username=b.user_id 
			join dashboard.auth_group c on c.id=b.group_id 
			WHERE c.name in (?) `, s)
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
		fmt.Println(s)
		rows, err := monitorDB.Query("SELECT b.email,b.wxwork FROM  dashboard.user b WHERE b.username in (?)", s)
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
		Logger.Error.Printf("sendAlarm: invalid type field in alarm_link table AnchorPoint=%s \n", al.AnchorPoint)
		return
	}
	if al.Channel&1 == 1 && len(chans.emails) > 0 {
		data := fmt.Sprintf("to=%s&subject=MonitorAlarm&content=%s", strings.Join(chans.emails, ","), ae.String())
		//Logger.Info.Println(data)
		ret, err := http.Post(conf.EmailURL, "application/x-www-form-urlencoded", strings.NewReader(data))
		if err != nil {
			Logger.Error.Println(err)
		}

		if err != nil {
			Logger.Error.Println(err)
		}
		ret.Body.Close()

	}
	if al.Channel&2 == 2 && len(chans.wchats) > 0 {
		data := fmt.Sprintf("to=%s&content=%s", strings.Join(chans.wchats, "|"), ae.String())
		//Logger.Info.Println(data)
		ret, err := http.Post(conf.WchatURL, "application/x-www-form-urlencoded", strings.NewReader(data))
		if err != nil {
			Logger.Error.Println(err)
		}
		//Logger.Info.Println(ret)

		if err != nil {
			Logger.Error.Println(err)
		}
		ret.Body.Close()

	}

}
