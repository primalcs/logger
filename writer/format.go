package writer

import (
	"time"

	"github.com/rybnov/logger/types"
)

func Format(level types.LogLevel, delimiter, tag, prefix, msg string, kvs ...string) string {
	out := types.LogLevels[level]
	if tag != "" {
		out += delimiter + tag
	}
	out += delimiter + prefix + msg
	if len(kvs) > 0 {
		out += delimiter
		for k, v := range kvs {
			out += v
			if k%2 == 0 {
				out += ":"
			} else {
				out += " "
			}
		}
	}
	return out
}

func LogTime(loc *time.Location, format, delimiter, msg string) string {
	t := time.Now().In(loc).Format(format)
	return t + delimiter + msg
}
