package suite

import (
	"context"
	"log"
	"net"
	"sso/internal/config"
	"strconv"
	"testing"

	ssov1 "github.com/kirill-dolgii/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T
	Cfg       config.Config
	ApiClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPCConfig.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("grpc server connection failed %v", err)
	}

	return ctx, &Suite{
		T:         t,
		Cfg:       *cfg,
		ApiClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPCConfig.Port))
}
