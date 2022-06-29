package pipeline

import (
	"container/ring"
	"sync"
)

type Buffer struct {
	m          sync.Mutex
	ringBuffer *ring.Ring
	cap        int
	pos        int
}

func NewRingBuffer(l int) *Buffer {
	r := ring.New(l)
	return &Buffer{ringBuffer: r, cap: r.Len(), pos: -1, m: sync.Mutex{}}
}

func (b *Buffer) Pos() int {
	b.m.Lock()
	defer b.m.Unlock()
	return b.pos
}

func (b *Buffer) Push(d int) {
	b.m.Lock()
	b.ringBuffer.Value = d
	b.ringBuffer = b.ringBuffer.Next()
	if b.pos < (b.cap - 1) {
		b.pos++
	}
	b.m.Unlock()
}

func (b *Buffer) Get() *ring.Ring {
	r := ring.New(b.pos + 1)
	b.m.Lock()
	for i := 0; i < (b.cap - (b.pos + 1)); i++ {
		b.ringBuffer = b.ringBuffer.Next()
	}
	for i := 0; i < (b.pos + 1); i++ {
		r.Value = b.ringBuffer.Value
		r = r.Next()
		b.ringBuffer = b.ringBuffer.Next()
	}
	b.pos = -1
	b.m.Unlock()
	return r
}
