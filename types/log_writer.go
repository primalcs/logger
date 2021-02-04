package types

// LogWriter is an interface that must be implemented for using with logger
type LogWriter interface {
	Write([]byte) (int, error)
	WriteForced(LogLevel, []byte) (int, error)
	Close() error
}
