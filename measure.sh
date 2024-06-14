#!/bin/sh

set -e

. script/prepare_smc.sh

go run ./script/throughput_benchmark
