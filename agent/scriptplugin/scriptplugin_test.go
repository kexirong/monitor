package scriptplugin

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kexirong/monitor/common"
	"github.com/kexirong/monitor/common/scheduler"
)

var scriptScheduled = scheduler.New()
var scriptPath string

const testJSON = `[{"id":1,"host_ip":"127.0.0.1","host_name":"kk-debian","plugin_name":"cpu","interval":1,"timeout":3,"file_name":"cpu.py","plugin_type":"python"},{"id":2,"host_ip":"127.0.0.1","host_name":"kk-debian","plugin_name":"cpus","interval":3,"timeout":3,"file_name":"cpus.py","plugin_type":"python"},{"id":5,"host_ip":"127.0.0.1","host_name":"kk-debian","plugin_name":"cpus1","interval":2,"timeout":3,"file_name":"cpus1.py","plugin_type":"python"},{"id":3,"host_ip":"127.0.0.1","host_name":"kk-debian","plugin_name":"cpus2","interval":5,"timeout":3,"file_name":"cpus2.py","plugin_type":"python"},{"id":4,"host_ip":"127.0.0.1","host_name":"kk-debian","plugin_name":"cpus3","interval":7,"timeout":3,"file_name":"cpus3.py","plugin_type":"python"}]`

func Test_scriptplugin(t *testing.T) {

	var scs []*common.ScriptConf

	json.Unmarshal([]byte(testJSON), &scs)

	downloaurl := fmt.Sprintf("http://%s/downloadsscript/", "127.0.0.1:5001")
	for _, sc := range scs {

		err := CheckDownloads(downloaurl, filepath.Join(scriptPath, sc.FileName), false)

		if err != nil {
			t.Error(err)
			continue
		}
		tasker := NewScripter(filepath.Join(scriptPath, sc.FileName),
			time.Duration(sc.Timeout)*time.Second)
		fmt.Println(tasker.Name())
		scriptScheduled.AddTask(time.Duration(sc.Interval)*time.Second, tasker)
		s := scriptScheduled.EcheTaskList()
		var tmp interface{}
		jsoniter.Unmarshal([]byte(s), &tmp)
		r, _ := jsoniter.MarshalIndent(tmp, "", "    ")
		fmt.Println("#######################################################################\n", string(r), "\n", "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

		fmt.Println(scriptScheduled.Len())
	}

	var callback = func(b []byte, err error) {

	}

	t.Log("activePluginScheduler staring")
	scriptScheduled.Star(callback)

}
