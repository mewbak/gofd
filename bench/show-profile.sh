#!/bin/bash
# in case you have run run-benchs.sh with the -p option,
# you can show the profile immediately afterwards with that script.
# if you have run e.g. 
# $ ./run-benchs.sh -p IvDomain Good.IvD.Copy
# then run now
# $ ./show-profile IvDomain


# We need to set manually the GOPATH, which is a couple of directories
# above starting from bench; what a hack....
export GOPATH=$(dirname $(dirname $(dirname $(dirname $(dirname $(pwd))))))
# now we may start go with a main and it will find the other packages

# configuration for benching
GOCMD=go
LOG_FOLD=logs
BIN_FOLD=bin

# go through options, accept -g <gocmd>
while getopts "g:" arg; do
case "$arg" in
        g) GOCMD="$OPTARG";;
        [?]) echo >&2 "Usage: $0 [-g gocmd (go)] bench_program"
            exit 1;;
    esac
done
shift $(($OPTIND-1))

$GOCMD tool pprof $BIN_FOLD/bench_$1 $LOG_FOLD/bench_$1.pprof

