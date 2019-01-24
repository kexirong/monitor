package judge

import (
	"math"
)

/*
函数签名: func (datas []float64, operator string, compare float64, limit int) bool
limit 为满足判定条件数据条数，非必用
datas 数据顺序为:旧到新,datas[len(datas)-1]为最新数据
*/

func checkCompare(value, compare float64, operator string) bool {
	switch operator {
	case ">":
		return value > compare
	case "<":
		return value < compare
	case "==":
		return math.Abs(value-compare) < 0.000001
	case "!=":
		return math.Abs(value-compare) > 0.000001
	default:
		return false
	}
}

func all(datas []float64, operator string, compare float64, limit int) bool {
	var count int
	for _, d := range datas {
		if checkCompare(d, compare, operator) {
			count++
		}
	}

	return count >= limit
}
func sum(datas []float64, operator string, compare float64, limit int) bool {
	var sum float64
	for _, d := range datas {
		sum += d
	}
	return checkCompare(sum, compare, operator)
}
func avg(datas []float64, operator string, compare float64, limit int) bool {
	var sum float64
	for _, d := range datas {
		sum += d
	}
	return checkCompare(sum/float64(len(datas)), compare, operator)
}

func diff(datas []float64, operator string, compare float64, limit int) bool {
	var count int
	var new = datas[len(datas)-1]
	for _, d := range datas[:len(datas)-1] {
		if checkCompare(new-d, compare, operator) {
			count++
		}
	}
	return count >= limit
}

func pDiff(datas []float64, operator string, compare float64, limit int) bool {
	var count int
	var new = datas[len(datas)-1]
	for _, d := range datas[:len(datas)-1] {
		if checkCompare((new-d)/d*100, compare, operator) {
			count++
		}
	}
	return count >= limit
}
