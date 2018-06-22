// Code generated by protoc-gen-go. DO NOT EDIT.
// source: wordtype/wordtype.proto

package wordtype // import "github.com/marthjod/binquiry-experimental/wordtype"

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

type WordType int32

const (
	WordType_Noun      WordType = 0
	WordType_Adjective WordType = 1
	WordType_Verb      WordType = 2
	WordType_Unknown   WordType = 3
)

var WordType_name = map[int32]string{
	0: "Noun",
	1: "Adjective",
	2: "Verb",
	3: "Unknown",
}
var WordType_value = map[string]int32{
	"Noun":      0,
	"Adjective": 1,
	"Verb":      2,
	"Unknown":   3,
}

func (x WordType) String() string {
	return proto.EnumName(WordType_name, int32(x))
}
func (WordType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_wordtype_2ec6d2cb36a47b26, []int{0}
}

func init() {
	proto.RegisterEnum("wordtype.WordType", WordType_name, WordType_value)
}

func init() { proto.RegisterFile("wordtype/wordtype.proto", fileDescriptor_wordtype_2ec6d2cb36a47b26) }

var fileDescriptor_wordtype_2ec6d2cb36a47b26 = []byte{
	// 157 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2f, 0xcf, 0x2f, 0x4a,
	0x29, 0xa9, 0x2c, 0x48, 0xd5, 0x87, 0x31, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x38, 0x60,
	0x7c, 0x2d, 0x2b, 0x2e, 0x8e, 0xf0, 0xfc, 0xa2, 0x94, 0x90, 0xca, 0x82, 0x54, 0x21, 0x0e, 0x2e,
	0x16, 0xbf, 0xfc, 0xd2, 0x3c, 0x01, 0x06, 0x21, 0x5e, 0x2e, 0x4e, 0xc7, 0x94, 0xac, 0xd4, 0xe4,
	0x92, 0xcc, 0xb2, 0x54, 0x01, 0x46, 0x90, 0x44, 0x58, 0x6a, 0x51, 0x92, 0x00, 0x93, 0x10, 0x37,
	0x17, 0x7b, 0x68, 0x5e, 0x76, 0x5e, 0x7e, 0x79, 0x9e, 0x00, 0xb3, 0x93, 0x49, 0x94, 0x51, 0x7a,
	0x66, 0x49, 0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0xae, 0x7e, 0x6e, 0x62, 0x51, 0x49, 0x46, 0x56,
	0x7e, 0x8a, 0x7e, 0x52, 0x66, 0x5e, 0x61, 0x69, 0x66, 0x51, 0xa5, 0x6e, 0x6a, 0x45, 0x41, 0x6a,
	0x51, 0x66, 0x6e, 0x6a, 0x5e, 0x49, 0x62, 0x0e, 0xdc, 0x05, 0x49, 0x6c, 0x60, 0x27, 0x18, 0x03,
	0x02, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x46, 0x8f, 0x14, 0x9d, 0x00, 0x00, 0x00,
}
