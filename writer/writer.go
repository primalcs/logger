package writer

import (
	"log/syslog"
	"sync"

	"github.com/primalcs/logger/file_writer"
	"github.com/primalcs/logger/nsq_writer"
	"github.com/primalcs/logger/ring_buffer"
	"github.com/primalcs/logger/syslog_writer"
	"github.com/primalcs/logger/types"
)

// Writer is a logging unit with constant connection type, address, prefixTag and priority
// contains ring buffer for unsent messages
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

func createConnection(connection, addr, prefix string, priority syslog.Priority) (types.LogWriter, error) {
	var d types.LogWriter
	var err error
	switch connection {
	case types.ConnectionTCP, types.ConnectionUDP, types.ConnectionLOCAL:
		if d, err = syslog_writer.NewSysLogWriter(connection, addr, prefix, priority); err != nil {
			return nil, err
		}
	case types.ConnectionFILE:
		if d, err = file_writer.NewFileWriter(addr); err != nil {
			return nil, err
		}
	case types.ConnectionNSQ:
		if d, err = nsq_writer.NewNSQWriter(addr, prefix); err != nil {
			return nil, err
		}
	}
	return d, nil
}

// NewWriter creates an logging instance with given parameters
func NewWriter(connection, addr, prefix string, priority syslog.Priority, bufferLen int) (*Writer, error) {
	conn, err := createConnection(connection, addr, prefix, priority)
	if err != nil {
		return nil, err
	}

	w := &Writer{
		logWriter:     conn,
		connection:    connection,
		priority:      priority,
		addr:          addr,
		prefixTag:     prefix,
		status:        types.WriterStatusOk,
		messageBuffer: ring_buffer.NewMessageBuffer(bufferLen),
	}
	return w, nil
}

// Reconnect tries to create new logger connection with old parameters;
// if succeed sends messages from ring buffer and erases them
func (w *Writer) Reconnect() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	conn, err := createConnection(w.connection, w.addr, w.prefixTag, w.priority)
	if err != nil {
		return err
	}

	w.logWriter = conn
	w.status = types.WriterStatusOk

	if err := w.processMessageBuffer(); err != nil {
		w.stop()
		return err
	}
	return nil
}

// Write writes message with parameters if writer status is OK
// otherwise adds a message to ring buffer
func (w *Writer) Write(lp types.LogParams, tp types.TimeParams, mp types.MsgParams, kvs ...string) error {
	m := Format(lp.Level, mp.Delimiter, mp.Tag, w.prefixTag, mp.Msg, kvs...)
	if lp.IsWithCaller {
		m = LogCaller(mp.Delimiter, m)
	}
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
		w.stop()
		writeFunc = w.writeAtStatusStopped
	}
	if err := writeFunc(lp, m); err != nil {
		return err
	}

	return nil
}

func (w *Writer) writeAtStatusOk(lp types.LogParams, m string) error {
	if lp.IsForced {
		if _, err := w.logWriter.WriteForced(lp.Level, []byte(m)); err != nil {
			return err
		}
	} else {
		if _, err := w.logWriter.Write([]byte(m)); err != nil {
			return err
		}
	}

	return nil
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

// Stop closes Writer connection and sets it's status to Stopped; uses mutex lock
func (w *Writer) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.stop()
}

func (w *Writer) stop() {
	w.logWriter.Close()
	w.status = types.WriterStatusStopped
}
