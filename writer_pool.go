package logger

import (
	"context"
	"sync"
	"time"
)

var defaultReconnectionTime = time.Minute

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
	if wType != TCP && wType != UDP && wType != Local {
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
	val.w.Close()
	delete(wp.writers, wType)

}

func (wp *writerPool) WriteAll(msg string) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	for _, v := range wp.writers {
		_, err := v.w.Write([]byte(msg))
		if err != nil {
			wp.connector.outerQ <- v
			continue
		}
	}
}
