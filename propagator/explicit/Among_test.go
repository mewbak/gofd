package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"testing"
)

func among_test(t *testing.T, xiinit [][]int, kinit []int,
	ninit []int, expxi [][]int, expn []int, expready bool) {
	Xi := make([]core.VarId, len(xiinit))
	for i := 0; i < len(xiinit); i++ {
		s := fmt.Sprintf("X%v", i)
		Xi[i] = core.CreateIntVarExValues(s, store, xiinit[i])
	}
	N := core.CreateIntVarExValues("N", store, ninit)
	store.AddPropagators(CreateAmong(Xi, kinit, N))
	ready := store.IsConsistent()
	y := fmt.Sprintf("ready: %v", ready)
	log(y)
	ready_test(t, "Among", ready, expready)
	if expready {
		domainEquals_test(t, "Among", N, expn)
		for i := 0; i < len(expxi); i++ {
			domainEquals_test(t, "Among", Xi[i], expxi[i])
		}
	}
}

func Test_Amonga(t *testing.T) {
	setup()
	defer teardown()
	log("Among 1: Xi:{{1,2,3},{1,2,3},{1,2,3}}, K:{1,2}, N:2")
	xi := make([][]int, 3)
	for i := 0; i < len(xi); i++ {
		xi[i] = []int{1, 2, 3}
	}

	k := []int{1, 2}
	n := []int{2}

	among_test(t, xi, k, n, xi, n, true)
}

func Test_Amongb(t *testing.T) {
	setup()
	defer teardown()
	log("Among 2 Xi:{{2,3},{4,5},{5,6}}, K:{2,3,4}, N:{0,1,2,3}")
	xi := make([][]int, 3)
	xi[0] = []int{2, 3}
	xi[1] = []int{4, 5}
	xi[2] = []int{5, 6}

	expxi := make([][]int, 3)
	expxi[0] = []int{2, 3}
	expxi[1] = []int{4, 5}
	expxi[2] = []int{5, 6}

	k := []int{2, 3, 4}
	n := []int{0, 1, 2, 3}
	expn := []int{1, 2}

	among_test(t, xi, k, n, expxi, expn, true)
}

func Test_Amongc(t *testing.T) {
	setup()
	defer teardown()
	log("Among 3: Xi:{{2,3},{4,5},{5,6}}, K:{2,3,4}, N:{1}")
	xi := make([][]int, 3)
	xi[0] = []int{2, 3}
	xi[1] = []int{4, 5}
	xi[2] = []int{5, 6}

	expxi := make([][]int, 3)
	expxi[0] = []int{2, 3}
	expxi[1] = []int{5}
	expxi[2] = []int{5, 6}

	k := []int{2, 3, 4}
	n := []int{1}

	among_test(t, xi, k, n, expxi, n, true)
}

func Test_Amongd(t *testing.T) {
	setup()
	defer teardown()
	log("Among 4: Xi:{{2,3},{4,5},{5,6}}, K:{2,3,4}, N:{2}")
	xi := make([][]int, 3)
	xi[0] = []int{2, 3}
	xi[1] = []int{4, 5}
	xi[2] = []int{5, 6}

	expxi := make([][]int, 3)
	expxi[0] = []int{2, 3}
	expxi[1] = []int{4}
	expxi[2] = []int{5, 6}

	k := []int{2, 3, 4}
	n := []int{2}

	among_test(t, xi, k, n, expxi, n, true)
}

func Test_Amonge(t *testing.T) {
	setup()
	defer teardown()
	log("Among 5: Xi:{{2,3},{4,5},{4,5},{5,6}}, K:{2,3,4}, N:{1}")
	xi := make([][]int, 4)
	xi[0] = []int{2, 3}
	xi[1] = []int{4, 5}
	xi[2] = []int{4, 5}
	xi[3] = []int{5, 6}

	expxi := make([][]int, 4)
	expxi[0] = []int{2, 3}
	expxi[1] = []int{5}
	expxi[2] = []int{5}
	expxi[3] = []int{5, 6}

	k := []int{2, 3, 4}
	n := []int{1}

	among_test(t, xi, k, n, expxi, n, true)
}
