package log

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger    zerolog.Logger
	conf      *Config
	outputer  outputer
	stackSkip int
}

func NewLogger(conf *Config) (*Logger, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	setGlobalSettings(conf)

	l := new(Logger)
	l.conf = conf
	l.outputer = newOutputer(l.conf, os.Stderr)
	l.logger = zerolog.New(l.outputer.output())

	if conf.Prefix != "" {
		l.outputer.setPrefix(l, conf.Prefix)
	}

	return l, nil
}

func setGlobalSettings(conf *Config) {
	level := getLevel(conf.LogLevel)

	zerolog.TimeFieldFormat = conf.TimestampLayout
	zerolog.SetGlobalLevel(level)
}

func (l *Logger) Printf(format string, v ...any) {
	l.logger.Printf(format, v...)
}

func (l *Logger) Print(v ...any) {
	l.logger.Print(v...)
}

func (l *Logger) withCaller(e *zerolog.Event) *zerolog.Event {
	if l.conf.EnableCaller {
		return e.Caller(2 + l.stackSkip)
	}

	return e
}

func (l *Logger) Fatal(v ...any) {
	l.withCaller(l.logger.Fatal()).Timestamp().Msg(fmt.Sprint(v...))
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.withCaller(l.logger.Fatal()).Timestamp().Msgf(format, fmt.Sprint(v...))
}

func (l *Logger) Panic(v ...any) {
	l.withCaller(l.logger.Panic()).Timestamp().Msg(fmt.Sprint(v...))
}

func (l *Logger) Panicf(format string, v ...any) {
	l.withCaller(l.logger.Panic()).Timestamp().Msgf(format, v...)
}

func (l *Logger) Debug(v ...any) {
	l.withCaller(l.logger.Debug()).Timestamp().Msg(fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...any) {
	l.withCaller(l.logger.Debug()).Timestamp().Msgf(format, v...)
}

func (l *Logger) Info(v ...any) {
	l.withCaller(l.logger.Info()).Timestamp().Msg(fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...any) {
	l.withCaller(l.logger.Info()).Timestamp().Msgf(format, v...)
}

func (l *Logger) Warn(v ...any) {
	l.withCaller(l.logger.Warn()).Timestamp().Msg(fmt.Sprint(v...))
}

func (l *Logger) Warnf(format string, v ...any) {
	l.withCaller(l.logger.Warn()).Timestamp().Msgf(format, v...)
}

func (l *Logger) Error(v ...any) {
	l.withCaller(l.logger.Error()).Timestamp().Msg(fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...any) {
	l.withCaller(l.logger.Error()).Timestamp().Msgf(format, v...)
}

func (l *Logger) WithField(key string, value any) *Logger {
	return &Logger{
		logger:   l.logger.With().Interface(key, value).Logger(),
		conf:     l.conf,
		outputer: l.outputer,
	}
}

// WithFields returns logger which contains passed
// fileds in context.
// Basic and safe type for fields is map[string]any.
// De-duplication is not supported!
//
// EXPERIMENTAL:
// for more performance and to reduce allocations can be
// used array(slice) of any elements with size [2 * elems count]
// as example: []any{"key","value","field_key","field_value"}.
// May be removed in next versions or with other logger.
func (l *Logger) WithFields(fields any) *Logger {
	return &Logger{
		logger:   l.logger.With().Fields(fields).Logger(),
		conf:     l.conf,
		outputer: l.outputer,
	}
}

func (l *Logger) Prefix() string {
	return l.outputer.prefix()
}

func (l *Logger) SetPrefix(prefix string) {
	l.outputer.setPrefix(l, prefix)
}

func (l *Logger) Output() io.Writer {
	return l.outputer.output()
}

func (l *Logger) SetOutput(w io.Writer) {
	l.outputer.setOutput(l, w)
}

func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}

	return l.WithField(Error, err.Error())
}

// WithContext adds trace_id and (in the future) other fields from context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	fieldsMap := map[string]interface{}{}
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		fieldsMap[TraceID] = spanCtx.TraceID()
	}

	if len(fieldsMap) == 0 {
		return l
	}

	return l.WithFields(fieldsMap)
}

// WithStackSkip The argument skip is the number of stack frames to ascend Skip
func (l *Logger) WithStackSkip(skip int) *Logger {
	return &Logger{
		logger:    l.logger,
		conf:      l.conf,
		outputer:  l.outputer,
		stackSkip: l.stackSkip + skip,
	}
}
