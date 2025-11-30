#!/bin/sh
set -e

just build
./build/out/nnc "$@"
