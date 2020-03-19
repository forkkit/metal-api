package grpc

import (
	"context"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
)

func newPartitionService(ds *datastore.RethinkStore) *partitionService {
	return &partitionService{
		ds: ds,
	}
}

type partitionService struct {
	ds *datastore.RethinkStore
}

func (s *partitionService) Create(ctx context.Context, req *v1.PartitionCreateRequest) (*v1.PartitionResponse, error) {
	return nil, nil
}

func (s *partitionService) Update(ctx context.Context, req *v1.PartitionUpdateRequest) (*v1.PartitionResponse, error) {
	return nil, nil
}

func (s *partitionService) Get(ctx context.Context, req *v1.PartitionGetRequest) (*v1.PartitionResponse, error) {
	return nil, nil
}

func (s *partitionService) Find(ctx context.Context, req *v1.PartitionFindRequest) (*v1.PartitionListResponse, error) {
	return nil, nil
}

func (s *partitionService) List(ctx context.Context, req *v1.PartitionListRequest) (*v1.PartitionListResponse, error) {
	return nil, nil
}
