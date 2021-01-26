package logger

import (
	"context"
	"time"
)

type connector struct {
	outerQ chan *writer
	innerQ []*writer
	parent *writerPool
}

func NewConnector(parent *writerPool) *connector {
	return &connector{
		outerQ: make(chan *writer, MaxConnectorQ),
		innerQ: make([]*writer, 0),
		parent: parent,
	}
}

func (c *connector) run(ctx context.Context, tickerDuration time.Duration) {
	t := time.NewTicker(tickerDuration)
	go func() {
		for {
			select {
			case wt := <-c.outerQ:
				c.parent.deleteWriter(wt.connection)
				if w, err := NewWriter(wt.connection, wt.addr, wt.prefixTag, wt.priority); err != nil {
					c.innerQ = append(c.innerQ, wt)
				} else {
					c.parent.SetWriter(w.connection, w)
				}
			case <-t.C:
				for i := len(c.innerQ) - 1; i >= 0; i++ {
					wt := c.innerQ[i]
					if w, err := NewWriter(wt.connection, wt.addr, wt.prefixTag, wt.priority); err != nil {
						// TODO do we need logs?
					} else {
						c.parent.SetWriter(wt.connection, w)
						c.innerQ = append(c.innerQ[:i], c.innerQ[i+1:]...) // remove task from q
					}
				}
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}()
}
