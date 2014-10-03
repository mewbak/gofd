package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func suminterval_test(t *testing.T, varsMapping []*VarMapping, sum *VarMapping,
	expready bool) {
	sumVar, varList := createSumtestVars(varsMapping, sum)
	p := CreateSum(store, sumVar, varList)
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

func Test_SumA(t *testing.T) {
	setup()
	defer teardown()
	log("SumA_intervals: V0:0..4, V1:0..4, V2:0..4, V3:11")

	v0 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	sum := CreateVarMapping([]int{11}, []int{11})
	suminterval_test(t, []*VarMapping{v0, v1, v2}, sum, true)
}

func Test_SumB(t *testing.T) {
	setup()
	defer teardown()
	log("SumB_intervals: X:0..4, Y:0..4, Q:0..4, Z:12")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	v3 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	sum := CreateVarMapping([]int{12}, []int{12})
	suminterval_test(t, []*VarMapping{v1, v2, v3}, sum, true)
}

func Test_SumC(t *testing.T) {
	setup()
	defer teardown()
	log("SumC_intervals: X:0..4, Y:0..4, Q:0..4, Z:13")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	v3 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	sum := CreateVarMapping([]int{13}, []int{})
	suminterval_test(t, []*VarMapping{v1, v2, v3}, sum, false)
}

func Test_SumD(t *testing.T) {
	setup()
	defer teardown()
	log("SumD_intervals: X:0..4, Y:0..4, Z:8")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	sum := CreateVarMapping([]int{8}, []int{8})
	suminterval_test(t, []*VarMapping{v1, v2}, sum, true)
}

func Test_SumE(t *testing.T) {
	setup()
	defer teardown()
	log("SumE_intervals: X:[[0,4][6,8]], Y:[[0,4][6,8]], Z:[[15,17]]")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4, 6, 7, 8}, []int{7, 8})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4, 6, 7, 8}, []int{7, 8})
	sum := CreateVarMapping([]int{15, 16, 17}, []int{15, 16})
	suminterval_test(t, []*VarMapping{v1, v2}, sum, true)
}

func Test_Sum_clone(t *testing.T) {
	setup()
	defer teardown()
	log("Sum_clone")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4, 6, 7, 8}, []int{7, 8})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4, 6, 7, 8}, []int{7, 8})
	sum := CreateVarMapping([]int{15, 16, 17}, []int{15, 16})
	sumVar, varList := createSumtestVars([]*VarMapping{v1, v2}, sum)
	constraint := CreateSum(store, sumVar, varList)

	clone_test(t, store, constraint)
}
