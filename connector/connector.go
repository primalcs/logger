package connector

import (
	"context"
	"time"

	"github.com/rybnov/logger/writer"

	"github.com/rybnov/logger/types"
)

type Connector struct {
	outerQ chan *writer.Writer
	innerQ []*writer.Writer
}

func NewConnector() *Connector {
	return &Connector{
		outerQ: make(chan *writer.Writer, types.MaxConnectorQ),
		innerQ: make([]*writer.Writer, 0),
	}
}

func (c *Connector) AddToQ(w *writer.Writer) {
	c.outerQ <- w
}

func (c *Connector) Run(ctx context.Context, tickerDuration time.Duration) {
	t := time.NewTicker(tickerDuration)
	go func() {
		for {
			select {
			case wt := <-c.outerQ:
				wt.Stop(true)
				if err := wt.Reconnect(); err != nil {
					c.innerQ = append(c.innerQ, wt)
				}
			case <-t.C:
				for i := len(c.innerQ) - 1; i >= 0; i++ {
					wt := c.innerQ[i]
					if err := wt.Reconnect(); err != nil {
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
