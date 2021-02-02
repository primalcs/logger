package logger

import (
	"log/syslog"
	"time"

	"github.com/rybnov/logger/listener"
	"github.com/rybnov/logger/writer"

	"github.com/rybnov/logger/types"
)

type Option func(*Logger) error

func WithTCPConnection(addr, tag string, priority syslog.Priority, bufferLen int) Option {
	return func(logger *Logger) error {
		w, err := writer.NewWriter(types.ConnectionTCP, addr, tag, priority, bufferLen)
		if err != nil {
			return err
		}
		logger.writers.AddWriter(w)
		return nil
	}
}

func WithUDPConnection(addr, tag string, priority syslog.Priority, bufferLen int) Option {
	return func(logger *Logger) error {
		w, err := writer.NewWriter(types.ConnectionUDP, addr, tag, priority, bufferLen)
		if err != nil {
			return err
		}
		logger.writers.AddWriter(w)
		return nil
	}
}

func WithLocalWriter(tag string, priority syslog.Priority, bufferLen int) Option {
	return func(logger *Logger) error {
		w, err := writer.NewWriter(types.ConnectionLOCAL, "", tag, priority, bufferLen)
		if err != nil {
			return err
		}
		logger.writers.AddWriter(w)
		return nil
	}
}

func WithFileWriter(addr string) Option {
	return func(logger *Logger) error {
		w, err := writer.NewWriter(types.ConnectionFILE, addr, "", 0, 1)
		if err != nil {
			return err
		}
		logger.writers.AddWriter(w)
		return nil
	}
}

func WithNSQWriter(addr, topic string) Option {
	return func(logger *Logger) error {
		w, err := writer.NewWriter(types.ConnectionNSQ, addr, topic, 0, 1)
		if err != nil {
			return err
		}
		logger.writers.AddWriter(w)
		return nil
	}
}

func WithLogLevel(level types.LogLevel) Option {
	return func(logger *Logger) error {
		logger.config.SetLogLevel(level)
		return nil
	}
}

func WithDelimiter(delimiter string) Option {
	return func(logger *Logger) error {
		logger.config.SetDelimiter(delimiter)
		return nil
	}
}

func WithHttpListener(port int) Option {
	return func(logger *Logger) error {
		go listener.NewListener(port, logger.config)
		return nil
	}
}

func WithTimeLog(format string, loc *time.Location) Option {
	return func(logger *Logger) error {
		logger.config.SetTimeOptions(format, loc)
		return nil
	}
}
