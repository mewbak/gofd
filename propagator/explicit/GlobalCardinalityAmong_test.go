package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"testing"
)

func globalCardinalityAmong_test(t *testing.T, xiinit [][]int, k []int, counters []int, expxi [][]int, expready bool) {
	Xi := make([]core.VarId, len(xiinit))
	for i := 0; i < len(xiinit); i++ {
		s := fmt.Sprintf("X%v", i)
		Xi[i] = core.CreateIntVarExValues(s, store, xiinit[i])
	}

	p := CreateGCCAmong(Xi, k, counters, store)
	store.AddPropagators(p...)
	ready := store.IsConsistent()
	y := fmt.Sprintf("ready: %v", ready)
	log(y)
	ready_test(t, "GlobalCardinality with Among", ready, expready)
	if expready {
		for i := 0; i < len(expxi); i++ {
			domainEquals_test(t, "GlobalCardinality with Among", Xi[i], expxi[i])
		}
	}
}

func Test_GlobalCardinalityAmonga(t *testing.T) {
	setup()
	defer teardown()
	log("GlobalCardinality with Among: Xi:{{1,2,3},{1,3,4},{2,5},{5,7}}, K:{2,3,7}, Counters:<1,2,1>")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2, 3}
	xi[1] = []int{1, 3, 4}
	xi[2] = []int{2, 5}
	xi[3] = []int{5, 7}

	expxi := make([][]int, 4)
	expxi[0] = []int{3}
	expxi[1] = []int{3}
	expxi[2] = []int{2}
	expxi[3] = []int{7}

	k := []int{2, 3, 7}

	counters := []int{1, 2, 1}

	globalCardinalityAmong_test(t, xi, k, counters, expxi, true)
}

func Test_GlobalCardinalityAmongb(t *testing.T) {
	setup()
	defer teardown()
	log("GlobalCardinality with Among: Xi:{{1,2,3},{1,2,3},{1,2,3},{5,7}}, K:<2,3>, Counters:<2,1>")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2, 3}
	xi[1] = []int{1, 2, 3}
	xi[2] = []int{1, 2, 3}
	xi[3] = []int{5, 7}

	expxi := make([][]int, 4)
	expxi[0] = []int{1, 2, 3}
	expxi[1] = []int{1, 2, 3}
	expxi[2] = []int{1, 2, 3}
	expxi[3] = []int{5, 7}

	k := []int{2, 3}

	counters := []int{2, 1}

	globalCardinalityAmong_test(t, xi, k, counters, expxi, true)
}

func Test_GlobalCardinalityAmongcs(t *testing.T) {
	setup()
	defer teardown()
	log("GlobalCardinality with Among: Xi:{{1,2,3},{1,2,3},{1,2,3},{5,7}}, K:<2,3>, Counters:<3,3>")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2, 3}
	xi[1] = []int{1, 2, 3}
	xi[2] = []int{1, 2, 3}
	xi[3] = []int{5, 7}

	expxi := make([][]int, 4)
	expxi[0] = []int{1, 2, 3}
	expxi[1] = []int{1, 2, 3}
	expxi[2] = []int{1, 2, 3}
	expxi[3] = []int{5, 7}

	k := []int{2, 3}

	counters := []int{3, 3}

	globalCardinalityAmong_test(t, xi, k, counters, expxi, false)
}
