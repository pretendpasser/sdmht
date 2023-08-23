#!/usr/bin/env sh

# Install proto3 from source
#  brew install autoconf automake libtool
#  git clone https://github.com/google/protobuf
#  ./autogen.sh ; ./configure ; make ; make install
#

# see https://grpc.io/docs/languages/go/quickstart/
# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
# export PATH="$PATH:$(go env GOPATH)/bin"

# See also
#  https://github.com/grpc/grpc-go/tree/master/examples

#protoc usersvc.proto --go_out=plugins=grpc:.
protoc -I. -I../../../.. conn.proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:.
