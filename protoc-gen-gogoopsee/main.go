package main

import (
	"github.com/gogo/protobuf/vanity"
	"github.com/gogo/protobuf/vanity/command"
	_ "github.com/opsee/protobuf/plugin/graphql"
)

func main() {
	req := command.Read()
	files := req.GetProtoFile()
	vanity.ForEachFile(files, vanity.TurnOnTestGenAll)
	vanity.ForEachFile(files, vanity.TurnOnEqualAll)
	vanity.ForEachFile(files, vanity.TurnOnPopulateAll)
	resp := command.Generate(req)
	command.Write(resp)
}
