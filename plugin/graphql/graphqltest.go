package graphql

import (
	"fmt"
	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/opsee/protobuf/gogogqlproto"
	"strings"
)

type test struct {
	*generator.Generator
	generator.PluginImports
}

// func init() {
// 	testgen.RegisterTestPlugin(NewTest)
// }

func NewTest() *test {
	return &test{}
}

func (p *test) Name() string {
	return "graphqltest"
}

func (p *test) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *test) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)
	used := false

	if gogogqlproto.GetGraphQLFile(file.FileDescriptorProto) != true {
		return
	}

	testingPkg := p.NewImport("testing")
	randPkg := p.NewImport("math/rand")
	timePkg := p.NewImport("time")

	for mi, message := range file.Messages() {
		if message.DescriptorProto.GetOptions().GetMapEntry() {
			continue
		}

		messageGQL := strings.TrimSpace(p.Comments(fmt.Sprintf("4,%d", mi)))
		ccTypeName := generator.CamelCaseSlice(message.TypeName())

		if gogoproto.HasTestGen(file.FileDescriptorProto, message.DescriptorProto) {
			used = true
			p.P(`func Test`, ccTypeName, `GraphQL(t *`, testingPkg.Use(), `.T) {`)
			p.In()
			p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(`, timePkg.Use(), `.Now().UnixNano()))`)
			p.P(`_ = NewPopulated`, ccTypeName, `(popr, false)`)
			p.P(`objdesc := "`, messageGQL, `"`)
			p.P(`pdesc := `, graphQLTypeVarName(ccTypeName), `.PrivateDescription`)
			p.P(`if pdesc != objdesc {`)
			p.In()
			p.P(`t.Fatalf("String want %v got %v", objdesc, pdesc)`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
		}

	}

	if used {
		p.P(`//These tests are generated by github.com/opsee/protobuf/plugin/graphql`)
	}
}
