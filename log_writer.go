package logger

type LogWriter interface {
	Emerg(m string) error
	Alert(m string) error
	Crit(m string) error
	Err(m string) error
	Warning(m string) error
	Notice(m string) error
	Info(m string) error
	Debug(m string) error
	Write([]byte) (int, error)
	Close() error
}
