package dango

import "testing"

func TestBits(t *testing.T) {
	const (
		S1 uint = 1 << iota
		S2
		S3
		S4
		S5
		S6
	)
	b := NewBits()
	if b.Value() != 0 {
		t.Errorf("Expect 0, got %d", b.Value())
	}
	b.Set(S2 | S3 | S5)
	if b.Value() != 22 {
		t.Errorf("Expect 22, got %d", b.Value())
	}
	b.Clear(S3)
	if b.Value() != 18 {
		t.Errorf("Expect 18, got %d", b.Value())
	}
	if !b.Has(S2) {
		t.Errorf("Expect True, got %t, %d != %d", S2 == S2&(S2|S5), S2, S2&(S2|S5))
	}
	if !b.Has(S2 | S3) {
		t.Errorf("Expect True, got False")
	}

	if b.HasAll(S2 | S3) {
		t.Errorf("Expect False, got True")
	}

}
