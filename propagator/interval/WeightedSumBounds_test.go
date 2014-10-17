package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"strconv"
	"testing"
)

func createWeightedSumtestVars(varsMapping []*VarMappingWeighted,
	sum *VarMapping) (core.VarId, []int, []core.VarId) {

	cnt := 0
	varList := make([]core.VarId, len(varsMapping))
	weightList := make([]int, len(varsMapping))
	for _, varMapping := range varsMapping {
		v := core.CreateIntVarIvValues("V"+strconv.Itoa(cnt),
			store, varMapping.initDomain)
		varMapping.intVar = v
		varList[cnt] = v
		weightList[cnt] = varMapping.weight
		cnt += 1
	}
	sumVar := core.CreateIntVarIvValues("SUM", store, sum.initDomain)
	sum.intVar = sumVar

	return sumVar, weightList, varList
}

func sumweightedBoundsinterval_test(t *testing.T, varsMapping []*VarMappingWeighted,
	sum *VarMapping, expready bool) {

	sumVar, weightList, varList := createWeightedSumtestVars(varsMapping,
		sum)
	p := CreateWeightedSumBounds(store, sumVar, weightList, varList...)
	store.AddPropagators(p)
	ready := store.IsConsistent()
	ready_test(t, "WeightedSumBounds_intervals", ready, expready)
	if expready {
		for _, varMapping := range varsMapping {
			expDomain := core.CreateIvDomainFromIntArr(varMapping.expDomain)
			DomainEquals_test(t, "WeightedSumBounds_intervals",
				varMapping.intVar, expDomain)
		}
		sumexpDomain := core.CreateIvDomainFromIntArr(sum.expDomain)
		DomainEquals_test(t, "WeightedSumBounds_intervals", sum.intVar, sumexpDomain)
	}
}

func Test_WeightedSumBoundsA(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumBoundsA_intervals: 1*X + 2*Y = Z, X:0..2, Y:0..2, Z:6")
	X := CreateVarMappingWeighted(1, []int{0, 1, 2}, []int{2})
	Y := CreateVarMappingWeighted(2, []int{0, 1, 2}, []int{2})
	sum := CreateVarMapping([]int{6}, []int{6})
	sumweightedBoundsinterval_test(t, []*VarMappingWeighted{X, Y}, sum, true)
}

func Test_WeightedSumBoundsB(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumBoundsB_intervals: 1*X + 2*Y + 3*Q = Z, X:0..2, Y:0..2, Q:0..2, Z:12")
	X := CreateVarMappingWeighted(1, []int{0, 1, 2}, []int{2})
	Y := CreateVarMappingWeighted(2, []int{0, 1, 2}, []int{2})
	Q := CreateVarMappingWeighted(3, []int{0, 1, 2}, []int{2})
	sum := CreateVarMapping([]int{12}, []int{12})
	sumweightedBoundsinterval_test(t, []*VarMappingWeighted{X, Y, Q}, sum, true)
}

func Test_WeightedSumBoundsC(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumBoundsC_intervals: 1*X + 2*Y + 3*Q = Z, X:0..2, Y:0..2, Q:0..2, Z:13")
	X := CreateVarMappingWeighted(1, []int{0, 1, 2}, []int{})
	Y := CreateVarMappingWeighted(2, []int{0, 1, 2}, []int{})
	Q := CreateVarMappingWeighted(3, []int{0, 1, 2}, []int{})
	sum := CreateVarMapping([]int{13}, []int{})
	sumweightedBoundsinterval_test(t, []*VarMappingWeighted{X, Y, Q}, sum, false)
}

func Test_WeightedSumBoundsD(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumBoundsD_intervals: 1*A + 1*B = Z, A:0..4, B:0..4, C:[5..10, 12..15]")
	A := CreateVarMappingWeighted(1, []int{0, 1, 2, 3, 4}, []int{1, 2, 3, 4})
	B := CreateVarMappingWeighted(1, []int{0, 1, 2, 3, 4}, []int{1, 2, 3, 4})
	sum := CreateVarMapping([]int{5, 6, 7, 8, 9, 10, 12, 13, 14, 15}, []int{5, 6, 7, 8})
	sumweightedBoundsinterval_test(t, []*VarMappingWeighted{A, B}, sum, true)
}

func Test_WeightedSumBoundsE(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumBoundsE_intervals: 1*A + 1*B = Z, A:0..4,9 , B:0..4,9, C:[5..10, 12..15]")
	A := CreateVarMappingWeighted(1, []int{0, 1, 2, 3, 4, 9}, []int{0, 1, 2, 3, 4, 9})
	B := CreateVarMappingWeighted(1, []int{0, 1, 2, 3, 4, 9}, []int{0, 1, 2, 3, 4, 9})
	sum := CreateVarMapping([]int{5, 6, 7, 8, 9, 10, 12, 13, 14, 15},
		[]int{5, 6, 7, 8, 9, 10, 12, 13, 14, 15})
	sumweightedBoundsinterval_test(t, []*VarMappingWeighted{A, B}, sum, true)
}

func Test_WeightedSumBoundsF(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumBoundsF_intervals: 1*A + 1*B = Z, A:0..4,11 , B:0..4,11 , C:[5..10, 12..15]")
	A := CreateVarMappingWeighted(1, []int{0, 1, 2, 3, 4, 11}, []int{0, 1, 2, 3, 4, 11})
	B := CreateVarMappingWeighted(1, []int{0, 1, 2, 3, 4, 11}, []int{0, 1, 2, 3, 4, 11})
	sum := CreateVarMapping([]int{5, 6, 7, 8, 12, 13, 14, 15},
		[]int{5, 6, 7, 8, 12, 13, 14, 15})
	sumweightedBoundsinterval_test(t, []*VarMappingWeighted{A, B}, sum, true)
}

func Test_WeightedSumBounds_clone(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumBounds_clone")

	X := CreateVarMappingWeighted(1, []int{0, 1, 2}, []int{2})
	Y := CreateVarMappingWeighted(2, []int{0, 1, 2}, []int{2})
	sum := CreateVarMapping([]int{6}, []int{6})
	sumVar, weightList, varList := createWeightedSumtestVars(
		[]*VarMappingWeighted{X, Y}, sum)
	constraint := CreateWeightedSumBounds(store, sumVar, weightList, varList...)

	clone_test(t, store, constraint)
}
