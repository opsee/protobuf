VANITY_DIR = opseeproto
EXAMPLE_DIR = examples
PROTO_DIR = proto

all: EXAMPLE_DIR

VANITY_DIR:
	docker run --rm -it -v $$(pwd):/build quay.io/opsee/build-go:go15 /bin/bash -c 'cd /build/$(VANITY_DIR) && make generate'

docker: clean VANITY_DIR
	docker build -t quay.io/opsee/build-go:gogoopsee .

EXAMPLE_DIR: docker
	docker run --rm -it -v $$(pwd):/gopath/src/github.com/opsee/protobuf quay.io/opsee/build-go:gogoopsee /bin/bash -c 'cd /gopath/src/github.com/opsee/protobuf/$(EXAMPLE_DIR) && make generate && go get -t ./... && go test -v ./...'

push:
	docker push quay.io/opsee/build-go:gogoopsee

clean:
	$(MAKE) -C $(EXAMPLE_DIR) clean

.PHONY:
	docker
	clean
	push
	all