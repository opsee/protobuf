package graphql

import (
	"fmt"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/opsee/protobuf/gogogqlproto"
	"strings"
)

type graphql struct {
	*generator.Generator
	generator.PluginImports
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
	// p.localName = generator.FileName(file)
	graphQLPkg := p.NewImport("github.com/graphql-go/graphql")
	fmtPkg := p.NewImport("fmt")

	for _, message := range file.Messages() {
		messageGQL := gogogqlproto.GetGraphQLMessage(message.DescriptorProto)

		if messageGQL == nil {
			continue
		}
		if message.DescriptorProto.GetOptions().GetMapEntry() {
			continue
		}
		if len(message.DescriptorProto.Field) == 0 {
			continue
		}

		ccTypeName := generator.CamelCaseSlice(message.TypeName())
		typeName := gogogqlproto.SnakeCase(strings.Join(message.TypeName(), "_"))
		// p.P(`func New`, ccTypeName, `GraphQLObject() *`, graphQLPkg.Use(), `.Object {`)
		// p.In()
		p.P(`var `, gogogqlproto.GraphQLTypeVarName(ccTypeName), ` = `, graphQLPkg.Use(), `.NewObject(`, graphQLPkg.Use(), `.ObjectConfig{`)
		p.In()
		p.P(`Name:        "`, typeName, `",`)
		p.P(`Description: "`, *messageGQL, `",`)
		p.P(`Fields: (`, graphQLPkg.Use(), `.FieldsThunk)(func() `, graphQLPkg.Use(), `.Fields {`)
		p.In()
		p.P(`return `, graphQLPkg.Use(), `.Fields{`)
		p.In()
		for _, field := range message.DescriptorProto.Field {
			// get type
			// get non null
			// get description
			// required := field.IsRequired()
			// repeated := field.IsRepeated()

			p.P(`"`, field.GetName(), `": &`, graphQLPkg.Use(), `.Field{`)
			p.In()
			p.P(`Type:        `, p.graphQLType(message, field, graphQLPkg.Use()), `,`)
			p.P(`Description: "foo field description",`)
			p.P(`Resolve: func(p `, graphQLPkg.Use(), `.ResolveParams) (interface{}, error) {`)
			p.In()
			p.P(`obj, ok := p.Source.(*`, ccTypeName, `)`)
			p.P(`if !ok {`)
			p.In()
			p.P(`return nil, `, fmtPkg.Use(), `.Errorf("field `, field.GetName(), ` not resolved")`)
			p.Out()
			p.P(`}`)
			p.P(`return obj.`, p.GetFieldName(message, field), `, nil`)
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
		// p.Out()
		// p.P(`}`)
	}
}

func (p *graphql) graphQLType(message *generator.Descriptor, field *descriptor.FieldDescriptorProto, pkgName string) string {
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
		mobj := p.ObjectNamed(field.GetTypeName())
		if mobj.PackageName() != message.PackageName() {
			gqltype = fmt.Sprint(pkgName, ".", "String")
			break
		}
		gqltype = gogogqlproto.GraphQLTypeVarName(p.TypeName(mobj))
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		panic("GETOUT")
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
