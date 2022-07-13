package handlers

import (
	"sync/atomic"
	"testing"
)

func TestAtomic(t *testing.T) {
	var ops *uint64 = 0
	atomic.AddUint64(ops, 1)
}
