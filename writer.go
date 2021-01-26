package logger

import (
	"log/syslog"
)

type writer struct {
	logWriter  LogWriter
	connection string
	priority   syslog.Priority
	addr       string
	prefixTag  string
}

func NewWriter(connection, addr, prefix string, priority syslog.Priority) (*writer, error) {
	var d LogWriter
	var err error
	switch connection {
	case TCP, UDP:
		if d, err = syslog.Dial(connection, addr, priority, prefix); err != nil {
			return nil, err
		}
	case LOCAL:
		if d, err = syslog.New(priority, prefix); err != nil {
			return nil, err
		}
	case FILE:
		if d, err = NewFileWriter(addr); err != nil {
			return nil, err
		}
	}

	w := &writer{
		logWriter:  d,
		connection: connection,
		priority:   priority,
		addr:       addr,
		prefixTag:  prefix,
	}
	return w, nil
}

func (w *writer) write(lp LogParams, tp TimeParams, mp MsgParams, kvs ...string) error {
	m := Format(lp.Level, mp.Delimiter, mp.Tag, w.prefixTag, mp.Msg, kvs...)
	if tp.Location != nil {
		m = LogTime(tp.Location, tp.Format, mp.Delimiter, m)
	}
	writeFunc := func(m string) error {
		_, err := w.logWriter.Write([]byte(m))
		return err
	}
	if lp.IsForced {
		switch lp.Level {
		case EMERG:
			writeFunc = w.logWriter.Emerg
		case ALERT:
			writeFunc = w.logWriter.Alert
		case CRIT:
			writeFunc = w.logWriter.Crit
		case ERR:
			writeFunc = w.logWriter.Err
		case WARN:
			writeFunc = w.logWriter.Warning
		case NOTIFY:
			writeFunc = w.logWriter.Notice
		case INFO:
			writeFunc = w.logWriter.Info
		case DEBUG:
			writeFunc = w.logWriter.Debug
		default:
			return nil
		}
	}
	if err := writeFunc(m); err != nil {
		return err
	}
	return nil
}
