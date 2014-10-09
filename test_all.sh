#!/bin/bash
# which packages to test
PACKAGES=core
PACKAGES=$PACKAGES:indexical:indexical/ixrange:indexical/ixterm
PACKAGES=$PACKAGES:propagator/interval:propagator/explicit:propagator/indexical
PACKAGES=$PACKAGES:propagator/reification
PACKAGES=$PACKAGES:demo:labeling

# which cmd to use to run go
GOCMD=go

PRECSECS=10000 
# compute a milli second precision time stamp
# cycling every PRECSECS seconds
function msecs {
stamp=$(($(date +%s)%$PRECSECS))$((date +%N) | cut -c 1-3)
echo $stamp
}

# compute a time difference between now and an older time stamp
function diffmsecs {
stamp=$1
now=$(msecs)
if [ $now -gt $1 ]
then 
	timediff=$((now-stamp))
else
	timediff=$(((PRECSECS+now)-stamp))
   fi
echo $timediff
}

function test_package {
    start=$(msecs)
    printf "%25s : " $1
    pushd $1 > /dev/null
    $GOCMD clean -i
    $GOCMD test -i
    printf "%s" $($GOCMD test 2>&1 | egrep "PASS|FAIL")
    used_time=$(diffmsecs $start)
    printf " %5d msecs" $used_time
    echo
    popd > /dev/null
}

# http://stackoverflow.com/questions/1527049/bash-join-elements-of-an-array
function join { 
    local IFS="$1"; shift; echo "$*"; 
}

while getopts "g:" arg; do
case "$arg" in
	g) GOCMD="$OPTARG";;
	[?]) print >&2 "Usage: $0 [-g gocmd (go)] [packages] "
	    exit 1;;
    esac
done
shift $(($OPTIND-1))

if [ -n "$1" ]; then
    PKGARGS=$(join ':' $*)
else 
    PKGARGS=$PACKAGES 
fi

gstart=$(msecs)
IFS=':' read -ra ARRPKG <<< "$PKGARGS"
echo "Testing all packages with : $GOCMD"
for PACKAGE in "${ARRPKG[@]}"; do
    test_package $PACKAGE
done
gused_time=$(diffmsecs $gstart)
secs=$(($gused_time / 1000)) 
decmsecs=$((($gused_time%1000)/10))
printf "Duration: %3d.%02d secs\n" $secs $decmsecs
