#!/bin/bash

echo "Building tpsg in debug mode..."
cd tpsg
go build -o ../bins/tpsg_debug
cd ..
echo "Debug build complete: bins/tpsg_debug"
