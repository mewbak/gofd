package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
)

type WeightedSum struct {
	vars             []core.VarId
	hvars            []core.VarId //helper-variables
	cs               []int
	resultVar        core.VarId
	outCh            chan<- *core.ChangeEvent
	inCh             <-chan *core.ChangeEntry
	varidToDomainMap map[core.VarId]*core.ExDomain
	id               core.PropId

	pseudoPropsXCY []*XmultCeqY
	pseudoPropsXYZ []*XPlusYEqZ
	finalProp      *XplusCeqY
}

func (this *WeightedSum) Start(store *core.Store) {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Dln("WeightedSum_'initial consistency check'")
	}
	evt := core.CreateChangeEvent()
	weightSumInitialCheck(this.varidToDomainMap,
		this.pseudoPropsXCY, this.pseudoPropsXYZ,
		this.finalProp, this.cs, evt)
	this.outCh <- evt // send changes to store

	for changeEntry := range this.inCh {
		if loggerDebug {
			core.GetLogger().Df("%s_'Incoming Change for %s'",
				this, core.GetNameRegistry().GetName(changeEntry.GetID()))
		}
		evt = core.CreateChangeEvent()
		varidChanged := changeEntry.GetID()
		changedDom := this.varidToDomainMap[varidChanged]
		changedDom.Removes(changeEntry.GetValues())
		weightedSumPropagate(varidChanged, this.varidToDomainMap,
			this.pseudoPropsXCY, this.pseudoPropsXYZ,
			this.finalProp, this.cs, evt)
		if loggerDebug {
			msg := "%s_propagate_'communicate change, evt-value: %s'"
			core.GetLogger().Df(msg, this, evt)
		}
		this.outCh <- evt // send changes to store
	}
}

// propagate check for changes. First look for X*C=Y propagators, then
// X+Y=Z and finally for the final propagator X=Y. Collect changes
func weightedSumPropagate(varid core.VarId,
	varidToDomainMap map[core.VarId]*core.ExDomain,
	pseudoPropsXCY []*XmultCeqY, pseudoPropsXYZ []*XPlusYEqZ,
	finalProp *XplusCeqY, cs []int, evt *core.ChangeEvent) {
	for i, prop := range pseudoPropsXCY {
		xDom := varidToDomainMap[prop.x]
		yDom := varidToDomainMap[prop.y]
		c := cs[i]
		if prop.x == varid {
			firstInMultSecondOut(xDom, c, yDom, prop.y, evt)
		} else if prop.y == varid {
			secondInMultFirstOut(yDom, c, xDom, prop.x, evt)
		}
	}
	sumPropagate(varid, varidToDomainMap, pseudoPropsXYZ, finalProp, evt)
}

// initialCheck check for changes. First look for X*C=Y propagators, then
// X+Y=Z and finally for the final propagator X=Y. Collect changes
func weightSumInitialCheck(varidToDomainMap map[core.VarId]*core.ExDomain,
	pseudoPropsXCY []*XmultCeqY, pseudoPropsXYZ []*XPlusYEqZ,
	finalProp *XplusCeqY, cs []int, evt *core.ChangeEvent) {

	for i, prop := range pseudoPropsXCY {
		xDom := varidToDomainMap[prop.x]
		yDom := varidToDomainMap[prop.y]
		c := cs[i]
		firstInMultSecondOut(xDom, c, yDom, prop.y, evt)
		secondInMultFirstOut(yDom, c, xDom, prop.x, evt)
	}
	sumInitialCheck(varidToDomainMap, pseudoPropsXYZ, finalProp, evt)
}

// Register generates auxiliary variables and makes pseudo structs
// and all vars will be registered at store and get domains and channels
func (this *WeightedSum) Register(store *core.Store) {
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
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(allvars, this.id)
	this.varidToDomainMap = core.GetVaridToExplicitDomainsMap(domains)

}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen
func (this *WeightedSum) SetID(propID core.PropId) {
	this.id = propID
}

func (this *WeightedSum) GetID() core.PropId {
	return this.id
}

func CreateWeightedSum(store *core.Store, resultVar core.VarId, cs []int,
	intVars ...core.VarId) *WeightedSum {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateWeightedSum-propagator")
	}
	prop := new(WeightedSum)
	prop.vars = intVars
	prop.resultVar = resultVar
	prop.cs = cs
	prop.pseudoPropsXCY = make([]*XmultCeqY, len(prop.vars))
	prop.hvars = make([]core.VarId, 0)
	for i, X := range prop.vars {
		H := core.CreateAuxIntVarExValues(store,
			core.ScalarSlice(prop.cs[i], store.GetDomain(X).Values_asSlice()))
		prop.pseudoPropsXCY[i] = CreateXmultCeqY(X, prop.cs[i], H)
		prop.hvars = append(prop.hvars, H)
	}
	prop.pseudoPropsXYZ = make([]*XPlusYEqZ, len(prop.pseudoPropsXCY)-1)
	H := prop.pseudoPropsXCY[0].y
	for i, X := range prop.pseudoPropsXCY[1:] {
		NewH := core.CreateAuxIntVarExFromTo(store,
			store.GetDomain(H).GetMin()+store.GetDomain(X.y).GetMin(),
			store.GetDomain(H).GetMax()+store.GetDomain(X.y).GetMax())
		prop.pseudoPropsXYZ[i] = CreateXplusYeqZ(H, X.y, NewH)
		H = NewH
		prop.hvars = append(prop.hvars, NewH)
	}
	prop.finalProp = CreateXeqY(H, prop.resultVar)
	return prop
}

func (this *WeightedSum) String() string {
	vars_str := make([]string, len(this.vars))
	for i, var_id := range this.vars {
		vars_str[i] = fmt.Sprintf("%v*%s",
			this.cs[i], core.GetNameRegistry().GetName(var_id))
	}
	return fmt.Sprintf("PROP_%d %s = %s",
		this.id, strings.Join(vars_str, "+"),
		core.GetNameRegistry().GetName(this.resultVar))
}

func (this *WeightedSum) Clone() core.Constraint {
	prop := new(WeightedSum)
	prop.vars = make([]core.VarId, len(this.vars))
	for i, single_var := range this.vars {
		prop.vars[i] = single_var
	}
	prop.resultVar = this.resultVar
	prop.cs = make([]int, len(this.cs))
	for i, c := range this.cs {
		prop.cs[i] = c
	}
	prop.pseudoPropsXCY = make([]*XmultCeqY, len(this.pseudoPropsXCY))
	for i, p := range this.pseudoPropsXCY {
		prop.pseudoPropsXCY[i] = p.Clone().(*XmultCeqY)
	}
	prop.pseudoPropsXYZ = make([]*XPlusYEqZ, len(this.pseudoPropsXYZ))
	for i, p := range this.pseudoPropsXYZ {
		prop.pseudoPropsXYZ[i] = p.Clone().(*XPlusYEqZ)
	}
	prop.hvars = make([]core.VarId, len(this.hvars))
	for i, single_var := range this.hvars {
		prop.hvars[i] = single_var
	}
	prop.finalProp = this.finalProp.Clone().(*XplusCeqY)
	return prop
}

// propagate-functions
// X*C=Y
// firstInMultSecondOut collect changes, when first variable has changed
// e.g. X*C=Y, then X is first variable
func firstInMultSecondOut(firstInDomain *core.ExDomain, c int,
	secondOutDomain *core.ExDomain, secondOutVarId core.VarId,
	evt *core.ChangeEvent) {
	var chEntry *core.ChangeEntry = nil
	for y_val := range secondOutDomain.Values {
		if c == 0 || !(y_val%c == 0) {
			if y_val != 0 {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(secondOutVarId)
				}
				chEntry.Add(y_val)
			}
		} else {
			if !firstInDomain.Contains(y_val / c) {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(secondOutVarId)
				}
				chEntry.Add(y_val)
			}
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

// secondInMultFirstOut collect changes, when second variable has changed
// e.g. X*C=Y, then Y is second variable
func secondInMultFirstOut(secondInDomain *core.ExDomain, c int,
	firstOutDomain *core.ExDomain, firstOutVarId core.VarId,
	evt *core.ChangeEvent) {
	var chEntry *core.ChangeEntry = nil
	for x_val := range firstOutDomain.Values {
		if !secondInDomain.Contains(x_val * c) {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(firstOutVarId)
			}
			chEntry.Add(x_val)
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

func (this *WeightedSum) GetAllVars() []core.VarId {
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

func (this *WeightedSum) GetVarIds() []core.VarId {
	return this.GetAllVars()
}

func (this *WeightedSum) GetDomains() []core.Domain {
	return core.ValuesOfMapVarIdToExDomain(this.GetAllVars(), this.varidToDomainMap)
}

func (this *WeightedSum) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *WeightedSum) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
