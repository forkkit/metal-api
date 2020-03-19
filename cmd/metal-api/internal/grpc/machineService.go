package grpc

import (
	"context"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/machine"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
)

func newMachineService(ds *datastore.RethinkStore) *machineService {
	return &machineService{
		ds: ds,
	}
}

type machineService struct {
	ds *datastore.RethinkStore
}

func (s *machineService) Create(ctx context.Context, req *v1.MachineCreateRequest) (*v1.MachineResponse, error) {
	return nil, nil
}

func (s *machineService) Update(ctx context.Context, req *v1.MachineUpdateRequest) (*v1.MachineResponse, error) {
	return nil, nil
}

func (s *machineService) Delete(ctx context.Context, req *v1.MachineDeleteRequest) (*v1.MachineResponse, error) {
	return nil, nil
}

func (s *machineService) Get(ctx context.Context, req *v1.MachineGetRequest) (*v1.MachineResponse, error) {
	return machine.FindMachine(s.ds, req.Identifiable.Id)
}

func (s *machineService) Find(ctx context.Context, req *v1.MachineFindRequest) (*v1.MachineListResponse, error) {
	mm, err := machine.FindMachines(s.ds, req.MachineSearchQuery)
	if err != nil {
		return nil, err
	}
	return &v1.MachineListResponse{
		Machines: mm,
	}, nil
}

func (s *machineService) List(ctx context.Context, req *v1.MachineListRequest) (*v1.MachineListResponse, error) {
	return nil, nil
}
