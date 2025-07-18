package grpcx

import (
	"fmt"
	"net"

	"github.com/boostgo/core/appx"
	"github.com/boostgo/core/grpcx/intercept"
	"github.com/boostgo/core/log"

	"google.golang.org/grpc"
)

type Registry func(server *grpc.Server)

func Run(port int, registries ...Registry) {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			intercept.Recovery,
			intercept.Tracer,
			intercept.Logging,
			intercept.Validate,
			intercept.ErrorHandling,
		),
	)

	if len(registries) > 0 {
		for _, registry := range registries {
			registry(server)
		}
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	appx.Tear(func() error {
		server.GracefulStop()
		return nil
	})

	log.
		Info().
		Int("port", port).
		Msg("gRPC server start")

	if err = server.Serve(l); err != nil {
		log.
			Error().
			Err(err).
			Msg("gRPC server serve")

		appx.Cancel()
	}
}
