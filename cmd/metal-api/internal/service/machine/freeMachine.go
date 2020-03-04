package machine

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"go.uber.org/zap"
	"net/http"
)

func (r machineResource) addFreeMachineRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.DELETE("/{id}/free").
		To(helper.Editor(r.freeMachine)).
		Operation("freeMachine").
		Doc("free a machine").
		Param(ws.PathParameter("id", "identifier of the machine").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "OK", v1.MachineResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r machineResource) freeMachine(request *restful.Request, response *restful.Response) {
	err := r.reinstallOrDeleteMachine(request, response, nil)
	helper.CheckError(request, response, utils.CurrentFuncName(), err)
}

func (r machineResource) releaseMachineNetworks(machine *metal.Machine, machineNetworks []*metal.MachineNetwork) error {
	for _, machineNetwork := range machineNetworks {
		for _, ipString := range machineNetwork.IPs {
			ip, err := r.DS.FindIPByID(ipString)
			if err != nil {
				return err
			}
			// ignore ips that were associated with the machine for allocation but the association is not present anymore at the ip
			if !ip.HasMachineId(machine.GetID()) {
				continue
			}
			// disassociate machine from ip
			newIP := *ip
			newIP.RemoveMachineId(machine.GetID())
			err = r.DS.UpdateIP(ip, &newIP)
			if err != nil {
				return err
			}
			// static ips should not be released automatically
			if ip.Type == metal.Static {
				continue
			}
			// ips that are associated to other machines will should not be released automatically
			if len(newIP.GetMachineIds()) > 0 {
				continue
			}
			// release and delete
			err = r.ipamer.ReleaseIP(*ip)
			if err != nil {
				return err
			}
			err = r.DS.DeleteIP(ip)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// reinstallOrDeleteMachine (re)installs the requested machine with given image by either allocating
// the machine if not yet allocated or not modifying any other allocation parameter than 'ImageID'
// and 'Reinstall' set to true.
// If the given image ID is nil, it deletes the machine instead.
func (r machineResource) reinstallOrDeleteMachine(request *restful.Request, response *restful.Response, imageID *string) error {
	id := request.PathParameter("id")
	m, err := r.DS.FindMachineByID(id)
	if err != nil {
		return err
	}

	if m.State.Value == metal.LockedState {
		return fmt.Errorf("machine is locked")
	}

	log := utils.Logger(request).Sugar()

	if m.Allocation != nil {
		old := *m

		if imageID == nil {
			// we drop networks of allocated machines from our database
			err = r.releaseMachineNetworks(m, m.Allocation.MachineNetworks)
			if err != nil {
				// TODO: Trigger network garbage collection
				// TODO: Check if all IPs in rethinkdb are in the IPAM and vice versa, cleanup if this is not the case
				// TODO: Check if there are network prefixes in the IPAM that are not in any of our networks
				log.Errorf("an error during releasing machine networks occurred, scheduled network garbage collection", "error", err)
				return err
			}

			m.Allocation = nil
			m.Tags = nil

			log.Infow("free machine", "machineID", id)
		} else {
			m.Allocation.ImageID = *imageID
			m.Allocation.Reinstall = true

			log.Infow("reinstall machine", "machineID", id, "imageID", *imageID)
		}

		err = r.DS.UpdateMachine(&old, m)
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return err
		}
	}

	// do the next steps in any case, so a client can call this function multiple times to
	// fire of the needed events

	sw, err := helper.SetVrfAtSwitches(r.DS, m, "")
	log.Infow("set VRF at switch", "machineID", id, "error", err)
	if err != nil {
		return err
	}

	deleteEvent := metal.MachineEvent{Type: metal.DELETE, Old: m}
	err = r.Publish(metal.TopicMachine.GetFQN(m.PartitionID), deleteEvent)
	log.Infow("published machine delete event", "machineID", id, "error", err)
	if err != nil {
		return err
	}

	switchEvent := metal.SwitchEvent{Type: metal.UPDATE, Machine: *m, Switches: sw}
	err = r.Publish(metal.TopicSwitch.GetFQN(m.PartitionID), switchEvent)
	log.Infow("published switch update event", "machineID", id, "error", err)
	if err != nil {
		return err
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, helper.MakeMachineResponse(m, r.DS, utils.Logger(request).Sugar()))
	if err != nil {
		log.Error("Failed to send response", zap.Error(err))
	}
	return err
}
