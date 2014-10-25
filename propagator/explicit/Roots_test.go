package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"testing"
)

func roots_test(t *testing.T, xiinit [][]int, sinit []int, tinit []int, expxi [][]int, exps []int, expready bool) {
	Xi := make([]core.VarId, len(xiinit))
	for i := 0; i < len(xiinit); i++ {
		s := fmt.Sprintf("X%v", i)
		Xi[i] = core.CreateIntVarExValues(s, store, xiinit[i])
	}

	S := core.CreateAuxIntVarExValues(store, sinit)
	T := core.CreateAuxIntVarExValues(store, tinit)

	store.AddPropagator(CreateRoots(Xi, S, T))
	ready := store.IsConsistent()
	y := fmt.Sprintf("ready: %v", ready)
	log(y)
	ready_test(t, "Roots", ready, expready)
	if expready {
		for i := 0; i < len(expxi); i++ {
			domainEquals_test(t, "Roots", Xi[i], expxi[i])
		}

		domainEquals_test(t, "Roots", S, exps)
	}
}

func Test_Rootsa(t *testing.T) {
	setup()
	defer teardown()
	log("Roots1: Xi:{{1,2,3},{1,3,4},{2,5},{5,7}}, S:{2,3,5}, T:{2,3,4}")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2, 3}
	xi[1] = []int{1, 3, 4}
	xi[2] = []int{2, 5}
	xi[3] = []int{5, 7}

	expxi := make([][]int, 4)
	expxi[0] = []int{1}
	expxi[1] = []int{3, 4}
	expxi[2] = []int{2}
	expxi[3] = []int{5, 7}

	s := []int{2, 3, 5}
	exps := []int{2, 3}

	tinit := []int{2, 3, 4}

	roots_test(t, xi, s, tinit, expxi, exps, true)
}

func Test_Rootsb(t *testing.T) {
	setup()
	defer teardown()
	log("Roots2: Xi:{{1,2,3},{1,3,4},{2,5},{5,7}}, S:{2,3,5}, T:{3,4}")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2, 3}
	xi[1] = []int{1, 3, 4}
	xi[2] = []int{2, 5}
	xi[3] = []int{5, 7}

	expxi := make([][]int, 4)
	expxi[0] = []int{1, 2}
	expxi[1] = []int{3, 4}
	expxi[2] = []int{2, 5}
	expxi[3] = []int{5, 7}

	s := []int{2, 3, 5}
	exps := []int{2}

	tinit := []int{3, 4}

	roots_test(t, xi, s, tinit, expxi, exps, true)
}

func Test_Rootsc(t *testing.T) {
	setup()
	defer teardown()
	log("Roots3: Xi:{{1,2,3},{1,3,4},{2,5},{5,7}}, S:{2,3,5}, T:{2,3,4,5,8}")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2, 3}
	xi[1] = []int{1, 3, 4}
	xi[2] = []int{2, 5}
	xi[3] = []int{5, 7}

	expxi := make([][]int, 4)
	expxi[0] = []int{1}
	expxi[1] = []int{3, 4}
	expxi[2] = []int{2, 5}
	expxi[3] = []int{7}

	s := []int{2, 3, 5}
	exps := []int{2, 3}

	tinit := []int{2, 3, 4, 5, 8}

	roots_test(t, xi, s, tinit, expxi, exps, true)
}

func Test_Rootsd(t *testing.T) {
	setup()
	defer teardown()
	log("Roots4: Xi:{{1,2,3},{1,3,4},{2,5},{5,7}}, S:{2,3,5}, T:{7,8}")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2, 3}
	xi[1] = []int{1, 3, 4}
	xi[2] = []int{2, 5}
	xi[3] = []int{5, 7}

	expxi := make([][]int, 4)
	expxi[0] = []int{1, 2, 3}
	expxi[1] = []int{}
	expxi[2] = []int{}
	expxi[3] = []int{5}

	s := []int{2, 3, 5}
	exps := []int{}

	tinit := []int{7, 8}

	roots_test(t, xi, s, tinit, expxi, exps, false)
}

func Test_Rootse(t *testing.T) {
	setup()
	defer teardown()
	log("Roots5: Xi:{{1,2,3},{1,3,4},{2,5},{5,7}}, S:{2}, T:{3,7,8}")
	xi := make([][]int, 4)
	xi[0] = []int{1, 2, 3}
	xi[1] = []int{1, 3, 4}
	xi[2] = []int{2, 5}
	xi[3] = []int{5, 7}

	expxi := make([][]int, 4)
	expxi[0] = []int{1, 2}
	expxi[1] = []int{3}
	expxi[2] = []int{2, 5}
	expxi[3] = []int{5}

	s := []int{2}
	exps := []int{2}

	tinit := []int{3, 7, 8}

	roots_test(t, xi, s, tinit, expxi, exps, true)
}
