package config

import (
	"sync"
	"time"

	"github.com/primalcs/logger/types"
)

// Config is a struct that encapsulates logger configs for thread-safe usage
type Config struct {
	mu           sync.RWMutex
	logLevel     types.LogLevel
	delimiter    string
	location     *time.Location
	timeFormat   string
	isWithCaller bool
}

// NewConfig creates new Config instance
func NewConfig() *Config {
	cfg := &Config{
		delimiter: types.DefaultDelimiter,
	}
	return cfg
}

// SetLogLevel sets logging level; thread-safe
func (c *Config) SetLogLevel(level types.LogLevel) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.logLevel != level {
		c.logLevel = level
	}
}

// GetLogLevel gets logging level; thread-safe
func (c *Config) GetLogLevel() types.LogLevel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logLevel
}

// SetDelimiter sets delimiter string; thread-safe
func (c *Config) SetDelimiter(d string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.delimiter != d {
		c.delimiter = d
	}
}

// GetDelimiter gets delimiter string; thread-safe
func (c *Config) GetDelimiter() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.delimiter
}

// SetTimeOptions sets time format and location config; thread-safe
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

// GetTimeOptions gets time format and location config; thread-safe
func (c *Config) GetTimeOptions() (string, *time.Location) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.timeFormat, c.location
}

// SetWithCaller sets config for calling runtime.Caller() for logging; thread-safe
func (c *Config) SetWithCaller(val bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isWithCaller != val {
		c.isWithCaller = val
	}
}

// GetWithCaller shows if the runtime.Caller() in needed; thread-safe
func (c *Config) GetWithCaller() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isWithCaller
}
