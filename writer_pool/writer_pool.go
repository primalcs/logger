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
	writers   map[string]*writer.Writer
	connector *connector.Connector
}

func NewWriterPool(ctx context.Context) *WriterPool {
	wp := &WriterPool{
		writers:   make(map[string]*writer.Writer),
		connector: connector.NewConnector(),
	}
	wp.connector.Run(ctx, types.DefaultReconnectionTime)
	return wp
}

func (wp *WriterPool) SetWriter(wType string, w *writer.Writer) {
	if wType != types.TCP && wType != types.UDP && wType != types.LOCAL {
		return
	}
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.writers[wType] = w
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
