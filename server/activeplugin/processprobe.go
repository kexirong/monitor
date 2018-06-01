package activeplugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/kexirong/monitor/common"

	"github.com/kexirong/monitor/common/packetparse"
)

//HTTPProbe built-in *http.Client
type ProcessProbe struct {
	client     *http.Client
	parttens   []string
	hostName   string
	ip         string
	pluginName string
}

//NewProcessProbe return *ProcessProbe built-in *http.Client
func NewProcessProbe(hostname, ip string) *ProcessProbe {
	return &ProcessProbe{
		client:     NewHTTPClient(),
		hostName:   hostname,
		ip:         ip,
		pluginName: "process_probe",
	}
}

func (p *ProcessProbe) Name() string {
	return fmt.Sprintf("[%s]%s", p.pluginName, p.hostName)
}

func (p *ProcessProbe) Do() ([]byte, error) {

	var tps []packetparse.TargetPacket
	var tp = packetparse.TargetPacket{
		HostName:  p.hostName,
		TimeStamp: packetparse.Nsecond2Unix(time.Now().UnixNano()),
		Plugin:    p.pluginName,
		Type:      "gauge",
		VlTags:    "pid|cpupercent|memroyused",
	}
	v := url.Values{}
	for _, pattern := range p.parttens {
		v.Add("pattern", pattern)
	}

	resp, err := p.client.Get("http://" + p.ip + ":5101/process?" + v.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	bbody, err := ioutil.ReadAll(resp.Body)
	var ret common.HttpResp
	var pl common.ProcessList
	ret.Result = &pl
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bbody, &ret)
	if err != nil {
		return nil, err
	}
	if ret.Code != 200 {
		return nil, errors.New(ret.Msg)
	}
	for _, p := range pl {
		if len(p.CmdLine) > 128 {
			tp.Instance = p.CmdLine[:125] + "..."
		} else {
			tp.Instance = p.CmdLine
		}
		tp.Value = append([]float64{}, float64(p.Pid))
		tp.Value = append(tp.Value, (p.CPUPercent))
		tp.Value = append(tp.Value, float64(p.MemoryUse))
		tps = append(tps, tp)
	}

	return json.Marshal(tps)
}

//AddJob args value must be has (partten),partten is a  regular expression
func (p *ProcessProbe) AddJob(param ...interface{}) error {
	if len(param) < 1 {
		return errors.New("invalid param")
	}
	var partten, _ = param[0].(string)
	_, err := regexp.Compile(partten)
	if err != nil {
		return err
	}
	p.parttens = append(p.parttens, partten)
	return nil
}

func (p *ProcessProbe) DeleteJob(param ...interface{}) error {
	if len(param) < 1 {
		return errors.New("invalid param")
	}
	var partten, _ = param[0].(string)
	for i := range p.parttens {
		if partten == p.parttens[i] {
			p.parttens = append(p.parttens[:i], p.parttens[i+1:]...)
			return nil
		}
	}
	return errors.New("not exist")
}
