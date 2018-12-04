package judge

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func Test_judge(t *testing.T) {
	var testValue = map[string]float64{"feild": 15, "feild1": 12, "feild2": 8, "feild3": 10, "feild4": 6}

	express1 := "diff(${feild1}+${feild3},3,3)>1.23"
	express2 := "avg(((${feild1}+(${feild}-${feild4})*${feild2})*${feild3}), 3) > -1.23"
	express3 := "avg(${feild1}/${feild2},3 , 2) > -1.23"
	_, _, _ = express1, express2, express3
	re1, err := regexp.Compile(`(?P<func>\w+)\((?P<formula>\(?[^,]+\)?),(?P<total>\d+),?(?P<limit>\d*)?\)(?P<operator>[<!=>]+)(?P<compare>[-]?\d+(?:\.\d+)?$)`)
	if err != nil {
		t.Error(err)
	}
	n1 := re1.SubexpNames()
	t.Log(n1)
	result1 := re1.FindStringSubmatch(express2)

	(strings.Replace(express2, " ", "", -1))
	var formula *binTreeNode
	for k, v := range result1 {
		fmt.Println(n1[k], ":", v)
		if n1[k] == "formula" {
			var err error
			formula, err = genBinTreeFormula(v)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	fmt.Println("genFormula: ", formula.toFormula(), "=", formula.calculate(testValue))
}
