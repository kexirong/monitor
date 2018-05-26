package activeplugin

//ping需要root权限（或加suid） 所以用exec.Command
import (
	"time"

	"github.com/kexirong/monitor/common"
)

//HostPinger  timeout 单位s，host为域名或者ip
func HostPinger(timeout int, host string) (string, error) {
	var args []string
	//ping 非root用户最低200ms间隔
	args = append(args, "-i 0.2", "-c 4", host)
	out, err := common.Command("ping", time.Second*time.Duration(timeout), args...)
	return string(out), err
}
