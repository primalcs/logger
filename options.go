package logger

import "log/syslog"

type Option func(*Logger) error

func WithTCPConnection(addr, tag string, priority syslog.Priority) Option {
	return func(logger *Logger) error {
		w, err := NewWriter(TCP, addr, tag, priority)
		if err != nil {
			return err
		}
		logger.writers.SetWriter(TCP, w)
		return nil
	}
}

func WithUDPConnection(addr, tag string, priority syslog.Priority) Option {
	return func(logger *Logger) error {
		w, err := NewWriter(UDP, addr, tag, priority)
		if err != nil {
			return err
		}
		logger.writers.SetWriter(UDP, w)
		return nil
	}
}

func WithLocalWriter(addr, tag string, priority syslog.Priority) Option {
	return func(logger *Logger) error {
		w, err := NewWriter(Local, addr, tag, priority)
		if err != nil {
			return err
		}
		logger.writers.SetWriter(Local, w)
		return nil
	}
}
