package scriptplugin

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kexirong/monitor/common"
)

func Test_scriptplugin(t *testing.T) {

	pp, err := Initialize("./")

	if err != nil {
		panic(err)
	}

	err = pp.InsertEntry("cpu.py", 5, 3)
	if err != nil {
		t.Log(err)
	}

	err = pp.InsertEntry("cpus.py", 2, 3)
	if err != nil {
		t.Log(err)
	}

	go func() {
		es := `[{"method":"add", "target":"cpus1.py","arg":{"interval":"1"}},{"method":"add", "target":"cpus2.py","arg":{"interval":"2"}},{"method":"add", "target":"cpus3.py","arg":{"interval":"2"}},{"method":"delete", "target":"cpus1.py"},{"method":"getlist" }]`
		var events []common.Event
		err := json.Unmarshal([]byte(es), &events)
		if err != nil {
			t.Log(err)
			return
		}
		for i := 0; i < len(events); i++ {
			nv := pp.AddEventAndWaitResult(events[i])
			events[i].Result = nv.Result
			if events[i].UniqueID != nv.UniqueID {
				events[i].Result = "server internal error"
			}
		}
		b, e := json.MarshalIndent(events, "", "    ")
		if e == nil {
			fmt.Println(string(b))
		} else {
			t.Log(e)
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
