package grpc

import (
	"context"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
)

func newSwitchService(ds *datastore.RethinkStore) *switchService {
	return &switchService{
		ds: ds,
	}
}

type switchService struct {
	ds *datastore.RethinkStore
}

func (s *switchService) Update(ctx context.Context, req *v1.SwitchRegisterRequest) (*v1.SwitchResponse, error) {
	return nil, nil
}

func (s *switchService) Get(ctx context.Context, req *v1.SwitchGetRequest) (*v1.SwitchResponse, error) {
	return nil, nil
}

func (s *switchService) Find(ctx context.Context, req *v1.SwitchFindRequest) (*v1.SwitchListResponse, error) {
	return nil, nil
}

func (s *switchService) List(ctx context.Context, req *v1.SwitchListRequest) (*v1.SwitchListResponse, error) {
	return nil, nil
}
