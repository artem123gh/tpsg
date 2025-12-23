#!/bin/bash
cd "$(dirname "$0")/tpsg"
go run . test-gkvs
