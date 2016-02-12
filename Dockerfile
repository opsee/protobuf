FROM quay.io/opsee/build-go:go15

RUN go get -u github.com/gogo/protobuf/proto && \
    go get -u github.com/gogo/protobuf/protoc-gen-gogo && \
    go get -u github.com/gogo/protobuf/gogoproto && \
    go get -u go.pedge.io/pb || true

COPY ./ /gopath/src/github.com/opsee/protobuf

RUN cd /gopath/src/github.com/opsee/protobuf && \
  go get ./opseeproto && \
  go get ./plugin/... && \
  go get ./protoc-gen-gogoopsee && \
  go install ./opseeproto && \
  go install ./plugin/... && \
  go install ./protoc-gen-gogoopsee
