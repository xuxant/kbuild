package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xuxant/kbuild/pkg/constants"
)

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type Level uint32

var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

const DefaultLogLevel = WarnLevel

type contextKey struct{}

var ContextKey = contextKey{}

type EventContext struct {
	Task    constants.Phase
	Subtask string
}

var logger = New()

func (level Level) String() string {
	switch level {
	case TraceLevel:
		return "trace"

	case DebugLevel:
		return "debug"

	case InfoLevel:
		return "info"

	case WarnLevel:
		return "warn"

	case ErrorLevel:
		return "error"

	case FatalLevel:
		return "fatal"

	case PanicLevel:
		return "panic"
	}
	return "unknown"
}

func Entry(ctx context.Context) *logrus.Entry {
	val := ctx.Value(ContextKey)
	if eventContext, ok := val.(EventContext); ok {
		return logger.WithFields(logrus.Fields{
			"task":    eventContext.Task,
			"subtask": eventContext.Subtask,
		})
	}
	return logger.WithField("task", constants.Build)

}

func New() *logrus.Logger {
	return logrus.New()
}
