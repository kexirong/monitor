package activeplugin

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kexirong/monitor/common/packetparse"
)

var TLSClient *http.Client

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}
	TLSClient = &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

}

//HTTPProbe built-in *http.Client
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
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				DisableCompression: true,
			},
			Timeout: 5 * time.Second,
		},
		hostName:   hostname,
		pluginName: "httpProbe",
	}
}

func (h *HTTPProbe) Name() string {
	return fmt.Sprintf("[%s]%s", h.pluginName, h.hostName)
}

func (h *HTTPProbe) Gather() ([]packetparse.TargetPacket, error) {

	var tps []packetparse.TargetPacket
	var tp = packetparse.TargetPacket{
		HostName:  h.hostName,
		TimeStamp: packetparse.Nsecond2Unix(time.Now().UnixNano()),
		Plugin:    h.pluginName,
		Type:      "bool",
		Value:     []float64{1},
	}

	for _, t := range h.target {
		var rsp []byte
		var err error
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
	return tps, nil
}

//AddJob args value must be has (method,url,contentType,reqdata)
func (h *HTTPProbe) AddJob(args ...interface{}) error {

	if len(args) != 4 {
		return errors.New("invalid args")
	}

	var method, _ = args[0].(string)
	var URL, _ = args[1].(string)
	var contentType, _ = args[2].(string)
	var reqdata, _ = args[3].(string)

	if method != "get" || method != "post" {
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

func (h *HTTPProbe) DeleteJob(target string) error {
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

/*






























 */

//Get ..
func Get(url string) (string, error) {
	rsp, err := TLSClient.Get(url)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	if 200 != rsp.StatusCode {
		return "", errors.New(rsp.Status)
	}
	byteData, err := ioutil.ReadAll(rsp.Body)

	return string(byteData), err

}

//Post ..
func Post(url string, contentType string, data string) (string, error) {
	body := strings.NewReader(data)
	rsp, err := TLSClient.Post(url, contentType, body)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	if 200 != rsp.StatusCode {
		return "", errors.New(rsp.Status)
	}
	byteData, err := ioutil.ReadAll(rsp.Body)

	return string(byteData), err

}
