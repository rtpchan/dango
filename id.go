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
	defer i.lock.Unlock()
	return i.current
}

func (i *ID) Reset() {
	i.lock.Lock()
	i.current = 0
	i.lock.Unlock()
}

// Up to the caller to ensure id is still unique
func (i *ID) SetCurrent(n int) {
	i.lock.Lock()
	i.current = n
	i.lock.Unlock()
}

func (i *ID) Current() int {
	return i.current
}
