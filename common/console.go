package common

//Console used for http console api json data
type Console struct {
	Category string  `json:"category"`
	Events   []Event `json:"events"`
}

//Event is include Console
type Event struct {
	Method string `json:"method"`
	Target string `json:"target"`
	Arg    string `json:"arg,omitempty"`
	Result string `json:"Result"`
}
