package grpc

import (
	"context"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
)

func newIPService(ds *datastore.RethinkStore) *ipService {
	return &ipService{
		ds: ds,
	}
}

type ipService struct {
	ds *datastore.RethinkStore
}

func (s *ipService) Allocate(ctx context.Context, req *v1.IPAllocateRequest) (*v1.IPResponse, error) {
	return nil, nil
}

func (s *ipService) Update(ctx context.Context, req *v1.IPUpdateRequest) (*v1.IPResponse, error) {
	return nil, nil
}

func (s *ipService) Delete(ctx context.Context, req *v1.IPDeleteRequest) (*v1.IPResponse, error) {
	return nil, nil
}

func (s *ipService) Get(ctx context.Context, req *v1.IPGetRequest) (*v1.IPResponse, error) {
	return nil, nil
}

func (s *ipService) Find(ctx context.Context, req *v1.IPFindRequest) (*v1.IPListResponse, error) {
	return nil, nil
}

func (s *ipService) List(ctx context.Context, req *v1.IPListRequest) (*v1.IPListResponse, error) {
	return nil, nil
}
