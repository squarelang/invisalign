#!/bin/bash

go build -o bin/invisalign
mkdir -p sandbox
cp *.go sandbox/
ls | grep .ivs | xargs -n 1 sh -c 'bin/invisalign $0 > sandbox/`echo $0 | sed -e "s/.ivs/.go/g"`'
ln -s ../test_cases sandbox/test_cases
go test sandbox/*.go
rm -rf sandbox
