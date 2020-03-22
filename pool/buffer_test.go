package pool

import (
	"testing"
)

func TestBufferPutGet(t *testing.T) {
	pool := NewBufferPool(3)

	b1 := pool.Get()
	b2 := pool.Get()
	b3 := pool.Get()

	pool.Put(b1)
	pool.Put(b2)
	pool.Put(b3)
}

func TestBufferOverflow(t *testing.T) {
	pool := NewBufferPool(3)

	b1 := pool.Get()
	b2 := pool.Get()
	b3 := pool.Get()
	b4 := pool.Get()

	pool.Put(b1)
	pool.Put(b2)
	pool.Put(b3)
	pool.Put(b4)
}
