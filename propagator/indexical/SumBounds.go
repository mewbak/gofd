package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical"
	"fmt"
	"strings"
)

// idea: X+Y+Q=Z, --> X+Y=H1, --> H1+Q=Z

type SumBounds struct {
	vars             []core.VarId
	hvars            []core.VarId //helper-variables
	resultVar        core.VarId
	outCh            chan<- *core.ChangeEvent
	inCh             <-chan *core.ChangeEntry
	varidToDomainMap map[core.VarId]*core.IvDomain
	id               core.PropId
	store            *core.Store
	iColl            *indexical.IndexicalCollection
	pseudoProps      []*XplusYeqZ_Rel //with pseudoFinalProp
}

func (this *SumBounds) GetIndexicalCollection() *indexical.IndexicalCollection {
	return this.iColl
}

func (this *SumBounds) Start() {
	indexical.InitProcessConstraint(this, true)
	indexical.ProcessConstraint(this, true)
}

func (this *SumBounds) GetAllVars() []core.VarId {
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
func (this *SumBounds) Register(store *core.Store) {
	allvars := this.GetAllVars()

	var domains map[core.VarId]core.Domain

	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(allvars, this.id)

	this.varidToDomainMap = core.GetVaridToIntervalDomains(domains)

	this.store = store

	this.iColl = indexical.CreateIndexicalCollection()

	this.iColl.AddIndexicalsAtPrio(MakeSumBoundsIndexicals(this.pseudoProps, this.varidToDomainMap), indexical.HIGHEST)
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen
func (this *SumBounds) SetID(propID core.PropId) {
	this.id = propID
}

func (this *SumBounds) GetID() core.PropId {
	return this.id
}

func CreateSumBounds(store *core.Store,
	resultVar core.VarId, intVars []core.VarId) core.Constraint {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateSumBounds-propagator")
	}
	prop := new(SumBounds)
	prop.vars = intVars
	prop.resultVar = resultVar

	if len(prop.vars) == 1 {
		prop.pseudoProps = make([]*XplusYeqZ_Rel, 1)
		prop.hvars = make([]core.VarId, 1)
		H := core.CreateAuxIntVarIvFromTo(store, 0, 0)
		prop.hvars[0] = H
		prop.pseudoProps[0] = CreateXplusYeqZ_Rel(prop.vars[0], H, prop.resultVar)
		return prop
	}

	prop.pseudoProps = make([]*XplusYeqZ_Rel, len(prop.vars)-1)
	prop.hvars = make([]core.VarId, 0)
	H := prop.vars[0]
	//exclusive... [1:len(prop.vars)-1] means, without the last two ones
	for i, X := range prop.vars[1 : len(prop.vars)-1] {
		hDom := store.GetDomain(H)
		xDom := store.GetDomain(X)

		NewH := core.CreateAuxIntVarIvFromTo(store,
			hDom.GetMin()+xDom.GetMin(),
			hDom.GetMax()+xDom.GetMax())
		prop.pseudoProps[i] = CreateXplusYeqZ_Rel(H, X, NewH)
		H = NewH
		prop.hvars = append(prop.hvars, NewH)
	}
	X := prop.vars[len(prop.vars)-1]
	prop.pseudoProps[len(prop.pseudoProps)-1] = CreateXplusYeqZ_Rel(H, X, prop.resultVar)
	return prop
}

func (this *SumBounds) Clone() core.Constraint {
	prop := new(SumBounds)
	prop.vars = make([]core.VarId, len(this.vars))
	for i, single_var := range this.vars {
		prop.vars[i] = single_var
	}
	prop.resultVar = this.resultVar

	prop.pseudoProps = make([]*XplusYeqZ_Rel, len(this.pseudoProps))
	for i, p := range this.pseudoProps {
		prop.pseudoProps[i] = p.Clone()
	}

	prop.hvars = make([]core.VarId, len(this.hvars))
	for i, single_var := range this.hvars {
		prop.hvars[i] = single_var
	}

	return prop
}

func (this *SumBounds) String() string {
	vars_str := make([]string, len(this.vars))
	for i, var_id := range this.vars {
		vars_str[i] = this.store.GetName(var_id)
	}
	return fmt.Sprintf("PROP_%d %s = %s",
		this.id, strings.Join(vars_str, "+"),
		this.store.GetName(this.resultVar))
}

func (this *SumBounds) GetVarIds() []core.VarId {
	return this.GetAllVars()
}

func (this *SumBounds) GetDomains() []core.Domain {
	return core.ValuesOfMapVarIdToIvDomain(this.GetAllVars(), this.varidToDomainMap)
}

func (this *SumBounds) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *SumBounds) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
