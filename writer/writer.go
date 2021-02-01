package writer

import (
	"log/syslog"
	"sync"

	"github.com/rybnov/logger/file_writer"

	"github.com/rybnov/logger/ring_buffer"

	"github.com/rybnov/logger/types"
)

type Writer struct {
	logWriter     types.LogWriter
	connection    string
	priority      syslog.Priority
	addr          string
	prefixTag     string
	status        types.WriterStatus
	mu            sync.RWMutex
	messageBuffer *ring_buffer.MessageBuffer
}

func NewWriter(connection, addr, prefix string, priority syslog.Priority, bufferLen int) (*Writer, error) {
	var d types.LogWriter
	var err error
	switch connection {
	case types.TCP, types.UDP:
		if d, err = syslog.Dial(connection, addr, priority, prefix); err != nil {
			return nil, err
		}
	case types.LOCAL:
		if d, err = syslog.New(priority, prefix); err != nil {
			return nil, err
		}
	case types.FILE:
		if d, err = file_writer.NewFileWriter(addr); err != nil {
			return nil, err
		}
	}

	w := &Writer{
		logWriter:     d,
		connection:    connection,
		priority:      priority,
		addr:          addr,
		prefixTag:     prefix,
		status:        types.WriterStatusOk,
		messageBuffer: ring_buffer.NewMessageBuffer(bufferLen),
	}
	return w, nil
}

func (w *Writer) Reconnect() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	var d types.LogWriter
	var err error
	switch w.connection {
	case types.TCP, types.UDP:
		if d, err = syslog.Dial(w.connection, w.addr, w.priority, w.prefixTag); err != nil {
			return err
		}
	case types.LOCAL:
		if d, err = syslog.New(w.priority, w.prefixTag); err != nil {
			return err
		}
	case types.FILE:
		if d, err = file_writer.NewFileWriter(w.addr); err != nil {
			return err
		}
	}

	w.logWriter = d
	w.status = types.WriterStatusOk

	if err := w.processMessageBuffer(); err != nil {
		w.Stop(false)
		return err
	}
	return nil
}

func (w *Writer) Write(lp types.LogParams, tp types.TimeParams, mp types.MsgParams, kvs ...string) error {
	m := Format(lp.Level, mp.Delimiter, mp.Tag, w.prefixTag, mp.Msg, kvs...)
	if tp.Location != nil {
		m = LogTime(tp.Location, tp.Format, mp.Delimiter, m)
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	status := w.status

	var writeFunc func(lp types.LogParams, m string) error
	switch status {
	case types.WriterStatusOk:
		writeFunc = w.writeAtStatusOk
	case types.WriterStatusStopped:
		writeFunc = w.writeAtStatusStopped
	default:
		w.Stop(false)
		writeFunc = w.writeAtStatusStopped
	}
	if err := writeFunc(lp, m); err != nil {
		return err
	}

	return nil
}

func (w *Writer) writeAtStatusOk(lp types.LogParams, m string) error {
	writeFunc := func(m string) error {
		_, err := w.logWriter.Write([]byte(m))
		return err
	}
	if lp.IsForced {
		switch lp.Level {
		case types.EMERG:
			writeFunc = w.logWriter.Emerg
		case types.ALERT:
			writeFunc = w.logWriter.Alert
		case types.CRIT:
			writeFunc = w.logWriter.Crit
		case types.ERR:
			writeFunc = w.logWriter.Err
		case types.WARN:
			writeFunc = w.logWriter.Warning
		case types.NOTIFY:
			writeFunc = w.logWriter.Notice
		case types.INFO:
			writeFunc = w.logWriter.Info
		case types.DEBUG:
			writeFunc = w.logWriter.Debug
		default:
			return nil
		}
	}

	return writeFunc(m)
}

func (w *Writer) writeAtStatusStopped(lp types.LogParams, m string) error {
	w.messageBuffer.AddCell(ring_buffer.NewCell(lp, m))
	return nil
}

func (w *Writer) processMessageBuffer() error {
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

func (w *Writer) Stop(lock bool) {
	if lock {
		w.mu.Lock()
		defer w.mu.Unlock()
	}

	w.logWriter.Close()
	w.status = types.WriterStatusStopped
}
