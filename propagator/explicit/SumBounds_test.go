package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"strconv"
	"testing"
)

func sumBounds_test(t *testing.T, varsMapping []*VarMapping, sum *VarMapping,
	expready bool) {
	cnt := 0
	varList := make([]core.VarId, len(varsMapping))
	for _, varMapping := range varsMapping {
		v := core.CreateIntVarExValues("V"+strconv.Itoa(cnt),
			store, varMapping.initDomain)
		varMapping.intVar = v
		varList[cnt] = v
		cnt += 1
	}
	sumVar := core.CreateIntVarExValues("SUM", store, sum.initDomain)
	sum.intVar = sumVar
	p := CreateSumBounds(store, sumVar, varList)
	store.AddPropagators(p...)
	ready := store.IsConsistent()
	ready_test(t, "Sum", ready, expready)
	if expready {
		for _, varMapping := range varsMapping {
			domainEquals_test(t, "Sum", varMapping.intVar,
				varMapping.expDomain)
		}
		domainEquals_test(t, "Sum", sum.intVar, sum.expDomain)
	}
}

func Test_SumBoundsA(t *testing.T) {
	setup()
	defer teardown()
	log("SumBoundsA: X:0..4, Y:0..4, Q:0..4, Z:11")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	v3 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	sum := CreateVarMapping([]int{11}, []int{11})
	sumBounds_test(t, []*VarMapping{v1, v2, v3}, sum, true)
}

func Test_SumBoundsB(t *testing.T) {
	setup()
	defer teardown()
	log("SumBoundsB: X:0..4, Y:0..4, Q:0..4, Z:12")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	v3 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	sum := CreateVarMapping([]int{12}, []int{12})
	sumBounds_test(t, []*VarMapping{v1, v2, v3}, sum, true)
}

func Test_SumBoundsC(t *testing.T) {
	setup()
	defer teardown()
	log("SumBoundsC: X:0..4, Y:0..4, Q:0..4, Z:13")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	v3 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	sum := CreateVarMapping([]int{13}, []int{})
	sumBounds_test(t, []*VarMapping{v1, v2, v3}, sum, false)
}
