package goplugin

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kexirong/monitor/common/packetparse"
)

//Goplugin interface
type Goplugin interface {
	Gather() ([]packetparse.Packet, error)
	Config(key string, value string) bool
}

func register(name string, plugin Goplugin) error {
	gopluginMap[name] = plugin
	return nil
}

var gopluginMap = map[string]Goplugin{}

type commonStruct struct {
	vltags   string
	valueMap map[string]int
	valueC   []string
	preValue procValue
}
type procValue map[string][]float64

func fsliced(fs1 []float64, fs2 []float64) ([]float64, error) {
	var leng int
	leng = len(fs1)
	if leng < 1 || leng != len(fs2) {
		return nil, fmt.Errorf("fsliced error: args len notequal or 0")
	}
	ret := make([]float64, leng)
	for i := range fs1 {
		ret[i] = fs2[i] - fs1[i]
	}
	return ret, nil

}

func readFileToStrings(filepath string, offset uint, n int) ([]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ret []string

	r := bufio.NewReader(f)
	for i := 0; i < n+int(offset) || n < 0; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if i < int(offset) {
			continue
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}

	return ret, nil

}
