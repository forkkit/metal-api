package network

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful"
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *networkResource) createNetwork(request *restful.Request, response *restful.Response) {
	var requestPayload v1.NetworkCreateRequest
	err := request.ReadEntity(&requestPayload)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	nw := requestPayload.Network
	nwi := requestPayload.NetworkImmutable

	if nw.ProjectID.GetValue() != "" {
		_, err = r.mdc.Project().Get(context.Background(), &mdmv1.ProjectGetRequest{Id: nw.ProjectID.GetValue()})
		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
	}

	if len(nwi.Prefixes) == 0 {
		// TODO: Should return a bad request 401
		if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("no prefixes given")) {
			return
		}
	}
	prefixes := metal.Prefixes{}
	// all Prefixes must be valid
	for _, p := range nwi.Prefixes {
		prefix, err := metal.NewPrefixFromCIDR(p)
		// TODO: Should return a bad request 401
		if err != nil {
			if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("given prefix %v is not a valid ip with mask: %v", p, err)) {
				return
			}
		}
		prefixes = append(prefixes, *prefix)
	}

	destPrefixes := metal.Prefixes{}
	for _, p := range nwi.DestinationPrefixes {
		prefix, err := metal.NewPrefixFromCIDR(p)
		if err != nil {
			if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("given prefix %v is not a valid ip with mask: %v", p, err)) {
				return
			}
		}
		destPrefixes = append(destPrefixes, *prefix)
	}

	allNws, err := r.ds.ListNetworks()
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
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
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	var partitionID string
	if nw.PartitionID != nil {
		partition, err := r.ds.FindPartition(nw.PartitionID.GetValue())
		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}

		if nwi.PrivateSuper {
			err := r.ds.FindNetwork(&v1.NetworkSearchQuery{
				PartitionID:  util.StringProto(partition.ID),
				PrivateSuper: util.BoolProto(true),
			}, &metal.Network{})
			if err != nil {
				if !metal.IsNotFound(err) {
					if service.CheckError(request, response, util.CurrentFuncName(), err) {
						return
					}
				}
			} else {
				if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("partition with id %q already has a private super network", partition.ID)) {
					return
				}
			}
		}
		if nwi.Underlay {
			err := r.ds.FindNetwork(&v1.NetworkSearchQuery{
				PartitionID:  util.StringProto(partition.ID),
				PrivateSuper: util.BoolProto(true),
			}, &metal.Network{})
			if err != nil {
				if !metal.IsNotFound(err) {
					if service.CheckError(request, response, util.CurrentFuncName(), err) {
						return
					}
				}
			} else {
				if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("partition with id %q already has an underlay network", partition.ID)) {
					return
				}
			}
		}
		partitionID = partition.ID
	}

	if (nwi.PrivateSuper || nwi.Underlay) && nwi.Nat {
		service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("private super or underlay network is not supposed to NAT"))
		return
	}

	if nwi.Vrf.GetValue() != 0 {
		_, err := r.ds.AcquireUniqueInteger(uint(nwi.Vrf.GetValue()))
		if err != nil {
			if !metal.IsConflict(err) {
				if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("could not acquire vrf: %v", err)) {
					return
				}
			}
			if !nwi.VrfShared.GetValue() {
				if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("cannot acquire a unique vrf id twice except vrfShared is set to true: %v", err)) {
					return
				}
			}
		}
	}

	network := &metal.Network{
		Base: metal.Base{
			ID:          nw.Common.Meta.Id,
			Name:        nw.Common.Name.GetValue(),
			Description: nw.Common.Description.GetValue(),
		},
		Prefixes:            prefixes,
		DestinationPrefixes: destPrefixes,
		PartitionID:         partitionID,
		ProjectID:           nw.ProjectID.GetValue(),
		Nat:                 nwi.Nat,
		PrivateSuper:        nwi.PrivateSuper,
		Underlay:            nwi.Underlay,
		Vrf:                 uint(nwi.Vrf.GetValue()),
	}

	for _, p := range network.Prefixes {
		err := r.ipamer.CreatePrefix(p)
		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
	}

	err = r.ds.CreateNetwork(network)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	usage := GetNetworkUsage(network, r.ipamer)
	err = response.WriteHeaderAndEntity(http.StatusCreated, helper.NewNetworkResponse(network, usage))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
