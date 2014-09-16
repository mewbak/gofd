package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"strconv"
	"testing"
)

func sumweighted_test(t *testing.T, varsMapping []*VarMappingWeighted,
	sum *VarMapping, expready bool) {
	cnt := 0
	varList := make([]core.VarId, len(varsMapping))
	weightList := make([]int, len(varsMapping))
	for _, varMapping := range varsMapping {
		v := core.CreateIntVarExValues("V"+strconv.Itoa(cnt),
			store, varMapping.initDomain)
		varMapping.intVar = v
		varList[cnt] = v
		weightList[cnt] = varMapping.weight
		cnt += 1
	}
	sumVar := core.CreateIntVarExValues("SUM", store, sum.initDomain)
	sum.intVar = sumVar
	p := CreateWeightedSum(store, sumVar, weightList, varList...)
	store.AddPropagators(p)
	ready := store.IsConsistent()
	ready_test(t, "WeightedSum", ready, expready)
	if expready {
		for _, varMapping := range varsMapping {
			domainEquals_test(t, "WeightedSum",
				varMapping.intVar, varMapping.expDomain)
		}
		domainEquals_test(t, "WeightedSum", sum.intVar, sum.expDomain)
	}
}

func Test_WeightedSumA(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumA: 1*X + 2*Y = Z, X:0..2, Y:0..2, Z:6")

	X := CreateVarMappingWeighted(1, []int{0, 1, 2}, []int{2})
	Y := CreateVarMappingWeighted(2, []int{0, 1, 2}, []int{2})
	sum := CreateVarMapping([]int{6}, []int{6})
	sumweighted_test(t, []*VarMappingWeighted{X, Y}, sum, true)
}

func Test_WeightedSumB(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumB: 1*X + 2*Y + 3*Q = Z, X:0..2, Y:0..2, Q:0..2, Z:12")

	X := CreateVarMappingWeighted(1, []int{0, 1, 2}, []int{2})
	Y := CreateVarMappingWeighted(2, []int{0, 1, 2}, []int{2})
	Q := CreateVarMappingWeighted(3, []int{0, 1, 2}, []int{2})
	sum := CreateVarMapping([]int{12}, []int{12})
	sumweighted_test(t, []*VarMappingWeighted{X, Y, Q}, sum, true)
}

func Test_WeightedSumC(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumC: 1*X + 2*Y + 3*Q = Z, X:0..2, Y:0..2, Q:0..2, Z:13")

	X := CreateVarMappingWeighted(1, []int{0, 1, 2}, []int{})
	Y := CreateVarMappingWeighted(2, []int{0, 1, 2}, []int{})
	Q := CreateVarMappingWeighted(3, []int{0, 1, 2}, []int{})
	sum := CreateVarMapping([]int{13}, []int{})
	sumweighted_test(t, []*VarMappingWeighted{X, Y, Q}, sum, false)
}
