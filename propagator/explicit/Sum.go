package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
)

// idea: X+Y+Q=Z, X+Y=H1, H1+Q=H2, H2=Z

type Sum struct {
	vars             []core.VarId
	hvars            []core.VarId //helper-variables
	resultVar        core.VarId
	outCh            chan<- *core.ChangeEvent
	inCh             <-chan *core.ChangeEntry
	varidToDomainMap map[core.VarId]*core.ExDomain
	id               core.PropId
	store            *core.Store
	pseudoProps      []*XPlusYEqZ
	finalProp        *XplusCeqY
}

func (this *Sum) Start(store *core.Store) {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Dln("Sum_'initial consistency check'")
	}
	evt := core.CreateChangeEvent()
	sumInitialCheck(this.varidToDomainMap,
		this.pseudoProps, this.finalProp, evt)
	this.outCh <- evt // send changes to store
	for changeEntry := range this.inCh {
		if loggerDebug {
			msg := "%s_'Incoming Change for %s'"
			core.GetLogger().Df(msg, this, store.GetName(changeEntry.GetID()))
		}
		evt = core.CreateChangeEvent()
		varidChanged := changeEntry.GetID()
		changedDom := this.varidToDomainMap[varidChanged]
		changedDom.Removes(changeEntry.GetValues())
		sumPropagate(varidChanged, this.varidToDomainMap, this.pseudoProps,
			this.finalProp, evt)
		if loggerDebug {
			msg := "%s_propagate_'communicate change, evt-value: %s'"
			core.GetLogger().Df(msg, this, evt)
		}
		this.outCh <- evt // send changes to store
	}
}

func sumPropagate(varid core.VarId,
	varidToDomainMap map[core.VarId]*core.ExDomain,
	pseudoProps []*XPlusYEqZ, finalProp *XplusCeqY, evt *core.ChangeEvent) {
	for _, prop := range pseudoProps {
		xDom := varidToDomainMap[prop.x]
		yDom := varidToDomainMap[prop.y]
		zDom := varidToDomainMap[prop.z]
		if prop.x == varid {
			firstInSecondInResultOut(xDom, yDom, zDom, prop.z, evt)
			firstInResultInSecondOut(xDom, zDom, yDom, prop.y, evt)
		} else if prop.y == varid {
			firstInSecondInResultOut(xDom, yDom, zDom, prop.z, evt)
			secondInResultInFirstOut(yDom, zDom, xDom, prop.x, evt)
		} else if prop.z == varid {
			firstInResultInSecondOut(xDom, zDom, yDom, prop.y, evt)
			secondInResultInFirstOut(yDom, zDom, xDom, prop.x, evt)
		}
	}
	xDom := varidToDomainMap[finalProp.x]
	yDom := varidToDomainMap[finalProp.y]
	if finalProp.x == varid {
		firstInSecondOut(xDom, yDom, finalProp.y, evt)
	} else if finalProp.y == varid {
		secondInFirstOut(yDom, xDom, finalProp.x, evt)
	}
}

// sumInitialCheck starts an initial consistency check
func sumInitialCheck(varidToDomainMap map[core.VarId]*core.ExDomain,
	pseudoProps []*XPlusYEqZ, finalProp *XplusCeqY, evt *core.ChangeEvent) {
	for _, prop := range pseudoProps {
		xDom := varidToDomainMap[prop.x]
		yDom := varidToDomainMap[prop.y]
		zDom := varidToDomainMap[prop.z]
		firstInSecondInResultOut(xDom, yDom, zDom, prop.z, evt)
		firstInResultInSecondOut(xDom, zDom, yDom, prop.y, evt)
		secondInResultInFirstOut(yDom, zDom, xDom, prop.x, evt)
	}
	xDom := varidToDomainMap[finalProp.x]
	yDom := varidToDomainMap[finalProp.y]
	firstInSecondOut(xDom, yDom, finalProp.y, evt)
	secondInFirstOut(yDom, xDom, finalProp.x, evt)
}

// Register generates auxiliary variables and makes pseudo structs
// and all vars will be registered at store and get domains and channels
func (this *Sum) Register(store *core.Store) {
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
	this.store = store
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen
func (this *Sum) SetID(propID core.PropId) {
	this.id = propID
}

func (this *Sum) GetID() core.PropId {
	return this.id
}

func CreateSum(store *core.Store,
	resultVar core.VarId, intVars []core.VarId) core.Constraint {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateSum-propagator")
	}
	prop := new(Sum)
	prop.vars = intVars
	prop.resultVar = resultVar
	prop.pseudoProps = make([]*XPlusYEqZ, len(prop.vars)-1)
	prop.hvars = make([]core.VarId, 0)
	H := prop.vars[0]
	for i, X := range prop.vars[1:] {
		NewH := core.CreateAuxIntVarExFromTo(store,
			store.GetDomain(H).GetMin()+store.GetDomain(X).GetMin(),
			store.GetDomain(H).GetMax()+store.GetDomain(X).GetMax())
		prop.pseudoProps[i] = CreateXplusYeqZ(H, X, NewH)
		H = NewH
		prop.hvars = append(prop.hvars, NewH)
	}
	prop.finalProp = CreateXeqY(H, prop.resultVar)
	return prop
}

func (this *Sum) Clone() core.Constraint {
	prop := new(Sum)
	prop.vars = make([]core.VarId, len(this.vars))
	for i, single_var := range this.vars {
		prop.vars[i] = single_var
	}
	prop.resultVar = this.resultVar
	prop.pseudoProps = make([]*XPlusYEqZ, len(this.pseudoProps))
	for i, p := range this.pseudoProps {
		prop.pseudoProps[i] = p.Clone().(*XPlusYEqZ)
	}
	prop.hvars = make([]core.VarId, len(this.hvars))
	for i, single_var := range this.hvars {
		prop.hvars[i] = single_var
	}
	prop.finalProp = this.finalProp.Clone().(*XplusCeqY)
	return prop
}

func (this *Sum) String() string {
	vars_str := make([]string, len(this.vars))
	for i, var_id := range this.vars {
		vars_str[i] = this.store.GetName(var_id)
	}
	return fmt.Sprintf("PROP_%d %s = %s",
		this.id, strings.Join(vars_str, "+"),
		this.store.GetName(this.resultVar))
}

//propagate-functions

//--------- X=Z ----------
// xinYout
func firstInSecondOut(firstInDomain *core.ExDomain, secondOutDomain *core.ExDomain,
	secondOutVarId core.VarId, evt *core.ChangeEvent) {
	var chEntry *core.ChangeEntry = nil
	for second_val := range secondOutDomain.Values {
		if !firstInDomain.Contains(second_val) {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(secondOutVarId)
			}
			chEntry.Add(second_val)
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

// yinXout
func secondInFirstOut(secondInDomain *core.ExDomain, firstOutDomain *core.ExDomain,
	firstOutVarId core.VarId, evt *core.ChangeEvent) {
	var chEntry *core.ChangeEntry = nil
	for first_val := range firstOutDomain.Values {
		if !secondInDomain.Contains(first_val) {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(firstOutVarId)
			}
			chEntry.Add(first_val)
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

//--------- X+Y=Z ----------

// xinYinZout
func firstInSecondInResultOut(firstInDomain *core.ExDomain,
	secondInDomain *core.ExDomain, resultOutDomain *core.ExDomain,
	resultOutVarId core.VarId, evt *core.ChangeEvent) {
	var chEntry *core.ChangeEntry = nil
	for result_val := range resultOutDomain.Values {
		match := false
		for first_val := range firstInDomain.Values {
			if secondInDomain.Contains(result_val - first_val) {
				match = true
				break
			}
		}
		if !match {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(resultOutVarId)
			}
			chEntry.Add(result_val)
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

// xinZinYout
func firstInResultInSecondOut(firstInDomain *core.ExDomain,
	resultInDomain *core.ExDomain, secondOutDomain *core.ExDomain,
	secondOutVarId core.VarId, evt *core.ChangeEvent) {
	var chEntry *core.ChangeEntry = nil
	for second_val := range secondOutDomain.Values {
		match := false
		for first_val := range firstInDomain.Values {
			if resultInDomain.Contains(second_val + first_val) {
				match = true
				break
			}
		}
		if !match {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(secondOutVarId)
			}
			chEntry.Add(second_val)
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

// yinZinXout
func secondInResultInFirstOut(secondInDomain *core.ExDomain,
	resultInDomain *core.ExDomain, firstOutDomain *core.ExDomain,
	firstOutVarId core.VarId, evt *core.ChangeEvent) {
	var chEntry *core.ChangeEntry = nil
	for first_val := range firstOutDomain.Values {
		match := false
		for second_val := range secondInDomain.Values {
			if resultInDomain.Contains(first_val + second_val) {
				match = true
				break
			}
		}
		if !match {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(firstOutVarId)
			}
			chEntry.Add(first_val)
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

func (this *Sum) GetAllVars() []core.VarId {
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

func (this *Sum) GetVarIds() []core.VarId {
	return this.GetAllVars()
}

func (this *Sum) GetDomains() []core.Domain {
	return core.ValuesOfMapVarIdToExDomain(this.GetAllVars(), this.varidToDomainMap)
}

func (this *Sum) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *Sum) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
