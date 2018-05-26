package activeplugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

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
		Type:      "bool",
	}
	v := url.Values{}
	for _, pattern := range p.parttens {
		v.Add("pattern", pattern)
	}

	resp, err := http.Get("http://" + p.ip + ":5101/process?" + v.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return json.Marshal(tps)
}

//AddJob args value must be has (partten),partten is a  regular expression
func (p *ProcessProbe) AddJob(param ...interface{}) error {
	if len(param) != 3 {
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
	if len(param) != 4 {
		return errors.New("invalid param")
	}
	var method, _ = param[0].(string)
	var URL, _ = param[1].(string)
	var target = fmt.Sprintf("[%s]%s", method, URL)
	for i := range h.target {
		if target == fmt.Sprintf("[%s]%s", h.target[i].Method, h.target[i].URL) {
			h.target = append(h.target[:i], h.target[i+1:]...)
			return nil
		}
	}
	return errors.New("not exist")
}

func (p *ProcessProbe) Get(url string) ([]byte, error) {
	rsp, err := h.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if 200 != rsp.StatusCode {
		return nil, errors.New(rsp.Status)
	}
	return ioutil.ReadAll(rsp.Body)
}

func (h *HTTPProbe) Post(url string, contentType string, data string) ([]byte, error) {
	body := strings.NewReader(data)
	rsp, err := h.client.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if 200 != rsp.StatusCode {
		return nil, errors.New(rsp.Status)
	}
	return ioutil.ReadAll(rsp.Body)
}
