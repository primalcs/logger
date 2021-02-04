package logger

import "github.com/primalcs/logger/types"

// Contains aliases for connection constants for easier access and usage
const (
	ConnectionTCP   = types.ConnectionTCP
	ConnectionUDP   = types.ConnectionUDP
	ConnectionLOCAL = types.ConnectionLOCAL
	ConnectionFILE  = types.ConnectionFILE
	ConnectionNSQ   = types.ConnectionNSQ
)

// Contains aliases for priority constants for easier access and usage
const (
	EMERG  = types.EMERG
	ALERT  = types.ALERT
	CRIT   = types.CRIT
	ERR    = types.ERR
	WARN   = types.WARN
	NOTIFY = types.NOTIFY
	INFO   = types.INFO
	DEBUG  = types.DEBUG
	FORCE  = types.FORCE // force one of levels above =  types.FORCE // force one of levels above
)

// Contains aliases for delimiter constants for easier access and usage
const (
	DelimiterV       = types.DelimiterV
	DelimiterH       = types.DelimiterH
	DelimiterA       = types.DelimiterA
	DefaultDelimiter = DelimiterV
)
