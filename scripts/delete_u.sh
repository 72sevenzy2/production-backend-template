#!/bin/bash

cd .. || exit 1
cd benchmarks || exit 1

go test -bench=BenchmarkDeleteUser -benchmem -count=20 > results/deleteUserResults.txt
benchstat results/deleteUserResults.txt

