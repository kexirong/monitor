package scriptplugin

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"monitor/common"
)

//Scripter   implement scheduler.Tasker interface
type Scripter struct {
	scriptFile string
	timeout    time.Duration
}

//NewScripter return *Scripter
func NewScripter(filePath string, timeout time.Duration) *Scripter {
	return &Scripter{
		scriptFile: filePath,
		timeout:    timeout,
	}
}

//Name scheduler.Tasker's method
func (s *Scripter) Name() string {
	return path.Base(s.scriptFile)
}

//Do scheduler.Tasker's method
func (s *Scripter) Do() ([]byte, error) {
	return common.Command(s.scriptFile, s.timeout)
}

//AddJob scheduler.Tasker's method
func (s *Scripter) AddJob(param ...interface{}) error {
	return nil
}

//DeleteJob scheduler.Tasker's method
func (s *Scripter) DeleteJob(param ...interface{}) error {
	return nil
}

func CheckDownloads(url, filePath string, check bool) error {
	if check && common.CheckFileIsExist(filePath) {
		return nil
	}
	res, err := http.Get(url + path.Base(filePath))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	robots, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(robots)
	file.Sync()
	return nil
}
