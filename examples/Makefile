generate:
	protoc --gogoopsee_out=plugins=grpc+graphql+bitflags,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:flags --proto_path=/gopath/src:flags flags/*.proto
	protoc --gogoopsee_out=plugins=grpc+graphql,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:dessert --proto_path=/gopath/src:dessert dessert/*.proto
	protoc --gogoopsee_out=plugins=grpc+graphql,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:. --proto_path=/gopath/src:. *.proto

clean:
	rm -f {,**/}*.pb.go
	rm -f {,**/}*pb_test.go

.PHONY:
	clean
	generate
