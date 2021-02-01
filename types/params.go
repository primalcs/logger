package types

import (
	"time"
)

type LogParams struct {
	IsForced     bool
	Level        LogLevel
	IsWithCaller bool
}

type TimeParams struct {
	Location *time.Location
	Format   string
}

type MsgParams struct {
	Delimiter string
	Tag       string
	Msg       string
}
