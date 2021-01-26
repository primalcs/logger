package logger

import "context"

type Logger struct {
	config  *config
	writers *writerPool
}

func NewLogger(ctx context.Context, opts ...Option) (*Logger, error) {
	l := &Logger{
		config:  NewConfig(),
		writers: NewWriterPool(ctx),
	}
	for _, opt := range opts {
		if err := opt(l); err != nil {
			return nil, err
		}
	}

	return l, nil
}
