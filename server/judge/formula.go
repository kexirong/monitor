package judge

import (
	"bytes"
	"errors"
	"strconv"
)

type binTreeNode struct {
	key    string
	flag   nodeType
	pLeft  *binTreeNode
	pRight *binTreeNode
}

func (b *binTreeNode) left() *binTreeNode {
	return b.pLeft
}

func (b *binTreeNode) right() *binTreeNode {
	return b.pRight
}

func (b *binTreeNode) toFormula() string {
	var s string
	if b != nil {
		s += b.pLeft.toFormula()
		s += b.key
		s += b.pRight.toFormula()
		if b.flag == operator {
			s = "(" + s + ")"
		}
	}
	return s
}

func (b *binTreeNode) calculate(value map[string]float64) float64 {
	var cul func(b *binTreeNode) float64
	cul = func(b *binTreeNode) float64 {
		switch b.flag {
		case operator:
			switch b.key {
			case "+":
				return cul(b.pLeft) + cul(b.pRight)
			case "-":
				return cul(b.pLeft) - cul(b.pRight)
			case "*":
				return cul(b.pLeft) * cul(b.pRight)
			case "/":
				return cul(b.pLeft) / cul(b.pRight)
			}
		case variable:
			return value[b.key]
		case digit:
			v, _ := strconv.ParseFloat(b.key, 32)
			return v
		}
		return 0
	}
	return cul(b)
}

func isLower(top, newTop string) bool {
	// 注意 a + b + c 的后缀表达式是 ab + c +，不是 abc + +
	switch top {
	case "+", "-":
		if newTop == "*" || newTop == "/" {
			return true
		}
	case "(":
		return true
	}
	return false
}

//formula 不支持 负数 小数
func genBinTreeFormula(formula string) (*binTreeNode, error) {
	bs := []byte(formula)
	stack1, stack2 := stack{}, stack{}
	//root := &binTreeNode{}
	l := len(bs)
	for i := 0; i < l; i++ {
		switch bs[i] {
		case '(':
			stack2.Push(&binTreeNode{
				key:  "(",
				flag: bracket,
			})
		case ')':
			for !stack2.IsEmpty() {
				item := stack2.Pop()
				if item.key == "(" {
					// 弹出 "("
					break
				}
				if v1, v2 := stack1.Pop(), stack1.Pop(); v1 != nil && v2 != nil {
					item.pLeft = v2
					item.pRight = v1
					stack1.Push(item)
				} else {
					return nil, errors.New("invalid formula format")
				}

			}
		case '$':
			if bs[i+1] != '{' {
				return nil, errors.New("invalid formula format")
			}
			n := bytes.IndexByte(bs[i+2:], '}')
			if n == -1 || n == i+2 {
				return nil, errors.New("invalid formula format")
			}

			stack1.Push(&binTreeNode{
				key:  string(bs[i+2 : i+2+n]),
				flag: variable,
			})
			i = i + 2 + n
		case '+', '-', '*', '/':
			for !stack2.IsEmpty() {
				top := stack2.Top()
				if top.key == "(" || isLower(top.key, string(bs[i])) {
					break
				}
				if v1, v2 := stack1.Pop(), stack1.Pop(); v1 != nil && v2 != nil {
					top.pLeft = v2
					top.pRight = v1
					stack1.Push(top)
					//stack2.Pop()
				} else {
					return nil, errors.New("invalid formula format")
				}
				stack2.Pop()
			}
			// 低优先级的运算符入栈
			stack2.Push(&binTreeNode{
				key:  string(bs[i]),
				flag: operator,
			})
		case ' ':
			continue
		default:
			n := 0
			for ; n < l-i; n++ {
				if bs[i+n] > 47 && bs[i+n] < 58 {
					continue
				}
				if n == 0 {
					return nil, errors.New("invalid formula format")
				}
				break
			}
			stack1.Push(&binTreeNode{
				key:  string(bs[i : i+n]),
				flag: digit,
			})
			i += n
		}
	}
	for !stack2.IsEmpty() {
		item := stack2.Pop()
		if v1, v2 := stack1.Pop(), stack1.Pop(); v1 != nil && v2 != nil {
			item.pLeft = v2
			item.pRight = v1
			stack1.Push(item)
		} else {
			return nil, errors.New("invalid formula format")
		}
	}

	if stack1.Len() != 1 {
		return nil, errors.New("invalid formula format")
	}
	r := stack1.Pop()

	return r, nil
}

type stack struct {
	items []*binTreeNode
}

func (s *stack) Push(item *binTreeNode) {
	s.items = append(s.items, item)
}

func (s *stack) Pop() *binTreeNode {
	item := s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]
	return item
}
func (s *stack) Len() int {
	return len(s.items)
}

func (s *stack) Top() *binTreeNode {
	return s.items[len(s.items)-1]
}

func (s *stack) IsEmpty() bool {
	return len(s.items) == 0
}
