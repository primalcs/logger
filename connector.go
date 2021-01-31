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
				c.parent.stopWriter(wt.connection)
				if err := wt.reconnect(wt.connection, wt.addr, wt.prefixTag, wt.priority); err != nil {
					c.innerQ = append(c.innerQ, wt)
				}
			case <-t.C:
				for i := len(c.innerQ) - 1; i >= 0; i++ {
					wt := c.innerQ[i]
					if err := wt.reconnect(wt.connection, wt.addr, wt.prefixTag, wt.priority); err != nil {
						// TODO do we need logs?
					} else {
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
