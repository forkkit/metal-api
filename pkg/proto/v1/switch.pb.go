// Code generated by protoc-gen-go. DO NOT EDIT.
// source: v1/switch.proto

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

type BGPFilter struct {
	CIDRs                []string                `protobuf:"bytes,1,rep,name=CIDRs,proto3" json:"CIDRs,omitempty"`
	VNIs                 []*wrappers.StringValue `protobuf:"bytes,2,rep,name=VNIs,proto3" json:"VNIs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *BGPFilter) Reset()         { *m = BGPFilter{} }
func (m *BGPFilter) String() string { return proto.CompactTextString(m) }
func (*BGPFilter) ProtoMessage()    {}
func (*BGPFilter) Descriptor() ([]byte, []int) {
	return fileDescriptor_4559081f66c40988, []int{0}
}

func (m *BGPFilter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BGPFilter.Unmarshal(m, b)
}
func (m *BGPFilter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BGPFilter.Marshal(b, m, deterministic)
}
func (m *BGPFilter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BGPFilter.Merge(m, src)
}
func (m *BGPFilter) XXX_Size() int {
	return xxx_messageInfo_BGPFilter.Size(m)
}
func (m *BGPFilter) XXX_DiscardUnknown() {
	xxx_messageInfo_BGPFilter.DiscardUnknown(m)
}

var xxx_messageInfo_BGPFilter proto.InternalMessageInfo

func (m *BGPFilter) GetCIDRs() []string {
	if m != nil {
		return m.CIDRs
	}
	return nil
}

func (m *BGPFilter) GetVNIs() []*wrappers.StringValue {
	if m != nil {
		return m.VNIs
	}
	return nil
}

type SwitchNic struct {
	MacAddress           string                `protobuf:"bytes,1,opt,name=macAddress,proto3" json:"macAddress,omitempty"`
	Name                 string                `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Vrf                  *wrappers.StringValue `protobuf:"bytes,3,opt,name=vrf,proto3" json:"vrf,omitempty"`
	BGPFilter            *BGPFilter            `protobuf:"bytes,4,opt,name=BGPFilter,proto3" json:"BGPFilter,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *SwitchNic) Reset()         { *m = SwitchNic{} }
func (m *SwitchNic) String() string { return proto.CompactTextString(m) }
func (*SwitchNic) ProtoMessage()    {}
func (*SwitchNic) Descriptor() ([]byte, []int) {
	return fileDescriptor_4559081f66c40988, []int{1}
}

func (m *SwitchNic) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SwitchNic.Unmarshal(m, b)
}
func (m *SwitchNic) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SwitchNic.Marshal(b, m, deterministic)
}
func (m *SwitchNic) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchNic.Merge(m, src)
}
func (m *SwitchNic) XXX_Size() int {
	return xxx_messageInfo_SwitchNic.Size(m)
}
func (m *SwitchNic) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchNic.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchNic proto.InternalMessageInfo

func (m *SwitchNic) GetMacAddress() string {
	if m != nil {
		return m.MacAddress
	}
	return ""
}

func (m *SwitchNic) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SwitchNic) GetVrf() *wrappers.StringValue {
	if m != nil {
		return m.Vrf
	}
	return nil
}

func (m *SwitchNic) GetBGPFilter() *BGPFilter {
	if m != nil {
		return m.BGPFilter
	}
	return nil
}

type SwitchConnection struct {
	Nic                  *SwitchNic            `protobuf:"bytes,1,opt,name=nic,proto3" json:"nic,omitempty"`
	MachineID            *wrappers.StringValue `protobuf:"bytes,2,opt,name=machineID,proto3" json:"machineID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *SwitchConnection) Reset()         { *m = SwitchConnection{} }
func (m *SwitchConnection) String() string { return proto.CompactTextString(m) }
func (*SwitchConnection) ProtoMessage()    {}
func (*SwitchConnection) Descriptor() ([]byte, []int) {
	return fileDescriptor_4559081f66c40988, []int{2}
}

func (m *SwitchConnection) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SwitchConnection.Unmarshal(m, b)
}
func (m *SwitchConnection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SwitchConnection.Marshal(b, m, deterministic)
}
func (m *SwitchConnection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchConnection.Merge(m, src)
}
func (m *SwitchConnection) XXX_Size() int {
	return xxx_messageInfo_SwitchConnection.Size(m)
}
func (m *SwitchConnection) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchConnection.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchConnection proto.InternalMessageInfo

func (m *SwitchConnection) GetNic() *SwitchNic {
	if m != nil {
		return m.Nic
	}
	return nil
}

func (m *SwitchConnection) GetMachineID() *wrappers.StringValue {
	if m != nil {
		return m.MachineID
	}
	return nil
}

// A switch that can register at the api
type Switch struct {
	Common               *Common      `protobuf:"bytes,1,opt,name=common,proto3" json:"common,omitempty"`
	RackID               string       `protobuf:"bytes,2,opt,name=rackID,proto3" json:"rackID,omitempty"`
	Nics                 []*SwitchNic `protobuf:"bytes,3,rep,name=nics,proto3" json:"nics,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Switch) Reset()         { *m = Switch{} }
func (m *Switch) String() string { return proto.CompactTextString(m) }
func (*Switch) ProtoMessage()    {}
func (*Switch) Descriptor() ([]byte, []int) {
	return fileDescriptor_4559081f66c40988, []int{3}
}

func (m *Switch) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Switch.Unmarshal(m, b)
}
func (m *Switch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Switch.Marshal(b, m, deterministic)
}
func (m *Switch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Switch.Merge(m, src)
}
func (m *Switch) XXX_Size() int {
	return xxx_messageInfo_Switch.Size(m)
}
func (m *Switch) XXX_DiscardUnknown() {
	xxx_messageInfo_Switch.DiscardUnknown(m)
}

var xxx_messageInfo_Switch proto.InternalMessageInfo

func (m *Switch) GetCommon() *Common {
	if m != nil {
		return m.Common
	}
	return nil
}

func (m *Switch) GetRackID() string {
	if m != nil {
		return m.RackID
	}
	return ""
}

func (m *Switch) GetNics() []*SwitchNic {
	if m != nil {
		return m.Nics
	}
	return nil
}

type SwitchRegisterRequest struct {
	Switch               *Switch  `protobuf:"bytes,1,opt,name=switch,proto3" json:"switch,omitempty"`
	PartitionID          string   `protobuf:"bytes,2,opt,name=partitionID,proto3" json:"partitionID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SwitchRegisterRequest) Reset()         { *m = SwitchRegisterRequest{} }
func (m *SwitchRegisterRequest) String() string { return proto.CompactTextString(m) }
func (*SwitchRegisterRequest) ProtoMessage()    {}
func (*SwitchRegisterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4559081f66c40988, []int{4}
}

func (m *SwitchRegisterRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SwitchRegisterRequest.Unmarshal(m, b)
}
func (m *SwitchRegisterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SwitchRegisterRequest.Marshal(b, m, deterministic)
}
func (m *SwitchRegisterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchRegisterRequest.Merge(m, src)
}
func (m *SwitchRegisterRequest) XXX_Size() int {
	return xxx_messageInfo_SwitchRegisterRequest.Size(m)
}
func (m *SwitchRegisterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchRegisterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchRegisterRequest proto.InternalMessageInfo

func (m *SwitchRegisterRequest) GetSwitch() *Switch {
	if m != nil {
		return m.Switch
	}
	return nil
}

func (m *SwitchRegisterRequest) GetPartitionID() string {
	if m != nil {
		return m.PartitionID
	}
	return ""
}

type SwitchGetRequest struct {
	Identifiable         *Identifiable `protobuf:"bytes,1,opt,name=identifiable,proto3" json:"identifiable,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *SwitchGetRequest) Reset()         { *m = SwitchGetRequest{} }
func (m *SwitchGetRequest) String() string { return proto.CompactTextString(m) }
func (*SwitchGetRequest) ProtoMessage()    {}
func (*SwitchGetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4559081f66c40988, []int{5}
}

func (m *SwitchGetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SwitchGetRequest.Unmarshal(m, b)
}
func (m *SwitchGetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SwitchGetRequest.Marshal(b, m, deterministic)
}
func (m *SwitchGetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchGetRequest.Merge(m, src)
}
func (m *SwitchGetRequest) XXX_Size() int {
	return xxx_messageInfo_SwitchGetRequest.Size(m)
}
func (m *SwitchGetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchGetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchGetRequest proto.InternalMessageInfo

func (m *SwitchGetRequest) GetIdentifiable() *Identifiable {
	if m != nil {
		return m.Identifiable
	}
	return nil
}

type SwitchFindRequest struct {
	FindCriteria         *FindCriteria `protobuf:"bytes,1,opt,name=findCriteria,proto3" json:"findCriteria,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *SwitchFindRequest) Reset()         { *m = SwitchFindRequest{} }
func (m *SwitchFindRequest) String() string { return proto.CompactTextString(m) }
func (*SwitchFindRequest) ProtoMessage()    {}
func (*SwitchFindRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4559081f66c40988, []int{6}
}

func (m *SwitchFindRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SwitchFindRequest.Unmarshal(m, b)
}
func (m *SwitchFindRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SwitchFindRequest.Marshal(b, m, deterministic)
}
func (m *SwitchFindRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchFindRequest.Merge(m, src)
}
func (m *SwitchFindRequest) XXX_Size() int {
	return xxx_messageInfo_SwitchFindRequest.Size(m)
}
func (m *SwitchFindRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchFindRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchFindRequest proto.InternalMessageInfo

func (m *SwitchFindRequest) GetFindCriteria() *FindCriteria {
	if m != nil {
		return m.FindCriteria
	}
	return nil
}

type SwitchResponse struct {
	Switch               *Switch             `protobuf:"bytes,1,opt,name=switch,proto3" json:"switch,omitempty"`
	Connections          []*SwitchConnection `protobuf:"bytes,2,rep,name=connections,proto3" json:"connections,omitempty"`
	PartitionResponse    *PartitionResponse  `protobuf:"bytes,3,opt,name=partitionResponse,proto3" json:"partitionResponse,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *SwitchResponse) Reset()         { *m = SwitchResponse{} }
func (m *SwitchResponse) String() string { return proto.CompactTextString(m) }
func (*SwitchResponse) ProtoMessage()    {}
func (*SwitchResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4559081f66c40988, []int{7}
}

func (m *SwitchResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SwitchResponse.Unmarshal(m, b)
}
func (m *SwitchResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SwitchResponse.Marshal(b, m, deterministic)
}
func (m *SwitchResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchResponse.Merge(m, src)
}
func (m *SwitchResponse) XXX_Size() int {
	return xxx_messageInfo_SwitchResponse.Size(m)
}
func (m *SwitchResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchResponse proto.InternalMessageInfo

func (m *SwitchResponse) GetSwitch() *Switch {
	if m != nil {
		return m.Switch
	}
	return nil
}

func (m *SwitchResponse) GetConnections() []*SwitchConnection {
	if m != nil {
		return m.Connections
	}
	return nil
}

func (m *SwitchResponse) GetPartitionResponse() *PartitionResponse {
	if m != nil {
		return m.PartitionResponse
	}
	return nil
}

type SwitchListRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SwitchListRequest) Reset()         { *m = SwitchListRequest{} }
func (m *SwitchListRequest) String() string { return proto.CompactTextString(m) }
func (*SwitchListRequest) ProtoMessage()    {}
func (*SwitchListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4559081f66c40988, []int{8}
}

func (m *SwitchListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SwitchListRequest.Unmarshal(m, b)
}
func (m *SwitchListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SwitchListRequest.Marshal(b, m, deterministic)
}
func (m *SwitchListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchListRequest.Merge(m, src)
}
func (m *SwitchListRequest) XXX_Size() int {
	return xxx_messageInfo_SwitchListRequest.Size(m)
}
func (m *SwitchListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchListRequest proto.InternalMessageInfo

type SwitchListResponse struct {
	Switches             []*SwitchResponse `protobuf:"bytes,1,rep,name=switches,proto3" json:"switches,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SwitchListResponse) Reset()         { *m = SwitchListResponse{} }
func (m *SwitchListResponse) String() string { return proto.CompactTextString(m) }
func (*SwitchListResponse) ProtoMessage()    {}
func (*SwitchListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4559081f66c40988, []int{9}
}

func (m *SwitchListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SwitchListResponse.Unmarshal(m, b)
}
func (m *SwitchListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SwitchListResponse.Marshal(b, m, deterministic)
}
func (m *SwitchListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchListResponse.Merge(m, src)
}
func (m *SwitchListResponse) XXX_Size() int {
	return xxx_messageInfo_SwitchListResponse.Size(m)
}
func (m *SwitchListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchListResponse proto.InternalMessageInfo

func (m *SwitchListResponse) GetSwitches() []*SwitchResponse {
	if m != nil {
		return m.Switches
	}
	return nil
}

func init() {
	proto.RegisterType((*BGPFilter)(nil), "v1.BGPFilter")
	proto.RegisterType((*SwitchNic)(nil), "v1.SwitchNic")
	proto.RegisterType((*SwitchConnection)(nil), "v1.SwitchConnection")
	proto.RegisterType((*Switch)(nil), "v1.Switch")
	proto.RegisterType((*SwitchRegisterRequest)(nil), "v1.SwitchRegisterRequest")
	proto.RegisterType((*SwitchGetRequest)(nil), "v1.SwitchGetRequest")
	proto.RegisterType((*SwitchFindRequest)(nil), "v1.SwitchFindRequest")
	proto.RegisterType((*SwitchResponse)(nil), "v1.SwitchResponse")
	proto.RegisterType((*SwitchListRequest)(nil), "v1.SwitchListRequest")
	proto.RegisterType((*SwitchListResponse)(nil), "v1.SwitchListResponse")
}

func init() { proto.RegisterFile("v1/switch.proto", fileDescriptor_4559081f66c40988) }

var fileDescriptor_4559081f66c40988 = []byte{
	// 614 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xdd, 0x6e, 0xd3, 0x4c,
	0x10, 0x55, 0xea, 0x7c, 0xd1, 0x97, 0x09, 0x85, 0x76, 0x69, 0xaa, 0x10, 0xa1, 0x12, 0x7c, 0x55,
	0x09, 0xd5, 0x26, 0x85, 0x82, 0xc4, 0x1d, 0x4d, 0xd5, 0x12, 0x09, 0x55, 0xd5, 0x46, 0xf4, 0x02,
	0x89, 0x8b, 0xcd, 0x7a, 0xe3, 0xae, 0x12, 0xff, 0xe0, 0xdd, 0xb8, 0x2f, 0xc3, 0x4b, 0xf0, 0x54,
	0xbc, 0x06, 0xda, 0x1f, 0xdb, 0x9b, 0x52, 0x54, 0xee, 0xbc, 0x67, 0xce, 0x9c, 0x99, 0x33, 0x33,
	0x09, 0x3c, 0x29, 0xc7, 0xa1, 0xb8, 0xe5, 0x92, 0xde, 0x04, 0x79, 0x91, 0xc9, 0x0c, 0x6d, 0x95,
	0xe3, 0xa1, 0x02, 0x69, 0x96, 0x24, 0x59, 0x6a, 0xc0, 0x21, 0x2a, 0xc7, 0x61, 0x4e, 0x0a, 0xc9,
	0x25, 0xaf, 0xb1, 0x7e, 0x39, 0x0e, 0x79, 0xc4, 0x52, 0xc9, 0x17, 0x9c, 0xcc, 0x57, 0xcc, 0xc2,
	0x07, 0x71, 0x96, 0xc5, 0x2b, 0x16, 0xea, 0xd7, 0x7c, 0xbd, 0x08, 0x6f, 0x0b, 0x92, 0xe7, 0xac,
	0x10, 0x26, 0xee, 0xcf, 0xa0, 0x7b, 0x7a, 0x71, 0x75, 0xce, 0x57, 0x92, 0x15, 0x68, 0x0f, 0xfe,
	0x9b, 0x4c, 0xcf, 0xb0, 0x18, 0xb4, 0x46, 0xde, 0x61, 0x17, 0x9b, 0x07, 0x7a, 0x0d, 0xed, 0xeb,
	0xcb, 0xa9, 0x18, 0x6c, 0x8d, 0xbc, 0xc3, 0xde, 0xf1, 0xf3, 0xc0, 0x28, 0x06, 0x95, 0x62, 0x30,
	0x93, 0x05, 0x4f, 0xe3, 0x6b, 0xb2, 0x5a, 0x33, 0xac, 0x99, 0xfe, 0x8f, 0x16, 0x74, 0x67, 0xda,
	0xc5, 0x25, 0xa7, 0xe8, 0x00, 0x20, 0x21, 0xf4, 0x63, 0x14, 0x15, 0x4c, 0x28, 0xe9, 0xd6, 0x61,
	0x17, 0x3b, 0x08, 0x42, 0xd0, 0x4e, 0x49, 0xc2, 0x06, 0x5b, 0x3a, 0xa2, 0xbf, 0x51, 0x00, 0x5e,
	0x59, 0x2c, 0x06, 0xde, 0xa8, 0xf5, 0x60, 0x49, 0x45, 0x44, 0xaf, 0x1c, 0x1b, 0x83, 0xb6, 0xce,
	0xda, 0x0e, 0xca, 0x71, 0x50, 0x83, 0xb8, 0x89, 0xfb, 0x19, 0xec, 0x98, 0xee, 0x26, 0x59, 0x9a,
	0x32, 0xaa, 0x86, 0x88, 0x5e, 0x80, 0x97, 0x72, 0xaa, 0xbb, 0xb3, 0xa9, 0xb5, 0x01, 0xac, 0x22,
	0xe8, 0x03, 0x74, 0x13, 0x42, 0x6f, 0x78, 0xca, 0xa6, 0x67, 0xba, 0xd5, 0x87, 0xfa, 0x6a, 0xe8,
	0x7e, 0x0c, 0x1d, 0xa3, 0x86, 0x7c, 0xe8, 0x98, 0x4d, 0xda, 0x4a, 0xa0, 0x2a, 0x4d, 0x34, 0x82,
	0x6d, 0x04, 0xed, 0x43, 0xa7, 0x20, 0x74, 0x69, 0xcb, 0x74, 0xb1, 0x7d, 0xa1, 0x97, 0xd0, 0x4e,
	0x39, 0x15, 0x03, 0x4f, 0xef, 0xe1, 0x4e, 0x8f, 0x3a, 0xe4, 0x7f, 0x83, 0xbe, 0x81, 0x30, 0x8b,
	0xb9, 0x50, 0xb6, 0xd9, 0xf7, 0x35, 0x13, 0x52, 0xd5, 0x35, 0x67, 0xe5, 0xd6, 0xb5, 0x54, 0x1b,
	0x41, 0x23, 0xe8, 0xd5, 0x47, 0x55, 0x17, 0x77, 0x21, 0xff, 0x53, 0x35, 0xb8, 0x0b, 0x26, 0x2b,
	0xe5, 0xb7, 0xf0, 0xc8, 0x3d, 0x3b, 0xab, 0xbf, 0xa3, 0xf4, 0xa7, 0x0e, 0x8e, 0x37, 0x58, 0xfe,
	0x14, 0x76, 0x8d, 0xd2, 0x39, 0x4f, 0x23, 0x47, 0x6a, 0xc1, 0xd3, 0x68, 0x52, 0x70, 0xc9, 0x0a,
	0x4e, 0x5c, 0xa9, 0x73, 0x07, 0xc7, 0x1b, 0x2c, 0xff, 0x67, 0x0b, 0x1e, 0x57, 0xa6, 0x45, 0x9e,
	0xa5, 0x82, 0xfd, 0x93, 0xdb, 0x77, 0xd0, 0xa3, 0xf5, 0xfa, 0xab, 0xe3, 0xde, 0x6b, 0x88, 0xcd,
	0x6d, 0x60, 0x97, 0x88, 0x26, 0xb0, 0x5b, 0x8f, 0xa4, 0x2a, 0x68, 0xef, 0xb4, 0xaf, 0xb2, 0xaf,
	0xee, 0x06, 0xf1, 0x9f, 0x7c, 0xff, 0x69, 0x65, 0xff, 0x33, 0x17, 0xd5, 0x24, 0xfd, 0x33, 0x40,
	0x2e, 0x68, 0xbd, 0x04, 0xf0, 0xbf, 0xe9, 0x98, 0x99, 0x9f, 0x65, 0xef, 0x18, 0x39, 0x6e, 0xaa,
	0x1a, 0x35, 0xe7, 0xf8, 0x57, 0x0b, 0xb6, 0x4d, 0x70, 0xc6, 0x8a, 0x92, 0x53, 0x86, 0xde, 0x43,
	0xe7, 0x4b, 0x1e, 0x11, 0xc9, 0xd0, 0x33, 0x37, 0x73, 0xe3, 0x40, 0x86, 0xf7, 0x88, 0xa2, 0x10,
	0xbc, 0x0b, 0x26, 0x91, 0x33, 0x94, 0x66, 0xef, 0xf7, 0x26, 0x9c, 0x40, 0x5b, 0x2d, 0x0a, 0xf5,
	0x9b, 0x98, 0xb3, 0xdf, 0xe1, 0x7e, 0x03, 0x6f, 0x58, 0x3c, 0x81, 0xb6, 0x7a, 0xbb, 0x69, 0xce,
	0x5c, 0xfe, 0x96, 0x76, 0x1a, 0x7e, 0x3d, 0x8a, 0xb9, 0xbc, 0x59, 0xcf, 0x03, 0x9a, 0x25, 0x61,
	0xc2, 0x24, 0x59, 0x1d, 0x09, 0x49, 0xe8, 0xd2, 0x7e, 0x93, 0x9c, 0x87, 0xf9, 0x32, 0x36, 0x7f,
	0x7d, 0x61, 0x39, 0x9e, 0x77, 0xf4, 0xd7, 0x9b, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xbb, 0xa6,
	0x9a, 0x21, 0x65, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SwitchServiceClient is the client API for SwitchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SwitchServiceClient interface {
	//    rpc Create(SwitchCreateRequest) returns (SwitchResponse);
	//    rpc Update(SwitchUpdateRequest) returns (SwitchResponse);
	Update(ctx context.Context, in *SwitchRegisterRequest, opts ...grpc.CallOption) (*SwitchResponse, error)
	//    rpc Delete(SwitchDeleteRequest) returns (SwitchResponse);
	Get(ctx context.Context, in *SwitchGetRequest, opts ...grpc.CallOption) (*SwitchResponse, error)
	Find(ctx context.Context, in *SwitchFindRequest, opts ...grpc.CallOption) (*SwitchListResponse, error)
	List(ctx context.Context, in *SwitchListRequest, opts ...grpc.CallOption) (*SwitchListResponse, error)
}

type switchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSwitchServiceClient(cc grpc.ClientConnInterface) SwitchServiceClient {
	return &switchServiceClient{cc}
}

func (c *switchServiceClient) Update(ctx context.Context, in *SwitchRegisterRequest, opts ...grpc.CallOption) (*SwitchResponse, error) {
	out := new(SwitchResponse)
	err := c.cc.Invoke(ctx, "/v1.SwitchService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *switchServiceClient) Get(ctx context.Context, in *SwitchGetRequest, opts ...grpc.CallOption) (*SwitchResponse, error) {
	out := new(SwitchResponse)
	err := c.cc.Invoke(ctx, "/v1.SwitchService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *switchServiceClient) Find(ctx context.Context, in *SwitchFindRequest, opts ...grpc.CallOption) (*SwitchListResponse, error) {
	out := new(SwitchListResponse)
	err := c.cc.Invoke(ctx, "/v1.SwitchService/Find", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *switchServiceClient) List(ctx context.Context, in *SwitchListRequest, opts ...grpc.CallOption) (*SwitchListResponse, error) {
	out := new(SwitchListResponse)
	err := c.cc.Invoke(ctx, "/v1.SwitchService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SwitchServiceServer is the server API for SwitchService service.
type SwitchServiceServer interface {
	//    rpc Create(SwitchCreateRequest) returns (SwitchResponse);
	//    rpc Update(SwitchUpdateRequest) returns (SwitchResponse);
	Update(context.Context, *SwitchRegisterRequest) (*SwitchResponse, error)
	//    rpc Delete(SwitchDeleteRequest) returns (SwitchResponse);
	Get(context.Context, *SwitchGetRequest) (*SwitchResponse, error)
	Find(context.Context, *SwitchFindRequest) (*SwitchListResponse, error)
	List(context.Context, *SwitchListRequest) (*SwitchListResponse, error)
}

// UnimplementedSwitchServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSwitchServiceServer struct {
}

func (*UnimplementedSwitchServiceServer) Update(ctx context.Context, req *SwitchRegisterRequest) (*SwitchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedSwitchServiceServer) Get(ctx context.Context, req *SwitchGetRequest) (*SwitchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedSwitchServiceServer) Find(ctx context.Context, req *SwitchFindRequest) (*SwitchListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (*UnimplementedSwitchServiceServer) List(ctx context.Context, req *SwitchListRequest) (*SwitchListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

func RegisterSwitchServiceServer(s *grpc.Server, srv SwitchServiceServer) {
	s.RegisterService(&_SwitchService_serviceDesc, srv)
}

func _SwitchService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SwitchRegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwitchServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.SwitchService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwitchServiceServer).Update(ctx, req.(*SwitchRegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwitchService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SwitchGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwitchServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.SwitchService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwitchServiceServer).Get(ctx, req.(*SwitchGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwitchService_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SwitchFindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwitchServiceServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.SwitchService/Find",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwitchServiceServer).Find(ctx, req.(*SwitchFindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwitchService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SwitchListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwitchServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.SwitchService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwitchServiceServer).List(ctx, req.(*SwitchListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SwitchService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.SwitchService",
	HandlerType: (*SwitchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Update",
			Handler:    _SwitchService_Update_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _SwitchService_Get_Handler,
		},
		{
			MethodName: "Find",
			Handler:    _SwitchService_Find_Handler,
		},
		{
			MethodName: "List",
			Handler:    _SwitchService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/switch.proto",
}