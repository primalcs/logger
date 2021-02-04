package writer_pool

import (
	"context"
	"sync"

	"github.com/primalcs/logger/types"
	"github.com/primalcs/logger/writer"
)

// WriterPool is a thread-safe set of writers and connector serving them
type WriterPool struct {
	mu        sync.RWMutex
	writers   []*writer.Writer
	connector *connector
}

// NewWriterPool creates WriterPool
func NewWriterPool(ctx context.Context) *WriterPool {
	wp := &WriterPool{
		writers:   make([]*writer.Writer, 0),
		connector: NewConnector(types.MaxConnectorQ),
	}
	wp.connector.Run(ctx, types.DefaultReconnectionTime)
	return wp
}

// AddWriter adds existing writer to pool
func (wp *WriterPool) AddWriter(w *writer.Writer) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.writers = append(wp.writers, w)
}

// WriteAll tells all writers to write a message with params; failed writers are sent to connector
func (wp *WriterPool) WriteAll(lp types.LogParams, tp types.TimeParams, mp types.MsgParams, kvs ...string) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	for _, v := range wp.writers {
		if err := v.Write(lp, tp, mp, kvs...); err != nil {
			wp.connector.AddToQ(v)
			continue
		}
	}
}
