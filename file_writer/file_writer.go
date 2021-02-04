package file_writer

import (
	"os"
	"strings"

	"github.com/primalcs/logger/types"
)

type fileWriter struct {
	file    *os.File
	counter int
}

// NewFileWriter creates new instance for writing log into file
func NewFileWriter(addr string) (*fileWriter, error) {
	p := strings.Split(addr, "/")
	path := strings.Join(p[:len(p)-1], "/")
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	f, err := os.Create(p[len(p)-1])
	if err != nil {
		return nil, err
	}
	return &fileWriter{
		file: f,
	}, nil
}

// WriteForced for fileWriter is the same as Write(). Needed to implement interface
func (f *fileWriter) WriteForced(_ types.LogLevel, ba []byte) (int, error) {
	return f.Write(ba)
}

// Write writes bytes into file and syncs every types.SyncFileAfterMessagesCount records
func (f *fileWriter) Write(ba []byte) (int, error) {
	f.counter++
	if f.counter == types.SyncFileAfterMessagesCount {
		f.counter = 0
		if err := f.file.Sync(); err != nil {
			return 0, err
		}
	}
	return f.file.Write(ba)
}

// Close syncs fileWriter and closes it
func (f *fileWriter) Close() error {
	if err := f.file.Sync(); err != nil {
		return err
	}
	return f.file.Close()
}
