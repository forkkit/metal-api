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

func (s *machineService) Get(ctx context.Context, req *v1.MachineGetRequest) (*v1.MachineResponse, error) {
	return machine.FindMachine(s.ds, req.Identifiable.Id)
}

func (s *machineService) Find(ctx context.Context, req *v1.MachineFindRequest) (*v1.MachineListResponse, error) {
	return toMachineListResponse(machine.FindMachines(s.ds, req.MachineSearchQuery))
}

func (s *machineService) List(ctx context.Context, req *v1.MachineListRequest) (*v1.MachineListResponse, error) {
	return toMachineListResponse(machine.ListMachines(s.ds))
}

func toMachineListResponse(machines []*v1.MachineResponse, err error) (*v1.MachineListResponse, error) {
	if err != nil {
		return nil, err
	}
	return &v1.MachineListResponse{
		Machines: machines,
	}, nil
}

func (s *machineService) IPMIReport(ctx context.Context, req *v1.MachineIPMIReportRequest) (*v1.MachineIPMIReportResponse, error) {
	return machine.IPMIReport(s.ds, req)
}

func (s *machineService) Update(ctx context.Context, req *v1.MachineUpdateRequest) (*v1.MachineResponse, error) {
	return nil, nil
}

func (s *machineService) Delete(ctx context.Context, req *v1.MachineDeleteRequest) (*v1.MachineResponse, error) {
	return nil, nil
}
