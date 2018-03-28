package main

import (
	"fmt"
	"testing"
)

func Test_scanalarmdb(t *testing.T) {
	avs := scanalarmdb()
	fmt.Println(len(avs))

	for _, av := range avs {
		fmt.Println(av.String())

	}
	ret, err := mysql.Exec(
		"update  alarm_queue SET stat = ? where stat=? and value=?",
		1,
		0,
		100)
	if err != nil {
		fmt.Println("1", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("2", err)
		return
	}
	fmt.Println(n)
}
