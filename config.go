package logger

import (
	"sync"
)

const (
	TCP   = "tcp"
	UDP   = "udp"
	Local = "local"
)

type config struct {
	mu       sync.RWMutex
	logLevel uint8
}

func NewConfig() *config {
	cfg := &config{}
	return cfg
}
