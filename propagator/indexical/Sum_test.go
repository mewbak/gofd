package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"strconv"
	"testing"
)

func icsum_ic_test(t *testing.T, varsMapping []*VarMapping, sum *VarMapping,
	expready bool) {
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
	p := CreateSum(store, sumVar, varList)
	store.AddPropagator(p)
	ready := store.IsConsistent()
	ready_test(t, "Sum_IC", ready, expready)
	if expready {
		for _, varMapping := range varsMapping {
			expDomain := core.CreateIvDomainFromIntArr(varMapping.expDomain)
			DomainEquals_test(t, "Sum_IC",
				varMapping.intVar, expDomain)
		}
		sumexpDomain := core.CreateIvDomainFromIntArr(sum.expDomain)
		DomainEquals_test(t, "Sum_IC", sum.intVar, sumexpDomain)
	}
}

func Test_Sum_IC2a(t *testing.T) {
	setup()
	defer teardown()
	log("Sum_ICa: V0:0..4, V1:0..4, V2:0..4, V3:11")

	v0 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	sum := CreateVarMapping([]int{11}, []int{11})
	icsum_ic_test(t, []*VarMapping{v0, v1, v2}, sum, true)
}

func Test_Sum_IC2b(t *testing.T) {
	setup()
	defer teardown()
	log("Sum_ICb: X:0..4, Y:0..4, Q:0..4, Z:12")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	v3 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	sum := CreateVarMapping([]int{12}, []int{12})
	icsum_ic_test(t, []*VarMapping{v1, v2, v3}, sum, true)
}

func Test_Sum_IC2c(t *testing.T) {
	setup()
	defer teardown()
	log("Sum_ICc: X:0..4, Y:0..4, Q:0..4, Z:13")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	v3 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	sum := CreateVarMapping([]int{13}, []int{})
	icsum_ic_test(t, []*VarMapping{v1, v2, v3}, sum, false)
}

func Test_Sum_ICd(t *testing.T) {
	setup()
	defer teardown()
	log("Sum_ICd: X:0..4, Y:0..4, Z:8")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	sum := CreateVarMapping([]int{8}, []int{8})
	icsum_ic_test(t, []*VarMapping{v1, v2}, sum, true)
}

func Test_Sum_ICe(t *testing.T) {
	setup()
	defer teardown()
	log("Sum_ICe: X:[[0,4][6,8]], Y:[[0,4][6,8]], Z:[[15,17]]")

	v1 := CreateVarMapping([]int{0, 1, 2, 3, 4, 6, 7, 8}, []int{7, 8})
	v2 := CreateVarMapping([]int{0, 1, 2, 3, 4, 6, 7, 8}, []int{7, 8})
	sum := CreateVarMapping([]int{15, 16, 17}, []int{15, 16})
	icsum_ic_test(t, []*VarMapping{v1, v2}, sum, true)
}
