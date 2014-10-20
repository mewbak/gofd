#!/bin/sh
# We need to set manually the GOPATH, which is a couple of directories
# above starting from bench; what a hack....
export GOPATH=$(dirname $(dirname $(dirname $(dirname $(dirname $(pwd))))))
# now we may start go with a main and it will find the other packages

set_local_test_environment(){
  echo "parameters for benching local environment"
  GO=go
  LOG_FOLD=logs
  BIN_FOLD=bin
}

set_compute_environment(){
  echo "parameters for benching compute environment"
  GO=go
  LOG_FOLD=logs
  BIN_FOLD=bin
}

clean(){
  /bin/rm -rf $LOG_FOLD
  /bin/rm -rf $BIN_FOLD
  mkdir $LOG_FOLD
  mkdir $BIN_FOLD
}

# preparation: setting up bench-environment, building executables
build() {  
  #-- basic --
  $GO build -o ./$BIN_FOLD/bench_Domain bench_Domain.go benchframework.go
  $GO build -o ./$BIN_FOLD/bench_GC1XeqC2Y bench_GC1XeqC2Y.go benchframework.go
  $GO build -o ./$BIN_FOLD/bench_GC1XplusC2YeqC3Z bench_GC1XplusC2YeqC3Z.go benchframework.go
  $GO build -o ./$BIN_FOLD/bench_IntVar bench_IntVar.go benchframework.go
  $GO build -o ./$BIN_FOLD/bench_Simple bench_Simple.go benchframework.go
  $GO build -o ./$BIN_FOLD/bench_Store bench_Store.go benchframework.go
  $GO build -o ./$BIN_FOLD/bench_XmultCeqY bench_XmultCeqY.go benchframework.go
  #-- intervals --
  $GO build -o ./$BIN_FOLD/bench_IvDomain bench_IvDomain.go benchframework.go
  #-- applications --
  $GO build -o ./$BIN_FOLD/bench_CarSequencing bench_CarSequencing.go benchframework.go
  $GO build -o ./$BIN_FOLD/bench_MagicSeries bench_MagicSeries.go benchframework.go
  $GO build -o ./$BIN_FOLD/bench_MysteryShopper bench_MysteryShopper.go benchframework.go
  $GO build -o ./$BIN_FOLD/bench_NQueens bench_NQueens.go benchframework.go
  $GO build -o ./$BIN_FOLD/bench_SMM bench_SMM.go benchframework.go
}


#- bench execution (function calls) -
set_local_test_environment
#set_compute_environment

clean
build

#- warm up -
echo "warm up"
$GO run warmup.go

#- benching-
./$BIN_FOLD/bench_Domain
./$BIN_FOLD/bench_GC1XeqC2Y
./$BIN_FOLD/bench_GC1XplusC2YeqC3Z
./$BIN_FOLD/bench_IntVar
./$BIN_FOLD/bench_Simple
./$BIN_FOLD/bench_Store
./$BIN_FOLD/bench_XmultCeqY
./$BIN_FOLD/bench_IvDomain
./$BIN_FOLD/bench_CarSequencing
./$BIN_FOLD/bench_MagicSeries
./$BIN_FOLD/bench_MysteryShopper
./$BIN_FOLD/bench_NQueens
./$BIN_FOLD/bench_SMM

