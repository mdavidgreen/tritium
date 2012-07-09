// Code generated by protoc-gen-go from "package.proto"
// DO NOT EDIT!

package proto

import proto1 "code.google.com/p/goprotobuf/proto"
import "math"

// Reference proto and math imports to suppress error if they are not otherwise used.
var _ = proto1.Marshal
var _ = math.Inf

type Type struct {
	Name             *string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Implements       *int32  `protobuf:"varint,2,opt,name=implements" json:"implements,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *Type) Reset()         { *this = Type{} }
func (this *Type) String() string { return proto1.CompactTextString(this) }
func (*Type) ProtoMessage()       {}

func (this *Type) GetName() string {
	if this != nil && this.Name != nil {
		return *this.Name
	}
	return ""
}

func (this *Type) GetImplements() int32 {
	if this != nil && this.Implements != nil {
		return *this.Implements
	}
	return 0
}

type Package struct {
	Name             *string     `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Functions        []*Function `protobuf:"bytes,2,rep,name=functions" json:"functions,omitempty"`
	Types            []*Type     `protobuf:"bytes,3,rep,name=types" json:"types,omitempty"`
	Dependencies     []string    `protobuf:"bytes,4,rep,name=dependencies" json:"dependencies,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (this *Package) Reset()         { *this = Package{} }
func (this *Package) String() string { return proto1.CompactTextString(this) }
func (*Package) ProtoMessage()       {}

func (this *Package) GetName() string {
	if this != nil && this.Name != nil {
		return *this.Name
	}
	return ""
}

func init() {
}