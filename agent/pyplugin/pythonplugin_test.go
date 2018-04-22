package pyplugin

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kexirong/monitor/common"
)

func Test_pyplugin(t *testing.T) {
	pp, err := Initialize("")

	if err != nil {

		panic(err)
	}
	err = pp.InsertEntry("cpu", 5)
	if err != nil {
		fmt.Println(err)
	}
	err = pp.InsertEntry("cpus", 2)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		es := `[{"method":"add", "target":"cpus1","arg":"1"},{"method":"add", "target":"cpus2","arg":"3"},{"method":"add", "target":"cpus3","arg":"2"},{"method":"delete", "target":"cpus1"},{"method":"getlist" }]`
		var events []common.Event
		err := json.Unmarshal([]byte(es), &events)
		if err != nil {
			fmt.Println(err)
			return
		}
		for i := 0; i < len(events); i++ {
			nv := pp.AddEventAndWaitResult(events[i])
			events[i].Result = nv.Result
			if events[i].Target != nv.Target || events[i].Method != nv.Method || events[i].Arg != nv.Arg {
				events[i].Result = "server internal error"
			}
		}
		b, e := json.MarshalIndent(events, "", "    ")
		if e == nil {
			fmt.Println(string(b))
		} else {
			fmt.Println(e)
		}

	}()
	for {
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`")
		fmt.Println(pp.foreche())
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`")
		pp.WaitAndEventDeal()

		fmt.Println(pp.Scheduler())

	}

}
