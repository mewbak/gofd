package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"testing"
)

func alldifferentAmong_test(t *testing.T, xiinit [][]int, expxi [][]int,
	expready bool) {
	Xi := make([]core.VarId, len(xiinit))
	for i := 0; i < len(xiinit); i++ {
		s := fmt.Sprintf("X%v", i)
		Xi[i] = core.CreateIntVarExValues(s, store, xiinit[i])
	}
	p := CreateAlldifferentAmong(Xi, store)
	store.AddPropagators(p...)
	ready := store.IsConsistent()
	y := fmt.Sprintf("ready: %v", ready)
	log(y)
	ready_test(t, "Alldifferent with Among", ready, expready)
	if expready {
		for i := 0; i < len(expxi); i++ {
			domainEquals_test(t, "Alldifferent with Among", Xi[i], expxi[i])
		}
	}
}

func Test_AlldifferentAmonga(t *testing.T) {
	setup()
	defer teardown()
	log("AlldifferentAmong: Xi:{{1,2,3},{1},{2}}")
	xi := make([][]int, 3)
	xi[0] = []int{1, 2, 3}
	xi[1] = []int{1}
	xi[2] = []int{2}

	expxi := make([][]int, 3)
	expxi[0] = []int{3}
	expxi[1] = []int{1}
	expxi[2] = []int{2}
	alldifferentAmong_test(t, xi, expxi, true)
}

func Test_AlldifferentAmongb(t *testing.T) {
	setup()
	defer teardown()
	log("AlldifferentAmong: Xi:{{1,2},{1,2},{3,4},{3}}")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2}
	xi[1] = []int{1, 2}
	xi[2] = []int{3, 4}
	xi[3] = []int{3}

	expxi := make([][]int, 4)
	expxi[0] = []int{1, 2}
	expxi[1] = []int{1, 2}
	expxi[2] = []int{4}
	expxi[3] = []int{3}
	alldifferentAmong_test(t, xi, expxi, true)
}

func Test_AlldifferentAmongc(t *testing.T) {
	setup()
	defer teardown()
	log("AlldifferentAmong: Xi:{{1,2,3},{3,4},{5,6}}")
	xi := make([][]int, 3)
	xi[0] = []int{1, 2, 3}
	xi[1] = []int{3, 4}
	xi[2] = []int{5, 6}

	expxi := make([][]int, 3)
	expxi[0] = []int{1, 2, 3}
	expxi[1] = []int{3, 4}
	expxi[2] = []int{5, 6}
	alldifferentAmong_test(t, xi, expxi, true)
}
