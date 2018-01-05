package activeplugin

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

/*func checkErr(err error){

}*/
func get(url string) error {

	rsp, err := TLSClient.Get(url)
	if err != nil {
		return err
	}
	if 200 != rsp.StatusCode {
		return errors.New(rsp.Status)
	}
	byteData, err := ioutil.ReadAll(rsp.Body)
	fmt.Println(string(byteData))
	return err

}
