#!/bin/bash

cd .. || exit 1
cd benchmarks || exit 1

go test -bench=BenchmarkCreateUser -benchmem -count=20 > results/createUserResults.txt
benchstat results/createUserResults.txt
