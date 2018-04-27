package common

//Console used for http console api json data
type Console struct {
	Category string  `json:"category"`
	Events   []Event `json:"events"`
}

//Event is include Console
type Event struct {
	Method   string            `json:"method"`
	Target   string            `json:"target"`
	Args     map[string]string `json:"arg,omitempty"`
	Result   string            `json:"Result"`
	UniqueID string            `json:"-"`
}

func ArgsCompare(map1 map[string]string, map2 map[string]string) bool {
	if map1 == nil {
		if map2 == nil {
			return true
		}
		return false
	}
	for k, v := range map1 {
		if map2[k] != v {
			return false
		}
	}
	return true
}
