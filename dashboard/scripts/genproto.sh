#!/bin/bash

set -e

PROTO_PATH=../modlib/internal/proto
PROTOC_GEN_TS_PATH="$(yarn bin protoc-gen-ts)"
OUT_DIR=./src/protobuf

mkdir -p "$OUT_DIR"

protoc \
    --plugin=protoc-gen-ts="$PROTOC_GEN_TS_PATH" \
    --js_out=import_style=commonjs,binary:"$OUT_DIR" \
    --ts_out="$OUT_DIR" \
    --proto_path="$PROTO_PATH" \
    feed.proto
