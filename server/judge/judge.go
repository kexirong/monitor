package judge

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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
const reExpr = `(?P<func>\w+)\((?P<formula>\(?[^,]+\)?),(?P<total>\d+),?(?P<limit>\d*)?\)(?P<operator>[<!=>]+)(?P<compare>[-]?\d+(?:\.\d+)?$)`
const storeCap = 12

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
	//fmt.Println(dt)
	if err == nil {
		dt.level = aj.Level
		dt.title = aj.Title
		judger.deters[aj.Express] = dt
		j.judgers[aj.AnchorPoint] = judger
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
func (j *Judge) DoJudge(tp *packetparse.TargetPacket) (result []*models.AlarmEvent) {
	//anchorPoint 格式，[*].cpu.cpu0:percent
	//var results []*models.AlarmEvent
	var point string
	if tp.Instance != "" {
		point = fmt.Sprintf("%s.%s:%s", tp.Plugin, tp.Instance, tp.Type)
	} else {
		point = fmt.Sprintf("%s:%s", tp.Plugin, tp.Type)
	}
	j.mu.Lock()
	anchorPoint := fmt.Sprintf("[%s].%s", tp.HostName, point)
	judger := j.judgers[anchorPoint]
	if judger == nil {
		anchorPoint = "[*]." + point
		judger = j.judgers[anchorPoint]
	}
	j.mu.Unlock()

	if judger != nil {
		data := make(map[string]float64)

		sl := strings.Split(tp.VlTags, "|")

		for idx, value := range tp.Value {
			data[sl[idx]] = value

		}
		result = judger.doJudge(data)
		for i := 0; i < len(result); i++ {
			result[i].HostName = tp.HostName
			result[i].AnchorPoint = anchorPoint
			result[i].Stat = 1
			result[i].Count = 1
			result[i].Message = tp.Message
		}
	}
	return
}

/*
warning 警告
serious 严重
disaster 灾难
*/
type judger struct {
	store *ringNode
	//level  models.Level
	//title  string
	deters map[string]*deter //[key]*deter
}

func newJuger() *judger {
	return &judger{
		store:  makeStore(storeCap),
		deters: make(map[string]*deter),
	}
}

type deter struct {
	level    models.Level
	title    string
	total    int
	limit    int
	operator string
	compare  float64
	arg      *binTreeNode
	method   func(datas []float64, operator string, compare float64, limit int) bool
}

func (j *judger) doJudge(data map[string]float64) (result []*models.AlarmEvent) {
	//result = make(map[string]float64)
	j.store = j.store.store(data)
	for key, deter := range j.deters {
		ds := j.store.pull(deter.total)
		if len(ds) < deter.total {
			continue
		}
		datas := make([]float64, deter.total)

		for i, d := range ds {
			datas[i] = deter.arg.calculate(d)
		}

		if deter.method(datas, deter.operator, deter.compare, deter.limit) {
			//result[key] = datas[len(datas)-1]
			result = append(result, &models.AlarmEvent{Rule: key, Value: datas[len(datas)-1], Title: deter.title, Level: deter.level})
		} //else {
		//fmt.Printf("deter.method(datas:%v, operator:%v, compare:%v, limit:%v)\n", datas, deter.operator, deter.compare, deter.limit)
		//}
	}
	return result
}

func parseExpress(express string, parser *regexp.Regexp) (*deter, error) {
	if parser == nil {
		return nil, errors.New("regexp is nil")
	}
	out := parser.FindStringSubmatch(strings.Replace(express, " ", "", -1))
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
				dt.method = all
			case "sum":
				dt.method = sum
			case "avg":
				dt.method = avg
			case "diff":
				dt.method = diff
			case "pdiff":
				dt.method = pDiff
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
				dt.limit = t
			}

		case "operator":
			if len(out[i]) == 0 {
				return nil, fmt.Errorf("parse operator failed %s", out[i])
			}
			switch out[i] {
			case ">", "<", "==", "!=":
				dt.operator = out[i]
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
	leng := r.len()
	c := r
	for i := 1; i < leng+1; i++ {
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
