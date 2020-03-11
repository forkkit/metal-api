// Code generated by protoc-gen-go. DO NOT EDIT.
// source: v1/firewall.proto

package v1

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
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

type FirewallCreateRequest struct {
	MachineAllocateRequest *MachineAllocateRequest `protobuf:"bytes,1,opt,name=machineAllocateRequest,proto3" json:"machineAllocateRequest,omitempty"`
	HA                     *wrappers.BoolValue     `protobuf:"bytes,2,opt,name=HA,proto3" json:"HA,omitempty"`
	XXX_NoUnkeyedLiteral   struct{}                `json:"-"`
	XXX_unrecognized       []byte                  `json:"-"`
	XXX_sizecache          int32                   `json:"-"`
}

func (m *FirewallCreateRequest) Reset()         { *m = FirewallCreateRequest{} }
func (m *FirewallCreateRequest) String() string { return proto.CompactTextString(m) }
func (*FirewallCreateRequest) ProtoMessage()    {}
func (*FirewallCreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2e3898ff840d6971, []int{0}
}

func (m *FirewallCreateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FirewallCreateRequest.Unmarshal(m, b)
}
func (m *FirewallCreateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FirewallCreateRequest.Marshal(b, m, deterministic)
}
func (m *FirewallCreateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FirewallCreateRequest.Merge(m, src)
}
func (m *FirewallCreateRequest) XXX_Size() int {
	return xxx_messageInfo_FirewallCreateRequest.Size(m)
}
func (m *FirewallCreateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FirewallCreateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FirewallCreateRequest proto.InternalMessageInfo

func (m *FirewallCreateRequest) GetMachineAllocateRequest() *MachineAllocateRequest {
	if m != nil {
		return m.MachineAllocateRequest
	}
	return nil
}

func (m *FirewallCreateRequest) GetHA() *wrappers.BoolValue {
	if m != nil {
		return m.HA
	}
	return nil
}

type FirewallFindRequest struct {
	MachineFindRequest   *MachineFindRequest `protobuf:"bytes,1,opt,name=machineFindRequest,proto3" json:"machineFindRequest,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *FirewallFindRequest) Reset()         { *m = FirewallFindRequest{} }
func (m *FirewallFindRequest) String() string { return proto.CompactTextString(m) }
func (*FirewallFindRequest) ProtoMessage()    {}
func (*FirewallFindRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2e3898ff840d6971, []int{1}
}

func (m *FirewallFindRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FirewallFindRequest.Unmarshal(m, b)
}
func (m *FirewallFindRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FirewallFindRequest.Marshal(b, m, deterministic)
}
func (m *FirewallFindRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FirewallFindRequest.Merge(m, src)
}
func (m *FirewallFindRequest) XXX_Size() int {
	return xxx_messageInfo_FirewallFindRequest.Size(m)
}
func (m *FirewallFindRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FirewallFindRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FirewallFindRequest proto.InternalMessageInfo

func (m *FirewallFindRequest) GetMachineFindRequest() *MachineFindRequest {
	if m != nil {
		return m.MachineFindRequest
	}
	return nil
}

type FirewallResponse struct {
	MachineResponse      *MachineResponse `protobuf:"bytes,1,opt,name=machineResponse,proto3" json:"machineResponse,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *FirewallResponse) Reset()         { *m = FirewallResponse{} }
func (m *FirewallResponse) String() string { return proto.CompactTextString(m) }
func (*FirewallResponse) ProtoMessage()    {}
func (*FirewallResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2e3898ff840d6971, []int{2}
}

func (m *FirewallResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FirewallResponse.Unmarshal(m, b)
}
func (m *FirewallResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FirewallResponse.Marshal(b, m, deterministic)
}
func (m *FirewallResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FirewallResponse.Merge(m, src)
}
func (m *FirewallResponse) XXX_Size() int {
	return xxx_messageInfo_FirewallResponse.Size(m)
}
func (m *FirewallResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FirewallResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FirewallResponse proto.InternalMessageInfo

func (m *FirewallResponse) GetMachineResponse() *MachineResponse {
	if m != nil {
		return m.MachineResponse
	}
	return nil
}

func init() {
	proto.RegisterType((*FirewallCreateRequest)(nil), "v1.FirewallCreateRequest")
	proto.RegisterType((*FirewallFindRequest)(nil), "v1.FirewallFindRequest")
	proto.RegisterType((*FirewallResponse)(nil), "v1.FirewallResponse")
}

func init() { proto.RegisterFile("v1/firewall.proto", fileDescriptor_2e3898ff840d6971) }

var fileDescriptor_2e3898ff840d6971 = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xb1, 0x4e, 0xeb, 0x30,
	0x14, 0x86, 0x95, 0xe8, 0xaa, 0x83, 0xef, 0xd0, 0xe2, 0x42, 0x29, 0x19, 0x10, 0xea, 0x84, 0x90,
	0x1a, 0x2b, 0x45, 0x0c, 0x0c, 0x0c, 0x2d, 0x52, 0xd5, 0x85, 0x01, 0x23, 0x31, 0x20, 0x31, 0xb8,
	0xe1, 0x34, 0xb5, 0xea, 0xc6, 0xc6, 0x71, 0xdc, 0x07, 0xe0, 0x21, 0x78, 0x5d, 0xd4, 0xc6, 0x06,
	0x0b, 0x85, 0xcd, 0x39, 0xff, 0xf9, 0x3f, 0x7d, 0xb1, 0x8c, 0x8e, 0x6c, 0x46, 0x56, 0x5c, 0xc3,
	0x8e, 0x09, 0x91, 0x2a, 0x2d, 0x8d, 0xc4, 0xb1, 0xcd, 0x92, 0x9e, 0xcd, 0xc8, 0x96, 0xe5, 0x6b,
	0x5e, 0x42, 0x33, 0x4d, 0xce, 0x0b, 0x29, 0x0b, 0x01, 0xe4, 0xf0, 0xb5, 0xac, 0x57, 0x64, 0xa7,
	0x99, 0x52, 0xa0, 0xab, 0x26, 0x1f, 0x7d, 0x46, 0xe8, 0x64, 0xee, 0x40, 0xf7, 0x1a, 0x98, 0x01,
	0x0a, 0xef, 0x35, 0x54, 0x06, 0x53, 0x34, 0x70, 0xa8, 0xa9, 0x10, 0x32, 0xff, 0x49, 0x86, 0xd1,
	0x45, 0x74, 0xf9, 0x7f, 0x92, 0xa4, 0x36, 0x4b, 0x1f, 0x5a, 0x37, 0xe8, 0x1f, 0x4d, 0x7c, 0x85,
	0xe2, 0xc5, 0x74, 0x18, 0xbb, 0x7e, 0xa3, 0x96, 0x7a, 0xb5, 0x74, 0x26, 0xa5, 0x78, 0x66, 0xa2,
	0x06, 0x1a, 0x2f, 0xa6, 0xa3, 0x57, 0xd4, 0xf7, 0x62, 0x73, 0x5e, 0xbe, 0x79, 0xc4, 0x1c, 0x61,
	0x07, 0x0f, 0xa6, 0x4e, 0x69, 0x10, 0x28, 0x05, 0x29, 0x6d, 0x69, 0x8c, 0x1e, 0x51, 0xcf, 0xe3,
	0x29, 0x54, 0x4a, 0x96, 0x15, 0xe0, 0x3b, 0xd4, 0x75, 0x9b, 0x7e, 0xe4, 0xc0, 0xfd, 0x00, 0xec,
	0x23, 0xfa, 0x7b, 0x77, 0xf2, 0x11, 0xa1, 0xae, 0x67, 0x3e, 0x81, 0xb6, 0x3c, 0x07, 0x7c, 0x8b,
	0x3a, 0xcd, 0xb5, 0xe2, 0xb3, 0x3d, 0xa3, 0xf5, 0xaa, 0x93, 0xe3, 0x30, 0xfa, 0xb6, 0xb9, 0x41,
	0xff, 0xf6, 0xc2, 0xf8, 0x34, 0x4c, 0x83, 0x5f, 0x68, 0xaf, 0xcd, 0xc8, 0xcb, 0xb8, 0xe0, 0x66,
	0x5d, 0x2f, 0xd3, 0x5c, 0x6e, 0xc9, 0x16, 0x0c, 0x13, 0xe3, 0xca, 0xb0, 0x7c, 0xe3, 0xce, 0x4c,
	0x71, 0xa2, 0x36, 0x45, 0xf3, 0x22, 0x88, 0xcd, 0x96, 0x9d, 0xc3, 0xe9, 0xfa, 0x2b, 0x00, 0x00,
	0xff, 0xff, 0x6b, 0x41, 0x4c, 0xcf, 0x54, 0x02, 0x00, 0x00,
}
