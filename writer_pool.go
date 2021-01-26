package logger

import (
	"context"
	"sync"
)

type writerPool struct {
	mu        sync.RWMutex
	writers   map[string]*writer
	connector *connector
}

func NewWriterPool(ctx context.Context) *writerPool {
	wp := &writerPool{
		writers: make(map[string]*writer),
	}
	c := NewConnector(wp)
	wp.connector = c
	c.run(ctx, defaultReconnectionTime)
	return wp
}

func (wp *writerPool) SetWriter(wType string, w *writer) {
	if wType != TCP && wType != UDP && wType != LOCAL {
		return
	}
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.writers[wType] = w
}

func (wp *writerPool) deleteWriter(wType string) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	val, ok := wp.writers[wType]
	if !ok {
		return
	}
	val.logWriter.Close()
	delete(wp.writers, wType)
}

func (wp *writerPool) WriteAll(lp LogParams, tp TimeParams, mp MsgParams, kvs ...string) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	for _, v := range wp.writers {
		if err := v.write(lp, tp, mp, kvs...); err != nil {
			wp.connector.outerQ <- v
			continue
		}
	}
}
