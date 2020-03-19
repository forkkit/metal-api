package grpc

import (
	"context"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
)

func newSizeService(ds *datastore.RethinkStore) *sizeService {
	return &sizeService{
		ds: ds,
	}
}

type sizeService struct {
	ds *datastore.RethinkStore
}

func (s *sizeService) Create(ctx context.Context, req *v1.SizeCreateRequest) (*v1.SizeResponse, error) {
	return nil, nil
}

func (s *sizeService) Update(ctx context.Context, req *v1.SizeUpdateRequest) (*v1.SizeResponse, error) {
	return nil, nil
}

func (s *sizeService) Delete(ctx context.Context, req *v1.SizeDeleteRequest) (*v1.SizeResponse, error) {
	return nil, nil
}

func (s *sizeService) Get(ctx context.Context, req *v1.SizeGetRequest) (*v1.SizeResponse, error) {
	return nil, nil
}

func (s *sizeService) Find(ctx context.Context, req *v1.SizeFindRequest) (*v1.SizeListResponse, error) {
	return nil, nil
}

func (s *sizeService) List(ctx context.Context, req *v1.SizeListRequest) (*v1.SizeListResponse, error) {
	return nil, nil
}
