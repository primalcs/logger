package logger

import "time"

type LogParams struct {
	IsForced bool
	Level    LogLevel
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
