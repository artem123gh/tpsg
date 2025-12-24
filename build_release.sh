#!/bin/bash

echo "Building tpsg in release mode..."
cd tpsg
go build -ldflags="-s -w" -o ../bins/tpsg_release ./cmd/tpsg
cd ..
echo "Release build complete: bins/tpsg_release"
