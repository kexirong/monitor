package queue

import "time"

//Queue interface
type Queue interface {
	Len() int
	Put([]byte) (bool, error)
	Get() ([]byte, bool, error)
	PutWait([]byte, ...time.Duration) error
	GetWait(...time.Duration) ([]byte, error)
}
