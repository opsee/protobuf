FROM quay.io/opsee/build-go:16

COPY ./ /gopath/src/github.com/opsee/protobuf

RUN cd /gopath/src/github.com/opsee/protobuf && \
  go get ./opseeproto && \
  go get ./plugin/... && \
  go get ./protoc-gen-gogoopsee && \
  go install ./opseeproto && \
  go install ./plugin/... && \
  go install ./protoc-gen-gogoopsee
