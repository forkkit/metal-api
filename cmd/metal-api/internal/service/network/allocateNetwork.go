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
	v12 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/proto/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *networkResource) allocateNetwork(request *restful.Request, response *restful.Response) {
	var requestPayload v12.NetworkAllocateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	var name string
	if requestPayload.Name != nil {
		name = *requestPayload.Name
	}
	var description string
	if requestPayload.Description != nil {
		description = *requestPayload.Description
	}
	var projectID string
	if requestPayload.ProjectID != nil {
		projectID = *requestPayload.ProjectID
	}
	var partitionID string
	if requestPayload.PartitionID != nil {
		partitionID = *requestPayload.PartitionID
	}

	if projectID == "" {
		if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("projectid should not be empty")) {
			return
		}
	}
	if partitionID == "" {
		if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("partitionid should not be empty")) {
			return
		}
	}

	project, err := r.mdc.Project().Get(context.Background(), &mdmv1.ProjectGetRequest{Id: projectID})
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	partition, err := r.ds.FindPartition(partitionID)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	var superNetwork metal.Network
	boolTrue := true
	err = r.ds.FindNetwork(&datastore.NetworkSearchQuery{PartitionID: &partition.ID, PrivateSuper: &boolTrue}, &superNetwork)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	nwSpec := &metal.Network{
		Base: metal.Base{
			Name:        name,
			Description: description,
		},
		PartitionID: partition.ID,
		ProjectID:   project.GetProject().GetMeta().GetId(),
		Labels:      requestPayload.Labels,
	}

	nw, err := createChildNetwork(r.ds, r.ipamer, nwSpec, &superNetwork, partition.PrivateNetworkPrefixLength)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	usage := helper.GetNetworkUsage(nw, r.ipamer)
	err = response.WriteHeaderAndEntity(http.StatusCreated, v12.NewNetworkResponse(nw, usage))
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
