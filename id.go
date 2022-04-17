// Generate unique ID

package dango

import "sync"

type ID struct {
	current int
	lock    sync.Mutex
}

func NewIDGenerator() *ID {
	return &ID{current: 0}
}

func (i *ID) NewID() int {
	i.lock.Lock()
	i.current++
	i.lock.Unlock()
	return i.current
}

func (i *ID) Reset() {
	i.lock.Lock()
	i.current = 0
	i.lock.Unlock()
}
