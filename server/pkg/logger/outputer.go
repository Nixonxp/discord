package log

import (
	"fmt"
	"io"

	"github.com/rs/zerolog"
)

type outputer interface {
	output() io.Writer
	setOutput(*Logger, io.Writer)

	prefix() string
	setPrefix(*Logger, string)
}

func newOutputer(conf *Config, w io.Writer) outputer {
	if conf.EncodingType == EncodingPlain {
		return newPlainOutputer(conf, w)
	}
	return newJSONOutputer(conf, w)
}

type plainOutputer struct {
	out           zerolog.ConsoleWriter
	messagePrefix string
}

func newPlainOutputer(conf *Config, w io.Writer) outputer {
	return &plainOutputer{
		out: zerolog.ConsoleWriter{
			TimeFormat: conf.TimestampLayout,
			Out:        w,
		},
	}
}

func (cw *plainOutputer) output() io.Writer {
	return cw.out
}

func (cw *plainOutputer) setOutput(l *Logger, w io.Writer) {
	cw.out.Out = w
	l.logger = l.logger.Output(&cw.out)
}

func (cw *plainOutputer) prefix() string {
	return cw.messagePrefix
}

func (cw *plainOutputer) setPrefix(l *Logger, prefix string) {
	cw.messagePrefix = prefix
	cw.out.FormatMessage = func(i any) string {
		return fmt.Sprintf("%s%s", prefix, i)
	}
	l.logger = l.logger.Output(&cw.out)
}

type jsonOutputer struct {
	out           io.Writer
	messagePrefix string
}

func newJSONOutputer(conf *Config, w io.Writer) outputer {
	return &jsonOutputer{
		out:           w,
		messagePrefix: conf.Prefix,
	}
}

func (j *jsonOutputer) output() io.Writer {
	return j.out
}

func (j *jsonOutputer) setOutput(l *Logger, w io.Writer) {
	j.out = w
	l.logger = l.logger.Output(j.out)
}

func (j *jsonOutputer) prefix() string {
	return j.messagePrefix
}

func (j *jsonOutputer) setPrefix(l *Logger, prefix string) {
	l.logger = l.logger.With().Str(PrefixFieldName, prefix).Logger()
}
