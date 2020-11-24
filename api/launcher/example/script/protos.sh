#!/bin/bash

ORIGINDIR=$(pwd)

cd protos

docker run --rm -v $(pwd):$(pwd) -w $(pwd) znly/protoc -I . --go_out=plugins=grpc:. *.proto

cd ${ORIGINDIR}
