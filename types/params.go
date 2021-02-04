package types

import (
	"time"
)

// LogParams contains log level parameters
type LogParams struct {
	IsForced     bool
	Level        LogLevel
	IsWithCaller bool
}

// TimeParams contains information about location and time format for logging
type TimeParams struct {
	Location *time.Location
	Format   string
}

// MsgParams contains parameters of message to be logged
type MsgParams struct {
	Delimiter string
	Tag       string
	Msg       string
}
