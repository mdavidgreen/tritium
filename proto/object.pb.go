// Code generated by protoc-gen-go.
// source: object.proto
// DO NOT EDIT!

package proto

import proto1 "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto1.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type ScriptObject struct {
	Name             *string      `protobuf:"bytes,1,opt,name=name,def=main" json:"name,omitempty"`
	Root             *Instruction `protobuf:"bytes,2,opt,name=root" json:"root,omitempty"`
	Functions        []*Function  `protobuf:"bytes,3,rep,name=functions" json:"functions,omitempty"`
	ScopeTypeId      *int32       `protobuf:"varint,4,opt,name=scope_type_id" json:"scope_type_id,omitempty"`
	Linked           *bool        `protobuf:"varint,5,opt,name=linked" json:"linked,omitempty"`
	Module           *string      `protobuf:"bytes,6,opt,name=module" json:"module,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *ScriptObject) Reset()         { *m = ScriptObject{} }
func (m *ScriptObject) String() string { return proto1.CompactTextString(m) }
func (*ScriptObject) ProtoMessage()    {}

const Default_ScriptObject_Name string = "main"

func (m *ScriptObject) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return Default_ScriptObject_Name
}

func (m *ScriptObject) GetRoot() *Instruction {
	if m != nil {
		return m.Root
	}
	return nil
}

func (m *ScriptObject) GetFunctions() []*Function {
	if m != nil {
		return m.Functions
	}
	return nil
}

func (m *ScriptObject) GetScopeTypeId() int32 {
	if m != nil && m.ScopeTypeId != nil {
		return *m.ScopeTypeId
	}
	return 0
}

func (m *ScriptObject) GetLinked() bool {
	if m != nil && m.Linked != nil {
		return *m.Linked
	}
	return false
}

func (m *ScriptObject) GetModule() string {
	if m != nil && m.Module != nil {
		return *m.Module
	}
	return ""
}

func init() {
}
