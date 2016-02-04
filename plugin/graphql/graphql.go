package graphql

import (
	"bytes"
	"fmt"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/opsee/protobuf/gogogqlproto"
	"strings"
	"unicode"
)

type graphql struct {
	*generator.Generator
	generator.PluginImports
	messages []*generator.Descriptor
	oneofs   map[*descriptor.OneofDescriptorProto]*oneof
}

type oneof struct {
	message      *generator.Descriptor
	fields       []*descriptor.FieldDescriptorProto
	messageIndex int
	oneofIndex   int
}

func init() {
	generator.RegisterPlugin(NewGraphQL())
}

func NewGraphQL() *graphql {
	return &graphql{}
}

func (p *graphql) Name() string {
	return "graphql"
}

func (p *graphql) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *graphql) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)
	p.messages = make([]*generator.Descriptor, 0)
	p.oneofs = make(map[*descriptor.OneofDescriptorProto]*oneof)

	if gogogqlproto.GetGraphQLFile(file.FileDescriptorProto) != true {
		return
	}

	graphQLPkg := p.NewImport("github.com/graphql-go/graphql")
	schemaPkg := p.NewImport("github.com/opsee/protobuf/gogogqlproto")
	fmtPkg := p.NewImport("fmt")

	for mi, message := range file.Messages() {
		if message.DescriptorProto.GetOptions().GetMapEntry() {
			continue
		}

		if len(message.DescriptorProto.Field) == 0 {
			continue
		}

		p.messages = append(p.messages, message)

		// generate the var declarations first
		ccTypeName := generator.CamelCaseSlice(message.TypeName())
		p.P(`var `, graphQLTypeVarName(ccTypeName), ` *`, graphQLPkg.Use(), `.Object`)

		for i, field := range message.DescriptorProto.OneofDecl {
			p.P(`var `, graphQLUnionVarName(message, field), ` *`, graphQLPkg.Use(), `.Union`)

			// collect the unions to make them easier to access in the file
			p.oneofs[field] = oneofFields(message, mi, i)
		}
	}

	p.P()
	p.P(`func init() {`)
	p.In()

	for mi, message := range p.messages {
		messageGQL := strings.TrimSpace(p.Comments(fmt.Sprintf("4,%d", mi)))
		ccTypeName := generator.CamelCaseSlice(message.TypeName())
		typeName := snakeCase(strings.Join(message.TypeName(), "_"))

		p.P(graphQLTypeVarName(ccTypeName), ` = `, graphQLPkg.Use(), `.NewObject(`, graphQLPkg.Use(), `.ObjectConfig{`)
		p.In()
		p.P(`Name:        "`, typeName, `",`)
		p.P(`Description: "`, messageGQL, `",`)
		p.P(`Fields: (`, graphQLPkg.Use(), `.FieldsThunk)(func() `, graphQLPkg.Use(), `.Fields {`)
		p.In()
		p.P(`return `, graphQLPkg.Use(), `.Fields{`)
		p.In()
		for fi, field := range message.DescriptorProto.Field {
			// skip defining a regular object field for unions, that comes next
			if field.OneofIndex != nil {
				continue
			}

			fieldGQL := strings.TrimSpace(p.Comments(fmt.Sprintf("4,%d,2,%d", mi, fi)))

			p.P(`"`, field.GetName(), `": &`, graphQLPkg.Use(), `.Field{`)
			p.In()
			p.P(`Type:        `, p.graphQLType(message, field, graphQLPkg.Use(), schemaPkg.Use()), `,`)
			p.P(`Description: "`, fieldGQL, `",`)
			p.P(`Resolve: func(p `, graphQLPkg.Use(), `.ResolveParams) (interface{}, error) {`)
			p.In()
			p.P(`switch obj := p.Source.(type) {`)
			p.P(`case *`, ccTypeName, `:`)
			p.In()
			p.P(`return obj.`, p.GetFieldName(message, field), `, nil`)
			p.Out()
			for _, oo := range p.oneofs {
				for _, oneof := range oo.fields {
					tname := p.TypeName(p.ObjectNamed(oneof.GetTypeName()))
					if tname == ccTypeName {
						p.P(`case *`, generator.CamelCaseSlice(oo.message.TypeName()), `_`, tname, `:`)
						p.In()
						p.P(`return obj.`, tname, `.`, p.GetFieldName(message, field), `, nil`)
						p.Out()
					}
				}
			}
			p.P(`}`)
			p.P(`return nil, `, fmtPkg.Use(), `.Errorf("field `, field.GetName(), ` not resolved")`)
			p.Out()
			p.P(`},`)
			p.Out()
			p.P(`},`)
		}
		for fi, field := range message.DescriptorProto.OneofDecl {
			fieldGQL := strings.TrimSpace(p.Comments(fmt.Sprintf("4,%d,8,%d", mi, fi)))

			p.P(`"`, field.GetName(), `": &`, graphQLPkg.Use(), `.Field{`)
			p.In()
			p.P(`Type:        `, graphQLUnionVarName(message, field), `,`)
			p.P(`Description: "`, fieldGQL, `",`)
			p.P(`Resolve: func(p `, graphQLPkg.Use(), `.ResolveParams) (interface{}, error) {`)
			p.In()
			p.P(`obj, ok := p.Source.(*`, ccTypeName, `)`)
			p.P(`if !ok {`)
			p.In()
			p.P(`return nil, `, fmtPkg.Use(), `.Errorf("field `, field.GetName(), ` not resolved")`)
			p.Out()
			p.P(`}`)
			p.P(`return obj.Get`, generator.CamelCase(field.GetName()), `(), nil`)
			p.Out()
			p.P(`},`)
			p.Out()
			p.P(`},`)

		}
		p.Out()
		p.P(`}`)
		p.Out()
		p.P(`}),`)
		p.Out()
		p.P(`})`)
	}

	// declare our unions last, since the types will have needed to be defined from all messages first
	for decl, oo := range p.oneofs {
		ccTypeName := generator.CamelCaseSlice(oo.message.TypeName())
		fieldGQL := strings.TrimSpace(p.Comments(fmt.Sprintf("4,%d,8,%d", oo.messageIndex, oo.oneofIndex)))

		p.P(graphQLUnionVarName(oo.message, decl), ` = `, graphQLPkg.Use(), `.NewUnion(`, graphQLPkg.Use(), `.UnionConfig{`)
		p.In()
		p.P(`Name:        "`, graphQLUnionName(oo.message, decl), `",`)
		p.P(`Description: "`, fieldGQL, `",`)
		p.P(`Types:       []*`, graphQLPkg.Use(), `.Object{`)
		p.In()
		for _, field := range oo.fields {
			p.P(graphQLTypeVarName(p.TypeName(p.ObjectNamed(field.GetTypeName()))), `,`)
		}
		p.Out()
		p.P(`},`)
		p.P(`ResolveType: func (value interface{}, info `, graphQLPkg.Use(), `.ResolveInfo) *`, graphQLPkg.Use(), `.Object {`)
		p.In()
		p.P(`switch value.(type) {`)
		for _, field := range oo.fields {
			tname := p.TypeName(p.ObjectNamed(field.GetTypeName()))
			p.P(`case *`, ccTypeName, `_`, tname, `:`)
			p.In()
			p.P(`return `, graphQLTypeVarName(tname))
			p.Out()
		}
		p.P(`}`)
		p.P(`return nil`)
		p.Out()
		p.P(`},`)
		p.Out()
		p.P(`})`)
	}

	p.Out()
	p.P(`}`)
}

func (p *graphql) graphQLType(message *generator.Descriptor, field *descriptor.FieldDescriptorProto, pkgName, schemaPkgName string) string {
	var gqltype string
	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE, descriptor.FieldDescriptorProto_TYPE_FLOAT:
		gqltype = fmt.Sprint(pkgName, ".", "Float")
	case descriptor.FieldDescriptorProto_TYPE_INT64, descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_INT32, descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_FIXED32, descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64, descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		gqltype = fmt.Sprint(pkgName, ".", "Int")
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		gqltype = fmt.Sprint(pkgName, ".", "Boolean")
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		gqltype = fmt.Sprint(pkgName, ".", "String")
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		panic("mapping a proto group type to graphql is unimplemented")
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		panic("mapping a proto enum type to graphql is unimplemented")
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		// TODO: fix this to be more robust about imported objects
		mobj := p.ObjectNamed(field.GetTypeName())
		if mobj.PackageName() != message.PackageName() {
			if field.GetTypeName() == "Timestamp" {
				gqltype = fmt.Sprint(schemaPkgName, ".", "Timestamp")
				break
			}

			gqltype = fmt.Sprint(schemaPkgName, ".", "ByteString")
			break
		}
		gqltype = graphQLTypeVarName(p.TypeName(mobj))
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		gqltype = fmt.Sprint(schemaPkgName, ".", "ByteString")
	default:
		panic("unknown proto field type")
	}

	if field.IsRepeated() {
		gqltype = fmt.Sprint(pkgName, ".NewList(", gqltype, ")")
	}

	if field.IsRequired() {
		gqltype = fmt.Sprint(pkgName, ".NewNonNull(", gqltype, ")")
	}

	return gqltype
}

func oneofFields(message *generator.Descriptor, messageIndex, oneofIndex int) *oneof {
	fields := make([]*descriptor.FieldDescriptorProto, 0)

	for _, field := range message.DescriptorProto.Field {
		if field.OneofIndex != nil && *field.OneofIndex == int32(oneofIndex) {
			fields = append(fields, field)
		}
	}

	return &oneof{message, fields, messageIndex, oneofIndex}
}

func graphQLTypeVarName(typeName string) string {
	return fmt.Sprint("GraphQL", typeName, "Type")
}

func graphQLUnionName(message *generator.Descriptor, oneof *descriptor.OneofDescriptorProto) string {
	return generator.CamelCaseSlice(message.TypeName()) + generator.CamelCase(oneof.GetName())
}

func graphQLUnionVarName(message *generator.Descriptor, oneof *descriptor.OneofDescriptorProto) string {
	return fmt.Sprint("GraphQL", graphQLUnionName(message, oneof), "Union")
}

func snakeCase(in string) string {
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
