package grpc

import (
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
}

func Server(custom grpc.UnaryServerInterceptor, p recovery.RecoveryHandlerFunc) *grpc.Server {
	recoveryOptions := []recovery.Option{
		recovery.WithRecoveryHandler(p),
	}

	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.KeepaliveEnforcementPolicy(kaep))
	grpcOptions = append(grpcOptions, grpc.KeepaliveParams(kasp))
	grpcOptions = append(grpcOptions, grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOptions...),
		custom,
	))

	return grpc.NewServer(grpcOptions...)
}
