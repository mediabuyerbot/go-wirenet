// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/session.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type OpenSessionRequest struct {
	Sid                  []byte   `protobuf:"bytes,1,opt,name=sid,proto3" json:"sid,omitempty"`
	Token                []byte   `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	Identification       []byte   `protobuf:"bytes,3,opt,name=identification,proto3" json:"identification,omitempty"`
	LocalStreamNames     []string `protobuf:"bytes,4,rep,name=local_stream_names,json=localStreamNames,proto3" json:"local_stream_names,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OpenSessionRequest) Reset()         { *m = OpenSessionRequest{} }
func (m *OpenSessionRequest) String() string { return proto.CompactTextString(m) }
func (*OpenSessionRequest) ProtoMessage()    {}
func (*OpenSessionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_387f10401efe34ae, []int{0}
}

func (m *OpenSessionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OpenSessionRequest.Unmarshal(m, b)
}
func (m *OpenSessionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OpenSessionRequest.Marshal(b, m, deterministic)
}
func (m *OpenSessionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OpenSessionRequest.Merge(m, src)
}
func (m *OpenSessionRequest) XXX_Size() int {
	return xxx_messageInfo_OpenSessionRequest.Size(m)
}
func (m *OpenSessionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OpenSessionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OpenSessionRequest proto.InternalMessageInfo

func (m *OpenSessionRequest) GetSid() []byte {
	if m != nil {
		return m.Sid
	}
	return nil
}

func (m *OpenSessionRequest) GetToken() []byte {
	if m != nil {
		return m.Token
	}
	return nil
}

func (m *OpenSessionRequest) GetIdentification() []byte {
	if m != nil {
		return m.Identification
	}
	return nil
}

func (m *OpenSessionRequest) GetLocalStreamNames() []string {
	if m != nil {
		return m.LocalStreamNames
	}
	return nil
}

type OpenSessionResponse struct {
	Sid                  []byte   `protobuf:"bytes,1,opt,name=sid,proto3" json:"sid,omitempty"`
	RemoteStreamNames    []string `protobuf:"bytes,2,rep,name=remote_stream_names,json=remoteStreamNames,proto3" json:"remote_stream_names,omitempty"`
	Err                  string   `protobuf:"bytes,3,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OpenSessionResponse) Reset()         { *m = OpenSessionResponse{} }
func (m *OpenSessionResponse) String() string { return proto.CompactTextString(m) }
func (*OpenSessionResponse) ProtoMessage()    {}
func (*OpenSessionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_387f10401efe34ae, []int{1}
}

func (m *OpenSessionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OpenSessionResponse.Unmarshal(m, b)
}
func (m *OpenSessionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OpenSessionResponse.Marshal(b, m, deterministic)
}
func (m *OpenSessionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OpenSessionResponse.Merge(m, src)
}
func (m *OpenSessionResponse) XXX_Size() int {
	return xxx_messageInfo_OpenSessionResponse.Size(m)
}
func (m *OpenSessionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_OpenSessionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_OpenSessionResponse proto.InternalMessageInfo

func (m *OpenSessionResponse) GetSid() []byte {
	if m != nil {
		return m.Sid
	}
	return nil
}

func (m *OpenSessionResponse) GetRemoteStreamNames() []string {
	if m != nil {
		return m.RemoteStreamNames
	}
	return nil
}

func (m *OpenSessionResponse) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func init() {
	proto.RegisterType((*OpenSessionRequest)(nil), "pb.OpenSessionRequest")
	proto.RegisterType((*OpenSessionResponse)(nil), "pb.OpenSessionResponse")
}

func init() { proto.RegisterFile("pb/session.proto", fileDescriptor_387f10401efe34ae) }

var fileDescriptor_387f10401efe34ae = []byte{
	// 201 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x6a, 0x84, 0x30,
	0x10, 0x86, 0x51, 0xdb, 0x82, 0x43, 0x29, 0x36, 0xf6, 0x90, 0xa3, 0x78, 0x28, 0x1e, 0x8a, 0x3d,
	0xf4, 0x3d, 0x5a, 0x88, 0x0f, 0x20, 0x51, 0xa7, 0x10, 0xaa, 0x99, 0x34, 0x93, 0x7d, 0x91, 0x7d,
	0xe2, 0x25, 0xf1, 0xb2, 0x2e, 0x7b, 0xfb, 0xe7, 0xff, 0xe0, 0xe3, 0x67, 0xa0, 0x72, 0xd3, 0x27,
	0x23, 0xb3, 0x21, 0xdb, 0x3b, 0x4f, 0x81, 0x44, 0xee, 0xa6, 0xf6, 0x9c, 0x81, 0xf8, 0x71, 0x68,
	0x87, 0x9d, 0x28, 0xfc, 0x3f, 0x21, 0x07, 0x51, 0x41, 0xc1, 0x66, 0x91, 0x59, 0x93, 0x75, 0xcf,
	0x2a, 0x46, 0xf1, 0x06, 0x8f, 0x81, 0xfe, 0xd0, 0xca, 0x3c, 0x75, 0xfb, 0x21, 0xde, 0xe1, 0xc5,
	0x2c, 0x68, 0x83, 0xf9, 0x35, 0xb3, 0x0e, 0x86, 0xac, 0x2c, 0x12, 0xbe, 0x69, 0xc5, 0x07, 0x88,
	0x95, 0x66, 0xbd, 0x8e, 0x1c, 0x3c, 0xea, 0x6d, 0xb4, 0x7a, 0x43, 0x96, 0x0f, 0x4d, 0xd1, 0x95,
	0xaa, 0x4a, 0x64, 0x48, 0xe0, 0x3b, 0xf6, 0xad, 0x81, 0xfa, 0xb0, 0x89, 0x1d, 0x59, 0xc6, 0x3b,
	0xa3, 0x7a, 0xa8, 0x3d, 0x6e, 0x14, 0xf0, 0xe8, 0xcd, 0x93, 0xf7, 0x75, 0x47, 0x57, 0xe2, 0x68,
	0x40, 0xef, 0xd3, 0xc6, 0x52, 0xc5, 0x38, 0x3d, 0xa5, 0x57, 0x7c, 0x5d, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x73, 0xe0, 0x37, 0x75, 0x1e, 0x01, 0x00, 0x00,
}
