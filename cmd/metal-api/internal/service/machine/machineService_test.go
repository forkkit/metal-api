package machine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"net/http"
	"net/http/httptest"
	"testing"

	goipam "github.com/metal-stack/go-ipam"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/ipam"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/testdata"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

type emptyPublisher struct {
	doPublish func(topic string, data interface{}) error
}

func (p *emptyPublisher) Publish(topic string, data interface{}) error {
	if p.doPublish != nil {
		return p.doPublish(topic, data)
	}
	return nil
}

func (p *emptyPublisher) CreateTopic(string) error {
	return nil
}

func TestGetMachines(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	machineService := NewMachineService(ds, &emptyPublisher{}, ipam.New(goipam.New()), nil)
	container := restful.NewContainer().Add(machineService)
	req := httptest.NewRequest("GET", "/v1/machine", nil)
	container = service.InjectViewer(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result []*v1.MachineResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Len(t, result, len(testdata.TestMachines))
	require.Equal(t, testdata.M1.ID, result[0].Common.Meta.Id)
	require.Equal(t, testdata.M1.ID, result[0].Machine.Common.Meta.Id)
	require.Equal(t, testdata.M1.Allocation.Name, result[0].Machine.Allocation.Name)
	require.Equal(t, testdata.Sz1.Name, result[0].Machine.SizeResponse.Size.Common.Name.GetValue())
	require.Equal(t, testdata.Partition1.Name, result[0].Machine.PartitionResponse.Partition.Common.Name.GetValue())
	require.Equal(t, testdata.M2.ID, result[1].Common.Meta.Id)
	require.Equal(t, testdata.M2.ID, result[1].Machine.Common.Meta.Id)
}

func TestRegisterMachine(t *testing.T) {
	data := []struct {
		name                 string
		uuid                 string
		partitionid          string
		numcores             int
		memory               int
		dbpartitions         []metal.Partition
		dbsizes              []metal.Size
		dbmachines           metal.Machines
		expectedStatus       int
		expectedErrorMessage string
		expectedSizeName     string
	}{
		{
			name:             "insert new",
			uuid:             "0",
			partitionid:      "0",
			dbpartitions:     []metal.Partition{testdata.Partition1},
			dbsizes:          []metal.Size{testdata.Sz1},
			numcores:         1,
			memory:           100,
			expectedStatus:   http.StatusCreated,
			expectedSizeName: testdata.Sz1.Name,
		},
		{
			name:             "insert existing",
			uuid:             "1",
			partitionid:      "1",
			dbpartitions:     []metal.Partition{testdata.Partition1},
			dbsizes:          []metal.Size{testdata.Sz1},
			dbmachines:       metal.Machines{testdata.M1},
			numcores:         1,
			memory:           100,
			expectedStatus:   http.StatusOK,
			expectedSizeName: testdata.Sz1.Name,
		},
		{
			name:                 "empty uuid",
			uuid:                 "",
			partitionid:          "1",
			dbpartitions:         []metal.Partition{testdata.Partition1},
			dbsizes:              []metal.Size{testdata.Sz1},
			expectedStatus:       http.StatusUnprocessableEntity,
			expectedErrorMessage: "uuid cannot be empty",
		},
		{
			name:                 "empty partition",
			uuid:                 "1",
			partitionid:          "",
			dbpartitions:         nil,
			dbsizes:              []metal.Size{testdata.Sz1},
			expectedStatus:       http.StatusNotFound,
			expectedErrorMessage: "no partition with id \"\" found",
		},
		{
			name:             "new with unknown size",
			uuid:             "0",
			partitionid:      "1",
			dbpartitions:     []metal.Partition{testdata.Partition1},
			dbsizes:          []metal.Size{testdata.Sz1},
			numcores:         2,
			memory:           100,
			expectedStatus:   http.StatusCreated,
			expectedSizeName: metal.UnknownSize.Name,
		},
	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {
			ds, mock := datastore.InitMockDB()
			mock.On(r.DB("mockdb").Table("partition").Get(test.partitionid)).Return(test.dbpartitions, nil)

			if len(test.dbmachines) > 0 {
				mock.On(r.DB("mockdb").Table("size").Get(test.dbmachines[0].SizeID)).Return([]metal.Size{testdata.Sz1}, nil)
				mock.On(r.DB("mockdb").Table("machine").Get(test.dbmachines[0].ID).Replace(r.MockAnything())).Return(testdata.EmptyResult, nil)
			} else {
				mock.On(r.DB("mockdb").Table("machine").Get("0")).Return(nil, nil)
				mock.On(r.DB("mockdb").Table("machine").Insert(r.MockAnything(), r.InsertOpts{
					Conflict: "replace",
				})).Return(testdata.EmptyResult, nil)
			}
			mock.On(r.DB("mockdb").Table("size").Get(metal.UnknownSize.ID)).Return([]metal.Size{*metal.UnknownSize}, nil)
			mock.On(r.DB("mockdb").Table("switch").Filter(r.MockAnything(), r.FilterOpts{})).Return([]metal.Switch{}, nil)
			mock.On(r.DB("mockdb").Table("event").Filter(r.MockAnything(), r.FilterOpts{})).Return([]metal.ProvisioningEventContainer{}, nil)
			mock.On(r.DB("mockdb").Table("event").Insert(r.MockAnything(), r.InsertOpts{})).Return(testdata.EmptyResult, nil)
			testdata.InitMockDBData(mock)

			registerRequest := &v1.MachineRegisterRequest{
				UUID:        test.uuid,
				PartitionID: test.partitionid,
				RackID:      "1",
				IPMI: &v1.MachineIPMI{
					Address:    testdata.IPMI1.Address,
					Interface:  testdata.IPMI1.Interface,
					MacAddress: testdata.IPMI1.MacAddress,
					Fru: &v1.MachineFru{
						ChassisPartNumber:   util.StringProto(testdata.IPMI1.Fru.ChassisPartNumber),
						ChassisPartSerial:   util.StringProto(testdata.IPMI1.Fru.ChassisPartSerial),
						BoardMfg:            util.StringProto(testdata.IPMI1.Fru.BoardMfg),
						BoardMfgSerial:      util.StringProto(testdata.IPMI1.Fru.BoardMfgSerial),
						BoardPartNumber:     util.StringProto(testdata.IPMI1.Fru.BoardPartNumber),
						ProductManufacturer: util.StringProto(testdata.IPMI1.Fru.ProductManufacturer),
						ProductPartNumber:   util.StringProto(testdata.IPMI1.Fru.ProductPartNumber),
						ProductSerial:       util.StringProto(testdata.IPMI1.Fru.ProductSerial),
					},
				},
				Hardware: &v1.MachineHardwareExtended{
					Base: &v1.MachineHardwareBase{
						CpuCores: uint32(test.numcores),
						Memory:   uint64(test.memory),
					},
				},
			}

			js, _ := json.Marshal(registerRequest)
			body := bytes.NewBuffer(js)
			machineService := NewMachineService(ds, &emptyPublisher{}, ipam.New(goipam.New()), nil)
			container := restful.NewContainer().Add(machineService)
			req := httptest.NewRequest("POST", "/v1/machine/register", body)
			req.Header.Add("Content-Type", "application/json")
			container = service.InjectEditor(container, req)
			w := httptest.NewRecorder()
			container.ServeHTTP(w, req)

			resp := w.Result()
			require.Equal(t, test.expectedStatus, resp.StatusCode, w.Body.String())

			if test.expectedStatus > 300 {
				var result httperrors.HTTPErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&result)

				require.Nil(t, err)
				require.Regexp(t, test.expectedErrorMessage, result.Message)
			} else {
				var result v1.MachineResponse
				err := json.NewDecoder(resp.Body).Decode(&result)

				require.Nil(t, err)
				expectedid := "0"
				if len(test.dbmachines) > 0 {
					expectedid = test.dbmachines[0].ID
				}
				require.Equal(t, expectedid, result.Machine.Common.Meta.Id)
				require.Equal(t, "1", result.Machine.RackID)
				require.Equal(t, test.expectedSizeName, result.Machine.SizeResponse.Size.Common.Name.GetValue())
				require.Equal(t, testdata.Partition1.Name, result.Machine.PartitionResponse.Partition.Common.Name.GetValue())
			}
		})
	}
}

func TestMachineIPMIReport(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	data := []struct {
		name           string
		input          v1.MachineIpmiReport
		output         v1.MachineIpmiReportResponse
		wantStatusCode int
	}{
		{
			name: "update machine1 ipmi address",
			input: v1.MachineIpmiReport{
				PartitionID:  testdata.M1.PartitionID,
				ActiveLeases: map[string]string{testdata.M1.ID: "192.167.0.1"},
			},
			output: v1.MachineIpmiReportResponse{
				UpdatedLeases: map[string]string{testdata.M1.ID: "192.167.0.1"},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "don't update machine with unkown mac",
			input: v1.MachineIpmiReport{
				PartitionID:  testdata.M1.PartitionID,
				ActiveLeases: map[string]string{"xyz": "192.167.0.1"},
			},
			output: v1.MachineIpmiReportResponse{
				CreatedLeases: map[string]string{"xyz": "192.167.0.1"},
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {
			machineService := NewMachineService(ds, &emptyPublisher{}, ipam.New(goipam.New()), nil)
			container := restful.NewContainer().Add(machineService)
			js, _ := json.Marshal(test.input)
			body := bytes.NewBuffer(js)
			req := httptest.NewRequest("POST", fmt.Sprintf("/v1/machine/ipmi"), body)
			req.Header.Add("Content-Type", "application/json")
			container = service.InjectEditor(container, req)
			w := httptest.NewRecorder()
			container.ServeHTTP(w, req)

			resp := w.Result()
			require.Equal(t, test.wantStatusCode, resp.StatusCode, w.Body.String())

			var result v1.MachineIpmiReportResponse
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.Nil(t, err)
			require.Equal(t, test.output, result)
		})
	}
}

func TestMachineFindIPMI(t *testing.T) {
	data := []struct {
		name           string
		machine        *metal.Machine
		wantStatusCode int
	}{
		{
			name:           "retrieve machine1 ipmi",
			machine:        &testdata.M1,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "retrieve machine2 ipmi",
			machine:        &testdata.M2,
			wantStatusCode: http.StatusOK,
		},
	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {
			ds, mock := datastore.InitMockDB()
			mock.On(r.DB("mockdb").Table("machine").Filter(r.MockAnything())).Return([]interface{}{*test.machine}, nil)
			testdata.InitMockDBData(mock)

			machineService := NewMachineService(ds, &emptyPublisher{}, ipam.New(goipam.New()), nil)
			container := restful.NewContainer().Add(machineService)

			query := v1.MachineSearchQuery{
				ID: util.StringProto(test.machine.ID),
			}
			js, _ := json.Marshal(query)
			body := bytes.NewBuffer(js)
			req := httptest.NewRequest("POST", "/v1/machine/ipmi/find", body)
			req.Header.Add("Content-Type", "application/json")
			container = service.InjectViewer(container, req)
			w := httptest.NewRecorder()
			container.ServeHTTP(w, req)

			resp := w.Result()
			require.Equal(t, test.wantStatusCode, resp.StatusCode, w.Body.String())

			var results []*v1.MachineIPMIResponse
			err := json.NewDecoder(resp.Body).Decode(&results)

			require.Nil(t, err)
			require.Len(t, results, 1)

			result := results[0]

			require.Equal(t, test.machine.IPMI.Address, result.IPMI.Address)
			require.Equal(t, test.machine.IPMI.Interface, result.IPMI.Interface)
			require.Equal(t, test.machine.IPMI.User, result.IPMI.User)
			require.Equal(t, test.machine.IPMI.Password, result.IPMI.Password)
			require.Equal(t, test.machine.IPMI.MacAddress, result.IPMI.MacAddress)

			require.Equal(t, test.machine.IPMI.Fru.ChassisPartNumber, result.IPMI.Fru.ChassisPartNumber.GetValue())
			require.Equal(t, test.machine.IPMI.Fru.ChassisPartSerial, result.IPMI.Fru.ChassisPartSerial.GetValue())
			require.Equal(t, test.machine.IPMI.Fru.BoardMfg, result.IPMI.Fru.BoardMfg.GetValue())
			require.Equal(t, test.machine.IPMI.Fru.BoardMfgSerial, result.IPMI.Fru.BoardMfgSerial.GetValue())
			require.Equal(t, test.machine.IPMI.Fru.BoardPartNumber, result.IPMI.Fru.BoardPartNumber.GetValue())
			require.Equal(t, test.machine.IPMI.Fru.ProductManufacturer, result.IPMI.Fru.ProductManufacturer.GetValue())
			require.Equal(t, test.machine.IPMI.Fru.ProductPartNumber, result.IPMI.Fru.ProductPartNumber.GetValue())
			require.Equal(t, test.machine.IPMI.Fru.ProductSerial, result.IPMI.Fru.ProductSerial.GetValue())
		})
	}
}

func TestFinalizeMachineAllocation(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	data := []struct {
		name           string
		machineID      string
		wantStatusCode int
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:           "finalize successfully",
			machineID:      "1",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "finalize unknown machine",
			machineID:      "999",
			wantStatusCode: http.StatusNotFound,
			wantErr:        true,
		},
		{
			name:           "finalize unallocated machine",
			machineID:      "3",
			wantStatusCode: http.StatusUnprocessableEntity,
			wantErr:        true,
			wantErrMessage: "the machine \"3\" is not allocated",
		},
	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {

			machineService := NewMachineService(ds, &emptyPublisher{}, ipam.New(goipam.New()), nil)
			container := restful.NewContainer().Add(machineService)

			finalizeRequest := v1.MachineFinalizeAllocationRequest{
				ConsolePassword: "blubber",
			}

			js, _ := json.Marshal(finalizeRequest)
			body := bytes.NewBuffer(js)
			req := httptest.NewRequest("POST", fmt.Sprintf("/v1/machine/%s/finalize-allocation", test.machineID), body)
			req.Header.Add("Content-Type", "application/json")
			container = service.InjectEditor(container, req)
			w := httptest.NewRecorder()
			container.ServeHTTP(w, req)

			resp := w.Result()
			require.Equal(t, test.wantStatusCode, resp.StatusCode, w.Body.String())

			if test.wantErr {
				var result httperrors.HTTPErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&result)

				require.Nil(t, err)
				require.Equal(t, test.wantStatusCode, result.StatusCode)
				if test.wantErrMessage != "" {
					require.Regexp(t, test.wantErrMessage, result.Message)
				}
			} else {
				var result v1.MachineResponse
				err := json.NewDecoder(resp.Body).Decode(&result)

				require.Nil(t, err)
				require.Equal(t, finalizeRequest.ConsolePassword, result.Machine.Allocation.ConsolePassword.GetValue())
			}
		})
	}
}

func TestSetMachineState(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	machineService := NewMachineService(ds, &emptyPublisher{}, ipam.New(goipam.New()), nil)
	container := restful.NewContainer().Add(machineService)

	stateRequest := v1.MachineState{
		Value:       string(metal.ReservedState),
		Description: "blubber",
	}
	js, _ := json.Marshal(stateRequest)
	body := bytes.NewBuffer(js)
	req := httptest.NewRequest("POST", "/v1/machine/1/state", body)
	req.Header.Add("Content-Type", "application/json")
	container = service.InjectEditor(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.MachineResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, "1", result.Machine.Common.Meta.Id)
	require.Equal(t, string(metal.ReservedState), result.Machine.State.Value)
	require.Equal(t, "blubber", result.Machine.State.Description)
}

func TestGetMachine(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	machineService := NewMachineService(ds, &emptyPublisher{}, ipam.New(goipam.New()), nil)
	container := restful.NewContainer().Add(machineService)
	req := httptest.NewRequest("GET", "/v1/machine/1", nil)
	container = service.InjectViewer(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.MachineResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.M1.ID, result.Machine.Common.Meta.Id)
	require.Equal(t, testdata.M1.Allocation.Name, result.Machine.Allocation.Name)
	require.Equal(t, testdata.Sz1.Name, result.Machine.SizeResponse.Size.Common.Name.GetValue())
	require.Equal(t, testdata.Img1.Name, result.Machine.Allocation.ImageResponse.Image.Common.Name.GetValue())
	require.Equal(t, testdata.Partition1.Name, result.Machine.PartitionResponse.Partition.Common.Name.GetValue())
}

func TestGetMachineNotFound(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	machineService := NewMachineService(ds, &emptyPublisher{}, ipam.New(goipam.New()), nil)
	container := restful.NewContainer().Add(machineService)
	req := httptest.NewRequest("GET", "/v1/machine/999", nil)
	container = service.InjectEditor(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusNotFound, resp.StatusCode, w.Body.String())
}

func TestFreeMachine(t *testing.T) {
	// TODO: Add tests for IPAM, verifying that networks are cleaned up properly

	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	pub := &emptyPublisher{}
	events := []string{"1-switch", "1-machine"}
	eventidx := 0
	pub.doPublish = func(topic string, data interface{}) error {
		require.Equal(t, events[eventidx], topic)
		eventidx++
		if eventidx == 0 {
			dv := data.(metal.MachineEvent)
			require.Equal(t, "1", dv.Old.ID)
		}
		return nil
	}

	machineService := NewMachineService(ds, pub, ipam.New(goipam.New()), nil)
	container := restful.NewContainer().Add(machineService)
	req := httptest.NewRequest("DELETE", "/v1/machine/1/free", nil)
	container = service.InjectEditor(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.MachineResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.M1.ID, result.Machine.Common.Meta.Id)
	require.Nil(t, result.Machine.Allocation)
	require.Empty(t, result.Machine.Tags)
}

func TestSearchMachine(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	mock.On(r.DB("mockdb").Table("machine").Filter(r.MockAnything())).Return([]interface{}{testdata.M1}, nil)
	testdata.InitMockDBData(mock)

	machineService := NewMachineService(ds, &emptyPublisher{}, ipam.New(goipam.New()), nil)
	container := restful.NewContainer().Add(machineService)
	findReq := &v1.MachineFindRequest{
		MachineSearchQuery: &v1.MachineSearchQuery{
			NicsMacAddresses: []*wrappers.StringValue{
				util.StringProto("1"),
			},
		},
	}
	requestJSON, err := json.Marshal(findReq)
	require.Nil(t, err)
	req := httptest.NewRequest("POST", "/v1/machine/find", bytes.NewBuffer(requestJSON))
	req.Header.Add("Content-Type", "application/json")
	container = service.InjectViewer(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var results []v1.MachineResponse
	err = json.NewDecoder(resp.Body).Decode(&results)

	require.Nil(t, err)
	require.Len(t, results, len(testdata.TestMachines))
	result := results[0]
	require.Equal(t, testdata.M1.ID, result.Machine.Common.Meta.Id)
	require.Equal(t, testdata.M1.Allocation.Name, result.Machine.Allocation.Name)
	require.Equal(t, testdata.Sz1.Name, result.Machine.SizeResponse.Size.Common.Name.GetValue())
	require.Equal(t, testdata.Img1.Name, result.Machine.Allocation.ImageResponse.Image.Common.Name.GetValue())
	require.Equal(t, testdata.Partition1.Name, result.Machine.PartitionResponse.Partition.Common.Name.GetValue())
}

func TestAddProvisioningEvent(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	machineService := NewMachineService(ds, &emptyPublisher{}, ipam.New(goipam.New()), nil)
	container := restful.NewContainer().Add(machineService)
	event := &v1.MachineProvisioningEvent{
		Event:   string(metal.ProvisioningEventPreparing),
		Message: util.StringProto("starting metal-hammer"),
	}
	js, _ := json.Marshal(event)
	body := bytes.NewBuffer(js)
	req := httptest.NewRequest("POST", "/v1/machine/1/event", body)
	container = service.InjectEditor(container, req)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.MachineRecentProvisioningEvents
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, "0", result.IncompleteProvisioningCycles)
	require.Len(t, result.Events, 1)
	if len(result.Events) > 0 {
		require.Equal(t, "starting metal-hammer", result.Events[0].Message.GetValue())
		require.Equal(t, string(metal.ProvisioningEventPreparing), result.Events[0].Event)
	}
}

func TestOnMachine(t *testing.T) {
	data := []struct {
		cmd      metal.MachineCommand
		endpoint string
		param    string
	}{
		{
			cmd:      metal.MachineOnCmd,
			endpoint: "on",
		},
		{
			cmd:      metal.MachineOffCmd,
			endpoint: "off",
		},
		{
			cmd:      metal.MachineResetCmd,
			endpoint: "reset",
		},
		{
			cmd:      metal.MachineBiosCmd,
			endpoint: "bios",
		},
		{
			cmd:      metal.ChassisIdentifyLEDOnCmd,
			endpoint: "chassis-identify-led-on",
		},
		{
			cmd:      metal.ChassisIdentifyLEDOnCmd,
			endpoint: "chassis-identify-led-on/test",
		},
		{
			cmd:      metal.ChassisIdentifyLEDOffCmd,
			endpoint: "chassis-identify-led-off/test",
		},
	}

	for _, d := range data {
		t.Run("cmd_"+d.endpoint, func(t *testing.T) {
			ds, mock := datastore.InitMockDB()
			testdata.InitMockDBData(mock)

			pub := &emptyPublisher{}
			pub.doPublish = func(topic string, data interface{}) error {
				require.Equal(t, "1-machine", topic)
				dv := data.(metal.MachineEvent)
				require.Equal(t, d.cmd, dv.Cmd.Command)
				require.Equal(t, "1", dv.Cmd.Target.ID)
				return nil
			}

			machineService := NewMachineService(ds, pub, ipam.New(goipam.New()), nil)

			js, _ := json.Marshal([]string{d.param})
			body := bytes.NewBuffer(js)
			container := restful.NewContainer().Add(machineService)
			req := httptest.NewRequest("POST", "/v1/machine/1/power/"+d.endpoint, body)
			container = service.InjectEditor(container, req)
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			container.ServeHTTP(w, req)

			resp := w.Result()
			require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
		})
	}
}

func TestParsePublicKey(t *testing.T) {
	pubKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDi4+MA0u/luzH2iaKnBTHzo+BEmV1MsdWtPtAps9ccD1vF94AqKtV6mm387ZhamfWUfD1b3Q5ftk56ekwZgHbk6PIUb/W4GrBD4uslTL2lzNX9v0Njo9DfapDKv4Tth6Qz5ldUb6z7IuyDmWqn3FbIPo4LOZxJ9z/HUWyau8+JMSpwIyzp2S0Gtm/pRXhbkZlr4h9jGApDQICPFGBWFEVpyOOjrS8JnEC8YzUszvbj5W1CH6Sn/DtxW0/CTAWwcjIAYYV8GlouWjjALqmjvpxO3F5kvQ1xR8IYrD86+cSCQSP4TpehftzaQzpY98fcog2YkEra+1GCY456cVSUhe1X"
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubKey))
	require.Nil(t, err)

	pubKey = ""
	_, _, _, _, err = ssh.ParseAuthorizedKey([]byte(pubKey))
	require.NotNil(t, err)

	pubKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDi4+MA0u/luzH2iaKnBTHzo+BEmV1MsdWtPtAps9ccD1vF94AqKtV6mm387ZhamfWUfD1b3Q5ftk56ekwZgHbk6PIUb/W4GrBD4uslTL2lzNX9v0Njo9DfapDKv4Tth6Qz5ldUb6z7IuyDmWqn3FbIPo4LOZxJ9z/HUWyau8+JMSpwIyzp2S0Gtm/pRXhbkZlr4h9jGApDQICPFGBWFEVpyOOjrS8JnEC8YzUszvbj5W1CH6Sn/DtxW0/CTAWwcjIAYYV8GlouWjjALqmjvpxO3F5kvQ1xR8IYrD86+cSCQSP4TpehftzaQzpY98fcog2YkEra+1GCY456cVSUhe1"
	_, _, _, _, err = ssh.ParseAuthorizedKey([]byte(pubKey))
	require.NotNil(t, err)

	pubKey = "AAAAB3NzaC1yc2EAAAADAQABAAABAQDi4+MA0u/luzH2iaKnBTHzo+BEmV1MsdWtPtAps9ccD1vF94AqKtV6mm387ZhamfWUfD1b3Q5ftk56ekwZgHbk6PIUb/W4GrBD4uslTL2lzNX9v0Njo9DfapDKv4Tth6Qz5ldUb6z7IuyDmWqn3FbIPo4LOZxJ9z/HUWyau8+JMSpwIyzp2S0Gtm/pRXhbkZlr4h9jGApDQICPFGBWFEVpyOOjrS8JnEC8YzUszvbj5W1CH6Sn/DtxW0/CTAWwcjIAYYV8GlouWjjALqmjvpxO3F5kvQ1xR8IYrD86+cSCQSP4TpehftzaQzpY98fcog2YkEra+1GCY456cVSUhe1X"
	_, _, _, _, err = ssh.ParseAuthorizedKey([]byte(pubKey))
	require.NotNil(t, err)
}

func Test_validateAllocationSpec(t *testing.T) {
	ass := assert.New(t)

	tests := []struct {
		spec     AllocationSpec
		isError  bool
		name     string
		expected string
	}{
		{
			spec: AllocationSpec{
				UUID:       "gopher-uuid",
				ProjectID:  "123",
				IsFirewall: false,
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID: "network",
					},
				},
				IPs: []string{"1.2.3.4"},
			},
			isError:  false,
			expected: "",
			name:     "auto acquire network and additional ip",
		},
		{
			spec: AllocationSpec{
				UUID:      "gopher-uuid",
				ProjectID: "123",
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     "network",
						AutoAcquireIP: util.BoolProto(true),
					},
				},
			},
			isError: false,
			name:    "good case (explicit network)",
		},
		{
			spec: AllocationSpec{
				UUID:       "gopher-uuid",
				ProjectID:  "123",
				IsFirewall: false,
			},
			isError:  false,
			expected: "",
			name:     "good case (no network)",
		},
		{
			spec: AllocationSpec{
				PartitionID: "42",
				ProjectID:   "123",
				SizeID:      "42",
			},
			isError: false,
			name:    "partition and size id for absent uuid",
		},
		{
			spec: AllocationSpec{
				PartitionID: "42",
				ProjectID:   "123",
			},
			isError:  true,
			expected: "when no machine id is given, a size id must be specified",
			name:     "missing size id",
		},
		{
			spec: AllocationSpec{
				SizeID:    "42",
				ProjectID: "123",
			},
			isError:  true,
			expected: "when no machine id is given, a partition id must be specified",
			name:     "missing partition id",
		},
		{
			spec:     AllocationSpec{},
			isError:  true,
			expected: "project id must be specified",
			name:     "absent project id",
		},
		{
			spec: AllocationSpec{
				UUID:       "gopher-uuid",
				ProjectID:  "123",
				IsFirewall: false,
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     "network",
						AutoAcquireIP: util.BoolProto(false),
					},
				},
			},
			isError:  true,
			expected: "missing ip(s) for network(s) without automatic ip allocation",
			name:     "missing ip definition for noauto network",
		},
		{
			spec: AllocationSpec{
				UUID:      "42",
				ProjectID: "123",
				IPs:       []string{"42"},
			},
			isError:  true,
			expected: `"42" is not a valid IP address`,
			name:     "illegal ip",
		},
		{
			spec: AllocationSpec{
				UUID:       "42",
				ProjectID:  "123",
				IsFirewall: true,
			},
			isError:  true,
			expected: "when no ip is given at least one auto acquire network must be specified",
			name:     "missing network/ ip in case of firewall",
		},
		{
			spec: AllocationSpec{
				UUID:       "42",
				ProjectID:  "123",
				SSHPubKeys: []string{"42"},
			},
			isError:  true,
			expected: `invalid public SSH key: 42`,
			name:     "invalid ssh",
		},
		{
			spec: AllocationSpec{
				UUID:       "gopher-uuid",
				ProjectID:  "123",
				IsFirewall: false,
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID: "network",
					},
				},
			},
			isError:  false,
			expected: "",
			name:     "implicit auto acquire network",
		},
	}

	for _, test := range tests {
		err := ValidateAllocationSpec(&test.spec)
		if test.isError {
			ass.Error(err, "Test: %s", test.name)
			ass.EqualError(err, test.expected, "Test: %s", test.name)
		} else {
			ass.NoError(err, "Test: %s", test.name)
		}
	}

}

func Test_makeMachineTags(t *testing.T) {
	type args struct {
		m        *metal.Machine
		networks AllocationNetworkMap
		userTags []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "All possible tags",
			args: args{
				m: &metal.Machine{
					Allocation: &metal.MachineAllocation{
						MachineNetworks: []*metal.MachineNetwork{
							{
								Private: true,
								ASN:     1203874,
							},
						},
					},
					RackID: "rack01",
					IPMI: metal.IPMI{
						Fru: metal.Fru{
							ChassisPartSerial: "chassis123",
						},
					},
				},
				networks: AllocationNetworkMap{
					"network-uuid-1": &AllocationNetwork{
						Network: &metal.Network{
							Labels: map[string]string{
								"external-network-label": "1",
							},
						},
					},
					"network-uuid-2": &AllocationNetwork{
						Network: &metal.Network{
							Labels: map[string]string{
								"private-network-label": "1",
							},
						},
						IsPrivate: true,
					},
				},
				userTags: []string{"usertag=something"},
			},
			want: []string{
				"external-network-label=1",
				"private-network-label=1",
				"usertag=something",
				"machine.metal-stack.io/network.primary.asn=1203874",
				"machine.metal-stack.io/rack=rack01",
				"machine.metal-stack.io/chassis=chassis123",
			},
		},
		{
			name: "private network tags higher precedence than external network tags",
			args: args{
				m: &metal.Machine{
					Allocation: &metal.MachineAllocation{
						MachineNetworks: []*metal.MachineNetwork{},
					},
				},
				networks: AllocationNetworkMap{
					"network-uuid-1": &AllocationNetwork{
						Network: &metal.Network{
							Labels: map[string]string{
								"override": "1",
							},
						},
					},
					"network-uuid-2": &AllocationNetwork{
						Network: &metal.Network{
							Labels: map[string]string{
								"override": "2",
							},
						},
						IsPrivate: true,
					},
				},
			},
			want: []string{
				"override=2",
			},
		},
		{
			name: "user tags higher precedence than network tags",
			args: args{
				m: &metal.Machine{
					Allocation: &metal.MachineAllocation{
						MachineNetworks: []*metal.MachineNetwork{},
					},
				},
				networks: AllocationNetworkMap{
					"network-uuid-1": &AllocationNetwork{
						Network: &metal.Network{
							Labels: map[string]string{
								"override": "1",
							},
						},
					},
					"network-uuid-2": &AllocationNetwork{
						Network: &metal.Network{
							Labels: map[string]string{
								"override": "2",
							},
						},
						IsPrivate: true,
					},
				},
				userTags: []string{"override=3"},
			},
			want: []string{
				"override=3",
			},
		},
		{
			name: "system tags higher precedence than user tags",
			args: args{
				m: &metal.Machine{
					Allocation: &metal.MachineAllocation{
						MachineNetworks: []*metal.MachineNetwork{
							{
								Private: true,
								ASN:     1203874,
							},
						},
					},
				},
				networks: AllocationNetworkMap{},
				userTags: []string{"machine.metal-stack.io/network.primary.asn=iamdoingsomethingevil"},
			},
			want: []string{
				"machine.metal-stack.io/network.primary.asn=1203874",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeMachineTags(tt.args.m, tt.args.networks, tt.args.userTags)

			for _, wantElement := range tt.want {
				require.Contains(t, got, wantElement, "tag not contained in result")
			}
			require.Len(t, got, len(tt.want))
		})
	}
}

func Test_gatherNetworksFromSpec(t *testing.T) {
	partitionSuperNetworks := metal.Networks{testdata.Partition1PrivateSuperNetwork, testdata.Partition2PrivateSuperNetwork}

	type mock struct {
		term     r.Term
		response interface{}
		err      error
	}
	tests := []struct {
		name                   string
		allocationSpec         *AllocationSpec
		partition              *metal.Partition
		partitionSuperNetworks metal.Networks
		mocks                  []mock
		want                   AllocationNetworkMap
		wantErr                bool
		errRegex               string
	}{
		{
			name: "no networks given",
			allocationSpec: &AllocationSpec{
				Networks: []v1.MachineAllocationNetwork{},
			},
			partition:              &testdata.Partition1,
			partitionSuperNetworks: partitionSuperNetworks,
			wantErr:                true,
			errRegex:               "no private network given",
		},
		{
			name: "private network given",
			allocationSpec: &AllocationSpec{
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     testdata.Partition1ExistingPrivateNetwork.ID,
						AutoAcquireIP: util.BoolProto(true),
					},
				},
				ProjectID: testdata.Partition1ExistingPrivateNetwork.ProjectID,
			},
			partition:              &testdata.Partition1,
			partitionSuperNetworks: partitionSuperNetworks,
			wantErr:                false,
			want: AllocationNetworkMap{
				testdata.Partition1ExistingPrivateNetwork.ID: &AllocationNetwork{
					Network:        &testdata.Partition1ExistingPrivateNetwork,
					MachineNetwork: &metal.MachineNetwork{},
					IPs:            []metal.IP{},
					Auto:           true,
					IsPrivate:      true,
				},
			},
		},
		{
			name: "private network given, but no auto acquisition and no ip provided",
			allocationSpec: &AllocationSpec{
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     testdata.Partition1ExistingPrivateNetwork.ID,
						AutoAcquireIP: util.BoolProto(false),
					},
				},
				ProjectID: testdata.Partition1ExistingPrivateNetwork.ProjectID,
			},
			partition:              &testdata.Partition1,
			partitionSuperNetworks: partitionSuperNetworks,
			wantErr:                true,
			errRegex:               "the private network has no auto ip acquisition, but no suitable ips were provided",
		},
		{
			name: "private network and internet network given",
			allocationSpec: &AllocationSpec{
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     testdata.Partition1ExistingPrivateNetwork.ID,
						AutoAcquireIP: util.BoolProto(true),
					},
					{
						NetworkID:     testdata.Partition1InternetNetwork.ID,
						AutoAcquireIP: util.BoolProto(true),
					},
				},
				ProjectID: testdata.Partition1ExistingPrivateNetwork.ProjectID,
			},
			partition:              &testdata.Partition1,
			partitionSuperNetworks: partitionSuperNetworks,
			wantErr:                false,
			want: AllocationNetworkMap{
				testdata.Partition1ExistingPrivateNetwork.ID: &AllocationNetwork{
					Network:        &testdata.Partition1ExistingPrivateNetwork,
					MachineNetwork: &metal.MachineNetwork{},
					IPs:            []metal.IP{},
					Auto:           true,
					IsPrivate:      true,
				},
				testdata.Partition1InternetNetwork.ID: &AllocationNetwork{
					Network:        &testdata.Partition1InternetNetwork,
					MachineNetwork: &metal.MachineNetwork{},
					IPs:            []metal.IP{},
					Auto:           true,
					IsPrivate:      false,
				},
			},
		},
		{
			name: "ip which does not belong to any related network given",
			allocationSpec: &AllocationSpec{
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     testdata.Partition1ExistingPrivateNetwork.ID,
						AutoAcquireIP: util.BoolProto(true),
					},
				},
				IPs:       []string{testdata.Partition2InternetIP.IPAddress},
				ProjectID: testdata.Partition1ExistingPrivateNetwork.ProjectID,
			},
			partition:              &testdata.Partition1,
			partitionSuperNetworks: partitionSuperNetworks,
			wantErr:                true,
			errRegex:               "given ip .* is not in any of the given networks",
		},
		{
			name: "private network and internet network with no auto acquired internet ip",
			allocationSpec: &AllocationSpec{
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     testdata.Partition1ExistingPrivateNetwork.ID,
						AutoAcquireIP: util.BoolProto(true),
					},
					{
						NetworkID:     testdata.Partition1InternetNetwork.ID,
						AutoAcquireIP: util.BoolProto(false),
					},
				},
				IPs:       []string{testdata.Partition1InternetIP.IPAddress},
				ProjectID: testdata.Partition1ExistingPrivateNetwork.ProjectID,
			},
			partition:              &testdata.Partition1,
			partitionSuperNetworks: partitionSuperNetworks,
			wantErr:                false,
			want: AllocationNetworkMap{
				testdata.Partition1ExistingPrivateNetwork.ID: &AllocationNetwork{
					Network:        &testdata.Partition1ExistingPrivateNetwork,
					MachineNetwork: &metal.MachineNetwork{},
					IPs:            []metal.IP{},
					Auto:           true,
					IsPrivate:      true,
				},
				testdata.Partition1InternetNetwork.ID: &AllocationNetwork{
					Network:        &testdata.Partition1InternetNetwork,
					MachineNetwork: &metal.MachineNetwork{},
					IPs:            []metal.IP{testdata.Partition1InternetIP},
					Auto:           false,
					IsPrivate:      false,
				},
			},
		},
		{
			name: "private of other network given",
			allocationSpec: &AllocationSpec{
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     testdata.Partition1ExistingPrivateNetwork.ID,
						AutoAcquireIP: util.BoolProto(true),
					},
				},
				ProjectID: "another-project",
			},
			partition:              &testdata.Partition1,
			partitionSuperNetworks: partitionSuperNetworks,
			wantErr:                true,
			errRegex:               "the given private network does not belong to the project, which is not allowed",
		},
		{
			name: "try to assign machine to private network of other partition",
			allocationSpec: &AllocationSpec{
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     testdata.Partition2ExistingPrivateNetwork.ID,
						AutoAcquireIP: util.BoolProto(true),
					},
				},
				ProjectID: testdata.Partition2ExistingPrivateNetwork.ProjectID,
			},
			partition:              &testdata.Partition1,
			partitionSuperNetworks: partitionSuperNetworks,
			wantErr:                true,
			errRegex:               "the private network must be in the partition where the machine is going to be placed",
		},
		{
			name: "try to assign machine to super network",
			allocationSpec: &AllocationSpec{
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     testdata.Partition1PrivateSuperNetwork.ID,
						AutoAcquireIP: util.BoolProto(true),
					},
				},
			},
			partition:              &testdata.Partition1,
			partitionSuperNetworks: partitionSuperNetworks,
			wantErr:                true,
			errRegex:               "private super networks are not allowed to be set explicitly",
		},
		{
			name: "try to assign machine to underlay network",
			allocationSpec: &AllocationSpec{
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     testdata.Partition1UnderlayNetwork.ID,
						AutoAcquireIP: util.BoolProto(true),
					},
				},
			},
			partition:              &testdata.Partition1,
			partitionSuperNetworks: partitionSuperNetworks,
			wantErr:                true,
			errRegex:               "underlay networks are not allowed to be set explicitly",
		},
		{
			name: "try to add machine to multiple private networks",
			allocationSpec: &AllocationSpec{
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID: testdata.Partition1ExistingPrivateNetwork.ID,
					},
					{
						NetworkID: testdata.Partition2ExistingPrivateNetwork.ID,
					},
				},
			},
			partition:              &testdata.Partition1,
			partitionSuperNetworks: partitionSuperNetworks,
			wantErr:                true,
			errRegex:               "multiple private networks provided, which is not allowed",
		},
		{
			name: "try to add the same network a couple of times",
			allocationSpec: &AllocationSpec{
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID: testdata.Partition1InternetNetwork.ID,
					},
					{
						NetworkID: testdata.Partition1InternetNetwork.ID,
					},
				},
			},
			partition:              &testdata.Partition1,
			partitionSuperNetworks: partitionSuperNetworks,
			wantErr:                true,
			errRegex:               "given network ids are not unique",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// init tests
			ds, mock := datastore.InitMockDB()
			for _, testMock := range tt.mocks {
				mock.On(testMock.term).Return(testMock.response, testMock.err)
			}
			testdata.InitMockDBData(mock)

			// run
			got, err := GatherNetworksFromSpec(ds, tt.allocationSpec, tt.partition, tt.partitionSuperNetworks)

			// verify
			if err != nil {
				if !tt.wantErr {
					t.Errorf("gatherNetworksFromSpec() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.errRegex != "" {
					require.Regexp(t, tt.errRegex, err)
				}
				return
			}

			require.Len(t, got, len(tt.want), "number of gathered networks is incorrect")
			for wantNetworkID, wantNetwork := range tt.want {
				require.Contains(t, got, wantNetworkID)
				gotNetwork := got[wantNetworkID]
				require.Equal(t, wantNetwork.IsPrivate, gotNetwork.IsPrivate)
				require.Equal(t, wantNetwork.Auto, gotNetwork.Auto)

				var gotIPs []string
				for _, gotIP := range gotNetwork.IPs {
					gotIPs = append(gotIPs, gotIP.IPAddress)
				}

				for _, wantIP := range wantNetwork.IPs {
					require.Contains(t, gotIPs, wantIP.IPAddress)
				}
			}
		})
	}
}