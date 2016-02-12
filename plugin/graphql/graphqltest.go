package graphql

import (
	"fmt"
	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/plugin/testgen"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/opsee/protobuf/opseeproto"
	"strings"
)

type test struct {
	*generator.Generator
}

func init() {
	testgen.RegisterTestPlugin(NewTest)
}

func NewTest(g *generator.Generator) testgen.TestPlugin {
	return &test{g}
}

func (p *test) Generate(imports generator.PluginImports, file *generator.FileDescriptor) bool {
	used := false

	if opseeproto.GetGraphQLFile(file.FileDescriptorProto) != true {
		return used
	}

	testingPkg := imports.NewImport("testing")
	randPkg := imports.NewImport("math/rand")
	timePkg := imports.NewImport("time")

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
			p.P(`pdesc := `, p.graphQLTypeVarName(message), `.PrivateDescription`)
			p.P(`if pdesc != objdesc {`)
			p.In()
			p.P(`t.Fatalf("String want %v got %v", objdesc, pdesc)`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
		}

	}
	return used
}

func (p *test) graphQLTypeVarName(obj generator.Object) string {
	return fmt.Sprint(p.DefaultPackageName(obj), "GraphQL", generator.CamelCaseSlice(obj.TypeName()), "Type")
}
