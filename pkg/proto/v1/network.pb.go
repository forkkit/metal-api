// Code generated by protoc-gen-go. DO NOT EDIT.
// source: v1/network.proto

package v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Network struct {
	Common               *Common               `protobuf:"bytes,1,opt,name=common,proto3" json:"common,omitempty"`
	PartitionID          *wrappers.StringValue `protobuf:"bytes,2,opt,name=partitionID,proto3" json:"partitionID,omitempty"`
	ProjectID            *wrappers.StringValue `protobuf:"bytes,3,opt,name=projectID,proto3" json:"projectID,omitempty"`
	Labels               map[string]string     `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Network) Reset()         { *m = Network{} }
func (m *Network) String() string { return proto.CompactTextString(m) }
func (*Network) ProtoMessage()    {}
func (*Network) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ef602c4c85062d, []int{0}
}

func (m *Network) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Network.Unmarshal(m, b)
}
func (m *Network) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Network.Marshal(b, m, deterministic)
}
func (m *Network) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Network.Merge(m, src)
}
func (m *Network) XXX_Size() int {
	return xxx_messageInfo_Network.Size(m)
}
func (m *Network) XXX_DiscardUnknown() {
	xxx_messageInfo_Network.DiscardUnknown(m)
}

var xxx_messageInfo_Network proto.InternalMessageInfo

func (m *Network) GetCommon() *Common {
	if m != nil {
		return m.Common
	}
	return nil
}

func (m *Network) GetPartitionID() *wrappers.StringValue {
	if m != nil {
		return m.PartitionID
	}
	return nil
}

func (m *Network) GetProjectID() *wrappers.StringValue {
	if m != nil {
		return m.ProjectID
	}
	return nil
}

func (m *Network) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

// a network which contains prefixes from which IP addresses can be allocated
type NetworkImmutable struct {
	Prefixes             []string              `protobuf:"bytes,1,rep,name=prefixes,proto3" json:"prefixes,omitempty"`
	DestinationPrefixes  []string              `protobuf:"bytes,2,rep,name=destinationPrefixes,proto3" json:"destinationPrefixes,omitempty"`
	Nat                  bool                  `protobuf:"varint,3,opt,name=nat,proto3" json:"nat,omitempty"`
	PrivateSuper         bool                  `protobuf:"varint,4,opt,name=privateSuper,proto3" json:"privateSuper,omitempty"`
	Underlay             bool                  `protobuf:"varint,5,opt,name=underlay,proto3" json:"underlay,omitempty"`
	Vrf                  *wrappers.UInt64Value `protobuf:"bytes,6,opt,name=vrf,proto3" json:"vrf,omitempty"`
	VrfShared            *wrappers.BoolValue   `protobuf:"bytes,7,opt,name=vrfShared,proto3" json:"vrfShared,omitempty"`
	ParentNetworkID      *wrappers.StringValue `protobuf:"bytes,8,opt,name=parentNetworkID,proto3" json:"parentNetworkID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *NetworkImmutable) Reset()         { *m = NetworkImmutable{} }
func (m *NetworkImmutable) String() string { return proto.CompactTextString(m) }
func (*NetworkImmutable) ProtoMessage()    {}
func (*NetworkImmutable) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ef602c4c85062d, []int{1}
}

func (m *NetworkImmutable) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkImmutable.Unmarshal(m, b)
}
func (m *NetworkImmutable) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkImmutable.Marshal(b, m, deterministic)
}
func (m *NetworkImmutable) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkImmutable.Merge(m, src)
}
func (m *NetworkImmutable) XXX_Size() int {
	return xxx_messageInfo_NetworkImmutable.Size(m)
}
func (m *NetworkImmutable) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkImmutable.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkImmutable proto.InternalMessageInfo

func (m *NetworkImmutable) GetPrefixes() []string {
	if m != nil {
		return m.Prefixes
	}
	return nil
}

func (m *NetworkImmutable) GetDestinationPrefixes() []string {
	if m != nil {
		return m.DestinationPrefixes
	}
	return nil
}

func (m *NetworkImmutable) GetNat() bool {
	if m != nil {
		return m.Nat
	}
	return false
}

func (m *NetworkImmutable) GetPrivateSuper() bool {
	if m != nil {
		return m.PrivateSuper
	}
	return false
}

func (m *NetworkImmutable) GetUnderlay() bool {
	if m != nil {
		return m.Underlay
	}
	return false
}

func (m *NetworkImmutable) GetVrf() *wrappers.UInt64Value {
	if m != nil {
		return m.Vrf
	}
	return nil
}

func (m *NetworkImmutable) GetVrfShared() *wrappers.BoolValue {
	if m != nil {
		return m.VrfShared
	}
	return nil
}

func (m *NetworkImmutable) GetParentNetworkID() *wrappers.StringValue {
	if m != nil {
		return m.ParentNetworkID
	}
	return nil
}

type NetworkUsage struct {
	AvailableIPs         uint64   `protobuf:"varint,1,opt,name=availableIPs,proto3" json:"availableIPs,omitempty"`
	UsedIPs              uint64   `protobuf:"varint,2,opt,name=usedIPs,proto3" json:"usedIPs,omitempty"`
	AvailablePrefixes    uint64   `protobuf:"varint,3,opt,name=availablePrefixes,proto3" json:"availablePrefixes,omitempty"`
	UsedPrefixes         uint64   `protobuf:"varint,4,opt,name=usedPrefixes,proto3" json:"usedPrefixes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkUsage) Reset()         { *m = NetworkUsage{} }
func (m *NetworkUsage) String() string { return proto.CompactTextString(m) }
func (*NetworkUsage) ProtoMessage()    {}
func (*NetworkUsage) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ef602c4c85062d, []int{2}
}

func (m *NetworkUsage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkUsage.Unmarshal(m, b)
}
func (m *NetworkUsage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkUsage.Marshal(b, m, deterministic)
}
func (m *NetworkUsage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkUsage.Merge(m, src)
}
func (m *NetworkUsage) XXX_Size() int {
	return xxx_messageInfo_NetworkUsage.Size(m)
}
func (m *NetworkUsage) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkUsage.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkUsage proto.InternalMessageInfo

func (m *NetworkUsage) GetAvailableIPs() uint64 {
	if m != nil {
		return m.AvailableIPs
	}
	return 0
}

func (m *NetworkUsage) GetUsedIPs() uint64 {
	if m != nil {
		return m.UsedIPs
	}
	return 0
}

func (m *NetworkUsage) GetAvailablePrefixes() uint64 {
	if m != nil {
		return m.AvailablePrefixes
	}
	return 0
}

func (m *NetworkUsage) GetUsedPrefixes() uint64 {
	if m != nil {
		return m.UsedPrefixes
	}
	return 0
}

// NetworkSearchQuery can be used to search networks.
type NetworkSearchQuery struct {
	ID                   *wrappers.StringValue   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name                 *wrappers.StringValue   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	PartitionID          *wrappers.StringValue   `protobuf:"bytes,3,opt,name=partitionID,proto3" json:"partitionID,omitempty"`
	ProjectID            *wrappers.StringValue   `protobuf:"bytes,4,opt,name=projectID,proto3" json:"projectID,omitempty"`
	Prefixes             []*wrappers.StringValue `protobuf:"bytes,5,rep,name=prefixes,proto3" json:"prefixes,omitempty"`
	DestinationPrefixes  []*wrappers.StringValue `protobuf:"bytes,6,rep,name=destinationPrefixes,proto3" json:"destinationPrefixes,omitempty"`
	Nat                  *wrappers.BoolValue     `protobuf:"bytes,7,opt,name=nat,proto3" json:"nat,omitempty"`
	PrivateSuper         *wrappers.BoolValue     `protobuf:"bytes,8,opt,name=privateSuper,proto3" json:"privateSuper,omitempty"`
	Underlay             *wrappers.BoolValue     `protobuf:"bytes,9,opt,name=underlay,proto3" json:"underlay,omitempty"`
	Vrf                  *wrappers.UInt64Value   `protobuf:"bytes,10,opt,name=vrf,proto3" json:"vrf,omitempty"`
	VrfShared            *wrappers.BoolValue     `protobuf:"bytes,11,opt,name=vrfShared,proto3" json:"vrfShared,omitempty"`
	ParentNetworkID      *wrappers.StringValue   `protobuf:"bytes,12,opt,name=parentNetworkID,proto3" json:"parentNetworkID,omitempty"`
	Labels               map[string]string       `protobuf:"bytes,13,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *NetworkSearchQuery) Reset()         { *m = NetworkSearchQuery{} }
func (m *NetworkSearchQuery) String() string { return proto.CompactTextString(m) }
func (*NetworkSearchQuery) ProtoMessage()    {}
func (*NetworkSearchQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ef602c4c85062d, []int{3}
}

func (m *NetworkSearchQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkSearchQuery.Unmarshal(m, b)
}
func (m *NetworkSearchQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkSearchQuery.Marshal(b, m, deterministic)
}
func (m *NetworkSearchQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkSearchQuery.Merge(m, src)
}
func (m *NetworkSearchQuery) XXX_Size() int {
	return xxx_messageInfo_NetworkSearchQuery.Size(m)
}
func (m *NetworkSearchQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkSearchQuery.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkSearchQuery proto.InternalMessageInfo

func (m *NetworkSearchQuery) GetID() *wrappers.StringValue {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *NetworkSearchQuery) GetName() *wrappers.StringValue {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *NetworkSearchQuery) GetPartitionID() *wrappers.StringValue {
	if m != nil {
		return m.PartitionID
	}
	return nil
}

func (m *NetworkSearchQuery) GetProjectID() *wrappers.StringValue {
	if m != nil {
		return m.ProjectID
	}
	return nil
}

func (m *NetworkSearchQuery) GetPrefixes() []*wrappers.StringValue {
	if m != nil {
		return m.Prefixes
	}
	return nil
}

func (m *NetworkSearchQuery) GetDestinationPrefixes() []*wrappers.StringValue {
	if m != nil {
		return m.DestinationPrefixes
	}
	return nil
}

func (m *NetworkSearchQuery) GetNat() *wrappers.BoolValue {
	if m != nil {
		return m.Nat
	}
	return nil
}

func (m *NetworkSearchQuery) GetPrivateSuper() *wrappers.BoolValue {
	if m != nil {
		return m.PrivateSuper
	}
	return nil
}

func (m *NetworkSearchQuery) GetUnderlay() *wrappers.BoolValue {
	if m != nil {
		return m.Underlay
	}
	return nil
}

func (m *NetworkSearchQuery) GetVrf() *wrappers.UInt64Value {
	if m != nil {
		return m.Vrf
	}
	return nil
}

func (m *NetworkSearchQuery) GetVrfShared() *wrappers.BoolValue {
	if m != nil {
		return m.VrfShared
	}
	return nil
}

func (m *NetworkSearchQuery) GetParentNetworkID() *wrappers.StringValue {
	if m != nil {
		return m.ParentNetworkID
	}
	return nil
}

func (m *NetworkSearchQuery) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

type NetworkCreateRequest struct {
	Network              *Network          `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	NetworkImmutable     *NetworkImmutable `protobuf:"bytes,2,opt,name=networkImmutable,proto3" json:"networkImmutable,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *NetworkCreateRequest) Reset()         { *m = NetworkCreateRequest{} }
func (m *NetworkCreateRequest) String() string { return proto.CompactTextString(m) }
func (*NetworkCreateRequest) ProtoMessage()    {}
func (*NetworkCreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ef602c4c85062d, []int{4}
}

func (m *NetworkCreateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkCreateRequest.Unmarshal(m, b)
}
func (m *NetworkCreateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkCreateRequest.Marshal(b, m, deterministic)
}
func (m *NetworkCreateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkCreateRequest.Merge(m, src)
}
func (m *NetworkCreateRequest) XXX_Size() int {
	return xxx_messageInfo_NetworkCreateRequest.Size(m)
}
func (m *NetworkCreateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkCreateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkCreateRequest proto.InternalMessageInfo

func (m *NetworkCreateRequest) GetNetwork() *Network {
	if m != nil {
		return m.Network
	}
	return nil
}

func (m *NetworkCreateRequest) GetNetworkImmutable() *NetworkImmutable {
	if m != nil {
		return m.NetworkImmutable
	}
	return nil
}

type NetworkUpdateRequest struct {
	Common               *Common  `protobuf:"bytes,1,opt,name=common,proto3" json:"common,omitempty"`
	Prefixes             []string `protobuf:"bytes,2,rep,name=prefixes,proto3" json:"prefixes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkUpdateRequest) Reset()         { *m = NetworkUpdateRequest{} }
func (m *NetworkUpdateRequest) String() string { return proto.CompactTextString(m) }
func (*NetworkUpdateRequest) ProtoMessage()    {}
func (*NetworkUpdateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ef602c4c85062d, []int{5}
}

func (m *NetworkUpdateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkUpdateRequest.Unmarshal(m, b)
}
func (m *NetworkUpdateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkUpdateRequest.Marshal(b, m, deterministic)
}
func (m *NetworkUpdateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkUpdateRequest.Merge(m, src)
}
func (m *NetworkUpdateRequest) XXX_Size() int {
	return xxx_messageInfo_NetworkUpdateRequest.Size(m)
}
func (m *NetworkUpdateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkUpdateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkUpdateRequest proto.InternalMessageInfo

func (m *NetworkUpdateRequest) GetCommon() *Common {
	if m != nil {
		return m.Common
	}
	return nil
}

func (m *NetworkUpdateRequest) GetPrefixes() []string {
	if m != nil {
		return m.Prefixes
	}
	return nil
}

type NetworkFindRequest struct {
	NetworkSearchQuery   *NetworkSearchQuery `protobuf:"bytes,1,opt,name=networkSearchQuery,proto3" json:"networkSearchQuery,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *NetworkFindRequest) Reset()         { *m = NetworkFindRequest{} }
func (m *NetworkFindRequest) String() string { return proto.CompactTextString(m) }
func (*NetworkFindRequest) ProtoMessage()    {}
func (*NetworkFindRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ef602c4c85062d, []int{6}
}

func (m *NetworkFindRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkFindRequest.Unmarshal(m, b)
}
func (m *NetworkFindRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkFindRequest.Marshal(b, m, deterministic)
}
func (m *NetworkFindRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkFindRequest.Merge(m, src)
}
func (m *NetworkFindRequest) XXX_Size() int {
	return xxx_messageInfo_NetworkFindRequest.Size(m)
}
func (m *NetworkFindRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkFindRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkFindRequest proto.InternalMessageInfo

func (m *NetworkFindRequest) GetNetworkSearchQuery() *NetworkSearchQuery {
	if m != nil {
		return m.NetworkSearchQuery
	}
	return nil
}

type NetworkAllocateRequest struct {
	Network              *Network `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkAllocateRequest) Reset()         { *m = NetworkAllocateRequest{} }
func (m *NetworkAllocateRequest) String() string { return proto.CompactTextString(m) }
func (*NetworkAllocateRequest) ProtoMessage()    {}
func (*NetworkAllocateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ef602c4c85062d, []int{7}
}

func (m *NetworkAllocateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkAllocateRequest.Unmarshal(m, b)
}
func (m *NetworkAllocateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkAllocateRequest.Marshal(b, m, deterministic)
}
func (m *NetworkAllocateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkAllocateRequest.Merge(m, src)
}
func (m *NetworkAllocateRequest) XXX_Size() int {
	return xxx_messageInfo_NetworkAllocateRequest.Size(m)
}
func (m *NetworkAllocateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkAllocateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkAllocateRequest proto.InternalMessageInfo

func (m *NetworkAllocateRequest) GetNetwork() *Network {
	if m != nil {
		return m.Network
	}
	return nil
}

type NetworkResponse struct {
	Network              *Network          `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	NetworkImmutable     *NetworkImmutable `protobuf:"bytes,2,opt,name=networkImmutable,proto3" json:"networkImmutable,omitempty"`
	Usage                *NetworkUsage     `protobuf:"bytes,3,opt,name=usage,proto3" json:"usage,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *NetworkResponse) Reset()         { *m = NetworkResponse{} }
func (m *NetworkResponse) String() string { return proto.CompactTextString(m) }
func (*NetworkResponse) ProtoMessage()    {}
func (*NetworkResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ef602c4c85062d, []int{8}
}

func (m *NetworkResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkResponse.Unmarshal(m, b)
}
func (m *NetworkResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkResponse.Marshal(b, m, deterministic)
}
func (m *NetworkResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkResponse.Merge(m, src)
}
func (m *NetworkResponse) XXX_Size() int {
	return xxx_messageInfo_NetworkResponse.Size(m)
}
func (m *NetworkResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkResponse proto.InternalMessageInfo

func (m *NetworkResponse) GetNetwork() *Network {
	if m != nil {
		return m.Network
	}
	return nil
}

func (m *NetworkResponse) GetNetworkImmutable() *NetworkImmutable {
	if m != nil {
		return m.NetworkImmutable
	}
	return nil
}

func (m *NetworkResponse) GetUsage() *NetworkUsage {
	if m != nil {
		return m.Usage
	}
	return nil
}

type NetworkListRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkListRequest) Reset()         { *m = NetworkListRequest{} }
func (m *NetworkListRequest) String() string { return proto.CompactTextString(m) }
func (*NetworkListRequest) ProtoMessage()    {}
func (*NetworkListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ef602c4c85062d, []int{9}
}

func (m *NetworkListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkListRequest.Unmarshal(m, b)
}
func (m *NetworkListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkListRequest.Marshal(b, m, deterministic)
}
func (m *NetworkListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkListRequest.Merge(m, src)
}
func (m *NetworkListRequest) XXX_Size() int {
	return xxx_messageInfo_NetworkListRequest.Size(m)
}
func (m *NetworkListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkListRequest proto.InternalMessageInfo

type NetworkListResponse struct {
	Networks             []*NetworkResponse `protobuf:"bytes,1,rep,name=networks,proto3" json:"networks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *NetworkListResponse) Reset()         { *m = NetworkListResponse{} }
func (m *NetworkListResponse) String() string { return proto.CompactTextString(m) }
func (*NetworkListResponse) ProtoMessage()    {}
func (*NetworkListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ef602c4c85062d, []int{10}
}

func (m *NetworkListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkListResponse.Unmarshal(m, b)
}
func (m *NetworkListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkListResponse.Marshal(b, m, deterministic)
}
func (m *NetworkListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkListResponse.Merge(m, src)
}
func (m *NetworkListResponse) XXX_Size() int {
	return xxx_messageInfo_NetworkListResponse.Size(m)
}
func (m *NetworkListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkListResponse proto.InternalMessageInfo

func (m *NetworkListResponse) GetNetworks() []*NetworkResponse {
	if m != nil {
		return m.Networks
	}
	return nil
}

func init() {
	proto.RegisterType((*Network)(nil), "v1.Network")
	proto.RegisterMapType((map[string]string)(nil), "v1.Network.LabelsEntry")
	proto.RegisterType((*NetworkImmutable)(nil), "v1.NetworkImmutable")
	proto.RegisterType((*NetworkUsage)(nil), "v1.NetworkUsage")
	proto.RegisterType((*NetworkSearchQuery)(nil), "v1.NetworkSearchQuery")
	proto.RegisterMapType((map[string]string)(nil), "v1.NetworkSearchQuery.LabelsEntry")
	proto.RegisterType((*NetworkCreateRequest)(nil), "v1.NetworkCreateRequest")
	proto.RegisterType((*NetworkUpdateRequest)(nil), "v1.NetworkUpdateRequest")
	proto.RegisterType((*NetworkFindRequest)(nil), "v1.NetworkFindRequest")
	proto.RegisterType((*NetworkAllocateRequest)(nil), "v1.NetworkAllocateRequest")
	proto.RegisterType((*NetworkResponse)(nil), "v1.NetworkResponse")
	proto.RegisterType((*NetworkListRequest)(nil), "v1.NetworkListRequest")
	proto.RegisterType((*NetworkListResponse)(nil), "v1.NetworkListResponse")
}

func init() { proto.RegisterFile("v1/network.proto", fileDescriptor_77ef602c4c85062d) }

var fileDescriptor_77ef602c4c85062d = []byte{
	// 867 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x56, 0xdd, 0x6e, 0xe3, 0x44,
	0x14, 0x56, 0x1c, 0xe7, 0xef, 0xa4, 0xd0, 0x30, 0x8d, 0x16, 0xcb, 0x42, 0xa8, 0xb2, 0x04, 0xea,
	0x45, 0xd7, 0xde, 0x04, 0xb4, 0x5b, 0x8a, 0xb4, 0xc0, 0x6e, 0x88, 0x14, 0x69, 0xb5, 0x5a, 0x1c,
	0x75, 0x2f, 0x10, 0x37, 0x93, 0x64, 0x92, 0x9a, 0x38, 0xb6, 0x19, 0x8f, 0x5d, 0x72, 0xc5, 0x53,
	0x70, 0xcb, 0x05, 0x6f, 0xc0, 0x63, 0xf0, 0x1e, 0x3c, 0x08, 0x9a, 0xf1, 0xd8, 0x1d, 0x27, 0x69,
	0x9b, 0xb6, 0x88, 0xbb, 0x99, 0x73, 0xbe, 0x6f, 0xce, 0x99, 0x33, 0xe7, 0x3b, 0x36, 0x74, 0xd2,
	0x9e, 0x13, 0x10, 0x76, 0x15, 0xd2, 0xa5, 0x1d, 0xd1, 0x90, 0x85, 0x48, 0x4b, 0x7b, 0xe6, 0x61,
	0xda, 0x73, 0xa6, 0xe1, 0x6a, 0x15, 0x06, 0x99, 0xd1, 0xfc, 0x74, 0x11, 0x86, 0x0b, 0x9f, 0x38,
	0x62, 0x37, 0x49, 0xe6, 0xce, 0x15, 0xc5, 0x51, 0x44, 0x68, 0x9c, 0xf9, 0xad, 0xdf, 0x35, 0x68,
	0xbc, 0xcd, 0x8e, 0x41, 0x16, 0xd4, 0x33, 0xae, 0x51, 0x39, 0xae, 0x9c, 0xb4, 0xfb, 0x60, 0xa7,
	0x3d, 0xfb, 0xb5, 0xb0, 0xb8, 0xd2, 0x83, 0x5e, 0x42, 0x3b, 0xc2, 0x94, 0x79, 0xcc, 0x0b, 0x83,
	0xd1, 0xc0, 0xd0, 0x04, 0xf0, 0x13, 0x3b, 0x8b, 0x62, 0xe7, 0x51, 0xec, 0x31, 0xa3, 0x5e, 0xb0,
	0x78, 0x8f, 0xfd, 0x84, 0xb8, 0x2a, 0x01, 0x9d, 0x43, 0x2b, 0xa2, 0xe1, 0xcf, 0x64, 0xca, 0x46,
	0x03, 0xa3, 0xba, 0x07, 0xfb, 0x1a, 0x8e, 0x1c, 0xa8, 0xfb, 0x78, 0x42, 0xfc, 0xd8, 0xd0, 0x8f,
	0xab, 0x27, 0xed, 0xfe, 0xc7, 0x3c, 0x3f, 0x99, 0xbc, 0xfd, 0x46, 0x78, 0xbe, 0x0f, 0x18, 0x5d,
	0xbb, 0x12, 0x66, 0x7e, 0x05, 0x6d, 0xc5, 0x8c, 0x3a, 0x50, 0x5d, 0x92, 0xb5, 0xb8, 0x5c, 0xcb,
	0xe5, 0x4b, 0xd4, 0x85, 0x5a, 0xca, 0xa3, 0x88, 0x7b, 0xb4, 0xdc, 0x6c, 0x73, 0xae, 0x9d, 0x55,
	0xac, 0x7f, 0x34, 0xe8, 0xc8, 0xa3, 0x47, 0xab, 0x55, 0xc2, 0xf0, 0xc4, 0x27, 0xc8, 0x84, 0x66,
	0x44, 0xc9, 0xdc, 0xfb, 0x95, 0xc4, 0x46, 0xe5, 0xb8, 0x7a, 0xd2, 0x72, 0x8b, 0x3d, 0x7a, 0x06,
	0x47, 0x33, 0x12, 0x33, 0x2f, 0xc0, 0xfc, 0xa6, 0xef, 0x72, 0x98, 0x26, 0x60, 0xbb, 0x5c, 0x3c,
	0x9d, 0x00, 0x33, 0x51, 0x84, 0xa6, 0xcb, 0x97, 0xc8, 0x82, 0x83, 0x88, 0x7a, 0x29, 0x66, 0x64,
	0x9c, 0x44, 0x84, 0x1a, 0xba, 0x70, 0x95, 0x6c, 0x3c, 0x87, 0x24, 0x98, 0x11, 0xea, 0xe3, 0xb5,
	0x51, 0x13, 0xfe, 0x62, 0x8f, 0x6c, 0xa8, 0xa6, 0x74, 0x6e, 0xd4, 0x6f, 0x28, 0xeb, 0xc5, 0x28,
	0x60, 0xcf, 0xbf, 0xcc, 0xca, 0xca, 0x81, 0xe8, 0x0c, 0x5a, 0x29, 0x9d, 0x8f, 0x2f, 0x31, 0x25,
	0x33, 0xa3, 0x21, 0x58, 0xe6, 0x16, 0xeb, 0x55, 0x18, 0xfa, 0xf2, 0x29, 0x0a, 0x30, 0x1a, 0xc2,
	0x61, 0x84, 0x29, 0x09, 0x58, 0x5e, 0xa3, 0x81, 0xd1, 0xdc, 0xe3, 0x31, 0x37, 0x49, 0xd6, 0x1f,
	0x15, 0x38, 0x90, 0xbb, 0x8b, 0x18, 0x2f, 0x08, 0x2f, 0x01, 0x4e, 0xb1, 0xe7, 0xf3, 0x7a, 0x8f,
	0xde, 0xc5, 0xe2, 0xb1, 0x74, 0xb7, 0x64, 0x43, 0x06, 0x34, 0x92, 0x98, 0xcc, 0xb8, 0x5b, 0x13,
	0xee, 0x7c, 0x8b, 0x4e, 0xe1, 0xa3, 0x02, 0x59, 0x3c, 0x41, 0x55, 0x60, 0xb6, 0x1d, 0x3c, 0x16,
	0x27, 0x16, 0x40, 0x3d, 0x8b, 0xa5, 0xda, 0xac, 0xbf, 0xeb, 0x80, 0x64, 0x82, 0x63, 0x82, 0xe9,
	0xf4, 0xf2, 0x87, 0x84, 0xd0, 0x35, 0x3a, 0x05, 0x6d, 0x34, 0x90, 0x32, 0xb9, 0xfd, 0xca, 0xda,
	0x68, 0x80, 0x9e, 0x81, 0x1e, 0xe0, 0x15, 0xd9, 0x4b, 0x2d, 0x02, 0xb9, 0x29, 0xb3, 0xea, 0xa3,
	0x64, 0xa6, 0xdf, 0x4f, 0x66, 0x67, 0x4a, 0x97, 0xd7, 0x84, 0xd0, 0x6e, 0xa7, 0x5e, 0x6b, 0xe0,
	0xed, 0x6e, 0x0d, 0xd4, 0xf7, 0x38, 0x64, 0xa7, 0x42, 0x4e, 0x33, 0x85, 0xdc, 0xdd, 0x99, 0x42,
	0x3d, 0x2f, 0x37, 0xd4, 0xd3, 0xbc, 0x93, 0x56, 0x56, 0xd6, 0x73, 0x45, 0x59, 0xad, 0x3b, 0xb9,
	0x5b, 0xaa, 0x83, 0x07, 0xa9, 0xae, 0xfd, 0x48, 0xd5, 0x1d, 0x3c, 0x40, 0x75, 0xe8, 0xbc, 0x18,
	0xa4, 0x1f, 0x88, 0xa7, 0xb1, 0x94, 0x41, 0xaa, 0x74, 0xf9, 0x7f, 0x3d, 0x53, 0x7f, 0x83, 0xae,
	0x0c, 0xf2, 0x9a, 0x12, 0xcc, 0x88, 0x4b, 0x7e, 0x49, 0x48, 0xcc, 0xd0, 0x67, 0xd0, 0x90, 0x5f,
	0x32, 0xa9, 0xa8, 0xb6, 0x92, 0x8f, 0x9b, 0xfb, 0xd0, 0xb7, 0xd0, 0x09, 0x36, 0x26, 0xb2, 0x54,
	0x54, 0x57, 0xc1, 0x17, 0x3e, 0x77, 0x0b, 0x6d, 0xbd, 0x2f, 0x12, 0xb8, 0x88, 0x66, 0x4a, 0x02,
	0xfb, 0x7c, 0xf8, 0xd4, 0xd9, 0xaf, 0x95, 0x67, 0xbf, 0xf5, 0x53, 0x31, 0x23, 0x86, 0x5e, 0x30,
	0xcb, 0x4f, 0x1d, 0x02, 0x0a, 0xb6, 0x6a, 0x2a, 0x23, 0x3c, 0xd9, 0x5d, 0x71, 0x77, 0x07, 0xc3,
	0xfa, 0x06, 0x9e, 0x48, 0xe4, 0x77, 0xbe, 0x1f, 0x4e, 0xef, 0x5d, 0x38, 0xeb, 0xcf, 0x0a, 0x1c,
	0xe6, 0x46, 0x12, 0x47, 0x61, 0x10, 0x93, 0xff, 0xad, 0xe6, 0xe8, 0x73, 0xa8, 0x25, 0x7c, 0xb2,
	0xcb, 0x19, 0xd6, 0x51, 0x68, 0x62, 0xe2, 0xbb, 0x99, 0xdb, 0xea, 0x16, 0x35, 0x7c, 0xe3, 0xc5,
	0x4c, 0xde, 0xd0, 0x1a, 0xc2, 0x51, 0xc9, 0x2a, 0xb3, 0x77, 0xa0, 0x29, 0x03, 0x65, 0x1f, 0xe2,
	0x76, 0xff, 0x48, 0x4d, 0x5f, 0xc2, 0xdc, 0x02, 0xd4, 0xff, 0x4b, 0x83, 0x0f, 0x8b, 0x72, 0xd3,
	0xd4, 0x9b, 0x12, 0xf4, 0x02, 0xea, 0x59, 0x1b, 0x22, 0x43, 0xe1, 0x96, 0x3a, 0xd3, 0xdc, 0x75,
	0x2a, 0x27, 0x66, 0xed, 0x53, 0x22, 0x96, 0x3a, 0x6a, 0x37, 0xf1, 0x6b, 0x68, 0xe6, 0x2f, 0x88,
	0x4c, 0x05, 0xb0, 0xf1, 0xac, 0x37, 0x45, 0xd5, 0x79, 0x73, 0x21, 0xb5, 0x73, 0x94, 0x6e, 0x33,
	0xd5, 0x9f, 0xa1, 0x52, 0xad, 0x5e, 0x80, 0xce, 0xf7, 0x25, 0xa2, 0x52, 0xe2, 0x1b, 0x89, 0xaf,
	0x9c, 0x1f, 0x9f, 0x2e, 0x3c, 0x76, 0x99, 0x4c, 0xec, 0x69, 0xb8, 0x72, 0x56, 0x84, 0x61, 0xff,
	0x69, 0xcc, 0xf0, 0x74, 0x29, 0xd7, 0x38, 0xf2, 0x9c, 0x68, 0xb9, 0xc8, 0x7e, 0x2d, 0x9d, 0xb4,
	0x37, 0xa9, 0x8b, 0xd5, 0x17, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0x9c, 0xbc, 0x52, 0x89, 0x9b,
	0x0a, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// NetworkServiceClient is the client API for NetworkService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NetworkServiceClient interface {
	Create(ctx context.Context, in *NetworkCreateRequest, opts ...grpc.CallOption) (*NetworkResponse, error)
	Update(ctx context.Context, in *NetworkUpdateRequest, opts ...grpc.CallOption) (*NetworkResponse, error)
	//    rpc Delete(NetworkDeleteRequest) returns (NetworkResponse);
	//    rpc Get(NetworkGetRequest) returns (NetworkResponse);
	Allocate(ctx context.Context, in *NetworkAllocateRequest, opts ...grpc.CallOption) (*NetworkResponse, error)
	Find(ctx context.Context, in *NetworkFindRequest, opts ...grpc.CallOption) (*NetworkListResponse, error)
	List(ctx context.Context, in *NetworkListRequest, opts ...grpc.CallOption) (*NetworkListResponse, error)
}

type networkServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNetworkServiceClient(cc grpc.ClientConnInterface) NetworkServiceClient {
	return &networkServiceClient{cc}
}

func (c *networkServiceClient) Create(ctx context.Context, in *NetworkCreateRequest, opts ...grpc.CallOption) (*NetworkResponse, error) {
	out := new(NetworkResponse)
	err := c.cc.Invoke(ctx, "/v1.NetworkService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkServiceClient) Update(ctx context.Context, in *NetworkUpdateRequest, opts ...grpc.CallOption) (*NetworkResponse, error) {
	out := new(NetworkResponse)
	err := c.cc.Invoke(ctx, "/v1.NetworkService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkServiceClient) Allocate(ctx context.Context, in *NetworkAllocateRequest, opts ...grpc.CallOption) (*NetworkResponse, error) {
	out := new(NetworkResponse)
	err := c.cc.Invoke(ctx, "/v1.NetworkService/Allocate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkServiceClient) Find(ctx context.Context, in *NetworkFindRequest, opts ...grpc.CallOption) (*NetworkListResponse, error) {
	out := new(NetworkListResponse)
	err := c.cc.Invoke(ctx, "/v1.NetworkService/Find", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkServiceClient) List(ctx context.Context, in *NetworkListRequest, opts ...grpc.CallOption) (*NetworkListResponse, error) {
	out := new(NetworkListResponse)
	err := c.cc.Invoke(ctx, "/v1.NetworkService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkServiceServer is the server API for NetworkService service.
type NetworkServiceServer interface {
	Create(context.Context, *NetworkCreateRequest) (*NetworkResponse, error)
	Update(context.Context, *NetworkUpdateRequest) (*NetworkResponse, error)
	//    rpc Delete(NetworkDeleteRequest) returns (NetworkResponse);
	//    rpc Get(NetworkGetRequest) returns (NetworkResponse);
	Allocate(context.Context, *NetworkAllocateRequest) (*NetworkResponse, error)
	Find(context.Context, *NetworkFindRequest) (*NetworkListResponse, error)
	List(context.Context, *NetworkListRequest) (*NetworkListResponse, error)
}

// UnimplementedNetworkServiceServer can be embedded to have forward compatible implementations.
type UnimplementedNetworkServiceServer struct {
}

func (*UnimplementedNetworkServiceServer) Create(ctx context.Context, req *NetworkCreateRequest) (*NetworkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedNetworkServiceServer) Update(ctx context.Context, req *NetworkUpdateRequest) (*NetworkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedNetworkServiceServer) Allocate(ctx context.Context, req *NetworkAllocateRequest) (*NetworkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Allocate not implemented")
}
func (*UnimplementedNetworkServiceServer) Find(ctx context.Context, req *NetworkFindRequest) (*NetworkListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (*UnimplementedNetworkServiceServer) List(ctx context.Context, req *NetworkListRequest) (*NetworkListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

func RegisterNetworkServiceServer(s *grpc.Server, srv NetworkServiceServer) {
	s.RegisterService(&_NetworkService_serviceDesc, srv)
}

func _NetworkService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.NetworkService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceServer).Create(ctx, req.(*NetworkCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.NetworkService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceServer).Update(ctx, req.(*NetworkUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkService_Allocate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkAllocateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceServer).Allocate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.NetworkService/Allocate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceServer).Allocate(ctx, req.(*NetworkAllocateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkService_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkFindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.NetworkService/Find",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceServer).Find(ctx, req.(*NetworkFindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.NetworkService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceServer).List(ctx, req.(*NetworkListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _NetworkService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.NetworkService",
	HandlerType: (*NetworkServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _NetworkService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _NetworkService_Update_Handler,
		},
		{
			MethodName: "Allocate",
			Handler:    _NetworkService_Allocate_Handler,
		},
		{
			MethodName: "Find",
			Handler:    _NetworkService_Find_Handler,
		},
		{
			MethodName: "List",
			Handler:    _NetworkService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/network.proto",
}