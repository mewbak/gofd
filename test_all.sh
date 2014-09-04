#!/bin/bash
# which packages to test
PACKAGES=core
PACKAGES=$PACKAGES:indexical:indexical/ixrange:indexical/ixterm
PACKAGES=$PACKAGES:propagator/interval:propagator/explicit:propagator/indexical
PACKAGES=$PACKAGES:propagator/reification
PACKAGES=$PACKAGES:demo:labeling

# which cmd to use to run go
GOCMD=go

function test_package {
    start=$(date +%s)
    printf "%25s : " $1
    pushd $1 > /dev/null
    $GOCMD clean -i
    $GOCMD test -i
    printf $($GOCMD test 2>&1 | egrep "PASS|FAIL")
    end=$(date +%s)
    used_time=$((end-start))
    printf "%3d sec" $used_time
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

gstart=$(date +%s)
IFS=':' read -ra ARRPKG <<< "$PKGARGS"
echo "Testing all packages with : $GOCMD"
for PACKAGE in "${ARRPKG[@]}"; do
    test_package $PACKAGE
done
gend=$(date +%s)
gused_time=$((gend-gstart))
echo "Duration: "$gused_time" sec"
