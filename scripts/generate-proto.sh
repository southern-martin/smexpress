#!/bin/bash
set -euo pipefail

PROTO_DIR="api/proto"
OUT_DIR="pkg/proto"

echo "Generating protobuf Go code..."

for proto_file in "$PROTO_DIR"/*.proto; do
    if [ -f "$proto_file" ]; then
        echo "  Processing $(basename "$proto_file")..."
        protoc \
            --go_out="$OUT_DIR" --go_opt=paths=source_relative \
            --go-grpc_out="$OUT_DIR" --go-grpc_opt=paths=source_relative \
            -I "$PROTO_DIR" \
            "$proto_file"
    fi
done

echo "Protobuf generation complete."
