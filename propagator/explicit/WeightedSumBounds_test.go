package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"strconv"
	"testing"
)

func sumweightedBounds_test(t *testing.T, varsMapping []*VarMappingWeighted,
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
	p := CreateWeightedSumBounds(store, sumVar, weightList, varList)
	store.AddPropagators(p...)
	ready := store.IsConsistent()
	ready_test(t, "WeightedSumBounds", ready, expready)
	if expready {
		for _, varMapping := range varsMapping {
			domainEquals_test(t, "WeightedSumBounds",
				varMapping.intVar, varMapping.expDomain)
		}
		domainEquals_test(t, "WeightedSumBounds", sum.intVar, sum.expDomain)
	}
}

func domainCheck2(t *testing.T, curDomain *core.ExDomain,
	expDomain *core.ExDomain, varname string) {
	if !curDomain.Equals(expDomain) {
		t.Errorf("WeightedSum: got Dom%s=%s, want Dom%s=%s\n",
			varname, curDomain.String(), varname, expDomain.String())
	}
}

func Test_GSumWeightedA(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumBoundsA: 1*X + 2*Y = Z, X:0..2, Y:0..2, Z:6")

	X := CreateVarMappingWeighted(1, []int{0, 1, 2}, []int{2})
	Y := CreateVarMappingWeighted(2, []int{0, 1, 2}, []int{2})
	sum := CreateVarMapping([]int{6}, []int{6})
	sumweightedBounds_test(t, []*VarMappingWeighted{X, Y}, sum, true)
}

func Test_GSumWeightedB(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumBoundsB: 1*X + 2*Y + 3*Q = Z, X:0..2, Y:0..2, Q:0..2, Z:12")

	X := CreateVarMappingWeighted(1, []int{0, 1, 2}, []int{2})
	Y := CreateVarMappingWeighted(2, []int{0, 1, 2}, []int{2})
	Q := CreateVarMappingWeighted(3, []int{0, 1, 2}, []int{2})
	sum := CreateVarMapping([]int{12}, []int{12})
	sumweightedBounds_test(t, []*VarMappingWeighted{X, Y, Q}, sum, true)
}

func Test_GSumWeightedC(t *testing.T) {
	setup()
	defer teardown()
	log("WeightedSumBoundsC: 1*X + 2*Y + 3*Q = Z, X:0..2, Y:0..2, Q:0..2, Z:13")

	X := CreateVarMappingWeighted(1, []int{0, 1, 2}, []int{})
	Y := CreateVarMappingWeighted(2, []int{0, 1, 2}, []int{})
	Q := CreateVarMappingWeighted(3, []int{0, 1, 2}, []int{})
	sum := CreateVarMapping([]int{13}, []int{})
	sumweightedBounds_test(t, []*VarMappingWeighted{X, Y, Q}, sum, false)
}
