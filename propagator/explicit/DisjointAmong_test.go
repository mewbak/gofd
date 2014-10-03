package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"testing"
)

func disjointAmong_test(t *testing.T, xiinit [][]int, yjinit [][]int,
	expxi [][]int, expyj [][]int, expready bool) {
	Xi := make([]core.VarId, len(xiinit))
	for i := 0; i < len(xiinit); i++ {
		s := fmt.Sprintf("X%v", i)
		Xi[i] = core.CreateIntVarExValues(s, store, xiinit[i])
	}

	Yj := make([]core.VarId, len(yjinit))
	for i := 0; i < len(yjinit); i++ {
		s := fmt.Sprintf("Y%v", i)
		Yj[i] = core.CreateIntVarExValues(s, store, yjinit[i])
	}

	p := CreateDisjointAmong(Xi, Yj, store)
	store.AddPropagators(p...)
	ready := store.IsConsistent()
	y := fmt.Sprintf("ready: %v", ready)
	log(y)
	ready_test(t, "Disjoint with Among", ready, expready)
	if expready {
		for i := 0; i < len(expxi); i++ {
			domainEquals_test(t, "Disjoint with Among", Xi[i], expxi[i])
		}

		for i := 0; i < len(expyj); i++ {
			domainEquals_test(t, "Disjoint with Among", Yj[i], expyj[i])
		}
	}
}

func Test_DisjointAmonga(t *testing.T) {
	setup()
	defer teardown()
	log("Disjoint with Among 1: Xi:{{1,2}, {2,3}, {4}, {5,6}}, Yj:{{2,7}, {5,8}, {7}}")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2}
	xi[1] = []int{2, 3}
	xi[2] = []int{4}
	xi[3] = []int{5, 6}

	yj := make([][]int, 3)
	yj[0] = []int{2, 7}
	yj[1] = []int{5, 8}
	yj[2] = []int{7}

	expxi := make([][]int, 4)
	expxi[0] = []int{1}
	expxi[1] = []int{3}
	expxi[2] = []int{4}
	expxi[3] = []int{6}

	expyj := make([][]int, 3)
	expyj[0] = []int{7}
	expyj[1] = []int{8}
	expyj[2] = []int{7}

	disjointAmong_test(t, xi, yj, expxi, expyj, true)
}

func Test_DisjointAmongb(t *testing.T) {
	setup()
	defer teardown()
	log("Disjoint with Among 2: Xi:{{1,2}, {2}, {4}, {5,6}}, Yj:{{2}, {5,8}, {7}}")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2}
	xi[1] = []int{2}
	xi[2] = []int{4}
	xi[3] = []int{5, 6}

	yj := make([][]int, 3)
	yj[0] = []int{2}
	yj[1] = []int{5, 8}
	yj[2] = []int{7}

	expxi := make([][]int, 4)
	expxi[0] = []int{1}
	expxi[1] = []int{}
	expxi[2] = []int{4}
	expxi[3] = []int{6}

	disjointAmong_test(t, xi, yj, expxi, yj, false)
}

func Test_DisjointAmongc(t *testing.T) {
	setup()
	defer teardown()
	log("Disjoint with Among 3: Xi:{{1,2}, {2,3}, {4}, {6,9}}, Yj:{{5,10,11}, {5,8}, {7}}")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2}
	xi[1] = []int{2, 3}
	xi[2] = []int{4}
	xi[3] = []int{6, 9}

	yj := make([][]int, 3)
	yj[0] = []int{5, 10, 11}
	yj[1] = []int{5, 8}
	yj[2] = []int{7}

	disjointAmong_test(t, xi, yj, xi, yj, true)
}
