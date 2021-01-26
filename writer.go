package logger

import "log/syslog"

type writer struct {
	w          *syslog.Writer
	connection string
	priority   syslog.Priority
	addr       string
	tag        string
}

func NewWriter(connection, addr, tag string, priority syslog.Priority) (*writer, error) {
	var d *syslog.Writer
	var err error
	switch connection {
	case TCP, UDP:
		if d, err = syslog.Dial(connection, addr, priority, tag); err != nil {
			return nil, err
		}
	case Local:
		if d, err = syslog.New(priority, tag); err != nil {
			return nil, err
		}
	}

	w := &writer{
		w:          d,
		connection: connection,
		priority:   priority,
		addr:       addr,
		tag:        tag,
	}
	return w, nil
}
