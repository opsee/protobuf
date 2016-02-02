VANITY_DIR = gogogqlproto
EXAMPLE_DIR = examples
PROTO_DIR = proto

VANITY_DIR:
	docker run --rm -it -v $$(pwd):/build quay.io/opsee/build-go:go15 /bin/bash -c 'cd /build/$(VANITY_DIR) && make generate'

docker: clean VANITY_DIR
	docker build -t quay.io/opsee/build-go:gogoopsee .

PROTO_DIR:
	docker run --rm -it -v $$(pwd):/build/src quay.io/opsee/build-go:gogoopsee /bin/bash -c 'cd /build/src/$(PROTO_DIR) && make && export GOPATH="$$GOPATH:/build" && go get -t ./... && go test -v ./...'

EXAMPLE_DIR: docker PROTO_DIR
	docker run --rm -it -v $$(pwd):/build/src quay.io/opsee/build-go:gogoopsee /bin/bash -c 'cd /build/src/$(EXAMPLE_DIR) && make generate && export GOPATH="$$GOPATH:/build" && go get -t ./... && go test -v ./...'

push:
	docker push quay.io/opsee/build-go:gogoopsee

all: EXAMPLE_DIR

clean:
	$(MAKE) -C $(EXAMPLE_DIR) clean

.PHONY:
	docker
	clean
	push
	all