package activeplugin

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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
