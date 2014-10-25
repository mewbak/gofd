package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
)

// idea: X+Y+Q=Z --> X+Y=H1 --> H1+Q=Z

// Sum represents the constraint X1+X2+...+Xn=Z
// Its propagate functions establish bounds consistency.
// The basic idea of Sum is to substitute the whole Sum equation to many
// X+Y=Z equations.
// i.e. Sum constraint X+Y+Q=Z results in X+Y=H1 and H1+Q=Z
type SumBounds struct {
	vars             []core.VarId
	hvars            []core.VarId //helper-variables
	resultVar        core.VarId
	outCh            chan<- *core.ChangeEvent
	inCh             <-chan *core.ChangeEntry
	varidToDomainMap map[core.VarId]*core.IvDomain
	id               core.PropId
	store            *core.Store
	pseudoProps      []*PseudoXplusYeqZ //with pseudoFinalProp
}

func (this *SumBounds) Start() {
	core.LogInitConsistency(this)
	// initial check
	evt := core.CreateChangeEvent()
	this.ivSumBoundsInitialCheck(evt)
	core.SendChangesToStore(evt, this)
	// process all changes
	for changeEntry := range this.inCh {
		core.LogIncomingChange(this, this.store, changeEntry)
		evt = core.CreateChangeEvent()
		varidChanged := changeEntry.GetID()
		changedDom := this.varidToDomainMap[varidChanged]
		changedDom.Removes(changeEntry.GetValues())
		this.ivSumBoundsPropagate(varidChanged, evt)
		core.SendChangesToStore(evt, this)
	}
}

func (this *SumBounds) ivSumBoundsPropagate(varid core.VarId, evt *core.ChangeEvent) {
	ivSumBoundsBoundsPropagate(varid, this.varidToDomainMap, this.pseudoProps, evt)
}

func ivSumBoundsBoundsPropagate(varid core.VarId,
	varidToDomainMap map[core.VarId]*core.IvDomain,
	pseudoProps []*PseudoXplusYeqZ, evt *core.ChangeEvent) {

	for _, prop := range pseudoProps {
		xDom := varidToDomainMap[prop.x]
		yDom := varidToDomainMap[prop.y]
		zDom := varidToDomainMap[prop.z]

		if prop.x == varid {
			resOutBounds(xDom, yDom, zDom, prop.z, evt)
			secondOutBounds(xDom, yDom, zDom, prop.y, evt)
		} else if prop.y == varid {
			resOutBounds(xDom, yDom, zDom, prop.z, evt)
			firstOutBounds(xDom, yDom, zDom, prop.x, evt)
		} else if prop.z == varid {
			secondOutBounds(xDom, yDom, zDom, prop.y, evt)
			firstOutBounds(xDom, yDom, zDom, prop.x, evt)
		}
	}
}

func ivSumBoundsBoundsInitialCheck(varidToDomainMap map[core.VarId]*core.IvDomain,
	pseudoProps []*PseudoXplusYeqZ, evt *core.ChangeEvent) {

	for _, prop := range pseudoProps {
		xDom := varidToDomainMap[prop.x]
		yDom := varidToDomainMap[prop.y]
		zDom := varidToDomainMap[prop.z]

		firstOutBounds(xDom, yDom, zDom, prop.x, evt)
		secondOutBounds(xDom, yDom, zDom, prop.y, evt)
		resOutBounds(xDom, yDom, zDom, prop.z, evt)
	}
}

func (this *SumBounds) ivSumBoundsInitialCheck(evt *core.ChangeEvent) {
	ivSumBoundsBoundsInitialCheck(this.varidToDomainMap, this.pseudoProps, evt)
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
		core.GetLogger().Dln("CreateSumBoundsBounds-propagator")
	}
	prop := new(SumBounds)
	prop.vars = intVars
	prop.resultVar = resultVar
	prop.pseudoProps = make([]*PseudoXplusYeqZ, len(prop.vars)-1)
	prop.hvars = make([]core.VarId, 0)
	H := prop.vars[0]
	//exclusive... [1:len(prop.vars)-1] means, without the last two ones
	for i, X := range prop.vars[1 : len(prop.vars)-1] {
		hDom := store.GetDomain(H)
		xDom := store.GetDomain(X)

		NewH := core.CreateAuxIntVarIvFromTo(store,
			hDom.GetMin()+xDom.GetMin(),
			hDom.GetMax()+xDom.GetMax())
		prop.pseudoProps[i] = CreatePseudoXplusYeqZ(H, X, NewH)
		H = NewH
		prop.hvars = append(prop.hvars, NewH)
	}
	X := prop.vars[len(prop.vars)-1]
	prop.pseudoProps[len(prop.pseudoProps)-1] = CreatePseudoXplusYeqZ(H, X, prop.resultVar)

	return prop
}

func (this *SumBounds) Clone() core.Constraint {
	prop := new(SumBounds)
	prop.vars = make([]core.VarId, len(this.vars))
	for i, single_var := range this.vars {
		prop.vars[i] = single_var
	}
	prop.resultVar = this.resultVar
	prop.pseudoProps = make([]*PseudoXplusYeqZ, len(this.pseudoProps))
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
