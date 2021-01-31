package logger

import (
	"log/syslog"
	"time"
)

type Option func(*Logger) error

func WithTCPConnection(addr, tag string, priority syslog.Priority, bufferLen int) Option {
	return func(logger *Logger) error {
		w, err := NewWriter(TCP, addr, tag, priority, bufferLen)
		if err != nil {
			return err
		}
		logger.writers.SetWriter(TCP, w)
		return nil
	}
}

func WithUDPConnection(addr, tag string, priority syslog.Priority, bufferLen int) Option {
	return func(logger *Logger) error {
		w, err := NewWriter(UDP, addr, tag, priority, bufferLen)
		if err != nil {
			return err
		}
		logger.writers.SetWriter(UDP, w)
		return nil
	}
}

func WithLocalWriter(tag string, priority syslog.Priority, bufferLen int) Option {
	return func(logger *Logger) error {
		w, err := NewWriter(LOCAL, "", tag, priority, bufferLen)
		if err != nil {
			return err
		}
		logger.writers.SetWriter(LOCAL, w)
		return nil
	}
}

func WithFileWriter(addr, tag string) Option {
	return func(logger *Logger) error {
		w, err := NewWriter(FILE, addr, tag, 0, 1)
		if err != nil {
			return err
		}
		logger.writers.SetWriter(FILE, w)
		return nil
	}
}

func WithLogLevel(level LogLevel) Option {
	return func(logger *Logger) error {
		logger.config.SetLogLevel(level)
		return nil
	}
}

func WithDelimiter(delimiter string) Option {
	return func(logger *Logger) error {
		logger.config.mu.Lock()
		defer logger.config.mu.Unlock()
		logger.config.delimiter = delimiter
		return nil
	}
}

func WithHttpListener(port int) Option {
	return func(logger *Logger) error {
		NewListener(port, logger)
		return nil
	}
}

func WithTimeLog(loc *time.Location, format string) Option {
	return func(logger *Logger) error {
		logger.config.mu.Lock()
		defer logger.config.mu.Unlock()
		logger.config.location = loc
		logger.config.timeFormat = format

		return nil
	}
}
