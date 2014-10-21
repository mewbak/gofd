#!/bin/bash
# run benchmark programs
# run from within the directory with ./run-benchs.sh
# Benchmark programs have the filename bench_<name>.go
# and are in the main package. The only other package 
# the need for compilation is benchframework.go (of course
# anything else can be imported and used). 
# For each run all benchmark programs are compiled to an 
# executable under BIN_FOLD and executed. 
# It is possible to restrict the compiled benchmark programs to the
# ones that contain the first command line parameter.
# It is possible to restrict the executed benchmark functions to the 
# ones whose name contains the second command line parameter. 
# Try to choose names that can easily be selected.

# Examples: 
# run all benchmarks
# $ ./run-benchs.sh 
# only run the benchmarks in bench_IntVar.go
# $ ./run-benchs.sh IntVar
# only run the benchmarks in bench_IvDomain.go and 
# of these only the Bad.ExD ones
# clean first as well
# $ ./run-benchs.sh -c IvDomain Bad.ExD

# We need to set manually the GOPATH, which is a couple of directories
# above starting from bench; what a hack....
export GOPATH=$(dirname $(dirname $(dirname $(dirname $(dirname $(pwd))))))
# now we may start go with a main and it will find the other packages

# configuration for benching
GOCMD=go
LOG_FOLD=logs
BIN_FOLD=bin

# ensure that the directories exist
mkdir -p $LOG_FOLD
mkdir -p $BIN_FOLD

# optional, realclean bin and log
clean() {	
	echo "CLEANING"
  /bin/rm -rf $LOG_FOLD
  /bin/rm -rf $BIN_FOLD
  mkdir $LOG_FOLD
  mkdir $BIN_FOLD
}

function build_bench {
	$GOCMD build -o ./$BIN_FOLD/bench_$1 benchframework.go bench_$1.go
}

# warm up, intention is to prevent the processor to be in deep sleepstates
function warmup {
	$GOCMD run warmup.go 
}

function run_bench {
	build_bench $1
    warmup
    ./$BIN_FOLD/bench_$1 "${*:2}"
}


# go through options, accept -c and -g <gocmd>
while getopts "cg:" arg; do
case "$arg" in
        g) GOCMD="$OPTARG";;
        c) GOCLEAN="cleaning bin and logs first"; clean;;
        [?]) echo >&2 "Usage: $0 [-g gocmd (go)] [-c] [packages]"
         echo >&2 "   add -c to clean bin and log first"
            exit 1;;
    esac
done
shift $(($OPTIND-1))

echo "Benching with : $GOCMD    "$GOCLEAN
ls bench_*.go | while read f; do 
BENCH=$(echo $f | sed -r 's/(^bench_|\.go$)//g')
if [[ $f == *$1* ]]; then # empty $1 automagically fits
echo "=== Benching $BENCH ==="
run_bench $BENCH "${*:2}"
fi
done
