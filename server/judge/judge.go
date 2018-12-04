package judge

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"github.com/kexirong/monitor/common/packetparse"
	"github.com/kexirong/monitor/server/models"
)

/*
func(feild,tatol,limit): func 判断方法，判断feild 字段（feild之间支持4则运算），tatol 数据条数，limit 符合条件条数(可缺省，相当于1)
all(${cpu0},3,1): 最新的3个点有一个满足阈值条件则报警
all(${cpu0},3,3): 最新的3个点都满足阈值条件则报警
sum(${cpu0},3): 对于最新的3个点，其和满足阈值条件则报警
avg(${cpu0},3): 对于最新的3个点，其平均值满足阈值条件则报警
diff(${cpu0},3,1): 拿最新push上来的点（被减数），与历史最新的3个点（3个减数）相减，得到3个差，只要有一个差满足阈值条件则报警
pdiff(${cpu0},3,1): 拿最新push上来的点，与历史最新的3个点相减，得到3个差，再将3个差值分别除以减数乘100（相当于增长率），只要有1个商值满足阈值则报警
*/

type nodeType uint

const reExpr = `(?P<func>\w+)\((?P<formula>\(?[^,]+\)?),(?P<total>\d+),?(?P<limit>\d*)?\)(?P<operator>[<!=>]+)(?P<compare>[-]?\d+(?:\.\d+)?$)`
const storeCap = 12

const (
	operator nodeType = 1 << iota
	digit
	variable
	bracket
)

//Judge 实现告警判断
type Judge struct {
	judgers map[string]*judger
	parser  *regexp.Regexp
	mu      *sync.Mutex
	//sender  io.Writer
}

//NewJudge 返回初始化过的 Judge指针
func NewJudge() *Judge {
	re, err := regexp.Compile(reExpr)
	if err != nil {
		panic(err)
	}
	return &Judge{
		judgers: make(map[string]*judger),
		parser:  re,
		mu:      new(sync.Mutex),
	}
}

//AddRule 添加告警判定规则
func (j *Judge) AddRule(aj *models.AlarmJudge) error {
	j.mu.Lock()

	judger := j.judgers[aj.AnchorPoint]
	if judger == nil {
		judger = newJuger()
	}
	dt, err := parseExpress(aj.Express, j.parser)
	if err == nil {
		judger.deters[aj.Express] = dt
		judger.level = aj.Level
	}
	j.mu.Unlock()
	return err
}

//DelRule 删除内存中判定规则
func (j *Judge) DelRule(aj *models.AlarmJudge) {
	j.mu.Lock()

	delete(j.judgers, aj.AnchorPoint)
	j.mu.Unlock()
}

//DoJudge 对数据包进行告警判定
func (j *Judge) DoJudge(tp *packetparse.TargetPacket) []*models.AlarmEvent {
	//[host].plugin.instance:type

	var results []*models.AlarmEvent
	j.mu.Lock()

	anchorPoint := fmt.Sprintf("[%s].%s.%s:%s", tp.HostName, tp.Plugin, tp.Instance, tp.Type)
	judger := j.judgers[anchorPoint]
	if judger == nil {
		anchorPoint = fmt.Sprintf("[*].%s.%s:%s", tp.Plugin, tp.Instance, tp.Type)
		judger = j.judgers[anchorPoint]
	}
	j.mu.Unlock()
	if judger != nil {
		ret := judger.doJudge()
		for k, v := range ret {
			var result = &models.AlarmEvent{
				HostName:    tp.HostName,
				AnchorPoint: anchorPoint,
				Rule:        k,
				Value:       v,
				Level:       judger.level,
				Message:     tp.Message,
			}
			results = append(results, result)
		}
	}
	return results
}

/*
warning 警告
serious 严重
disaster 灾难
*/
type judger struct {
	store  *ringNode
	level  models.Level
	deters map[string]*deter //[key]*deter
}

func newJuger() *judger {
	return &judger{
		store:  makeStore(storeCap),
		deters: make(map[string]*deter),
	}
}

type deter struct {
	total    int
	limit    int
	operator string
	compare  float64
	arg      *binTreeNode
	method   func(datas []float64, operator string, compare float64, limit int) bool
}

func (j *judger) doJudge() map[string]float64 {
	var result = make(map[string]float64)
	for key, deter := range j.deters {
		ds := j.store.pull(deter.total)
		datas := make([]float64, deter.total)
		for i, d := range ds {
			if d == nil {
				break
			}
			datas[i] = deter.arg.calculate(d)
		}
		if len(datas) < deter.total {
			continue
		}
		if deter.method(datas, deter.operator, deter.compare, deter.limit) {
			result[key] = datas[len(datas)-1]
		}
	}
	return result
}

func parseExpress(express string, parser *regexp.Regexp) (*deter, error) {
	if parser == nil {
		return nil, errors.New("regexp is nil")
	}
	out := parser.FindStringSubmatch(string(bytes.Trim([]byte(express), " ")))
	l := len(out)
	if l == 0 {
		return nil, errors.New("express parse failed")
	}
	var dt deter
	exname := parser.SubexpNames()
	for i := 1; i < l; i++ {
		switch exname[i] {
		case "func":
			switch out[i] {
			case "all":
				dt.method = nil
			case "sum":
				dt.method = nil
			case "avg":
				dt.method = nil
			case "diff":
				dt.method = nil
			case "pdiff":
				dt.method = nil
			default:
				return nil, fmt.Errorf("unknow judge method %s", out[i])
			}
		case "formula":
			arg, err := genBinTreeFormula(out[i])
			if err != nil {
				return nil, err
			}
			dt.arg = arg
		case "total":
			if len(out[i]) == 0 {
				return nil, fmt.Errorf("parse total failed %s", out[i])
			}
			t, err := strconv.Atoi(out[i])
			if err != nil {
				return nil, fmt.Errorf("total Atoi failed %s", out[i])
			}
			if t > storeCap {
				t = storeCap
			}
			dt.total = t
		case "limit":
			if len(out[i]) != 0 {
				t, err := strconv.Atoi(out[i])
				if err != nil {
					return nil, fmt.Errorf("total Atoi failed %s", out[i])
				}
				dt.total = t
			}

		case "operator":
			if len(out[i]) == 0 {
				return nil, fmt.Errorf("parse operator failed %s", out[i])
			}
			switch out[i] {
			case ">", "<", "==", "!=":
			default:
				return nil, fmt.Errorf("%s operator invalid ", out[i])
			}
		case "compare":
			if len(out[i]) == 0 {
				return nil, fmt.Errorf("parse compare failed %s", out[i])
			}
			t, err := strconv.ParseFloat(out[i], 32)
			if err != nil {
				return nil, fmt.Errorf("compare ParseFloat failed %s", out[i])
			}
			dt.compare = float64(t)

		}
	}
	return &dt, nil
}

type ringNode struct {
	data  map[string]float64
	pNext *ringNode
}

func makeStore(n int) *ringNode {
	if n <= 0 {
		n = 1
	}
	r := new(ringNode)
	p := r
	for i := 1; i < n; i++ {
		p.pNext = &ringNode{}
		p = p.pNext
	}
	p.pNext = r
	return r
}

func (r *ringNode) next() *ringNode {
	return r.pNext
}

func (r *ringNode) store(data map[string]float64) *ringNode {
	p := r.pNext
	p.data = data
	return p
}

func (r *ringNode) pull(n int) []map[string]float64 {
	var data []map[string]float64
	leng := r.len() + 1
	c := r
	for i := 1; i < leng; i++ {
		c = c.pNext
		if i > leng-n {
			data = append(data, c.data)
		}
	}
	return data
}

func (r *ringNode) len() int {
	n := 0
	if r != nil {
		n = 1
		for p := r.pNext; p != r; p = p.pNext {
			n++
		}
	}
	return n
}
