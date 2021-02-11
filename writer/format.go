package writer

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/primalcs/logger/types"
)

// Format creates a message with giver parameters
func Format(level types.LogLevel, delimiter, tag, prefix, msg string, kvs ...string) string {
	out := types.LogLevels[level]
	if tag != "" {
		out += delimiter + tag
	}
	if prefix != "" {
		out += delimiter + prefix
	}
	out += delimiter + msg
	if len(kvs) > 0 {
		out += delimiter
		key := kvs[0]
		mp := map[string]string{key: ""}

		for i := 1; i < len(kvs); i++ {
			if i%2 != 0 {
				mp[key] = kvs[i]
			} else {
				key = kvs[i]
			}
		}
		ba, err := json.Marshal(mp)
		if err != nil {
			// TODO process err
		} else {
			out += string(ba)
		}
	}
	return out
}

// LogTime adds timestamp to the beginning of the message
func LogTime(loc *time.Location, format, delimiter, msg string) string {
	t := time.Now().In(loc).Format(format)
	return t + delimiter + msg
}

// LogCaller adds function name and line number before message
func LogCaller(delimiter, msg string) string {
	_, fn, ln, ok := runtime.Caller(types.LogCallerSkipLevels)
	if !ok {
		return msg
	}
	call := fmt.Sprintf("file: %s; line: %d", fn, ln)
	return call + delimiter + msg
}
