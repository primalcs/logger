package logger

import (
	"log/syslog"
	"time"

	"github.com/primalcs/logger/listener"
	"github.com/primalcs/logger/types"
	"github.com/primalcs/logger/writer"
)

// Option implements options pattern for Logger
type Option func(*Logger) error

// WithTCPConnection adds a tcp syslog writer to Logger
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

// WithUDPConnection adds a udp syslog writer to Logger
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

// WithLocalWriter adds a local syslog writer to Logger
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

// WithFileWriter creates a log file at specified address
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

// WithNSQWriter creates a simple connection to NSQ
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

// WithLogLevel specifies the maximum allowed log level
func WithLogLevel(level types.LogLevel) Option {
	return func(logger *Logger) error {
		logger.config.SetLogLevel(level)
		return nil
	}
}

// WithDelimiter specifies message delimiter
func WithDelimiter(delimiter string) Option {
	return func(logger *Logger) error {
		logger.config.SetDelimiter(delimiter)
		return nil
	}
}

// WithHttpListener creates new http-server for configuring logger and run it in a new goroutine
func WithHttpListener(port int) Option {
	return func(logger *Logger) error {
		go listener.NewListener(port, logger.config)
		return nil
	}
}

// WithTimeLog specifies time format and location for logs
func WithTimeLog(format string, loc *time.Location) Option {
	return func(logger *Logger) error {
		logger.config.SetTimeOptions(format, loc)
		return nil
	}
}

// WithCaller adds runtime.Caller() info to logs
func WithCaller() Option {
	return func(logger *Logger) error {
		logger.config.SetWithCaller(true)
		return nil
	}
}
