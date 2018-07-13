// Code generated by protoc-gen-go.
// source: message_proto.proto
// DO NOT EDIT!

/*
Package gridworker is a generated protocol buffer package.

It is generated from these files:
	message_proto.proto

It has these top-level messages:
	DynamicType
	MessageProto
*/
package gridworker

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type DynamicTypeType int32

const (
	DynamicTypeType_int        DynamicTypeType = 0
	DynamicTypeType_float      DynamicTypeType = 1
	DynamicTypeType_string     DynamicTypeType = 2
	DynamicTypeType_bytes      DynamicTypeType = 3
	DynamicTypeType_bool       DynamicTypeType = 4
	DynamicTypeType_json_slice DynamicTypeType = 5
	DynamicTypeType_json_map   DynamicTypeType = 6
)

var DynamicTypeType_name = map[int32]string{
	0: "int",
	1: "float",
	2: "string",
	3: "bytes",
	4: "bool",
	5: "json_slice",
	6: "json_map",
}
var DynamicTypeType_value = map[string]int32{
	"int":        0,
	"float":      1,
	"string":     2,
	"bytes":      3,
	"bool":       4,
	"json_slice": 5,
	"json_map":   6,
}

func (x DynamicTypeType) Enum() *DynamicTypeType {
	p := new(DynamicTypeType)
	*p = x
	return p
}
func (x DynamicTypeType) String() string {
	return proto.EnumName(DynamicTypeType_name, int32(x))
}
func (x *DynamicTypeType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(DynamicTypeType_value, data, "DynamicTypeType")
	if err != nil {
		return err
	}
	*x = DynamicTypeType(value)
	return nil
}
func (DynamicTypeType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type DynamicType struct {
	Type             *DynamicTypeType `protobuf:"varint,1,req,name=type,enum=gridworker.DynamicTypeType" json:"type,omitempty"`
	Int              *int64           `protobuf:"varint,2,opt,name=int" json:"int,omitempty"`
	Float            *float64         `protobuf:"fixed64,3,opt,name=float" json:"float,omitempty"`
	Str              *string          `protobuf:"bytes,4,opt,name=str" json:"str,omitempty"`
	Bytes            []byte           `protobuf:"bytes,5,opt,name=bytes" json:"bytes,omitempty"`
	Bool             *bool            `protobuf:"varint,6,opt,name=bool" json:"bool,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *DynamicType) Reset()                    { *m = DynamicType{} }
func (m *DynamicType) String() string            { return proto.CompactTextString(m) }
func (*DynamicType) ProtoMessage()               {}
func (*DynamicType) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DynamicType) GetType() DynamicTypeType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return DynamicTypeType_int
}

func (m *DynamicType) GetInt() int64 {
	if m != nil && m.Int != nil {
		return *m.Int
	}
	return 0
}

func (m *DynamicType) GetFloat() float64 {
	if m != nil && m.Float != nil {
		return *m.Float
	}
	return 0
}

func (m *DynamicType) GetStr() string {
	if m != nil && m.Str != nil {
		return *m.Str
	}
	return ""
}

func (m *DynamicType) GetBytes() []byte {
	if m != nil {
		return m.Bytes
	}
	return nil
}

func (m *DynamicType) GetBool() bool {
	if m != nil && m.Bool != nil {
		return *m.Bool
	}
	return false
}

type MessageProto struct {
	ReferenceID      *string                 `protobuf:"bytes,1,opt,name=referenceID" json:"referenceID,omitempty"`
	Command          *string                 `protobuf:"bytes,2,opt,name=Command" json:"Command,omitempty"`
	Done             *bool                   `protobuf:"varint,3,opt,name=done" json:"done,omitempty"`
	Arguments        map[string]*DynamicType `protobuf:"bytes,4,rep,name=arguments" json:"arguments,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *MessageProto) Reset()                    { *m = MessageProto{} }
func (m *MessageProto) String() string            { return proto.CompactTextString(m) }
func (*MessageProto) ProtoMessage()               {}
func (*MessageProto) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MessageProto) GetReferenceID() string {
	if m != nil && m.ReferenceID != nil {
		return *m.ReferenceID
	}
	return ""
}

func (m *MessageProto) GetCommand() string {
	if m != nil && m.Command != nil {
		return *m.Command
	}
	return ""
}

func (m *MessageProto) GetDone() bool {
	if m != nil && m.Done != nil {
		return *m.Done
	}
	return false
}

func (m *MessageProto) GetArguments() map[string]*DynamicType {
	if m != nil {
		return m.Arguments
	}
	return nil
}

func init() {
	proto.RegisterType((*DynamicType)(nil), "gridworker.DynamicType")
	proto.RegisterType((*MessageProto)(nil), "gridworker.MessageProto")
	proto.RegisterEnum("gridworker.DynamicTypeType", DynamicTypeType_name, DynamicTypeType_value)
}

func init() { proto.RegisterFile("message_proto.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 355 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x50, 0x4d, 0x8b, 0xdb, 0x30,
	0x14, 0xac, 0xfc, 0x91, 0xd8, 0xcf, 0x21, 0x15, 0x6a, 0xa1, 0xa2, 0xbd, 0x88, 0x5c, 0x6a, 0x0a,
	0x75, 0x21, 0xa7, 0xd2, 0x5b, 0x69, 0x72, 0xe8, 0xa1, 0x50, 0x44, 0x7b, 0x0e, 0x4e, 0xac, 0x18,
	0x37, 0xb6, 0x64, 0x24, 0xa5, 0x8b, 0xff, 0xce, 0xfe, 0xce, 0x3d, 0x2c, 0x92, 0x13, 0xe2, 0x5d,
	0xd8, 0x83, 0xcd, 0xcc, 0xd3, 0x48, 0x6f, 0x66, 0xe0, 0x4d, 0x27, 0x8c, 0x29, 0x6b, 0xb1, 0xeb,
	0xb5, 0xb2, 0xaa, 0xf0, 0x7f, 0x02, 0xb5, 0x6e, 0xaa, 0x3b, 0xa5, 0x4f, 0x42, 0xaf, 0xee, 0x11,
	0x64, 0x9b, 0x41, 0x96, 0x5d, 0x73, 0xf8, 0x33, 0xf4, 0x82, 0x7c, 0x81, 0xc8, 0x0e, 0xbd, 0xa0,
	0x88, 0x05, 0xf9, 0x72, 0xfd, 0xa1, 0xb8, 0x49, 0x8b, 0x89, 0xcc, 0x7d, 0xdc, 0x0b, 0x09, 0x86,
	0xb0, 0x91, 0x96, 0x06, 0x0c, 0xe5, 0x21, 0x77, 0x90, 0xbc, 0x85, 0xf8, 0xd8, 0xaa, 0xd2, 0xd2,
	0x90, 0xa1, 0x1c, 0xf1, 0x91, 0x38, 0x9d, 0xb1, 0x9a, 0x46, 0x0c, 0xe5, 0x29, 0x77, 0xd0, 0xe9,
	0xf6, 0x83, 0x15, 0x86, 0xc6, 0x0c, 0xe5, 0x0b, 0x3e, 0x12, 0x42, 0x20, 0xda, 0x2b, 0xd5, 0xd2,
	0x19, 0x43, 0x79, 0xc2, 0x3d, 0x5e, 0x3d, 0x20, 0x58, 0xfc, 0x1a, 0x83, 0xfc, 0xf6, 0x09, 0x18,
	0x64, 0x5a, 0x1c, 0x85, 0x16, 0xf2, 0x20, 0x7e, 0x6e, 0x28, 0xf2, 0x8f, 0x4e, 0x47, 0x84, 0xc2,
	0xfc, 0x87, 0xea, 0xba, 0x52, 0x56, 0xde, 0x5a, 0xca, 0xaf, 0xd4, 0x2d, 0xa8, 0x94, 0x14, 0xde,
	0x5d, 0xc2, 0x3d, 0x26, 0x5b, 0x48, 0x4b, 0x5d, 0x9f, 0x3b, 0x21, 0xad, 0xa1, 0x11, 0x0b, 0xf3,
	0x6c, 0xfd, 0x71, 0x1a, 0x7d, 0xba, 0xbc, 0xf8, 0x7e, 0x55, 0x6e, 0xa5, 0xd5, 0x03, 0xbf, 0xdd,
	0x7c, 0xff, 0x17, 0x96, 0x4f, 0x0f, 0x5d, 0xea, 0x93, 0x18, 0x2e, 0x06, 0x1d, 0x24, 0x9f, 0x21,
	0xfe, 0x5f, 0xb6, 0x67, 0xe1, 0x6d, 0x65, 0xeb, 0x77, 0x2f, 0x34, 0xcc, 0x47, 0xd5, 0xb7, 0xe0,
	0x2b, 0xfa, 0x54, 0xc1, 0xeb, 0x67, 0xdd, 0x93, 0xb9, 0x6f, 0x1d, 0xbf, 0x22, 0xe9, 0xa5, 0x6c,
	0x8c, 0x08, 0xc0, 0xcc, 0x58, 0xdd, 0xc8, 0x1a, 0x07, 0x6e, 0xec, 0xeb, 0xc4, 0x21, 0x49, 0xc6,
	0x42, 0x71, 0x44, 0x96, 0x00, 0xff, 0x8c, 0x92, 0x3b, 0xd3, 0x36, 0x07, 0x81, 0x63, 0xb2, 0x80,
	0xc4, 0xf3, 0xae, 0xec, 0xf1, 0xec, 0x31, 0x00, 0x00, 0xff, 0xff, 0xa0, 0x57, 0x06, 0x29, 0x2b,
	0x02, 0x00, 0x00,
}