#!/usr/bin/env bash

DIR_BASENAME=$(basename `pwd`)
go test -c -o "${DIR_BASENAME}.test"

"./${DIR_BASENAME}.test" -test.run=XXX -test.bench=. -test.cpuprofile=prof.out
