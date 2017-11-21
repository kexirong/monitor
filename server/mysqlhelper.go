package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// 为了美观 sql 放这里

//('10.1.1.107',3306,'monitor','monitor','monitor')
var mysql *sql.DB

func judgemapGet() judgeMap {
	judgemap := make(judgeMap)
	rows, err := mysql.Query("SELECT * FROM alarm_judge")
	checkErr(err)
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

func alarmInsert(av alarmValue) error {

	_, err := mysql.Exec(
		"INSERT alarm_queue SET hostname=?,alarmname=?,alarmele=?,value=?,message=?,time=?,level=?",
		av.HostName,
		av.Plugin,
		av.Instance,
		av.Value,
		av.Message,
		av.Time,
		av.Level)

	
	return err

}

/*
	stmt, err := mysql.Prepare("INSERT AlarmQueue SET username=?,departname=?,created=?")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "部门", "2012-12-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
	//更新数据
	stmt, err = mysql.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	//查询数据

	//删除数据
	stmt, err = mysql.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)*/
