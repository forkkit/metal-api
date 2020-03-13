package machine

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/sw"
	"github.com/metal-stack/metal-api/pkg/util"
	"go.uber.org/zap"
	"net/http"
)

func (r *machineResource) freeMachine(request *restful.Request, response *restful.Response) {
	err := r.reinstallOrDeleteMachine(request, response, nil)
	helper.CheckError(request, response, util.CurrentFuncName(), err)
}

func (r *machineResource) releaseMachineNetworks(machine *metal.Machine, machineNetworks []*metal.MachineNetwork) error {
	for _, machineNetwork := range machineNetworks {
		for _, ipString := range machineNetwork.IPs {
			ip, err := r.ds.FindIPByID(ipString)
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
			err = r.ds.UpdateIP(ip, &newIP)
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
			err = r.ds.DeleteIP(ip)
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
func (r *machineResource) reinstallOrDeleteMachine(request *restful.Request, response *restful.Response, imageID *string) error {
	id := request.PathParameter("id")
	m, err := r.ds.FindMachineByID(id)
	if err != nil {
		return err
	}

	if m.State.Value == metal.LockedState {
		return fmt.Errorf("machine is locked")
	}

	log := util.Logger(request).Sugar()

	// do the next steps in any case, so a client can call this function multiple times to
	// fire of the needed events

	ss, err := sw.SetVrfAtSwitches(r.ds, m, "")
	log.Infow("set VRF at switch", "machineID", id, "error", err)
	if err != nil {
		return err
	}

	switchEvent := metal.SwitchEvent{Type: metal.UPDATE, Machine: *m, Switches: ss}
	err = r.Publish(metal.TopicSwitch.GetFQN(m.PartitionID), switchEvent)
	log.Infow("published switch update event", "machineID", id, "error", err)
	if err != nil {
		log.Errorw("failed to publish switch update event", "machineID", id, "error", err)
	}

	deleteEvent := metal.MachineEvent{Type: metal.DELETE, Old: m}
	err = r.Publish(metal.TopicMachine.GetFQN(m.PartitionID), deleteEvent)
	log.Infow("published machine delete event", "machineID", id, "error", err)
	if err != nil {
		return err
	}

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

		err = r.ds.UpdateMachine(&old, m)
		if helper.CheckError(request, response, util.CurrentFuncName(), err) {
			return err
		}

		if imageID != nil {
			err = PublishMachineCmd(log, m, r, metal.MachineReinstall)
			if err != nil {
				log.Errorw("unable to publish ’Reinstall' command", "machineID", m.ID, "error", err)
			}
		}
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, MakeResponse(m, r.ds, util.Logger(request).Sugar()))
	if err != nil {
		log.Error("Failed to send response", zap.Error(err))
	}
	return err
}
