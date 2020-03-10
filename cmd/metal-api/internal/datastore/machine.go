package datastore

import (
	"context"
	"fmt"

	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
)

// FindMachineByID returns a machine for a given id.
func (rs *RethinkStore) FindMachineByID(id string) (*metal.Machine, error) {
	var m metal.Machine
	err := rs.findEntityByID(rs.machineTable(), &m, id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// FindMachine returns a machine by the given query, fails if there is no record or multiple records found.
func (rs *RethinkStore) FindMachine(q *v1.MachineSearchQuery, ms *metal.Machine) error {
	return rs.findEntity(q.GenerateTerm(rs), &ms)
}

// SearchMachines returns the result of the machines search request query.
func (rs *RethinkStore) SearchMachines(q *v1.MachineSearchQuery, ms *metal.Machines) error {
	return rs.searchEntities(q.GenerateTerm(rs), ms)
}

// ListMachines returns all machines.
func (rs *RethinkStore) ListMachines() (metal.Machines, error) {
	ms := make(metal.Machines, 0)
	err := rs.listEntities(rs.machineTable(), &ms)
	return ms, err
}

// CreateMachine creates a new machine in the database as "unallocated new machines".
// If the given machine has an allocation, the function returns an error because
// allocated machines cannot be created. If there is already a machine with the
// given ID in the database it will be replaced the the given machine.
// CreateNetwork creates a new network.
func (rs *RethinkStore) CreateMachine(m *metal.Machine) error {
	if m.Allocation != nil {
		return fmt.Errorf("a machine cannot be created when it is allocated: %q: %+v", m.ID, *m.Allocation)
	}
	return rs.createEntity(rs.machineTable(), m)
}

// DeleteMachine removes a machine from the database.
func (rs *RethinkStore) DeleteMachine(m *metal.Machine) error {
	return rs.deleteEntity(rs.machineTable(), m)
}

// UpdateMachine replaces a machine in the database if the 'changed' field of
// the old value equals the 'changed' field of the recored in the database.
func (rs *RethinkStore) UpdateMachine(oldMachine *metal.Machine, newMachine *metal.Machine) error {
	return rs.updateEntity(rs.machineTable(), newMachine, oldMachine)
}

// InsertWaitingMachine adds a machine to the wait table.
func (rs *RethinkStore) InsertWaitingMachine(m *metal.Machine) error {
	// does not prohibit concurrent wait calls for the same UUID
	return rs.upsertEntity(rs.waitTable(), m)
}

// RemoveWaitingMachine removes a machine from the wait table.
func (rs *RethinkStore) RemoveWaitingMachine(m *metal.Machine) error {
	return rs.deleteEntity(rs.waitTable(), m)
}

// UpdateWaitingMachine updates a machine in the wait table with the given machine
func (rs *RethinkStore) UpdateWaitingMachine(m *metal.Machine) error {
	_, err := rs.waitTable().Get(m.ID).Update(m).RunWrite(rs.session)
	return err
}

// WaitForMachineAllocation listens on changes on the wait table for a given machine and returns the changed machine.
func (rs *RethinkStore) WaitForMachineAllocation(ctx context.Context, m *metal.Machine) (*metal.Machine, error) {
	type responseType struct {
		NewVal metal.Machine `rethinkdb:"new_val" json:"new_val"`
		OldVal metal.Machine `rethinkdb:"old_val" json:"old_val"`
	}
	var response responseType
	err := rs.listenForEntityChange(ctx, rs.waitTable(), m, response)
	if err != nil {
		return nil, err
	}

	if response.NewVal.ID == "" {
		// the machine was taken out of the wait table and not allocated
		return nil, fmt.Errorf("machine %q was taken out of the wait table", m.ID)
	}

	// the machine was really allocated!
	return &response.NewVal, nil
}

// FindAvailableMachine returns an available machine that momentarily also sits in the wait table.
func (rs *RethinkStore) FindAvailableMachine(partitionid, sizeid string) (*metal.Machine, error) {
	q := *rs.waitTable()
	q = q.Filter(map[string]interface{}{
		"allocation":  nil,
		"partitionid": partitionid,
		"sizeid":      sizeid,
		"state": map[string]string{
			"value": string(metal.AvailableState),
		},
	})

	var available metal.Machines
	err := rs.searchEntities(&q, &available)
	if err != nil {
		return nil, err
	}

	if len(available) < 1 {
		return nil, fmt.Errorf("no machine available")
	}

	// we actually return the machine from the machine table, not from the wait table
	// otherwise we will get in trouble with update operations on the machine table because
	// we have mixed timestamps with the entity from the wait table...
	m, err := rs.FindMachineByID(available[0].ID)
	if err != nil {
		return nil, err
	}

	return m, nil
}
