package config

import (
	"sync"
	"time"

	"github.com/rybnov/logger/types"
)

type Config struct {
	mu         sync.RWMutex
	logLevel   types.LogLevel
	delimiter  string
	location   *time.Location
	timeFormat string
}

func NewConfig() *Config {
	cfg := &Config{
		delimiter: types.DefaultDelimiter,
	}
	return cfg
}

func (c *Config) SetLogLevel(level types.LogLevel) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.logLevel != level {
		c.logLevel = level
	}
}

func (c *Config) GetLogLevel() types.LogLevel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logLevel
}

func (c *Config) SetDelimiter(d string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.delimiter != d {
		c.delimiter = d
	}
}

func (c *Config) GetDelimiter() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.delimiter
}

func (c *Config) SetTimeOptions(format string, location *time.Location) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.timeFormat != format && format != "" {
		c.timeFormat = format
	}
	if location != nil && c.location != location {
		c.location = location
	}
}

func (c *Config) GetTimeOptions() (string, *time.Location) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.timeFormat, c.location
}
