package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"strconv"
	"testing"
)

func createSumtestVars(varsMapping []*VarMapping,
	sum *VarMapping) (core.VarId, []core.VarId) {

	cnt := 0
	varList := make([]core.VarId, len(varsMapping))
	for _, varMapping := range varsMapping {
		v := core.CreateIntVarIvValues("V"+strconv.Itoa(cnt),
			store, varMapping.initDomain)
		varMapping.intVar = v
		varList[cnt] = v
		cnt += 1
	}
	sumVar := core.CreateIntVarIvValues("SUM", store, sum.initDomain)
	sum.intVar = sumVar

	return sumVar, varList
}

func sumBoundsinterval_test(t *testing.T, varsMapping []*VarMapping, sum *VarMapping,
	expready bool) {
	sumVar, varList := createSumtestVars(varsMapping, sum)
	p := CreateSumBounds(store, sumVar, varList)
	store.AddPropagator(p)
	ready := store.IsConsistent()
	ready_test(t, "Sum_intervals", ready, expready)
	if expready {
		for _, varMapping := range varsMapping {
			expDomain := core.CreateIvDomainFromIntArr(varMapping.expDomain)
			DomainEquals_test(t, "Sum_intervals",
				varMapping.intVar, expDomain)
		}
		sumexpDomain := core.CreateIvDomainFromIntArr(sum.expDomain)
		DomainEquals_test(t, "Sum_intervals", sum.intVar, sumexpDomain)
	}
}

func Test_SumBoundsA(t *testing.T) {
	setup()
	defer teardown()
	log("SumBoundsA_intervals: V0:0..4, V1:0..4, V2:0..4, V3:11")

	v0 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	sum := CreateVarMapping([]int{11}, []int{11})
	sumBoundsinterval_test(t, []*VarMapping{v0, v1, v2}, sum, true)
}

func Test_SumBoundsB(t *testing.T) {
	setup()
	defer teardown()
	log("SumBoundsB_intervals: X:0..4, Y:0..4, Q:0..4, Z:12")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	v3 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	sum := CreateVarMapping([]int{12}, []int{12})
	sumBoundsinterval_test(t, []*VarMapping{v1, v2, v3}, sum, true)
}

func Test_SumBoundsC(t *testing.T) {
	setup()
	defer teardown()
	log("SumBoundsC_intervals: X:0..4, Y:0..4, Q:0..4, Z:13")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	v3 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	sum := CreateVarMapping([]int{13}, []int{})
	sumBoundsinterval_test(t, []*VarMapping{v1, v2, v3}, sum, false)
}

func Test_SumBoundsD(t *testing.T) {
	setup()
	defer teardown()
	log("SumBoundsD_intervals: X:0..4, Y:0..4, Z:8")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	sum := CreateVarMapping([]int{8}, []int{8})
	sumBoundsinterval_test(t, []*VarMapping{v1, v2}, sum, true)
}

func Test_SumBoundsE(t *testing.T) {
	setup()
	defer teardown()
	log("SumBoundsE_intervals: X:[[0,4][6,8]], Y:[[0,4][6,8]], Z:[[15,17]]")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4, 6, 7, 8}, []int{7, 8})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4, 6, 7, 8}, []int{7, 8})
	sum := CreateVarMapping([]int{15, 16, 17}, []int{15, 16})
	sumBoundsinterval_test(t, []*VarMapping{v1, v2}, sum, true)
}

func Test_SumBounds_clone(t *testing.T) {
	setup()
	defer teardown()
	log("SumBounds_clone")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4, 6, 7, 8}, []int{7, 8})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4, 6, 7, 8}, []int{7, 8})
	sum := CreateVarMapping([]int{15, 16, 17}, []int{15, 16})
	sumVar, varList := createSumtestVars([]*VarMapping{v1, v2}, sum)
	constraint := CreateSumBounds(store, sumVar, varList)

	clone_test(t, store, constraint)
}
