package gogogqlproto

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

func GetGraphQLField(field *descriptor.FieldDescriptorProto) *string {
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, E_GraphqlField)
		if err == nil && v.(*string) != nil {
			return (v.(*string))
		}
	}
	return nil
}

func GetGraphQLMessage(message *descriptor.DescriptorProto) *string {
	if message.Options != nil {
		v, err := proto.GetExtension(message.Options, E_GraphqlMessage)
		if err == nil && v.(*string) != nil {
			return (v.(*string))
		}
	}
	return nil
}
