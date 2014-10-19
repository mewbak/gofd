#!/bin/bash
# which packages to test
PACKAGES=core
PACKAGES=$PACKAGES:indexical:indexical/ixrange:indexical/ixterm
PACKAGES=$PACKAGES:propagator/interval:propagator/explicit
PACKAGES=$PACKAGES:propagator/indexical
PACKAGES=$PACKAGES:propagator/reification
PACKAGES=$PACKAGES:demo:labeling

# which cmd to use to run go
GOCMD=go
# do not clean and install by default
GOCLEAN=

PRECSECS=10000
# compute a milli second precision time stamp
# cycling every PRECSECS seconds
function msecs {
stamp=$(($(date +%s)%$PRECSECS))$((date +%N) | cut -c 1-3)
echo $stamp
}

# compute a time difference between now and an older time stamp
# take cycling once (after PRECSECS) into account
function diffmsecs {
stamp=$1
now=$(msecs)
if [ "$now" -gt "$stamp" ]
then
	timediff=$((now-stamp))
else
	timediff=$((((PRECSECS*1000)+now)-stamp))
fi
echo $timediff
}

function test_package {
start=$(msecs)
printf "%25s : " $1
pushd $1 > /dev/null
if [ -n "$GOCLEAN" ]; then
    $GOCMD clean -i
fi
$GOCMD test -i
printf "%s" $($GOCMD test 2>&1 | egrep "PASS|FAIL")
printf_usedtime $(diffmsecs $start)
echo
popd > /dev/null
}

# formatted print of used time in msecs resolution
# with prefix $2 and postfix $3
function printf_usedtime {
ut_secs=$(($1/1000))
ut_msecs=$((($1%1000)))
printf $2"%3d.%03ds"$3 $ut_secs $ut_msecs
}

# http://stackoverflow.com/questions/1527049/bash-join-elements-of-an-array
function join {
    local IFS="$1"; shift; echo "$*";
}

# go through options, accept -c and -g <gocmd>
while getopts "cg:" arg; do
case "$arg" in
	g) GOCMD="$OPTARG";;
	c) GOCLEAN="do some cleaning first";;
	[?]) echo >&2 "Usage: $0 [-g gocmd (go)] [-c] [packages]"
         echo >&2 "   add -c to clean first, use with changing gocmd"
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
echo "Testing all packages with : $GOCMD    "$GOCLEAN
for PACKAGE in "${ARRPKG[@]}"; do
    test_package $PACKAGE
done
printf_usedtime $(diffmsecs $gstart) "Duration:" "\n"
