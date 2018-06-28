package activeplugin

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kexirong/monitor/common/packetparse"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//HTTPProbe built-in *http.Client,implement scheduler.Tasker interface
type HTTPProbe struct {
	client     *http.Client
	target     []target
	hostName   string
	pluginName string
}

type target struct {
	Method      string
	URL         string
	ContentType string
	ReqData     string
}

//NewHTTPProbe return *HTTPProbe built-in *http.Client
func NewHTTPProbe(hostname string) *HTTPProbe {
	return &HTTPProbe{
		client:     NewHTTPClient(),
		hostName:   hostname,
		pluginName: "http_probe",
	}
}

func (h *HTTPProbe) Name() string {
	return fmt.Sprintf("[%s]%s", h.pluginName, h.hostName)
}

func (h *HTTPProbe) Do() ([]byte, error) {

	var tps []packetparse.TargetPacket
	var tp = packetparse.TargetPacket{
		HostName:  h.hostName,
		TimeStamp: packetparse.Nsecond2Unix(time.Now().UnixNano()),
		Plugin:    h.pluginName,
		Type:      "bool",
		VlTags:    "state",
	}

	for _, t := range h.target {
		var rsp []byte
		var err error
		tp.Value = []float64{0}
		switch t.Method {
		case "get":
			u, _ := url.Parse(t.URL)
			u.RawQuery = t.ReqData
			rsp, err = h.Get(u.String())
		case "post":
			if t.ContentType == "" {
				t.ContentType = "text/xml"
			}
			rsp, err = h.Post(t.URL, t.ContentType, t.ReqData)
		}
		tp.Instance = t.URL
		if "ok" == string(rsp) {
			tp.Value = []float64{1}
		}
		if err != nil {
			tp.Message = err.Error()
		}
		tps = append(tps, tp)
	}
	return json.Marshal(tps)
}

var httpprobeRegex = regexp.MustCompile(`\[(\w+)\](https?://\S+)`)

//AddJob args value must be has (target,contentType,reqdata)
func (h *HTTPProbe) AddJob(param ...interface{}) error {
	fmt.Println("httpprobe.go 91 len(param): ", len(param), param)
	if len(param) != 3 {
		return errors.New("invalid param")
	}
	var t, _ = param[0].(string)
	starget := httpprobeRegex.FindStringSubmatch(t)
	if len(starget) != 3 {
		return errors.New("parse target error")
	}
	var method = starget[1]
	var URL = starget[2]
	var contentType, _ = param[1].(string)
	var reqdata, _ = param[2].(string)

	if method != "get" && method != "post" {
		return errors.New("invalid Method")
	}
	if contentType == "" {
		contentType = "text/xml"
	}

	if _, err := url.ParseRequestURI(URL); err != nil {
		return err
	}

	h.target = append(h.target, target{
		Method:      method,
		URL:         URL,
		ContentType: contentType,
		ReqData:     reqdata,
	})
	return nil
}

func (h *HTTPProbe) DeleteJob(param ...interface{}) error {
	if len(param) < 1 {
		return errors.New("invalid param")
	}
	var target, _ = param[0].(string)

	for i := range h.target {
		if target == fmt.Sprintf("[%s]%s", h.target[i].Method, h.target[i].URL) {
			h.target = append(h.target[:i], h.target[i+1:]...)
			return nil
		}
	}
	return errors.New("not exist")
}

func (h *HTTPProbe) Get(url string) ([]byte, error) {
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

func NewHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DisableCompression: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*2)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 2,
		},
	}
}
