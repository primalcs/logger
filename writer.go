package logger

import (
	"log/syslog"
	"sync"
)

type writer struct {
	logWriter     LogWriter
	connection    string
	priority      syslog.Priority
	addr          string
	prefixTag     string
	status        WriterStatus
	mu            sync.RWMutex
	messageBuffer *MessageBuffer
}

func NewWriter(connection, addr, prefix string, priority syslog.Priority, bufferLen int) (*writer, error) {
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
		logWriter:     d,
		connection:    connection,
		priority:      priority,
		addr:          addr,
		prefixTag:     prefix,
		status:        WriterStatusOk,
		messageBuffer: NewMessageBuffer(bufferLen),
	}
	return w, nil
}

func (w *writer) reconnect(connection, addr, prefix string, priority syslog.Priority) error {
	var d LogWriter
	var err error
	switch connection {
	case TCP, UDP:
		if d, err = syslog.Dial(connection, addr, priority, prefix); err != nil {
			return err
		}
	case LOCAL:
		if d, err = syslog.New(priority, prefix); err != nil {
			return err
		}
	case FILE:
		if d, err = NewFileWriter(addr); err != nil {
			return err
		}
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logWriter = d
	w.connection = connection
	w.priority = priority
	w.addr = addr
	w.prefixTag = prefix
	w.status = WriterStatusOk

	if err := w.processMessageBuffer(); err != nil {
		w.stop(false)
		return err
	}
	return nil
}

func (w *writer) write(lp LogParams, tp TimeParams, mp MsgParams, kvs ...string) error {
	m := Format(lp.Level, mp.Delimiter, mp.Tag, w.prefixTag, mp.Msg, kvs...)
	if tp.Location != nil {
		m = LogTime(tp.Location, tp.Format, mp.Delimiter, m)
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	status := w.status

	var writeFunc func(lp LogParams, m string) error
	switch status {
	case WriterStatusOk:
		writeFunc = w.writeAtStatusOk
	case WriterStatusStopped:
		writeFunc = w.writeAtStatusStopped
	default:
		w.stop(false)
		writeFunc = w.writeAtStatusStopped
	}
	if err := writeFunc(lp, m); err != nil {
		return err
	}

	return nil
}

func (w *writer) writeAtStatusOk(lp LogParams, m string) error {
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

	return writeFunc(m)
}

func (w *writer) writeAtStatusStopped(lp LogParams, m string) error {
	w.messageBuffer.AddCell(NewCell(lp, m))
	return nil
}

func (w *writer) processMessageBuffer() error {
	for {
		cell, ptr, ok := w.messageBuffer.GetOldestCell()
		if !ok {
			break
		}
		if err := w.writeAtStatusOk(cell.LogParams, cell.Message); err != nil {
			return err
		}
		w.messageBuffer.EraseCell(ptr)
	}

	return nil
}

func (w *writer) stop(lock bool) {
	if lock {
		w.mu.Lock()
		defer w.mu.Unlock()
	}

	w.logWriter.Close()
	w.status = WriterStatusStopped
}
