package grpc

import (
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

const DefaultGRPCPort = 50005

func CreateServer(ds *datastore.RethinkStore) (*grpc.Server, error) {
	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			// Problem: https://github.com/grpc/grpc-go/issues/2160
			// Solution: https://stackoverflow.com/questions/52993259/problem-with-grpc-setup-getting-an-intermittent-rpc-unavailable-error/54703234#54703234
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Minute,
			Timeout:           5 * time.Minute,
		}),
		grpc.ReadBufferSize(128),              // defaults to 32 //TODO via config
		grpc.ConnectionTimeout(5*time.Minute), //TODO via config
		grpc.MaxConcurrentStreams(2000),       //TODO via config
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			// https://github.com/grpc/grpc-go/issues/2443
			MinTime:             1 * time.Minute,
			PermitWithoutStream: true,
		}),
	)

	v1.RegisterMachineServiceServer(s, newMachineService(ds))
	v1.RegisterFirewallServiceServer(s, newFirewallService(ds))
	v1.RegisterImageServiceServer(s, newImageService(ds))
	v1.RegisterIPServiceServer(s, newIPService(ds))
	v1.RegisterNetworkServiceServer(s, newNetworkService(ds))
	v1.RegisterPartitionServiceServer(s, newPartitionService(ds))
	v1.RegisterSizeServiceServer(s, newSizeService(ds))
	v1.RegisterSwitchServiceServer(s, newSwitchService(ds))

	return s, nil
}
