package pipeline

import (
	"container/ring"
	"log"
	"os"
	"sync"
)

type Buffer struct {
	m          sync.Mutex
	ringBuffer *ring.Ring
	cap        int
	pos        int
	elog       *log.Logger
	ilog       *log.Logger
}

func NewRingBuffer(l int) *Buffer {
	r := ring.New(l)
	return &Buffer{
		ringBuffer: r,
		cap:        r.Len(),
		pos:        -1,
		m:          sync.Mutex{},
		ilog:       log.New(os.Stdout, "Buffer INFO\t", log.Ldate|log.Ltime),
		elog:       log.New(os.Stdout, "Buffer ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
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
	b.ilog.Printf("В буффер добавлены данные: %v", d)
}

func (b *Buffer) Get() *ring.Ring {
	r := ring.New(b.pos + 1)
	b.m.Lock()
	for i := 0; i < (b.cap - (b.pos + 1)); i++ {
		b.ringBuffer = b.ringBuffer.Next()
	}
	b.ilog.Printf("Из буффера получены данные: \n")
	for i := 0; i < (b.pos + 1); i++ {
		r.Value = b.ringBuffer.Value
		b.ilog.Printf("\t%v", r.Value)
		r = r.Next()
		b.ringBuffer = b.ringBuffer.Next()
	}
	b.pos = -1
	b.m.Unlock()
	return r
}
