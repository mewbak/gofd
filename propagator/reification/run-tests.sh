#!/bin/sh
outfile=gotest.out
echo $PWD
export GOPATH=$(dirname $(dirname $(dirname $(dirname $(dirname $(dirname $PWD))))))
echo $GOPATH
go test -v | tee $outfile
go2xunit -fail -input $outfile -output tests.xml
rm -f $outfile
