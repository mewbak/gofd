package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"testing"
)

func amongSimple_test(t *testing.T, xiinit [][]int, kinit []int,
	ninit []int, expxi [][]int, expn []int, expready bool) {
	Xi := make([]core.VarId, len(xiinit))
	for i := 0; i < len(xiinit); i++ {
		s := fmt.Sprintf("X%v", i)
		Xi[i] = core.CreateIntVarExValues(s, store, xiinit[i])
	}
	N := core.CreateIntVarExValues("N", store, ninit)
	store.AddPropagators(CreateAmongSimple(Xi, kinit, N))
	ready := store.IsConsistent()
	y := fmt.Sprintf("ready: %v", ready)
	log(y)
	ready_test(t, "AmongSimple", ready, expready)
	if expready {
		domainEquals_test(t, "AmongSimple", N, expn)
		for i := 0; i < len(expxi); i++ {
			domainEquals_test(t, "AmongSimple", Xi[i], expxi[i])
		}
	}
}

func Test_AmongSimplea(t *testing.T) {
	setup()
	defer teardown()
	log("AmongSimple: Xi:{{1},{2},{3}}, K:{1,2}, N:{2}")
	xi := make([][]int, 3)
	xi[0] = []int{1}
	xi[1] = []int{2}
	xi[2] = []int{3}
	k := []int{1, 2}
	n := []int{2}
	amongSimple_test(t, xi, k, n, xi, n, true)
}

func Test_AmongSimpleb(t *testing.T) {
	setup()
	defer teardown()
	log("AmongSimple: Xi:{{2},{4},{5}}, K:{2,3,4}, N:{1}")
	xi := make([][]int, 3)
	xi[0] = []int{2}
	xi[1] = []int{4}
	xi[2] = []int{5}
	k := []int{2, 3, 4}
	n := []int{1}
	expn := []int{}
	amongSimple_test(t, xi, k, n, xi, expn, false)
}

func Test_AmongSimplec(t *testing.T) {
	setup()
	defer teardown()
	log("AmongSimple: Xi:{{2,3},{4,5},{5,6}}, K:{2,3,4}, N:{1}")
	xi := make([][]int, 3)
	xi[0] = []int{2, 3}
	xi[1] = []int{4, 5}
	xi[2] = []int{5, 6}
	k := []int{2, 3, 4}
	n := []int{1}
	amongSimple_test(t, xi, k, n, xi, n, true)
}

func Test_AmongSimpled(t *testing.T) {
	setup()
	defer teardown()
	log("AmongSimple: Xi:{{2,3},{4,5},{5,6}}, K:{2,3,4}, N:{4}")
	xi := make([][]int, 3)
	xi[0] = []int{2, 3}
	xi[1] = []int{4, 5}
	xi[2] = []int{5, 6}
	k := []int{2, 3, 4}
	n := []int{2}
	amongSimple_test(t, xi, k, n, xi, n, true)
}
