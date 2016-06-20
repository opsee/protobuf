package bitflags

import (
	"bytes"
	"unicode"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/opsee/protobuf/opseeproto"
)

type plugin struct {
	*generator.Generator
	generator.PluginImports
	messages []*generator.Descriptor
}

func NewBitflags() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return "bitflags"
}

func (p *plugin) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *plugin) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)
	p.messages = make([]*generator.Descriptor, 0)
	driverPkg := p.NewImport("database/sql/driver")

	for _, message := range file.Messages() {
		if !opseeproto.IsBitflags(file.FileDescriptorProto, message.DescriptorProto) {
			continue
		}
		p.messages = append(p.messages, message)
		baseTypeName := generator.CamelCaseSlice(message.TypeName())

		// UInt64()
		// returns a bitflags uint64 representation of the structure
		p.P(`func (this *`, baseTypeName, `) UInt64() uint64 {`)
		p.In()
		p.P(`b := uint64(0)`)
		for bit, field := range message.Field {
			fieldname := p.GetFieldName(message, field)

			p.P(`if this.`, fieldname, ` {`)
			p.In()
			p.P(`b |= uint64(1) << uint64(`, bit, `)`)
			p.Out()
			p.P(`}`)
		}
		p.P()
		p.P(`return b`)
		p.P(`}`)

		// HighFlags() returns fields in struct set to 1
		p.P(`func (this *`, baseTypeName, `) HighFlags() []string {`)
		p.In()
		p.P(`var b []string`)
		for _, field := range message.Field {
			fieldname := p.GetFieldName(message, field)

			p.P(`if this.`, fieldname, ` {`)
			p.In()
			p.P(`b = append(b, "`, snakeCase(fieldname), `")`)
			p.Out()
			p.P(`}`)
		}
		p.P(`return b`)
		p.P(`}`)
		p.P()

		// LowFlags() returns fields in struct set to 0
		p.P(`func (this *`, baseTypeName, `) LowFlags() []string {`)
		p.In()
		p.P(`var b []string`)
		for _, field := range message.Field {
			fieldname := p.GetFieldName(message, field)

			p.P(`if !this.`, fieldname, ` {`)
			p.In()
			p.P(`b = append(b, "`, snakeCase(fieldname), `")`)
			p.Out()
			p.P(`}`)
		}
		p.P(`return b`)
		p.P(`}`)
		p.P()

		// Sets a flag or returns error
		p.P(`func (this *`, baseTypeName, `) SetFlag(flag string) error {`)
		p.In()
		p.P(`switch flag {`)
		p.In()
		for _, field := range message.Field {
			fieldname := p.GetFieldName(message, field)
			p.P(`case "`, snakeCase(fieldname), `":`)
			p.In()
			p.P(`this.`, fieldname, `= true `)
			p.Out()
		}
		p.P(`default:`)
		p.In()
		p.P(`return fmt.Errorf("invalid flag: %v", flag)`)
		p.Out()
		p.P(`}`)
		p.Out()
		p.P(`return nil`)
		p.P(`}`)

		// Sets a flag or returns error
		p.P(`func (this *`, baseTypeName, `) ClearFlag(flag string) error {`)
		p.In()
		p.P(`switch flag {`)
		p.In()
		for _, field := range message.Field {
			fieldname := p.GetFieldName(message, field)
			p.P(`case "`, snakeCase(fieldname), `":`)
			p.In()
			p.P(`this.`, fieldname, `= false `)
			p.Out()
		}
		p.P(`default:`)
		p.In()
		p.P(`return fmt.Errorf("invalid flag: %v", flag)`)
		p.Out()
		p.P(`}`)
		p.Out()
		p.P(`return nil`)
		p.P(`}`)

		// SetFlags(...string) []error
		// sets a number of flags which correspond to fields in the struct
		p.P(`func (this *`, baseTypeName, `) SetFlags(flags ...string) []error {`)
		p.In()
		p.P(`var errs []error`)
		p.P(`for _, f := range flags {`)
		p.In()
		p.P(`if err := this.SetFlag(f); err != nil {`)
		p.In()
		p.P(`errs = append(errs, err)`)
		p.Out()
		p.P(`}`)
		p.Out()
		p.P(`}`)
		p.P(`return errs`)
		p.Out()
		p.P(`}`)

		// ClearFlags(...string) []error
		// sets a number of flags which correspond to fields in the struct
		p.P(`func (this *`, baseTypeName, `) ClearFlags(flags ...string) []error {`)
		p.In()
		p.P(`var errs []error`)
		p.P(`for _, f := range flags {`)
		p.In()
		p.P(`if err := this.ClearFlag(f); err != nil {`)
		p.In()
		p.P(`errs = append(errs, err)`)
		p.Out()
		p.P(`}`)
		p.Out()
		p.P(`}`)
		p.P(`return errs`)
		p.Out()
		p.P(`}`)

		// Returns the value of a flag or false if the flag does not exist
		p.P(`func (this *`, baseTypeName, `) TestFlag(flag string) bool {`)
		p.In()
		p.P(`switch flag {`)
		p.In()
		for _, field := range message.Field {
			fieldname := p.GetFieldName(message, field)
			p.P(`case "`, snakeCase(fieldname), `":`)
			p.In()
			p.P(`return this.`, fieldname)
			p.Out()
		}
		p.P(`}`)
		p.Out()
		p.P(`return false`)
		p.P(`}`)

		// TestFlags(...string) []error
		// Returns Flag1 AND Flag2 AND ...
		p.P(`func (this *`, baseTypeName, `) TestFlags(flags ...string) bool {`)
		p.In()
		p.P(`for _, f := range flags {`)
		p.In()
		p.P(`if !this.TestFlag(f) {`)
		p.In()
		p.P(`return false`)
		p.Out()
		p.P(`}`)
		p.Out()
		p.P(`}`)
		p.P(`return true`)
		p.Out()
		p.P(`}`)

		// FromInt64(b uint64)
		// TODO(dan) should return error if overflow
		p.P(`func (this *`, baseTypeName, `) FromUInt64(b uint64) error {`)
		p.In()
		p.P(`bb := b`)
		for i, field := range message.Field {
			p.P(`bb = b`)
			fieldname := p.GetFieldName(message, field)
			p.P(`if bb&(uint64(1)<<uint(`, i, `)) > 0 {`)
			p.In()
			p.P(`this.`, fieldname, ` = true`)
			p.Out()
			p.P(`} else {`)
			p.In()
			p.P(`this.`, fieldname, ` = false`)
			p.P(`}`)
			p.Out()
		}
		p.P()
		p.P(`return nil`)
		p.P(`}`)

		// returns a bitflags uint64 representation of the structure from database
		p.P(`func (this *`, baseTypeName, `) Scan(i interface{}) error {`)
		p.In()
		p.P(`switch v := i.(type) {`)
		types := []string{"int", "int32", "int64", "float32", "float64"}
		for _, t := range types {
			p.P(`case `, t, `:`)
			p.In()
			p.P(`return this.FromUInt64(uint64(v))`)
			p.Out()
		}
		p.P(`}`)
		p.P()
		p.P(`return fmt.Errorf("invalid type: %T", i)`)
		p.P(`}`)

		// Represent in the database as int64()
		// allows the type to be encoded by database
		p.P(`func (this *`, baseTypeName, `) Value() (`, driverPkg.Use(), `.Value, error) {`)
		p.In()
		p.P(`return int64(this.UInt64()), nil`)
		p.P(`}`)

	}
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
func init() {
	generator.RegisterPlugin(NewBitflags())
}
