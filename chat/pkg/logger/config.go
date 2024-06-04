package log

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type Config struct {
	TimestampLayout string
	Prefix          string
	EncodingType    string
	LogLevel        string
	EnableCaller    bool
}

func NewDefaultConfig() *Config {
	return &Config{
		TimestampLayout: time.RFC3339,
		Prefix:          "",
		EncodingType:    EncodingJSON,
		LogLevel:        LevelDebugString,
		EnableCaller:    false,
	}
}

func (conf *Config) Validate() error {
	return validator.New().Struct(conf)
}
