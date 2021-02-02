package writer_pool

import (
	"context"
	"sync"

	"github.com/rybnov/logger/types"

	"github.com/rybnov/logger/connector"

	"github.com/rybnov/logger/writer"
)

type WriterPool struct {
	mu        sync.RWMutex
	writers   []*writer.Writer
	connector *connector.Connector
}

func NewWriterPool(ctx context.Context) *WriterPool {
	wp := &WriterPool{
		writers:   make([]*writer.Writer, 0),
		connector: connector.NewConnector(),
	}
	wp.connector.Run(ctx, types.DefaultReconnectionTime)
	return wp
}

func (wp *WriterPool) AddWriter(w *writer.Writer) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.writers = append(wp.writers, w)
}

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
