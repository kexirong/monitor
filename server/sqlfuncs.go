package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

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
