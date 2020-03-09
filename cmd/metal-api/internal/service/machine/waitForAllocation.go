package machine

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v12 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/proto/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// The MachineAllocation contains the allocated machine or an error.
type MachineAllocation struct {
	Machine *metal.Machine
	Err     error
}

// An Allocation is a queue of allocated machines. You can read the machines
// to get the next allocated one.
type Allocation <-chan MachineAllocation

// An Allocator is a callback for some piece of code if this wants to read
// allocated machines.
type Allocator func(Allocation) error

func (r *machineResource) waitForAllocation(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	ctx, cancel := context.WithCancel(request.Request.Context())
	log := utils.Logger(request)

	// after leaving waiting, stop listening for machine table changes in the background
	defer cancel()

	err := r.wait(ctx, id, log.Sugar(), func(alloc Allocation) error {
		select {
		case <-time.After(waitForServerTimeout):
			err := response.WriteHeaderAndEntity(http.StatusGatewayTimeout, httperrors.NewHTTPError(http.StatusGatewayTimeout, fmt.Errorf("server timeout")))
			if err != nil {
				zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
				return nil
			}
		case a := <-alloc:
			if a.Err != nil {
				log.Sugar().Errorw("allocation returned an error", "error", a.Err)
				return a.Err
			}

			s, p, i, ec := helper.FindMachineReferencedEntities(a.Machine, r.ds, log.Sugar())
			err := response.WriteHeaderAndEntity(http.StatusOK, v12.NewMachineResponse(a.Machine, s, p, i, ec))
			if err != nil {
				zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("client timeout")
		}
		return nil
	})
	if err != nil {
		helper.SendError(log, response, utils.CurrentFuncName(), httperrors.InternalServerError(err))
	}
}

// Wait inserts the machine with the given ID in the waittable, so
// this machine is ready for allocation. After this, this function waits
// for an update of this record in the waittable, which is a signal that
// this machine is allocated. This allocation will be signaled via the
// given allocator in a separate goroutine. The allocator is a function
// which will receive a channel and the caller has to select on this
// channel to get a result. Using a channel allows the caller of this
// function to implement timeouts to not wait forever.
// The user of this function will block until this machine is allocated.
func (r *machineResource) wait(ctx context.Context, id string, logger *zap.SugaredLogger, allocator Allocator) error {
	m, err := r.ds.FindMachineByID(id)
	if err != nil {
		return err
	}
	a := make(chan MachineAllocation, 1)

	// the machine IS already allocated, so notify this allocation back.
	if m.Allocation != nil {
		go func() {
			a <- MachineAllocation{Machine: m}
		}()
		return allocator(a)
	}

	err = r.ds.InsertWaitingMachine(m)
	if err != nil {
		return err
	}
	defer func() {
		err := r.ds.RemoveWaitingMachine(m)
		if err != nil {
			logger.Errorw("could not remove machine from wait table", "error", err)
		}
	}()

	go func() {
		changedMachine, err := r.ds.WaitForMachineAllocation(ctx, m)
		if err != nil {
			logger.Errorw("WaitForMachineAllocation returned an error", "error", err)
			a <- MachineAllocation{Err: err}
		} else {
			a <- MachineAllocation{Machine: changedMachine}
		}
		close(a)
	}()

	return allocator(a)
}
