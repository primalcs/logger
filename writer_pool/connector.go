package writer_pool

import (
	"context"
	"time"

	"github.com/primalcs/logger/writer"
)

type connector struct {
	outerQ chan *writer.Writer
	innerQ []*writer.Writer
}

// NewConnector creates new instance of connector with qSize-len buffered chan for queue
func NewConnector(qSize int) *connector {
	return &connector{
		outerQ: make(chan *writer.Writer, qSize),
		innerQ: make([]*writer.Writer, 0),
	}
}

// AddToQ adds disconnected writer in queue (buffered chan) for reconnection
func (c *connector) AddToQ(w *writer.Writer) {
	c.outerQ <- w
}

// Run starts the connector loop in a new go-routine;
// loop reconnects writers from AddToQ, if fail - sends writer to inner queue
// every tickerDuration innerQ is checked and tried to reconnect
func (c *connector) Run(ctx context.Context, tickerDuration time.Duration) {
	t := time.NewTicker(tickerDuration)
	go func() {
		for {
			select {
			case wt, ok := <-c.outerQ:
				if !ok {
					return
				}
				wt.Stop()
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
				c.close()
				return
			default:
				time.Sleep(time.Millisecond * 1000) // TODO make configurable
			}
		}
	}()
}

func (c *connector) close() {
	close(c.outerQ)
}
