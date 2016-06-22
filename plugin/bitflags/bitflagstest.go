package bitflags

import (
	"github.com/gogo/protobuf/plugin/testgen"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/opsee/protobuf/opseeproto"
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
	for _, message := range file.Messages() {
		if opseeproto.IsBitflags(file.FileDescriptorProto, message.DescriptorProto) {
			return true
		}
	}
	return false
}
