#!/bin/bash

cd .. || exit 1
cd benchmarks || exit 1

go test -bench=BenchmarkGetUser -benchmem -count=20 > results/getUserResults.txt
benchstat results/getUserResults.txt

