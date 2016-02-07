generate:
	protoc --gogoopsee_out=import_path=github.com/opsee/protobuf/examples/dessert,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:./dessert --proto_path=/gopath/src:./dessert dessert/*.proto
	protoc --gogoopsee_out=import_path=github.com/opsee/protobuf/examples,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:. --proto_path=/gopath/src:. *.proto

clean:
	rm -f *.pb.go
	rm -f *pb_test.go

.PHONY:
	clean
	generate
