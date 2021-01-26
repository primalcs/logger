package logger

import (
	"sync"
	"time"
)

type config struct {
	mu         sync.RWMutex
	logLevel   LogLevel
	delimiter  string
	location   *time.Location
	timeFormat string
}

func NewConfig() *config {
	cfg := &config{
		delimiter: defaultDelimiter,
	}
	return cfg
}

func (c *config) SetLogLevel(level LogLevel) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.logLevel != level {
		c.logLevel = level
	}
}

func (c *config) GetLogLevel() LogLevel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logLevel
}
