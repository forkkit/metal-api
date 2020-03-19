package grpc

import (
	"context"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
)

func newFirewallService(ds *datastore.RethinkStore) *firewallService {
	return &firewallService{
		ds: ds,
	}
}

type firewallService struct {
	ds *datastore.RethinkStore
}

func (s *firewallService) Create(ctx context.Context, req *v1.FirewallCreateRequest) (*v1.FirewallResponse, error) {
	return nil, nil
}

func (s *firewallService) Find(ctx context.Context, req *v1.FirewallFindRequest) (*v1.FirewallResponse, error) {
	return nil, nil
}
