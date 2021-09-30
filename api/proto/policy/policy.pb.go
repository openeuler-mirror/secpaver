// Code generated by protoc-gen-go. DO NOT EDIT.
// source: policy.proto

package policy

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// common request message
type Req struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Engine               string   `protobuf:"bytes,2,opt,name=engine,proto3" json:"engine,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Req) Reset()         { *m = Req{} }
func (m *Req) String() string { return proto.CompactTextString(m) }
func (*Req) ProtoMessage()    {}
func (*Req) Descriptor() ([]byte, []int) {
	return fileDescriptor_ac3b897852294d6a, []int{0}
}

func (m *Req) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Req.Unmarshal(m, b)
}
func (m *Req) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Req.Marshal(b, m, deterministic)
}
func (m *Req) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Req.Merge(m, src)
}
func (m *Req) XXX_Size() int {
	return xxx_messageInfo_Req.Size(m)
}
func (m *Req) XXX_DiscardUnknown() {
	xxx_messageInfo_Req.DiscardUnknown(m)
}

var xxx_messageInfo_Req proto.InternalMessageInfo

func (m *Req) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Req) GetEngine() string {
	if m != nil {
		return m.Engine
	}
	return ""
}

// common ack message
type Ack struct {
	Level                string   `protobuf:"bytes,1,opt,name=level,proto3" json:"level,omitempty"`
	Status               string   `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ack) Reset()         { *m = Ack{} }
func (m *Ack) String() string { return proto.CompactTextString(m) }
func (*Ack) ProtoMessage()    {}
func (*Ack) Descriptor() ([]byte, []int) {
	return fileDescriptor_ac3b897852294d6a, []int{1}
}

func (m *Ack) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ack.Unmarshal(m, b)
}
func (m *Ack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ack.Marshal(b, m, deterministic)
}
func (m *Ack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ack.Merge(m, src)
}
func (m *Ack) XXX_Size() int {
	return xxx_messageInfo_Ack.Size(m)
}
func (m *Ack) XXX_DiscardUnknown() {
	xxx_messageInfo_Ack.DiscardUnknown(m)
}

var xxx_messageInfo_Ack proto.InternalMessageInfo

func (m *Ack) GetLevel() string {
	if m != nil {
		return m.Level
	}
	return ""
}

func (m *Ack) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type PolicyInfo struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Status               string   `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PolicyInfo) Reset()         { *m = PolicyInfo{} }
func (m *PolicyInfo) String() string { return proto.CompactTextString(m) }
func (*PolicyInfo) ProtoMessage()    {}
func (*PolicyInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ac3b897852294d6a, []int{2}
}

func (m *PolicyInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PolicyInfo.Unmarshal(m, b)
}
func (m *PolicyInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PolicyInfo.Marshal(b, m, deterministic)
}
func (m *PolicyInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PolicyInfo.Merge(m, src)
}
func (m *PolicyInfo) XXX_Size() int {
	return xxx_messageInfo_PolicyInfo.Size(m)
}
func (m *PolicyInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PolicyInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PolicyInfo proto.InternalMessageInfo

func (m *PolicyInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PolicyInfo) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

// ack message of ListPolicy rpc
type ListPolicyAck struct {
	PolicyInfos          []*PolicyInfo `protobuf:"bytes,1,rep,name=policyInfos,proto3" json:"policyInfos,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ListPolicyAck) Reset()         { *m = ListPolicyAck{} }
func (m *ListPolicyAck) String() string { return proto.CompactTextString(m) }
func (*ListPolicyAck) ProtoMessage()    {}
func (*ListPolicyAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_ac3b897852294d6a, []int{3}
}

func (m *ListPolicyAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListPolicyAck.Unmarshal(m, b)
}
func (m *ListPolicyAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListPolicyAck.Marshal(b, m, deterministic)
}
func (m *ListPolicyAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListPolicyAck.Merge(m, src)
}
func (m *ListPolicyAck) XXX_Size() int {
	return xxx_messageInfo_ListPolicyAck.Size(m)
}
func (m *ListPolicyAck) XXX_DiscardUnknown() {
	xxx_messageInfo_ListPolicyAck.DiscardUnknown(m)
}

var xxx_messageInfo_ListPolicyAck proto.InternalMessageInfo

func (m *ListPolicyAck) GetPolicyInfos() []*PolicyInfo {
	if m != nil {
		return m.PolicyInfos
	}
	return nil
}

// PolicyZipFile zip file message
type PolicyZipFile struct {
	Filename             string   `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PolicyZipFile) Reset()         { *m = PolicyZipFile{} }
func (m *PolicyZipFile) String() string { return proto.CompactTextString(m) }
func (*PolicyZipFile) ProtoMessage()    {}
func (*PolicyZipFile) Descriptor() ([]byte, []int) {
	return fileDescriptor_ac3b897852294d6a, []int{4}
}

func (m *PolicyZipFile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PolicyZipFile.Unmarshal(m, b)
}
func (m *PolicyZipFile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PolicyZipFile.Marshal(b, m, deterministic)
}
func (m *PolicyZipFile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PolicyZipFile.Merge(m, src)
}
func (m *PolicyZipFile) XXX_Size() int {
	return xxx_messageInfo_PolicyZipFile.Size(m)
}
func (m *PolicyZipFile) XXX_DiscardUnknown() {
	xxx_messageInfo_PolicyZipFile.DiscardUnknown(m)
}

var xxx_messageInfo_PolicyZipFile proto.InternalMessageInfo

func (m *PolicyZipFile) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *PolicyZipFile) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// ack message of ExportPolicy rpc
type ExportPolicyAck struct {
	File                 *PolicyZipFile `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ExportPolicyAck) Reset()         { *m = ExportPolicyAck{} }
func (m *ExportPolicyAck) String() string { return proto.CompactTextString(m) }
func (*ExportPolicyAck) ProtoMessage()    {}
func (*ExportPolicyAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_ac3b897852294d6a, []int{5}
}

func (m *ExportPolicyAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExportPolicyAck.Unmarshal(m, b)
}
func (m *ExportPolicyAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExportPolicyAck.Marshal(b, m, deterministic)
}
func (m *ExportPolicyAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExportPolicyAck.Merge(m, src)
}
func (m *ExportPolicyAck) XXX_Size() int {
	return xxx_messageInfo_ExportPolicyAck.Size(m)
}
func (m *ExportPolicyAck) XXX_DiscardUnknown() {
	xxx_messageInfo_ExportPolicyAck.DiscardUnknown(m)
}

var xxx_messageInfo_ExportPolicyAck proto.InternalMessageInfo

func (m *ExportPolicyAck) GetFile() *PolicyZipFile {
	if m != nil {
		return m.File
	}
	return nil
}

func init() {
	proto.RegisterType((*Req)(nil), "policy.Req")
	proto.RegisterType((*Ack)(nil), "policy.Ack")
	proto.RegisterType((*PolicyInfo)(nil), "policy.PolicyInfo")
	proto.RegisterType((*ListPolicyAck)(nil), "policy.ListPolicyAck")
	proto.RegisterType((*PolicyZipFile)(nil), "policy.PolicyZipFile")
	proto.RegisterType((*ExportPolicyAck)(nil), "policy.ExportPolicyAck")
}

func init() { proto.RegisterFile("policy.proto", fileDescriptor_ac3b897852294d6a) }

var fileDescriptor_ac3b897852294d6a = []byte{
	// 309 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x4f, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0x89, 0xad, 0xc5, 0x4e, 0x52, 0x0a, 0x83, 0x7f, 0x42, 0x4f, 0x25, 0xa7, 0x88, 0x58,
	0x6a, 0xda, 0x83, 0x07, 0x41, 0x0a, 0x56, 0x28, 0x28, 0xc8, 0x82, 0x17, 0x6f, 0x31, 0x4e, 0xcb,
	0xd2, 0x75, 0x93, 0x26, 0xab, 0xe8, 0x47, 0xf7, 0x26, 0xd9, 0x4d, 0xd2, 0x44, 0x2a, 0x7a, 0xdb,
	0xb7, 0x79, 0xbf, 0x7d, 0x6f, 0x86, 0x80, 0x93, 0xc4, 0x82, 0x47, 0x9f, 0xa3, 0x24, 0x8d, 0x55,
	0x8c, 0x1d, 0xa3, 0xbc, 0x0b, 0x68, 0x31, 0xda, 0x20, 0x42, 0x5b, 0x86, 0xaf, 0xe4, 0x5a, 0x43,
	0xcb, 0xef, 0x32, 0x7d, 0xc6, 0x63, 0xe8, 0x90, 0x5c, 0x71, 0x49, 0xee, 0x9e, 0xbe, 0x2d, 0x94,
	0x37, 0x81, 0xd6, 0x2c, 0x5a, 0xe3, 0x21, 0xec, 0x0b, 0x7a, 0x27, 0x51, 0x30, 0x46, 0xe4, 0x50,
	0xa6, 0x42, 0xf5, 0x96, 0x95, 0x90, 0x51, 0xde, 0x25, 0xc0, 0x83, 0x4e, 0x5c, 0xc8, 0x65, 0xfc,
	0x5b, 0xdc, 0x4e, 0x72, 0x0e, 0xbd, 0x3b, 0x9e, 0x29, 0x43, 0xe7, 0xc1, 0x53, 0xb0, 0x93, 0xea,
	0xa9, 0xcc, 0xb5, 0x86, 0x2d, 0xdf, 0x0e, 0x70, 0x54, 0x8c, 0xb7, 0x4d, 0x61, 0x75, 0x9b, 0x77,
	0x0d, 0x3d, 0xf3, 0xe9, 0x89, 0x27, 0xb7, 0x5c, 0x10, 0x0e, 0xe0, 0x60, 0xc9, 0x05, 0xd5, 0x7a,
	0x54, 0x3a, 0xef, 0xf7, 0x12, 0xaa, 0x50, 0x37, 0x71, 0x98, 0x3e, 0x7b, 0x57, 0xd0, 0x9f, 0x7f,
	0x24, 0x71, 0x5a, 0x6b, 0x72, 0x0a, 0xed, 0x1c, 0xd1, 0xb8, 0x1d, 0x1c, 0x35, 0x2b, 0x14, 0x39,
	0x4c, 0x5b, 0x82, 0x2f, 0x0b, 0xba, 0xe6, 0xfe, 0x7e, 0x95, 0xe2, 0x18, 0x60, 0x3b, 0x13, 0xda,
	0x25, 0xc8, 0x68, 0x33, 0xa8, 0x5e, 0x69, 0x0e, 0x7d, 0x06, 0xbd, 0x85, 0xcc, 0x54, 0x28, 0xc4,
	0x2e, 0xa8, 0x12, 0xb3, 0x68, 0x3d, 0xb6, 0xf0, 0x1c, 0xfa, 0x8f, 0x92, 0xff, 0xdb, 0xee, 0x83,
	0x73, 0x43, 0x82, 0x14, 0xfd, 0xe5, 0xc5, 0x29, 0x38, 0xf5, 0x1d, 0x34, 0x9d, 0x27, 0xa5, 0xf8,
	0xb1, 0xa6, 0xe7, 0x8e, 0xfe, 0xe5, 0x26, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x35, 0x67, 0x80,
	0xd7, 0x82, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PolicyMgrClient is the client API for PolicyMgr service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PolicyMgrClient interface {
	ListPolicy(ctx context.Context, in *Req, opts ...grpc.CallOption) (*ListPolicyAck, error)
	InstallPolicy(ctx context.Context, in *Req, opts ...grpc.CallOption) (PolicyMgr_InstallPolicyClient, error)
	UninstallPolicy(ctx context.Context, in *Req, opts ...grpc.CallOption) (PolicyMgr_UninstallPolicyClient, error)
	DeletePolicy(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Ack, error)
	ExportPolicy(ctx context.Context, in *Req, opts ...grpc.CallOption) (*ExportPolicyAck, error)
}

type policyMgrClient struct {
	cc *grpc.ClientConn
}

func NewPolicyMgrClient(cc *grpc.ClientConn) PolicyMgrClient {
	return &policyMgrClient{cc}
}

func (c *policyMgrClient) ListPolicy(ctx context.Context, in *Req, opts ...grpc.CallOption) (*ListPolicyAck, error) {
	out := new(ListPolicyAck)
	err := c.cc.Invoke(ctx, "/policy.PolicyMgr/ListPolicy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyMgrClient) InstallPolicy(ctx context.Context, in *Req, opts ...grpc.CallOption) (PolicyMgr_InstallPolicyClient, error) {
	stream, err := c.cc.NewStream(ctx, &_PolicyMgr_serviceDesc.Streams[0], "/policy.PolicyMgr/InstallPolicy", opts...)
	if err != nil {
		return nil, err
	}
	x := &policyMgrInstallPolicyClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PolicyMgr_InstallPolicyClient interface {
	Recv() (*Ack, error)
	grpc.ClientStream
}

type policyMgrInstallPolicyClient struct {
	grpc.ClientStream
}

func (x *policyMgrInstallPolicyClient) Recv() (*Ack, error) {
	m := new(Ack)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *policyMgrClient) UninstallPolicy(ctx context.Context, in *Req, opts ...grpc.CallOption) (PolicyMgr_UninstallPolicyClient, error) {
	stream, err := c.cc.NewStream(ctx, &_PolicyMgr_serviceDesc.Streams[1], "/policy.PolicyMgr/UninstallPolicy", opts...)
	if err != nil {
		return nil, err
	}
	x := &policyMgrUninstallPolicyClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PolicyMgr_UninstallPolicyClient interface {
	Recv() (*Ack, error)
	grpc.ClientStream
}

type policyMgrUninstallPolicyClient struct {
	grpc.ClientStream
}

func (x *policyMgrUninstallPolicyClient) Recv() (*Ack, error) {
	m := new(Ack)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *policyMgrClient) DeletePolicy(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/policy.PolicyMgr/DeletePolicy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyMgrClient) ExportPolicy(ctx context.Context, in *Req, opts ...grpc.CallOption) (*ExportPolicyAck, error) {
	out := new(ExportPolicyAck)
	err := c.cc.Invoke(ctx, "/policy.PolicyMgr/ExportPolicy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PolicyMgrServer is the server API for PolicyMgr service.
type PolicyMgrServer interface {
	ListPolicy(context.Context, *Req) (*ListPolicyAck, error)
	InstallPolicy(*Req, PolicyMgr_InstallPolicyServer) error
	UninstallPolicy(*Req, PolicyMgr_UninstallPolicyServer) error
	DeletePolicy(context.Context, *Req) (*Ack, error)
	ExportPolicy(context.Context, *Req) (*ExportPolicyAck, error)
}

// UnimplementedPolicyMgrServer can be embedded to have forward compatible implementations.
type UnimplementedPolicyMgrServer struct {
}

func (*UnimplementedPolicyMgrServer) ListPolicy(ctx context.Context, req *Req) (*ListPolicyAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPolicy not implemented")
}
func (*UnimplementedPolicyMgrServer) InstallPolicy(req *Req, srv PolicyMgr_InstallPolicyServer) error {
	return status.Errorf(codes.Unimplemented, "method InstallPolicy not implemented")
}
func (*UnimplementedPolicyMgrServer) UninstallPolicy(req *Req, srv PolicyMgr_UninstallPolicyServer) error {
	return status.Errorf(codes.Unimplemented, "method UninstallPolicy not implemented")
}
func (*UnimplementedPolicyMgrServer) DeletePolicy(ctx context.Context, req *Req) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePolicy not implemented")
}
func (*UnimplementedPolicyMgrServer) ExportPolicy(ctx context.Context, req *Req) (*ExportPolicyAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportPolicy not implemented")
}

func RegisterPolicyMgrServer(s *grpc.Server, srv PolicyMgrServer) {
	s.RegisterService(&_PolicyMgr_serviceDesc, srv)
}

func _PolicyMgr_ListPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyMgrServer).ListPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/policy.PolicyMgr/ListPolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyMgrServer).ListPolicy(ctx, req.(*Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyMgr_InstallPolicy_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Req)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PolicyMgrServer).InstallPolicy(m, &policyMgrInstallPolicyServer{stream})
}

type PolicyMgr_InstallPolicyServer interface {
	Send(*Ack) error
	grpc.ServerStream
}

type policyMgrInstallPolicyServer struct {
	grpc.ServerStream
}

func (x *policyMgrInstallPolicyServer) Send(m *Ack) error {
	return x.ServerStream.SendMsg(m)
}

func _PolicyMgr_UninstallPolicy_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Req)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PolicyMgrServer).UninstallPolicy(m, &policyMgrUninstallPolicyServer{stream})
}

type PolicyMgr_UninstallPolicyServer interface {
	Send(*Ack) error
	grpc.ServerStream
}

type policyMgrUninstallPolicyServer struct {
	grpc.ServerStream
}

func (x *policyMgrUninstallPolicyServer) Send(m *Ack) error {
	return x.ServerStream.SendMsg(m)
}

func _PolicyMgr_DeletePolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyMgrServer).DeletePolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/policy.PolicyMgr/DeletePolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyMgrServer).DeletePolicy(ctx, req.(*Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyMgr_ExportPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyMgrServer).ExportPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/policy.PolicyMgr/ExportPolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyMgrServer).ExportPolicy(ctx, req.(*Req))
	}
	return interceptor(ctx, in, info, handler)
}

var _PolicyMgr_serviceDesc = grpc.ServiceDesc{
	ServiceName: "policy.PolicyMgr",
	HandlerType: (*PolicyMgrServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListPolicy",
			Handler:    _PolicyMgr_ListPolicy_Handler,
		},
		{
			MethodName: "DeletePolicy",
			Handler:    _PolicyMgr_DeletePolicy_Handler,
		},
		{
			MethodName: "ExportPolicy",
			Handler:    _PolicyMgr_ExportPolicy_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "InstallPolicy",
			Handler:       _PolicyMgr_InstallPolicy_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "UninstallPolicy",
			Handler:       _PolicyMgr_UninstallPolicy_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "policy.proto",
}