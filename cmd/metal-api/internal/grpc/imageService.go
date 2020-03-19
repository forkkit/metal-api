package grpc

import (
	"context"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
)

func newImageService(ds *datastore.RethinkStore) *imageService {
	return &imageService{
		ds: ds,
	}
}

type imageService struct {
	ds *datastore.RethinkStore
}

func (s *imageService) Create(ctx context.Context, req *v1.ImageCreateRequest) (*v1.ImageResponse, error) {
	return nil, nil
}

func (s *imageService) Update(ctx context.Context, req *v1.ImageUpdateRequest) (*v1.ImageResponse, error) {
	return nil, nil
}

func (s *imageService) Delete(ctx context.Context, req *v1.ImageDeleteRequest) (*v1.ImageResponse, error) {
	return nil, nil
}

func (s *imageService) Get(ctx context.Context, req *v1.ImageGetRequest) (*v1.ImageResponse, error) {
	return nil, nil
}

func (s *imageService) Find(ctx context.Context, req *v1.ImageFindRequest) (*v1.ImageListResponse, error) {
	return nil, nil
}

func (s *imageService) List(ctx context.Context, req *v1.ImageListRequest) (*v1.ImageListResponse, error) {
	return nil, nil
}
