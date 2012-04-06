package proto

import (
	pb "goprotobuf.googlecode.com/hg/proto"
)

func (obj *ScriptObject) Imports() []string {
	list := make([]string, 0)
	obj.Root.Iterate(func(ins *Instruction) {
		if *ins.Type == Instruction_IMPORT {
			list = append(list, pb.GetString(ins.Value))
		}
		//println(*ins.Type)
	})
	return list
}
