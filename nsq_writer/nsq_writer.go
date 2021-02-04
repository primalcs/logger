package nsq_writer

import (
	"github.com/nsqio/go-nsq"
	"github.com/primalcs/logger/types"
)

type nsqWriter struct {
	producer *nsq.Producer
	topic    string
}

// New NSQWriter configures
func NewNSQWriter(addr, topic string) (*nsqWriter, error) {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer(addr, config)
	if err != nil {
		return nil, err
	}
	return &nsqWriter{
		producer: p,
		topic:    topic,
	}, err
}

// WriteForced for nsqWriter is the same as Write(). Needed to implement interface
func (n *nsqWriter) WriteForced(_ types.LogLevel, ba []byte) (int, error) {
	return n.Write(ba)
}

// Write publishes bytes to predefined topic of nsq
func (n *nsqWriter) Write(ba []byte) (int, error) {
	if err := n.producer.Publish(n.topic, ba); err != nil {
		return 0, err
	}

	return len(ba), nil
}

// Close stops the nsq producer
func (n *nsqWriter) Close() error {
	n.producer.Stop()
	return nil
}
