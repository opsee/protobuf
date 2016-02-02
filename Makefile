PROTO_DIR = gogogqlproto
EXAMPLE_DIR = examples

generate:
	docker run --rm -it -v $$(pwd):/build quay.io/opsee/build-go:go15 /bin/bash -c 'cd /build/$(PROTO_DIR) && make generate'

build: clean generate
	docker build -t quay.io/opsee/build-go:gogoopsee .

test: build
	docker run --rm -it -v $$(pwd):/build/src quay.io/opsee/build-go:gogoopsee /bin/bash -c 'cd /build/src/$(EXAMPLE_DIR) && make generate && export GOPATH="$$GOPATH:/build" && go get -t ./... && go test -v ./...'

push: test
	docker push quay.io/opsee/build-go:gogoopsee

clean:
	$(MAKE) -C $(EXAMPLE_DIR) clean

.PHONY:
	generate
	test
	build
	clean
	push