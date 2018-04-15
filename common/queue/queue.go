package queue

import "time"

//Queue interface
//Get retuen value and remove
type Queue interface {
	Len() int
	Put(interface{}) (bool, error)
	Get() (interface{}, bool, error)
	PutWait(interface{}, ...time.Duration) error
	GetWait(...time.Duration) (interface{}, error)
}
