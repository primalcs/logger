package logger

import (
	"context"
	"log/syslog"
	"time"

	"github.com/primalcs/logger/config"
	"github.com/primalcs/logger/types"
	"github.com/primalcs/logger/writer_pool"
)

// Logger is the main logging structure of this package
type Logger struct {
	config  *config.Config
	writers *writer_pool.WriterPool
}

// NewLogger creates a new logging instance with options
func NewLogger(ctx context.Context, opts ...Option) (*Logger, error) {
	l := &Logger{
		config:  config.NewConfig(),
		writers: writer_pool.NewWriterPool(ctx),
	}
	for _, opt := range opts {
		if err := opt(l); err != nil {
			return nil, err
		}
	}

	return l, nil
}

// NewDefaultLogger creates new logger with default options:
// delimiter: " | "
// fileWriter in "/var/log/logger/"
// httpListener on 8080 port
// localWriter with ring buffer length=1
// loglevel Debug
// and time format RFC3339 in UTC location
// mostly for testing purposes
func NewDefaultLogger() (*Logger, error) {
	opts := []Option{
		WithDelimiter(types.DelimiterV),
		WithFileWriter("/var/log/logger/"),
		WithHttpListener(8080),
		WithLocalWriter("local", syslog.LOG_DEBUG|syslog.LOG_SYSLOG, 1),
		WithLogLevel(types.DEBUG),
		WithTimeLog(time.RFC3339, time.UTC),
	}

	l, err := NewLogger(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Log logs message to all defined writers
func (lg *Logger) Log(level types.LogLevel, tag, msg string, kvs ...string) {
	var isForced bool
	maxLvl := lg.config.GetLogLevel()
	// ignore log
	if level < types.FORCE && level > maxLvl {
		return
	}
	// log is forced
	if level > types.DEBUG {
		isForced = true
		level &= 0b111
	}

	ft, loc := lg.config.GetTimeOptions()
	tp := types.TimeParams{
		Location: loc,
		Format:   ft,
	}

	lg.writers.WriteAll(
		types.LogParams{
			IsForced:     isForced,
			Level:        level,
			IsWithCaller: lg.config.GetWithCaller(),
		},
		tp,
		types.MsgParams{
			Delimiter: lg.config.GetDelimiter(),
			Tag:       tag,
			Msg:       msg,
		},
		kvs...)
}

// AddOptions options add new options to existing logger
// should be used before first call of logger
func (lg *Logger) AddOptions(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(lg); err != nil {
			return err
		}
	}
	return nil
}
