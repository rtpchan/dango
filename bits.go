// bits.go provide bitmask functions with mutex lock
// Example:
// const {
// 	state1 = 1 << iota
// 	state2
// 	state3
// }

// b = NewBits()
// b.Set(state2|state3)

package dango

import (
	"fmt"
	"sync"
)

// Convert bit position to uint, right most bit is bit 1,
func PosToUint(p uint) uint {
	return 1 << (p - 1)
}

type Bits struct {
	bits uint // 32 or 64 bit depends on platform
	mu   sync.Mutex
}

func NewBits() *Bits {
	return &Bits{}
}

func (b *Bits) Set(n uint) {
	b.mu.Lock()
	b.bits = b.bits | n
	b.mu.Unlock()
}

func (b *Bits) Clear(n uint) {
	b.mu.Lock()
	b.bits = b.bits &^ n
	b.mu.Unlock()
}

func (b *Bits) Toggle(n uint) {
	b.mu.Lock()
	b.bits = b.bits ^ n
	b.mu.Unlock()
}

// Has matches any bit position
func (b *Bits) Has(n uint) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	matched := b.bits & n
	return matched != 0
}

// HasAll matches all bit positions
func (b *Bits) HasAll(n uint) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	matched := b.bits & n
	return matched == n
}

func (b *Bits) Value() uint {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.bits
}

func (b *Bits) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return fmt.Sprint(b.bits)
}
