package main

import (
	"fmt"
	"testing"

	"github.com/kexirong/monitor/common"
	"github.com/kexirong/monitor/server/models"
)

func Test_httplisent(t *testing.T) {
	startHTTPsrv()
}

func Test_aj_json(t *testing.T) {
	var req common.HttpReq
	var aj = &models.AlarmJudge{}

	req.Cause = aj
	err := json.Unmarshal([]byte(`{"method":"add","cause":{"anchor_point":"qwe","express":"qwe","level":"warning"}} `), &req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(aj)

}
