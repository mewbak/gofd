package indexical

import (
	"bitbucket.org/gofd/gofd/core"
)

type VarMappingWeighted struct {
	intVar     core.VarId
	weight     int
	initDomain []int
	expDomain  []int
}

func CreateVarMappingWeighted(w int, initD []int, expD []int) *VarMappingWeighted {
	vm := new(VarMappingWeighted)
	vm.weight = w
	vm.initDomain = initD
	vm.expDomain = expD
	return vm
}

type VarMapping struct {
	intVar     core.VarId
	initDomain []int
	expDomain  []int
}

func CreateVarMapping(initD []int, expD []int) *VarMapping {
	vm := new(VarMapping)
	vm.initDomain = initD
	vm.expDomain = expD
	return vm
}
