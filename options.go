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
		w, err := writer.NewWriter(types.TCP, addr, tag, priority, bufferLen)
		if err != nil {
			return err
		}
		logger.writers.SetWriter(types.TCP, w)
		return nil
	}
}

func WithUDPConnection(addr, tag string, priority syslog.Priority, bufferLen int) Option {
	return func(logger *Logger) error {
		w, err := writer.NewWriter(types.UDP, addr, tag, priority, bufferLen)
		if err != nil {
			return err
		}
		logger.writers.SetWriter(types.UDP, w)
		return nil
	}
}

func WithLocalWriter(tag string, priority syslog.Priority, bufferLen int) Option {
	return func(logger *Logger) error {
		w, err := writer.NewWriter(types.LOCAL, "", tag, priority, bufferLen)
		if err != nil {
			return err
		}
		logger.writers.SetWriter(types.LOCAL, w)
		return nil
	}
}

func WithFileWriter(addr, tag string) Option {
	return func(logger *Logger) error {
		w, err := writer.NewWriter(types.FILE, addr, tag, 0, 1)
		if err != nil {
			return err
		}
		logger.writers.SetWriter(types.FILE, w)
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
		listener.NewListener(port, logger.config)
		return nil
	}
}

func WithTimeLog(format string, loc *time.Location) Option {
	return func(logger *Logger) error {
		logger.config.SetTimeOptions(format, loc)
		return nil
	}
}
