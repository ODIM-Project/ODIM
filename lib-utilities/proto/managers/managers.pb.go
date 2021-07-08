// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/managers/managers.proto

package managers

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

type ManagerRequest struct {
	SessionToken         string   `protobuf:"bytes,1,opt,name=sessionToken,proto3" json:"sessionToken,omitempty"`
	ManagerID            string   `protobuf:"bytes,2,opt,name=managerID,proto3" json:"managerID,omitempty"`
	URL                  string   `protobuf:"bytes,3,opt,name=URL,proto3" json:"URL,omitempty"`
	ResourceID           string   `protobuf:"bytes,4,opt,name=resourceID,proto3" json:"resourceID,omitempty"`
	RequestBody          []byte   `protobuf:"bytes,5,opt,name=RequestBody,proto3" json:"RequestBody,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ManagerRequest) Reset()         { *m = ManagerRequest{} }
func (m *ManagerRequest) String() string { return proto.CompactTextString(m) }
func (*ManagerRequest) ProtoMessage()    {}
func (*ManagerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_17e14ce9c60587aa, []int{0}
}

func (m *ManagerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ManagerRequest.Unmarshal(m, b)
}
func (m *ManagerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ManagerRequest.Marshal(b, m, deterministic)
}
func (m *ManagerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagerRequest.Merge(m, src)
}
func (m *ManagerRequest) XXX_Size() int {
	return xxx_messageInfo_ManagerRequest.Size(m)
}
func (m *ManagerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ManagerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ManagerRequest proto.InternalMessageInfo

func (m *ManagerRequest) GetSessionToken() string {
	if m != nil {
		return m.SessionToken
	}
	return ""
}

func (m *ManagerRequest) GetManagerID() string {
	if m != nil {
		return m.ManagerID
	}
	return ""
}

func (m *ManagerRequest) GetURL() string {
	if m != nil {
		return m.URL
	}
	return ""
}

func (m *ManagerRequest) GetResourceID() string {
	if m != nil {
		return m.ResourceID
	}
	return ""
}

func (m *ManagerRequest) GetRequestBody() []byte {
	if m != nil {
		return m.RequestBody
	}
	return nil
}

type ManagerResponse struct {
	StatusCode           int32             `protobuf:"varint,1,opt,name=statusCode,proto3" json:"statusCode,omitempty"`
	StatusMessage        string            `protobuf:"bytes,2,opt,name=statusMessage,proto3" json:"statusMessage,omitempty"`
	Body                 []byte            `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	Header               map[string]string `protobuf:"bytes,5,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ManagerResponse) Reset()         { *m = ManagerResponse{} }
func (m *ManagerResponse) String() string { return proto.CompactTextString(m) }
func (*ManagerResponse) ProtoMessage()    {}
func (*ManagerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_17e14ce9c60587aa, []int{1}
}

func (m *ManagerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ManagerResponse.Unmarshal(m, b)
}
func (m *ManagerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ManagerResponse.Marshal(b, m, deterministic)
}
func (m *ManagerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagerResponse.Merge(m, src)
}
func (m *ManagerResponse) XXX_Size() int {
	return xxx_messageInfo_ManagerResponse.Size(m)
}
func (m *ManagerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ManagerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ManagerResponse proto.InternalMessageInfo

func (m *ManagerResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *ManagerResponse) GetStatusMessage() string {
	if m != nil {
		return m.StatusMessage
	}
	return ""
}

func (m *ManagerResponse) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *ManagerResponse) GetHeader() map[string]string {
	if m != nil {
		return m.Header
	}
	return nil
}

func init() {
	proto.RegisterType((*ManagerRequest)(nil), "ManagerRequest")
	proto.RegisterType((*ManagerResponse)(nil), "ManagerResponse")
	proto.RegisterMapType((map[string]string)(nil), "ManagerResponse.HeaderEntry")
}

func init() { proto.RegisterFile("proto/managers/managers.proto", fileDescriptor_17e14ce9c60587aa) }

var fileDescriptor_17e14ce9c60587aa = []byte{
	// 355 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x41, 0x4f, 0xea, 0x40,
	0x14, 0x85, 0x5f, 0x81, 0x92, 0xc7, 0x85, 0xf7, 0xc0, 0x51, 0x93, 0x86, 0xa0, 0x69, 0x1a, 0x17,
	0xac, 0x6a, 0x44, 0x17, 0x40, 0x5c, 0x09, 0x44, 0x49, 0x64, 0xd3, 0xa8, 0xfb, 0x81, 0xde, 0x60,
	0xa5, 0x76, 0x70, 0xee, 0xd4, 0x84, 0x3f, 0xe4, 0xca, 0x9f, 0xe4, 0x8f, 0x31, 0x1d, 0x0a, 0x14,
	0x56, 0xdd, 0x9d, 0xfb, 0xcd, 0x9c, 0x39, 0x27, 0x33, 0x03, 0x67, 0x4b, 0x29, 0x94, 0xb8, 0x7c,
	0xe7, 0x11, 0x9f, 0xa3, 0xa4, 0xad, 0x70, 0x35, 0x77, 0xbe, 0x0c, 0xf8, 0x3f, 0x59, 0x23, 0x0f,
	0x3f, 0x62, 0x24, 0xc5, 0x1c, 0xa8, 0x11, 0x12, 0x05, 0x22, 0x7a, 0x12, 0x0b, 0x8c, 0x2c, 0xc3,
	0x36, 0xda, 0x15, 0x6f, 0x8f, 0xb1, 0x16, 0x54, 0xd2, 0x83, 0xc6, 0x43, 0xab, 0xa0, 0x37, 0xec,
	0x00, 0x6b, 0x40, 0xf1, 0xd9, 0x7b, 0xb4, 0x8a, 0x9a, 0x27, 0x92, 0x9d, 0x03, 0x48, 0x24, 0x11,
	0xcb, 0x19, 0x8e, 0x87, 0x56, 0x49, 0x2f, 0x64, 0x08, 0xb3, 0xa1, 0x9a, 0xc6, 0xdf, 0x09, 0x7f,
	0x65, 0x99, 0xb6, 0xd1, 0xae, 0x79, 0x59, 0xe4, 0xfc, 0x18, 0x50, 0xdf, 0x16, 0xa5, 0xa5, 0x88,
	0x08, 0x93, 0x53, 0x49, 0x71, 0x15, 0xd3, 0x40, 0xf8, 0xa8, 0x7b, 0x9a, 0x5e, 0x86, 0xb0, 0x0b,
	0xf8, 0xb7, 0x9e, 0x26, 0x48, 0xc4, 0xe7, 0x98, 0x36, 0xdd, 0x87, 0x8c, 0x41, 0x69, 0x9a, 0x84,
	0x96, 0x74, 0xa8, 0xd6, 0xec, 0x06, 0xca, 0xaf, 0xc8, 0x7d, 0x94, 0x96, 0x69, 0x17, 0xdb, 0xd5,
	0x4e, 0xcb, 0x3d, 0xc8, 0x76, 0x1f, 0xf4, 0xf2, 0x28, 0x52, 0x72, 0xe5, 0xa5, 0x7b, 0x9b, 0x3d,
	0xa8, 0x66, 0x70, 0x72, 0x0d, 0x0b, 0x5c, 0xa5, 0xf7, 0x97, 0x48, 0x76, 0x02, 0xe6, 0x27, 0x0f,
	0xe3, 0x4d, 0x91, 0xf5, 0xd0, 0x2f, 0x74, 0x8d, 0xce, 0x77, 0x01, 0xfe, 0xa6, 0x11, 0xc4, 0x6e,
	0xe1, 0xf4, 0x1e, 0xd5, 0x66, 0x1c, 0x88, 0x30, 0xc4, 0x99, 0x0a, 0x44, 0xc4, 0xea, 0xee, 0xfe,
	0x5b, 0x35, 0x1b, 0x87, 0xbd, 0x9c, 0x3f, 0xec, 0x0a, 0x60, 0xe7, 0xce, 0x67, 0xe9, 0xc3, 0x71,
	0x26, 0xd0, 0x4b, 0xdf, 0x25, 0x9f, 0xb7, 0x07, 0xec, 0x25, 0x90, 0x2a, 0xe6, 0xe1, 0x04, 0xfd,
	0x80, 0x8f, 0x23, 0x42, 0xa9, 0xf2, 0x59, 0xbb, 0x70, 0x94, 0xb5, 0x8e, 0xde, 0x70, 0x96, 0xcf,
	0x39, 0x2d, 0xeb, 0xdf, 0x7b, 0xfd, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x07, 0xce, 0x78, 0xb5, 0xde,
	0x02, 0x00, 0x00,
}
