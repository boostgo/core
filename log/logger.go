// Package log is wrap over zerolog library.
// Features:
// - Setting trace_id key from provided context.
// - Fatal level log calls global context cancel (not panic or os.Exit()).
// - Printing errorx.Error.
// - More key-value pairs.
package log

import (
	"context"
	"github.com/rs/zerolog"

	"github.com/boostgo/core/appx"
	"github.com/boostgo/core/log/logx"
)

// Debug print log on debug level.
//
// Provided context use trace id
func Debug() Event {
	l := logx.Logger()
	e := logx.Extractor()
	return newEvent(l.Debug(), e)
}

// Info print log on info level.
//
// Provided context use trace id
func Info() Event {
	l := logx.Logger()
	e := logx.Extractor()
	return newEvent(l.Info(), e)
}

// Warn print log on warning level.
//
// Provided context use trace id
func Warn() Event {
	l := logx.Logger()
	e := logx.Extractor()
	return newEvent(l.Warn(), e)
}

// Error print log on error level.
//
// Provided context use trace id
func Error() Event {
	l := logx.Logger()
	e := logx.Extractor()
	return newEvent(l.Error(), e)
}

// Fatal print log on error level but with bool fatal=true.
// Provided context use trace id.
//
// Call AppCancel function
func Fatal() Event {
	defer appx.Cancel()
	l := logx.Logger()
	e := logx.Extractor()
	return newEvent(l.Error().Bool("fatal", true), e)
}

func With(ctx context.Context) zerolog.Context {
	l := logx.Logger()
	return l.With().Ctx(ctx)
}

// Logger is wrap over zerolog logger
type Logger interface {
	// Debug print log on debug level.
	// Provided context use trace id
	Debug() Event
	// Info print log on info level.
	// Provided context use trace id
	Info() Event
	// Warn print log on warning level.
	// Provided context use trace id
	Warn() Event
	// Error print log on error level.
	// Provided context use trace id
	Error() Event
	// Fatal print log on error level but with bool fatal=true.
	// Provided context use trace id.
	//
	// Call life.Cancel() method which call graceful shutdown
	Fatal() Event
}

type wrapper struct {
	namespace string
	ctx       context.Context
}

// Namespace creates Logger implementation with namespace
func Namespace(namespace string) Logger {
	return Context(context.Background(), namespace)
}

// Context creates Logger implementation with context & namespace
func Context(ctx context.Context, namespace string) Logger {
	return &wrapper{
		ctx:       ctx,
		namespace: namespace,
	}
}

func (logger *wrapper) Debug() Event {
	return Debug().
		Ctx(logger.ctx).
		Namespace(logger.namespace)
}

func (logger *wrapper) Info() Event {
	return Info().
		Ctx(logger.ctx).
		Namespace(logger.namespace)
}

func (logger *wrapper) Warn() Event {
	return Warn().
		Ctx(logger.ctx).
		Namespace(logger.namespace)
}

func (logger *wrapper) Error() Event {
	return Error().
		Ctx(logger.ctx).
		Namespace(logger.namespace)
}

func (logger *wrapper) Fatal() Event {
	return Fatal().
		Ctx(logger.ctx).
		Namespace(logger.namespace)
}
