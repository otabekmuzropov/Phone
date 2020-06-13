#!/bin/bash
CURRENT_DIR=$1
for x in $(find ${CURRENT_DIR}/protos/* -type d); do
  protoc -I=${x} -I /usr/local/include --go_out=plugins=grpc:${CURRENT_DIR} ${x}/*.proto
done
