package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xplusyeqz_ic_test(t *testing.T, vars []*VarMapping,
	expready bool) {

	x := core.CreateIntVarIvValues("X", store, vars[0].initDomain)
	y := core.CreateIntVarIvValues("X", store, vars[1].initDomain)
	z := core.CreateIntVarIvValues("X", store, vars[2].initDomain)

	p := CreateXplusYeqZ(x, y, z)
	store.AddPropagator(p)
	ready := store.IsConsistent()
	ready_test(t, "XpusYeqZ_IC", ready, expready)
	if expready {
		DomainEquals_test(t, "XpusYeqZ_IC", x, core.CreateIvDomainFromIntArr(vars[0].expDomain))
		DomainEquals_test(t, "XpusYeqZ_IC", y, core.CreateIvDomainFromIntArr(vars[1].expDomain))
		DomainEquals_test(t, "XpusYeqZ_IC", z, core.CreateIvDomainFromIntArr(vars[2].expDomain))
	}
}

func Test_XplusYeqZ_ICa(t *testing.T) {
	setup()
	defer teardown()
	log("XpusYeqZ_ICa: X:0..4, Y:0..4, Z:7")

	x := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	y := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{3, 4})
	z := CreateVarMapping([]int{7}, []int{7})
	xplusyeqz_ic_test(t, []*VarMapping{x, y, z}, true)
}

func Test_XplusYeqZ_ICb(t *testing.T) {
	setup()
	defer teardown()
	log("XpusYeqZ_ICb: X:0..4, Y:0..4, Z:9")

	x := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	y := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{})
	z := CreateVarMapping([]int{9}, []int{})
	xplusyeqz_ic_test(t, []*VarMapping{x, y, z}, false)
}

func Test_XplusYeqZ_ICc(t *testing.T) {
	setup()
	defer teardown()
	log("XpusYeqZ_ICc: X:0..4, Y:0..4, Z:8")

	x := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	y := CreateVarMapping([]int{0, 1, 2, 3, 4}, []int{4})
	z := CreateVarMapping([]int{8}, []int{8})
	xplusyeqz_ic_test(t, []*VarMapping{x, y, z}, true)
}

func Test_XplusYeqZ_ICd(t *testing.T) {
	setup()
	defer teardown()
	log("XpusYeqZ_ICd: X:[[0,4][6,8]], Y:[[0,4][6,8]], Z:[[15,17]]")

	x := CreateVarMapping([]int{0, 1, 2, 3, 4, 6, 7, 8}, []int{7, 8})
	y := CreateVarMapping([]int{0, 1, 2, 3, 4, 6, 7, 8}, []int{7, 8})
	z := CreateVarMapping([]int{15, 16, 17}, []int{15, 16})
	xplusyeqz_ic_test(t, []*VarMapping{x, y, z}, true)
}
