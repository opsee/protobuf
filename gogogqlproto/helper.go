package gogogqlproto

import (
	"bytes"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"unicode"
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

func GraphQLTypeVarName(typeName string) string {
	return fmt.Sprint("GraphQL", typeName, "Type")
}

func SnakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)
	out := bytes.NewBuffer(make([]byte, 0, length))

	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out.WriteRune('_')
		}
		out.WriteRune(unicode.ToLower(runes[i]))
	}

	return out.String()
}
