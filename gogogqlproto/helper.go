package gogogqlproto

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

func GetGraphQLFile(file *descriptor.FileDescriptorProto) bool {
	if file.Options != nil {
		v, err := proto.GetExtension(file.Options, E_Graphql)
		if err == nil && v.(*bool) != nil {
			return (*v.(*bool))
		}
	}
	return false
}
