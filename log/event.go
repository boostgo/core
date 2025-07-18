package log

import (
	"context"
	"errors"
	"time"

	"github.com/boostgo/core/convert"
	"github.com/boostgo/core/errorx"
	"github.com/boostgo/core/log/logx"
	"github.com/boostgo/core/trace"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Event represents a log event. It is instanced by one of the level method of
// Logger and finalized by the Msg or Msgf method.
//
// Notice: this is wrap over zerolog.Event
type Event struct {
	inner     *zerolog.Event
	extractor logx.ExtractorFunc
}

func newEvent(inner *zerolog.Event, extractor logx.ExtractorFunc) Event {
	return Event{
		inner:     inner,
		extractor: extractor,
	}
}

func (e Event) Send() {
	e.inner.Send()
}

func (e Event) Ctx(ctx context.Context) Event {
	if ctx == nil {
		return e
	}

	e.inner.Ctx(ctx)

	traceID := trace.Get(ctx)
	if traceID != "" {
		e.Str("trace_id", traceID)
	}

	if e.extractor != nil {
		e.extractor(ctx, e.inner)
	}

	return e
}

func (e Event) Any(key string, object any) Event {
	e.inner.Interface(key, object)
	return e
}

func (e Event) Arr(key string, args ...any) Event {
	if args == nil {
		return e
	}

	stringArgs := make([]string, len(args))
	for i, arg := range args {
		stringArgs[i] = convert.String(arg)
	}
	e.inner.Strs(key, stringArgs)
	return e
}

func (e Event) Err(err error) Event {
	if err == nil {
		return e
	}

	var converted *errorx.Error
	ok := errors.As(err, &converted)

	if !ok {
		e.inner.Err(err)
	} else {
		if converted.Data() != nil {
			e.Obj("context", converted.Data())
		}

		e.inner.Err(errors.New(converted.Error()))
	}

	return e
}

func (e Event) Errs(key string, errors []error) Event {
	e.inner.Errs(key, errors)
	return e
}

func (e Event) Msg(message string) Event {
	e.inner.Msg(message)
	return e
}

func (e Event) Msgf(format string, args ...any) Event {
	e.inner.Msgf(format, args...)
	return e
}

func (e Event) Str(key, value string) Event {
	e.inner.Str(key, value)
	return e
}

func (e Event) UUID(key string, id uuid.UUID) Event {
	e.inner.Str(key, id.String())
	return e
}

func (e Event) Strs(key string, values []string) Event {
	e.inner.Strs(key, values)
	return e
}

func (e Event) Int(key string, value int) Event {
	e.inner.Int(key, value)
	return e
}

func (e Event) Int8(key string, value int8) Event {
	e.inner.Int8(key, value)
	return e
}

func (e Event) Int16(key string, value int16) Event {
	e.inner.Int16(key, value)
	return e
}

func (e Event) Int32(key string, value int32) Event {
	e.inner.Int32(key, value)
	return e
}

func (e Event) Int64(key string, value int64) Event {
	e.inner.Int64(key, value)
	return e
}

func (e Event) Ints(key string, values []int) Event {
	e.inner.Ints(key, values)
	return e
}

func (e Event) Ints8(key string, values []int8) Event {
	e.inner.Ints8(key, values)
	return e
}

func (e Event) Ints16(key string, values []int16) Event {
	e.inner.Ints16(key, values)
	return e
}

func (e Event) Ints32(key string, values []int32) Event {
	e.inner.Ints32(key, values)
	return e
}

func (e Event) Ints64(key string, values []int64) Event {
	e.inner.Ints64(key, values)
	return e
}

func (e Event) Uint(key string, value uint) Event {
	e.inner.Uint(key, value)
	return e
}

func (e Event) Uint8(key string, value uint8) Event {
	e.inner.Uint8(key, value)
	return e
}

func (e Event) Uint16(key string, value uint16) Event {
	e.inner.Uint16(key, value)
	return e
}

func (e Event) Uint32(key string, value uint32) Event {
	e.inner.Uint32(key, value)
	return e
}

func (e Event) Uint64(key string, value uint64) Event {
	e.inner.Uint64(key, value)
	return e
}

func (e Event) Uints(key string, value []uint) Event {
	e.inner.Uints(key, value)
	return e
}

func (e Event) Uints8(key string, value []uint8) Event {
	e.inner.Uints8(key, value)
	return e
}

func (e Event) Uints16(key string, value []uint16) Event {
	e.inner.Uints16(key, value)
	return e
}

func (e Event) Uints32(key string, value []uint32) Event {
	e.inner.Uints32(key, value)
	return e
}

func (e Event) Uints64(key string, value []uint64) Event {
	e.inner.Uints64(key, value)
	return e
}

func (e Event) Float32(key string, value float32) Event {
	e.inner.Float32(key, value)
	return e
}

func (e Event) Floats32(key string, values []float32) Event {
	e.inner.Floats32(key, values)
	return e
}

func (e Event) Float64(key string, value float64) Event {
	e.inner.Float64(key, value)
	return e
}

func (e Event) Floats64(key string, values []float64) Event {
	e.inner.Floats64(key, values)
	return e
}

func (e Event) Bool(key string, value bool) Event {
	e.inner.Bool(key, value)
	return e
}

func (e Event) Time(key string, value time.Time) Event {
	e.inner.Time(key, value)
	return e
}

func (e Event) Times(key string, value []time.Time) Event {
	e.inner.Times(key, value)
	return e
}

func (e Event) Duration(key string, value time.Duration) Event {
	e.inner.Dur(key, value)
	return e
}

func (e Event) Durations(key string, value []time.Duration) Event {
	e.inner.Durs(key, value)
	return e
}

func (e Event) Obj(key string, obj any) Event {
	e.Any(key, convert.String(obj))
	return e
}

func (e Event) Object(key string, obj zerolog.LogObjectMarshaler) Event {
	e.inner.Object(key, obj)
	return e
}

func (e Event) Bytes(key string, bytes []byte) Event {
	e.inner.Bytes(key, bytes)
	return e
}

func (e Event) Type(key string, obj any) Event {
	e.inner.Type(key, obj)
	return e
}

func (e Event) Namespace(namespace string) Event {
	if namespace == "" {
		return e
	}

	e.Str("namespace", namespace)
	return e
}
