package log

import (
	"github.com/rs/zerolog"
)

const (
	LevelDebug uint8 = iota
	LevelInfo
	LevelWarn
	LevelError
)

const (
	LevelDebugString = "DEBUG"
	LevelInfoString  = "INFO"
	LevelWarnString  = "WARN"
	LevelErrorString = "ERROR"
)

const (
	EncodingJSON  = "JSON"
	EncodingPlain = "PLAIN"
)

const PrefixFieldName = "prefix"

//nolint:gochecknoglobals // used as constant
var levelsMapping = map[uint8]zerolog.Level{
	LevelDebug: zerolog.DebugLevel,
	LevelInfo:  zerolog.InfoLevel,
	LevelWarn:  zerolog.WarnLevel,
	LevelError: zerolog.ErrorLevel,
}

//nolint:gochecknoglobals // used as constant
var envLevelsMapping = map[string]uint8{
	LevelDebugString: LevelDebug,
	LevelInfoString:  LevelInfo,
	LevelWarnString:  LevelWarn,
	LevelErrorString: LevelError,
}

func getLevel(level string) zerolog.Level {
	return levelsMapping[envLevelsMapping[level]]
}
