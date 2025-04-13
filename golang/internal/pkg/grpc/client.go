package grpc

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

var kacp = keepalive.ClientParameters{
	Time:                time.Second * 5,
	Timeout:             time.Second,
	PermitWithoutStream: true,
}

func Connection(ctx context.Context, dns string, inter grpc.UnaryClientInterceptor) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	ch := byte(0)

	if len(dns) > 0 {
		ch = dns[0]
	}

	if ch < '0' || ch > '9' {
		resolver.SetDefaultScheme("dns")
		opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	}

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithKeepaliveParams(kacp))
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			inter,
		),
	))
	hostName, _ := os.Hostname()
	opts = append(opts, grpc.WithUserAgent(hostName))

	gc, err := grpc.DialContext(ctx, dns, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial grpc connection")
	}

	return gc, nil
}
