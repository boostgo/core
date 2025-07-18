package intercept

import (
	"context"

	"github.com/boostgo/core/errorx"
	"github.com/boostgo/core/grpcx/errs"
	"github.com/boostgo/core/log"
	"github.com/boostgo/core/trace"
	"github.com/boostgo/core/validator"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Interceptor func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)

func Recovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	var response any
	var err error

	if err = errorx.Try(func() error {
		response, err = handler(ctx, req)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return response, nil
}

func Tracer(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	traceID := trace.Get(ctx)
	if traceID == "" {
		ctx = trace.Set(ctx)
	}

	return handler(ctx, req)
}

func Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	const (
		start      = "start"
		stageError = "error"
		end        = "end"
	)

	log.
		Info().
		Str("stage", start).
		Obj("request", req).
		Msg(info.FullMethod)

	response, err := handler(ctx, req)
	if err != nil {
		log.
			Error().
			Str("stage", stageError).
			Err(err).
			Msg(info.FullMethod)

		return nil, err
	}

	log.
		Info().
		Str("stage", end).
		Obj("response", response).
		Msg(info.FullMethod)

	return response, nil
}

func ErrorHandling(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	response, err := handler(ctx, req)
	if err != nil {
		return nil, status.Error(errs.Code(err), err.Error())
	}

	return response, nil
}

func Validate(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if err := validator.Get().Struct(req); err != nil {
		return nil, status.Error(errs.Code(err), err.Error())
	}

	return handler(ctx, req)
}
