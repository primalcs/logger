package logger

import (
	"context"
	"log/syslog"
	"time"
)

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

func NewDefaultLogger() (*Logger, error) {
	opts := []Option{
		WithDelimiter(DelimiterV),
		WithFileWriter("/var/log/logger/", "file"),
		WithHttpListener(8080),
		WithLocalWriter("local", syslog.LOG_DEBUG|syslog.LOG_SYSLOG, 1),
		WithLogLevel(DEBUG),
		WithTimeLog(time.UTC, time.RFC3339),
	}

	l, err := NewLogger(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (lg *Logger) Log(level LogLevel, tag, msg string, kvs ...string) {
	var isForced bool
	maxLvl := lg.config.GetLogLevel()
	// ignore log
	if level < FORCE && level > maxLvl {
		return
	}
	// log is forced
	if level > DEBUG {
		isForced = true
		level = level & 0b111
	}
	lg.writers.WriteAll(
		LogParams{
			IsForced: isForced,
			Level:    level,
		},
		TimeParams{
			Location: lg.config.location,
			Format:   lg.config.timeFormat,
		},
		MsgParams{
			Delimiter: lg.config.delimiter,
			Tag:       tag,
			Msg:       msg,
		},
		kvs...)
}

func (lg *Logger) AddOptions(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(lg); err != nil {
			return err
		}
	}
	return nil
}
