package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"fmt"
	"strings"
)

type WeightedSumBounds struct {
	vars             []core.VarId
	hvars            []core.VarId //helper-variables
	cs               []int
	resultVar        core.VarId
	outCh            chan<- *core.ChangeEvent
	inCh             <-chan *core.ChangeEntry
	varidToDomainMap map[core.VarId]*core.IvDomain
	id               core.PropId
	store            *core.Store
	iColl            *indexical.IndexicalCollection
	pseudoPropsXCY   []*XmultCeqY_Rel
	pseudoPropsXYZ   []*XplusYeqZ_Rel
}

func (this *WeightedSumBounds) GetIndexicalCollection() *indexical.IndexicalCollection {
	return this.iColl
}

func (this *WeightedSumBounds) Start() {
	indexical.InitProcessConstraint(this, true)
	indexical.ProcessConstraint(this, true)
}

func (this *WeightedSumBounds) GetAllVars() []core.VarId {
	allvars := make([]core.VarId, len(this.vars)+len(this.hvars)+1)
	i := 0
	for _, v := range this.vars {
		allvars[i] = v
		i++
	}
	for _, v := range this.hvars {
		allvars[i] = v
		i++
	}
	allvars[i] = this.resultVar
	return allvars
}

// Register generates auxiliary variables and makes pseudo structs
// and all vars will be registered at store and get domains and channels
func (this *WeightedSumBounds) Register(store *core.Store) {
	allvars := this.GetAllVars()
	var domains map[core.VarId]core.Domain

	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(allvars, this.id)

	this.varidToDomainMap = core.GetVaridToIntervalDomains(domains)

	this.store = store

	this.iColl = indexical.CreateIndexicalCollection()

	hvarXCY := make([]core.VarId, len(this.pseudoPropsXCY))
	for i, hxcy := range this.pseudoPropsXCY {
		hvarXCY[i] = hxcy.y
	}

	this.iColl.AddIndexicalsAtPrio(MakeSumBoundsIndexicals(this.pseudoPropsXYZ, this.varidToDomainMap), indexical.HIGHEST)

	//X*C=Y
	for _, p := range this.pseudoPropsXCY {
		//mul: Y in X*C

		xDomR := ixrange.CreateDomRange(p.x, this.varidToDomainMap[p.x])
		cDomR := ixrange.CreateDomRange(-1, core.CreateIvDomainFromTo(p.c, p.c))
		multR := ixrange.CreateMultRange(xDomR, cDomR)

		this.iColl.CreateAndAddIndexical(p.y, this.varidToDomainMap[p.y], multR, indexical.HIGH)

		//div: X in Y/C
		yDomR := ixrange.CreateDomRange(p.y, this.varidToDomainMap[p.y])
		cDomR = ixrange.CreateDomRange(-1, core.CreateIvDomainFromTo(p.c, p.c))
		divR := ixrange.CreateDivRange(yDomR, cDomR)
		this.iColl.CreateAndAddIndexical(p.x, this.varidToDomainMap[p.x], divR, indexical.HIGH)
	}
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen
func (this *WeightedSumBounds) SetID(propID core.PropId) {
	this.id = propID
}

func (this *WeightedSumBounds) GetID() core.PropId {
	return this.id
}

func CreateWeightedSumBounds(store *core.Store, resultVar core.VarId, cs []int,
	intVars ...core.VarId) *WeightedSumBounds {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateWeightedSumBounds-propagator")
	}
	prop := new(WeightedSumBounds)
	prop.vars = intVars
	prop.resultVar = resultVar
	prop.cs = cs
	prop.pseudoPropsXCY = make([]*XmultCeqY_Rel, len(prop.vars))
	prop.hvars = make([]core.VarId, 0)
	for i, X := range prop.vars {
		H := core.CreateAuxIntVarIvValues(store,
			core.ScalarSlice(prop.cs[i], store.GetDomain(X).Values_asSlice()))
		prop.pseudoPropsXCY[i] = CreateXmultCeqY_Rel(X, prop.cs[i], H)
		prop.hvars = append(prop.hvars, H)
	}
	prop.pseudoPropsXYZ = make([]*XplusYeqZ_Rel, len(prop.pseudoPropsXCY)-1)
	H := prop.pseudoPropsXCY[0].y
	newHVars := make([]core.VarId, 0)
	for i, p := range prop.pseudoPropsXCY[1 : len(prop.vars)-1] {
		NewH := core.CreateAuxIntVarIvFromTo(store,
			store.GetDomain(H).GetMin()+store.GetDomain(p.y).GetMin(),
			store.GetDomain(H).GetMax()+store.GetDomain(p.y).GetMax())
		prop.pseudoPropsXYZ[i] = CreateXplusYeqZ_Rel(H, p.y, NewH)
		H = NewH
		newHVars = append(newHVars, NewH)
	}
	X := prop.hvars[len(prop.hvars)-1]
	prop.hvars = append(prop.hvars, newHVars...)
	prop.pseudoPropsXYZ[len(prop.pseudoPropsXYZ)-1] = CreateXplusYeqZ_Rel(H, X, prop.resultVar)

	return prop
}

func (this *WeightedSumBounds) String() string {
	vars_str := make([]string, len(this.vars))
	for i, var_id := range this.vars {
		vars_str[i] = fmt.Sprintf("%v*%s",
			this.cs[i], this.store.GetName(var_id))
	}
	return fmt.Sprintf("PROP_%d %s = %s",
		this.id, strings.Join(vars_str, "+"),
		this.store.GetName(this.resultVar))
}

func (this *WeightedSumBounds) Clone() core.Constraint {
	prop := new(WeightedSumBounds)
	prop.vars = make([]core.VarId, len(this.vars))
	for i, single_var := range this.vars {
		prop.vars[i] = single_var
	}
	prop.resultVar = this.resultVar
	prop.cs = make([]int, len(this.cs))
	for i, c := range this.cs {
		prop.cs[i] = c
	}
	prop.pseudoPropsXCY = make([]*XmultCeqY_Rel, len(this.pseudoPropsXCY))
	for i, p := range this.pseudoPropsXCY {
		prop.pseudoPropsXCY[i] = p.Clone()
	}
	prop.pseudoPropsXYZ = make([]*XplusYeqZ_Rel, len(this.pseudoPropsXYZ))
	for i, p := range this.pseudoPropsXYZ {
		prop.pseudoPropsXYZ[i] = p.Clone()
	}
	prop.hvars = make([]core.VarId, len(this.hvars))
	for i, single_var := range this.hvars {
		prop.hvars[i] = single_var
	}

	return prop
}

func (this *WeightedSumBounds) GetVarIds() []core.VarId {
	return this.GetAllVars()
}

func (this *WeightedSumBounds) GetDomains() []core.Domain {
	return core.ValuesOfMapVarIdToIvDomain(this.GetAllVars(), this.varidToDomainMap)
}

func (this *WeightedSumBounds) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *WeightedSumBounds) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
