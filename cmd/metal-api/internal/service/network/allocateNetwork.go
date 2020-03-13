package network

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful"
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/ipam"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *networkResource) allocateNetwork(request *restful.Request, response *restful.Response) {
	var requestPayload v1.NetworkAllocateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	nw := requestPayload.Network

	if nw.ProjectID.GetValue() == "" {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("projectid should not be empty")) {
			return
		}
	}
	if nw.PartitionID.GetValue() == "" {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("partitionid should not be empty")) {
			return
		}
	}

	project, err := r.mdc.Project().Get(context.Background(), &mdmv1.ProjectGetRequest{Id: nw.ProjectID.GetValue()})
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	partition, err := r.ds.FindPartition(nw.PartitionID.GetValue())
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	var superNetwork metal.Network
	err = r.ds.FindNetwork(&v1.NetworkSearchQuery{
		PartitionID:  util.StringProto(partition.ID),
		PrivateSuper: util.BoolProto(true),
	}, &superNetwork)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	nwSpec := &metal.Network{
		Base: metal.Base{
			Name:        nw.Common.Name.GetValue(),
			Description: nw.Common.Description.GetValue(),
		},
		PartitionID: partition.ID,
		ProjectID:   project.GetProject().GetMeta().GetId(),
		Labels:      nw.Labels,
	}

	network, err := createChildNetwork(r.ds, r.ipamer, nwSpec, &superNetwork, int(partition.PrivateNetworkPrefixLength))
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	usage := GetNetworkUsage(network, r.ipamer)
	err = response.WriteHeaderAndEntity(http.StatusCreated, NewNetworkResponse(network, usage))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}

func createChildNetwork(ds *datastore.RethinkStore, ipamer ipam.IPAMer, nwSpec *metal.Network, parent *metal.Network, childLength int) (*metal.Network, error) {
	vrf, err := ds.AcquireRandomUniqueInteger()
	if err != nil {
		return nil, fmt.Errorf("Could not acquire a vrf: %v", err)
	}

	childPrefix, err := createChildPrefix(parent.Prefixes, childLength, ipamer)
	if err != nil {
		return nil, err
	}

	if childPrefix == nil {
		return nil, fmt.Errorf("could not allocate child prefix in parent Network: %s", parent.ID)
	}

	nw := &metal.Network{
		Base: metal.Base{
			Name:        nwSpec.Name,
			Description: nwSpec.Description,
		},
		Prefixes:            metal.Prefixes{*childPrefix},
		DestinationPrefixes: metal.Prefixes{},
		PartitionID:         parent.PartitionID,
		ProjectID:           nwSpec.ProjectID,
		Nat:                 parent.Nat,
		PrivateSuper:        false,
		Underlay:            false,
		Vrf:                 vrf,
		ParentNetworkID:     parent.ID,
		Labels:              nwSpec.Labels,
	}

	err = ds.CreateNetwork(nw)
	if err != nil {
		return nil, err
	}

	return nw, nil
}

func createChildPrefix(parentPrefixes metal.Prefixes, childLength int, ipamer ipam.IPAMer) (*metal.Prefix, error) {
	var errors []error
	var err error
	var childPrefix *metal.Prefix
	for _, parentPrefix := range parentPrefixes {
		childPrefix, err = ipamer.AllocateChildPrefix(parentPrefix, childLength)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		if childPrefix != nil {
			break
		}
	}
	if childPrefix == nil {
		if len(errors) > 0 {
			return nil, fmt.Errorf("cannot allocate free child prefix in ipam: %v", errors)
		}
		return nil, fmt.Errorf("cannot allocate free child prefix in one of the given parent prefixes in ipam: %v", parentPrefixes)
	}

	return childPrefix, nil
}
