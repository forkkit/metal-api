// +build integration

package machine

import (
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"net"
	"net/http"
	"testing"

	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPrivateSuperCidr = "192.168.0.0/20"

func TestMachineAllocationIntegration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	te := CreateTestEnvironment(t)
	defer te.Teardown()

	// Empty DB with empty alloc request
	machine := v1.MachineAllocateRequest{}

	// Register a machine
	mrr := v1.MachineRegisterRequest{
		UUID:        "test-uuid",
		PartitionID: "test-partition",
		RackID:      "test-rack",
		Hardware: &v1.MachineHardwareExtended{
			Base: &v1.MachineHardwareBase{
				CpuCores: 8,
				Memory:   1500,
				Disks: []*v1.MachineBlockDevice{
					{
						Name: "sda",
						Size: 2500,
					},
				},
			},
			Nics: []*v1.MachineNicExtended{
				{
					MachineNic: &v1.MachineNic{
						Name:       "eth0",
						MacAddress: "aa:aa:aa:aa:aa:aa",
					},
					Neighbors: []*v1.MachineNicExtended{
						{
							MachineNic: &v1.MachineNic{
								Name:       "swp1",
								MacAddress: "bb:aa:aa:aa:aa:aa",
							},
						},
					},
				},
			},
		},
	}

	var registeredMachine v1.MachineResponse
	status := te.MachineRegister(t, mrr, &registeredMachine)
	require.Equal(http.StatusCreated, status)
	require.NotNil(registeredMachine)
	assert.Equal(mrr.PartitionID, registeredMachine.Machine.PartitionResponse.Partition.Common.Meta.Id)
	assert.Equal(mrr.RackID, registeredMachine.Machine.RackID)
	assert.Equal("test-size", registeredMachine.Machine.SizeResponse.Size.Common.Meta.Id)
	assert.Len(mrr.Hardware.Nics, 1)
	assert.Equal(mrr.Hardware.Nics[0].MachineNic.MacAddress, registeredMachine.Machine.Hardware.Nics[0].MacAddress)

	go te.MachineWait("test-uuid")

	// DB contains at least a machine which is allocatable
	machine = v1.MachineAllocateRequest{
		ImageID:     "test-image",
		PartitionID: "test-partition",
		ProjectID:   te.PrivateNetwork.Network.ProjectID.GetValue(),
		SizeID:      "test-size",
	}

	var allocatedMachine v1.MachineResponse
	status = te.MachineAllocate(t, machine, &allocatedMachine)
	require.Equal(http.StatusOK, status)
	require.NotNil(allocatedMachine)
	require.NotNil(allocatedMachine.Machine.Allocation)
	require.NotNil(allocatedMachine.Machine.Allocation.ImageResponse.Image)
	assert.Equal(machine.ImageID, allocatedMachine.Machine.Allocation.ImageResponse.Image.Common.Meta.Id)
	assert.Equal(machine.ProjectID, allocatedMachine.Machine.Allocation.ProjectID)
	assert.Len(allocatedMachine.Machine.Allocation.MachineNetworks, 1)
	assert.True(allocatedMachine.Machine.Allocation.MachineNetworks[0].Private)
	assert.NotEmpty(allocatedMachine.Machine.Allocation.MachineNetworks[0].Vrf)
	assert.GreaterOrEqual(allocatedMachine.Machine.Allocation.MachineNetworks[0].Vrf, datastore.IntegerPoolRangeMin)
	assert.LessOrEqual(allocatedMachine.Machine.Allocation.MachineNetworks[0].Vrf, datastore.IntegerPoolRangeMax)
	assert.GreaterOrEqual(allocatedMachine.Machine.Allocation.MachineNetworks[0].ASN, metal.ASNBase)
	assert.Len(allocatedMachine.Machine.Allocation.MachineNetworks[0].IPs, 1)
	_, ipnet, _ := net.ParseCIDR(testPrivateSuperCidr)
	ip := net.ParseIP(allocatedMachine.Machine.Allocation.MachineNetworks[0].IPs[0])
	assert.True(ipnet.Contains(ip), "%s must be within %s", ip, ipnet)

	// Free machine for next test
	status = te.MachineFree(t, "test-uuid", &v1.MachineResponse{})
	require.Equal(http.StatusOK, status)

	go te.MachineWait("test-uuid")

	// DB contains at least a machine which is allocatable
	machine = v1.MachineAllocateRequest{
		ImageID:     "test-image",
		PartitionID: "test-partition",
		ProjectID:   te.PrivateNetwork.Network.ProjectID.GetValue(),
		SizeID:      "test-size",
		Networks: []*v1.MachineAllocationNetwork{
			{
				NetworkID: te.PrivateNetwork.Network.Common.Meta.Id,
			},
		},
	}

	allocatedMachine = v1.MachineResponse{}
	status = te.MachineAllocate(t, machine, &allocatedMachine)
	require.Equal(http.StatusOK, status)
	require.NotNil(allocatedMachine)
	require.NotNil(allocatedMachine.Machine.Allocation)
	require.NotNil(allocatedMachine.Machine.Allocation.ImageResponse.Image)
	assert.Equal(machine.ImageID, allocatedMachine.Machine.Allocation.ImageResponse.Image.Common.Meta.Id)
	assert.Equal(machine.ProjectID, allocatedMachine.Machine.Allocation.ProjectID)
	assert.Len(allocatedMachine.Machine.Allocation.MachineNetworks, 1)
	assert.True(allocatedMachine.Machine.Allocation.MachineNetworks[0].Private)
	assert.NotEmpty(allocatedMachine.Machine.Allocation.MachineNetworks[0].Vrf)
	assert.GreaterOrEqual(allocatedMachine.Machine.Allocation.MachineNetworks[0].Vrf, datastore.IntegerPoolRangeMin)
	assert.LessOrEqual(allocatedMachine.Machine.Allocation.MachineNetworks[0].Vrf, datastore.IntegerPoolRangeMax)
	assert.GreaterOrEqual(allocatedMachine.Machine.Allocation.MachineNetworks[0].ASN, metal.ASNBase)
	assert.Len(allocatedMachine.Machine.Allocation.MachineNetworks[0].IPs, 1)
	_, ipnet, _ = net.ParseCIDR(te.PrivateNetwork.NetworkImmutable.Prefixes[0])
	ip = net.ParseIP(allocatedMachine.Machine.Allocation.MachineNetworks[0].IPs[0])
	assert.True(ipnet.Contains(ip), "%s must be within %s", ip, ipnet)

	// Check if allocated machine created a machine <-> switch connection
	var foundSwitch v1.SwitchResponse
	status = te.SwitchGet(t, "test-switch01", &foundSwitch)
	require.Equal(http.StatusOK, status)
	require.NotNil(foundSwitch)
	require.Equal("test-switch01", foundSwitch.Switch.Common.Meta.Id)

	require.Len(foundSwitch.Connections, 1)
	require.Equal("swp1", foundSwitch.Connections[0].Nic.Name, "we expected exactly one connection from one allocated machine->switch.swp1")
	require.Equal("bb:aa:aa:aa:aa:aa", foundSwitch.Connections[0].Nic.MacAddress)
	require.Equal("test-uuid", foundSwitch.Connections[0].MachineID, "the allocated machine ID must be connected to swp1")

	require.Len(foundSwitch.Switch.Nics, 1)
	require.NotNil(foundSwitch.Switch.Nics[0].BGPFilter)
	require.Len(foundSwitch.Switch.Nics[0].BGPFilter.CIDRs, 1, "on this switch port, only the cidrs from the allocated machine are allowed.")
	require.Equal(allocatedMachine.Machine.Allocation.MachineNetworks[0].Prefixes[0], foundSwitch.Switch.Nics[0].BGPFilter.CIDRs[0], "exactly the prefixes of the allocated machine must be allowed on this switch port")
	require.Empty(foundSwitch.Switch.Nics[0].BGPFilter.VNIs, "to this switch port a machine with no evpn connections, so no vni filter")

	// Free machine for next test
	status = te.MachineFree(t, "test-uuid", &v1.MachineResponse{})
	require.Equal(http.StatusOK, status)

	// Check on the switch that connections still exists, but filters are nil,
	// this ensures that the freeMachine call executed and reset the machine<->switch configuration items.
	status = te.SwitchGet(t, "test-switch01", &foundSwitch)
	require.Equal(http.StatusOK, status)
	require.NotNil(foundSwitch)
	require.Equal("test-switch01", foundSwitch.Switch.Common.Meta.Id)

	require.Len(foundSwitch.Connections, 1, "machine is free for further allocations, but still connected to this switch port")
	require.Equal("swp1", foundSwitch.Connections[0].Nic.Name, "we expected exactly one connection from one allocated machine->switch.swp1")
	require.Equal("bb:aa:aa:aa:aa:aa", foundSwitch.Connections[0].Nic.MacAddress)
	require.Equal("test-uuid", foundSwitch.Connections[0].MachineID, "the allocated machine ID must be connected to swp1")

	require.Len(foundSwitch.Switch.Nics, 1)
	require.Nil(foundSwitch.Switch.Nics[0].BGPFilter, "no machine allocated anymore")

}
