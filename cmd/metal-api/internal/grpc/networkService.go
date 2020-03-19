package grpc

import (
	"context"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
)

func newNetworkService(ds *datastore.RethinkStore) *networkService {
	return &networkService{
		ds: ds,
	}
}

type networkService struct {
	ds *datastore.RethinkStore
}

func (s *networkService) Create(ctx context.Context, req *v1.NetworkCreateRequest) (*v1.NetworkResponse, error) {
	return nil, nil
}

func (s *networkService) Update(ctx context.Context, req *v1.NetworkUpdateRequest) (*v1.NetworkResponse, error) {
	return nil, nil
}

func (s *networkService) Allocate(ctx context.Context, req *v1.NetworkAllocateRequest) (*v1.NetworkResponse, error) {
	return nil, nil
}

func (s *networkService) Find(ctx context.Context, req *v1.NetworkFindRequest) (*v1.NetworkListResponse, error) {
	return nil, nil
}

func (s *networkService) List(ctx context.Context, req *v1.NetworkListRequest) (*v1.NetworkListResponse, error) {
	return nil, nil
}
