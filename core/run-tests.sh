#!/bin/sh
outfile=gotest.out
go test -v | tee $outfile
go2xunit -fail -input $outfile -output tests.xml
rm -f $outfile
