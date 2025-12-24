#!/bin/bash

echo "Building tpsg in debug mode..."
cd tpsg
go build -o ../bins/tpsg_debug ./cmd/tpsg
cd ..
echo "Debug build complete: bins/tpsg_debug"
