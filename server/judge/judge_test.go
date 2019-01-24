package judge

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func Test_Regexp(t *testing.T) {
	expresss := []string{
		"diff(${feild1}+${feild3},3,3)>1.23",
		"avg(((${feild1}+(${feild}-${feild4})*${feild2})*${feild3}), 3) > -1.23",
		"all(100-${idle},10,10)>1",
	}

	re, err := regexp.Compile(reExpr)

	if err != nil {
		t.Fatal(err)
	}
	n := re.SubexpNames()

	for _, express := range expresss {
		express = strings.Replace(express, " ", "", -1)
		fmt.Println(express)
		result := re.FindStringSubmatch(express)
		if len(result) == 0 {
			t.Error("express parse failed")
		}
		out := "\n"
		for k, v := range result {
			if k == 0 {
				out += v
			} else {
				out = out + "\n" + n[k] + ":" + v
			}

		}
		t.Log(out)
	}

}
func Test_genBinTreeFormula(t *testing.T) {
	formulas := []string{
		"100-${idle}",
		"((${feild1}+(${feild}-${feild4})*${feild2})*${feild3})",
		"${feild1}+${feild3}",
	}

	for _, formula := range formulas {
		v, err := genBinTreeFormula(formula)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(v.toFormula())
		}
	}

}

func Test_parseExpress(t *testing.T) {
	var testValue = map[string]float64{"feild": 15, "feild1": 12, "feild2": 8, "feild3": 10, "feild4": 6}

	expresss := []string{
		"diff(${feild1}+${feild3},3,3)>1.23",
		"avg(((${feild1}+(${feild}-${feild4})*${feild2})*${feild3}), 3) > 1.23",
		"all(100-${idle},10,10)>1",
	}

	re, err := regexp.Compile(reExpr)
	if err != nil {
		t.Fatal(err)
	}

	for _, express := range expresss {
		dt, err := parseExpress(express, re)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(dt)
			t.Log(dt.arg.toFormula(), "=>", dt.arg.calculate(testValue))
		}
	}

}
