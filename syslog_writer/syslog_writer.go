package syslog_writer

import (
	"log/syslog"

	"github.com/primalcs/logger/types"
)

type sysLogWriter struct {
	slw *syslog.Writer
}

// NewSysLogWriter creates a new connection to
func NewSysLogWriter(connection, addr, prefix string, priority syslog.Priority) (*sysLogWriter, error) {
	var d *syslog.Writer
	var err error
	switch connection {
	case types.ConnectionTCP, types.ConnectionUDP:
		if d, err = syslog.Dial(connection, addr, priority, prefix); err != nil {
			return nil, err
		}
	case types.ConnectionLOCAL:
		if d, err = syslog.New(priority, prefix); err != nil {
			return nil, err
		}
	}
	return &sysLogWriter{slw: d}, nil
}

// Write performs standard syslog writing
func (s *sysLogWriter) Write(ba []byte) (int, error) {
	return s.Write(ba)
}

// WriteForced ignores current syslog level and forces writing with level specified
func (s *sysLogWriter) WriteForced(level types.LogLevel, ba []byte) (int, error) {
	writeFunc := func(m string) error {
		_, err := s.slw.Write(ba)
		return err
	}
	switch level {
	case types.EMERG:
		writeFunc = s.slw.Emerg
	case types.ALERT:
		writeFunc = s.slw.Alert
	case types.CRIT:
		writeFunc = s.slw.Crit
	case types.ERR:
		writeFunc = s.slw.Err
	case types.WARN:
		writeFunc = s.slw.Warning
	case types.NOTIFY:
		writeFunc = s.slw.Notice
	case types.INFO:
		writeFunc = s.slw.Info
	case types.DEBUG:
		writeFunc = s.slw.Debug
	default:
		return 0, nil
	}
	err := writeFunc(string(ba))
	if err != nil {
		return 0, err
	}
	return len(ba), nil
}

// Close closes connection to syslog server
func (s *sysLogWriter) Close() error {
	return s.Close()
}
