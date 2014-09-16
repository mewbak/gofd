package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"testing"
)

func commonAmong_test(t *testing.T, xiinit [][]int, yjinit [][]int,
	ninit []int, minit []int, expxi [][]int, expyj [][]int, expn []int, expm []int, expready bool) {
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

	N := core.CreateIntVarExValues("N", store, ninit)
	M := core.CreateIntVarExValues("M", store, minit)
	p := CreateCommonAmong(Xi, Yj, N, M, store)
	store.AddPropagators(p...)
	ready := store.IsConsistent()
	y := fmt.Sprintf("ready: %v", ready)
	log(y)
	ready_test(t, "Common with Among", ready, expready)
	if expready {
		domainEquals_test(t, "Common with Among", N, expn)
		for i := 0; i < len(expxi); i++ {
			domainEquals_test(t, "Common with Among", Xi[i], expxi[i])
		}

		for i := 0; i < len(expyj); i++ {
			domainEquals_test(t, "Common with Among", Yj[i], expyj[i])
		}
	}
}

func Test_CommonAmonga(t *testing.T) {
	setup()
	defer teardown()
	log("Common with Among: Xi:{{1,2}, {2,3}, {4}, {5,6}}, Yj:{{2}, {3,4}, {7}}, N:{3}, M:{2}")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2}
	xi[1] = []int{2, 3}
	xi[2] = []int{4}
	xi[3] = []int{5, 6}

	yj := make([][]int, 3)
	yj[0] = []int{2}
	yj[1] = []int{3, 4}
	yj[2] = []int{7}

	expxi := make([][]int, 4)
	expxi[0] = []int{2}
	expxi[1] = []int{2, 3}
	expxi[2] = []int{4}
	expxi[3] = []int{5, 6}

	expyj := make([][]int, 3)
	expyj[0] = []int{2}
	expyj[1] = []int{3, 4}
	expyj[2] = []int{7}

	n := []int{3}
	expn := []int{3}

	m := []int{2}
	expm := []int{2}

	commonAmong_test(t, xi, yj, n, m, expxi, expyj, expn, expm, true)
}

func Test_CommonAmongb(t *testing.T) {
	setup()
	defer teardown()
	log("Common with Among: Xi:{{1,2}, {2,3}, {4}, {5,6}}, Yj:{{2}, {3,4}, {7}}, N:{0,...,4}, M:{0,...,3}")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2}
	xi[1] = []int{2, 3}
	xi[2] = []int{4}
	xi[3] = []int{5, 6}

	yj := make([][]int, 3)
	yj[0] = []int{2}
	yj[1] = []int{3, 4}
	yj[2] = []int{7}

	expxi := make([][]int, 4)
	expxi[0] = []int{1, 2}
	expxi[1] = []int{2, 3}
	expxi[2] = []int{4}
	expxi[3] = []int{5, 6}

	expyj := make([][]int, 3)
	expyj[0] = []int{2}
	expyj[1] = []int{3, 4}
	expyj[2] = []int{7}

	n := []int{0, 1, 2, 3, 4}
	expn := []int{2, 3}

	m := []int{0, 1, 2, 3}
	expm := []int{1, 2}

	commonAmong_test(t, xi, yj, n, m, expxi, expyj, expn, expm, true)
}
