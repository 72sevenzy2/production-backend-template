#!/bin/bash

cd .. || exit 1
cd benchmarks || exit 1

go test -bench=BenchmarkEditUser -benchmem -count=20 > results/editUserResults.txt
benchstat results/editUserResults.txt

