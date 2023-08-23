#!/bin/bash

protoc -I. -I../../../.. account.proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:.