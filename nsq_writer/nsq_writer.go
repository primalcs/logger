package nsq_writer

import (
	"github.com/nsqio/go-nsq"
)

type nsqWriter struct {
	producer *nsq.Producer
	topic    string
}

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

func (n *nsqWriter) Emerg(m string) error   { _, err := n.Write([]byte(m)); return err }
func (n *nsqWriter) Alert(m string) error   { _, err := n.Write([]byte(m)); return err }
func (n *nsqWriter) Crit(m string) error    { _, err := n.Write([]byte(m)); return err }
func (n *nsqWriter) Err(m string) error     { _, err := n.Write([]byte(m)); return err }
func (n *nsqWriter) Warning(m string) error { _, err := n.Write([]byte(m)); return err }
func (n *nsqWriter) Notice(m string) error  { _, err := n.Write([]byte(m)); return err }
func (n *nsqWriter) Info(m string) error    { _, err := n.Write([]byte(m)); return err }
func (n *nsqWriter) Debug(m string) error   { _, err := n.Write([]byte(m)); return err }

func (n *nsqWriter) Write(ba []byte) (int, error) {
	if err := n.producer.Publish(n.topic, ba); err != nil {
		return 0, err
	}

	return len(ba), nil
}

func (n *nsqWriter) Close() error {
	n.producer.Stop()
	return nil
}
