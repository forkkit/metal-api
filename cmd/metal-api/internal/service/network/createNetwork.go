package network

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful"
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r networkResource) createNetwork(request *restful.Request, response *restful.Response) {
	var requestPayload v1.NetworkCreateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	var id string
	if requestPayload.ID != nil {
		id = *requestPayload.ID
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
	var vrf uint
	if requestPayload.Vrf != nil {
		vrf = *requestPayload.Vrf
	}
	vrfShared := false
	if requestPayload.VrfShared != nil {
		vrfShared = *requestPayload.VrfShared
	}

	privateSuper := requestPayload.PrivateSuper
	underlay := requestPayload.Underlay
	nat := requestPayload.Nat

	if projectID != "" {
		_, err = r.mdc.Project().Get(context.Background(), &mdmv1.ProjectGetRequest{Id: projectID})
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
	}

	if len(requestPayload.Prefixes) == 0 {
		// TODO: Should return a bad request 401
		if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("no prefixes given")) {
			return
		}
	}
	prefixes := metal.Prefixes{}
	// all Prefixes must be valid
	for _, p := range requestPayload.Prefixes {
		prefix, err := metal.NewPrefixFromCIDR(p)
		// TODO: Should return a bad request 401
		if err != nil {
			if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("given prefix %v is not a valid ip with mask: %v", p, err)) {
				return
			}
		}
		prefixes = append(prefixes, *prefix)
	}

	destPrefixes := metal.Prefixes{}
	for _, p := range requestPayload.DestinationPrefixes {
		prefix, err := metal.NewPrefixFromCIDR(p)
		if err != nil {
			if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("given prefix %v is not a valid ip with mask: %v", p, err)) {
				return
			}
		}
		destPrefixes = append(destPrefixes, *prefix)
	}

	allNws, err := r.DS.ListNetworks()
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	existingPrefixes := metal.Prefixes{}
	existingPrefixesMap := make(map[string]bool)
	for _, nw := range allNws {
		for _, p := range nw.Prefixes {
			_, ok := existingPrefixesMap[p.String()]
			if !ok {
				existingPrefixes = append(existingPrefixes, p)
				existingPrefixesMap[p.String()] = true
			}
		}
	}

	err = r.ipamer.PrefixesOverlapping(existingPrefixes, prefixes)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	var partitionID string
	if requestPayload.PartitionID != nil {
		partition, err := r.DS.FindPartition(*requestPayload.PartitionID)
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}

		if privateSuper {
			boolTrue := true
			err := r.DS.FindNetwork(&datastore.NetworkSearchQuery{PartitionID: &partition.ID, PrivateSuper: &boolTrue}, &metal.Network{})
			if err != nil {
				if !metal.IsNotFound(err) {
					if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
						return
					}
				}
			} else {
				if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("partition with id %q already has a private super network", partition.ID)) {
					return
				}
			}
		}
		if underlay {
			boolTrue := true
			err := r.DS.FindNetwork(&datastore.NetworkSearchQuery{PartitionID: &partition.ID, Underlay: &boolTrue}, &metal.Network{})
			if err != nil {
				if !metal.IsNotFound(err) {
					if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
						return
					}
				}
			} else {
				if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("partition with id %q already has an underlay network", partition.ID)) {
					return
				}
			}
		}
		partitionID = partition.ID
	}

	if (privateSuper || underlay) && nat {
		helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("private super or underlay network is not supposed to NAT"))
		return
	}

	if vrf != 0 {
		_, err := r.DS.AcquireUniqueInteger(vrf)
		if err != nil {
			if !metal.IsConflict(err) {
				if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("could not acquire vrf: %v", err)) {
					return
				}
			}
			if !vrfShared {
				if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("cannot acquire a unique vrf id twice except vrfShared is set to true: %v", err)) {
					return
				}
			}
		}
	}

	nw := &metal.Network{
		Base: metal.Base{
			ID:          id,
			Name:        name,
			Description: description,
		},
		Prefixes:            prefixes,
		DestinationPrefixes: destPrefixes,
		PartitionID:         partitionID,
		ProjectID:           projectID,
		Nat:                 nat,
		PrivateSuper:        privateSuper,
		Underlay:            underlay,
		Vrf:                 vrf,
	}

	for _, p := range nw.Prefixes {
		err := r.ipamer.CreatePrefix(p)
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
	}

	err = r.DS.CreateNetwork(nw)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	usage := helper.GetNetworkUsage(nw, r.ipamer)
	err = response.WriteHeaderAndEntity(http.StatusCreated, v1.NewNetworkResponse(nw, usage))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
