package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
)

// idea: X+Y+Q=Z --> X+Y=H1 --> H1+Q=Z

// Sum represents the constraint X_1+X_2+...+Xn=Z
// Its propagate functions establish arc consistency (with bounds and arc
// algorithms).
// The basic idea of Sum is to substitute the whole Sum equation to many
// X+Y=Z equations.
// i.e. Sum constraint X+Y+Q=Z results in X+Y=H1 and H1+Q=Z

type Sum struct {
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

func (this *Sum) Start(store *core.Store) {
	core.LogInitConsistency(this)

	// initial check
	evt := core.CreateChangeEvent()
	this.ivsumInitialCheck(evt)
	core.SendChangesToStore(evt, this)

	for changeEntry := range this.inCh {
		core.LogIncomingChange(this, store, changeEntry)

		evt = core.CreateChangeEvent()
		varidChanged := changeEntry.GetID()
		changedDom := this.varidToDomainMap[varidChanged]
		changedDom.Removes(changeEntry.GetValues())
		this.ivsumPropagate(varidChanged, evt)
		core.SendChangesToStore(evt, this)
	}
}

func (this *Sum) ivsumPropagate(varid core.VarId, evt *core.ChangeEvent) {
	ivsumBoundsPropagate(varid, this.varidToDomainMap, this.pseudoProps, evt)
	ivsumArcPropagate(varid, this.varidToDomainMap, this.pseudoProps, evt)
}

func ivsumBoundsPropagate(varid core.VarId,
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

func ivsumArcPropagate(varid core.VarId,
	varidToDomainMap map[core.VarId]*core.IvDomain,
	pseudoProps []*PseudoXplusYeqZ, evt *core.ChangeEvent) {
	for _, prop := range pseudoProps {
		xDom := varidToDomainMap[prop.x]
		yDom := varidToDomainMap[prop.y]
		zDom := varidToDomainMap[prop.z]

		//Kanten
		if prop.x == varid {
			ivfirstInSecondInResultOut(xDom, yDom, zDom, prop.z, evt)
			ivfirstInResultInSecondOut(xDom, zDom, yDom, prop.y, evt)
		} else if prop.y == varid {
			ivfirstInSecondInResultOut(xDom, yDom, zDom, prop.z, evt)
			ivsecondInResultInFirstOut(yDom, zDom, xDom, prop.x, evt)
		} else if prop.z == varid {
			ivfirstInResultInSecondOut(xDom, zDom, yDom, prop.y, evt)
			ivsecondInResultInFirstOut(yDom, zDom, xDom, prop.x, evt)
		}
	}
}

func ivsumBoundsInitialCheck(varidToDomainMap map[core.VarId]*core.IvDomain,
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

func (this *Sum) ivsumInitialCheck(evt *core.ChangeEvent) {
	ivsumBoundsInitialCheck(this.varidToDomainMap, this.pseudoProps, evt)
	ivsumArcInitialCheck(this.varidToDomainMap, this.pseudoProps, evt)
}

// ivsumArcInitialCheck starts an initial consistency check
func ivsumArcInitialCheck(varidToDomainMap map[core.VarId]*core.IvDomain,
	pseudoProps []*PseudoXplusYeqZ, evt *core.ChangeEvent) {
	for _, prop := range pseudoProps {
		xDom := varidToDomainMap[prop.x]
		yDom := varidToDomainMap[prop.y]
		zDom := varidToDomainMap[prop.z]

		//Kanten
		ivfirstInSecondInResultOut(xDom, yDom, zDom, prop.z, evt)
		ivfirstInResultInSecondOut(xDom, zDom, yDom, prop.y, evt)
		ivsecondInResultInFirstOut(yDom, zDom, xDom, prop.x, evt)
	}
}

// Register generates auxiliary variables and makes pseudo structs
// and all vars will be registered at store and get domains and channels
func (this *Sum) Register(store *core.Store) {
	allvars := this.GetAllVars()
	var domains map[core.VarId]core.Domain

	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(allvars, this.id)

	this.varidToDomainMap = core.GetVaridToIntervalDomains(domains)

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

func (this *Sum) Clone() core.Constraint {
	prop := new(Sum)
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

//--------- X+Y=Z ----------

//--- bounds ---

func resOutBounds(firstInDomain *core.IvDomain,
	secondInDomain *core.IvDomain, resultOutDomain *core.IvDomain,
	resultOutVarId core.VarId, evt *core.ChangeEvent) {

	//Z in min(Y)+min(X)..max(Y)+max(X)
	if !secondInDomain.IsEmpty() && !firstInDomain.IsEmpty() {
		minY := secondInDomain.GetParts()[0].From
		minX := firstInDomain.GetParts()[0].From
		maxY := secondInDomain.GetParts()[len(secondInDomain.GetParts())-1].To
		maxX := firstInDomain.GetParts()[len(firstInDomain.GetParts())-1].To

		lowerBound := minY + minX
		upperBound := maxY + maxX

		i1 := core.CreateIvDomPart(core.NEG_INFINITY, lowerBound-1)
		i2 := core.CreateIvDomPart(upperBound+1, core.INFINITY)

		removingD := core.CreateIvDomainDomParts([]*core.IvDomPart{i1, i2})
		chEntry := core.CreateChangeEntryWithValues(resultOutVarId, removingD)
		evt.AddChangeEntry(chEntry)
		resultOutDomain.Removes(removingD)
	}
}

func secondOutBounds(firstInDomain *core.IvDomain,
	secondOutDomain *core.IvDomain, resultInDomain *core.IvDomain,
	secondOutVarId core.VarId, evt *core.ChangeEvent) {

	//Y in min(Z)-max(X)..max(Z)-min(X)
	if !resultInDomain.IsEmpty() && !firstInDomain.IsEmpty() {
		minZ := resultInDomain.GetParts()[0].From
		minX := firstInDomain.GetParts()[0].From
		maxZ := resultInDomain.GetParts()[len(resultInDomain.GetParts())-1].To
		maxX := firstInDomain.GetParts()[len(firstInDomain.GetParts())-1].To

		lowerBound := minZ - maxX
		upperBound := maxZ - minX

		i1 := core.CreateIvDomPart(core.NEG_INFINITY, lowerBound-1)
		i2 := core.CreateIvDomPart(upperBound+1, core.INFINITY)

		removingD := core.CreateIvDomainDomParts([]*core.IvDomPart{i1, i2})
		chEntry := core.CreateChangeEntryWithValues(secondOutVarId, removingD)
		evt.AddChangeEntry(chEntry)
		secondOutDomain.Removes(removingD)
	}
}

func firstOutBounds(firstOutDomain *core.IvDomain,
	secondInDomain *core.IvDomain, resultInDomain *core.IvDomain,
	firstOutVarId core.VarId, evt *core.ChangeEvent) {

	//X in min(Z)-max(Y)..max(Z)-min(Y)
	if !resultInDomain.IsEmpty() && !secondInDomain.IsEmpty() {
		minZ := resultInDomain.GetParts()[0].From
		minY := secondInDomain.GetParts()[0].From
		maxZ := resultInDomain.GetParts()[len(resultInDomain.GetParts())-1].To
		maxY := secondInDomain.GetParts()[len(secondInDomain.GetParts())-1].To

		lowerBound := minZ - maxY
		upperBound := maxZ - minY

		i1 := core.CreateIvDomPart(core.NEG_INFINITY, lowerBound-1)
		i2 := core.CreateIvDomPart(upperBound+1, core.INFINITY)

		removingD := core.CreateIvDomainDomParts([]*core.IvDomPart{i1, i2})
		chEntry := core.CreateChangeEntryWithValues(firstOutVarId, removingD)
		evt.AddChangeEntry(chEntry)
		firstOutDomain.Removes(removingD)

	}
}

//--- arc ---

// xinYinZout
func ivfirstInSecondInResultOut(firstInDomain *core.IvDomain,
	secondInDomain *core.IvDomain, resultOutDomain *core.IvDomain,
	resultOutVarId core.VarId, evt *core.ChangeEvent) {

	rparts := resultOutDomain.GetParts()
	if len(rparts) == 0 {
		return
	}
	beginRparts := rparts[0]
	endRparts := rparts[len(rparts)-1]

	translatedParts := make([]*core.IvDomPart, 0)
	for _, firstInP := range firstInDomain.GetParts() {
		for _, secondInP := range secondInDomain.GetParts() {
			tPart := firstInP.ADD(secondInP)
			//only take parts, which are relevant (are not out of bounds of resultOutDomain)
			if tPart.LT_DP(beginRparts) || tPart.GT_DP(endRparts) {
				continue
			}
			translatedParts = append(translatedParts, tPart)
		}
	}

	var removingD core.Domain

	if len(translatedParts) != 0 {
		//A+B = U			(U: Union)
		//removingD: outD diff V   (Difference)
		unionD := core.CreateIvDomainUnion(translatedParts)
		removingD = resultOutDomain.DifferenceWithIvDomain(unionD)
		if removingD.IsEmpty() {
			return
		}
	} else {
		removingD = resultOutDomain.Copy()
	}
	chEntry := core.CreateChangeEntryWithValues(resultOutVarId, removingD)
	evt.AddChangeEntry(chEntry)
}

// xinZinYout
func ivfirstInResultInSecondOut(firstInDomain *core.IvDomain,
	resultInDomain *core.IvDomain, secondOutDomain *core.IvDomain,
	secondOutVarId core.VarId, evt *core.ChangeEvent) {

	rparts := secondOutDomain.GetParts()
	if len(rparts) == 0 {
		return
	}
	beginRparts := rparts[0]
	endRparts := rparts[len(rparts)-1]

	translatedParts := make([]*core.IvDomPart, 0)
	for _, resultInP := range resultInDomain.GetParts() {
		for _, firstInP := range firstInDomain.GetParts() {
			tPart := resultInP.SUBTRACT(firstInP)
			if tPart.LT_DP(beginRparts) || tPart.GT_DP(endRparts) {
				continue
			}
			translatedParts = append(translatedParts, tPart)
		}
	}

	var removingD core.Domain

	if len(translatedParts) != 0 {
		unionD := core.CreateIvDomainUnion(translatedParts)
		removingD = secondOutDomain.DifferenceWithIvDomain(unionD)
		if removingD.IsEmpty() {
			return
		}
	} else {
		removingD = secondOutDomain.Copy()
	}
	chEntry := core.CreateChangeEntryWithValues(secondOutVarId, removingD)
	evt.AddChangeEntry(chEntry)
}

// yinZinXout
func ivsecondInResultInFirstOut(secondInDomain *core.IvDomain,
	resultInDomain *core.IvDomain, firstOutDomain *core.IvDomain,
	firstOutVarId core.VarId, evt *core.ChangeEvent) {

	rparts := firstOutDomain.GetParts()
	if len(rparts) == 0 {
		return
	}
	beginRparts := rparts[0]
	endRparts := rparts[len(rparts)-1]

	translatedParts := make([]*core.IvDomPart, 0)
	for _, resultInP := range resultInDomain.GetParts() {
		for _, secondInP := range secondInDomain.GetParts() {
			tPart := resultInP.SUBTRACT(secondInP)
			if tPart.LT_DP(beginRparts) || tPart.GT_DP(endRparts) {
				continue
			}
			translatedParts = append(translatedParts, tPart)
		}
	}

	var removingD core.Domain

	if len(translatedParts) != 0 {
		unionD := core.CreateIvDomainUnion(translatedParts)
		removingD = firstOutDomain.DifferenceWithIvDomain(unionD)
		if removingD.IsEmpty() {
			return
		}
	} else {
		removingD = firstOutDomain.Copy()
	}
	chEntry := core.CreateChangeEntryWithValues(firstOutVarId, removingD)
	evt.AddChangeEntry(chEntry)
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
	return core.ValuesOfMapVarIdToIvDomain(this.GetAllVars(), this.varidToDomainMap)
}

func (this *Sum) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *Sum) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}

// XplusYneqZ represents the propagator for the constraint X + Y == Z
type PseudoXplusYeqZ struct {
	x, y, z core.VarId
}

func (this *PseudoXplusYeqZ) Clone() *PseudoXplusYeqZ {
	prop := new(PseudoXplusYeqZ)
	prop.x, prop.y, prop.z = this.x, this.y, this.z
	return prop
}

func CreatePseudoXplusYeqZ(x core.VarId, y core.VarId, z core.VarId) *PseudoXplusYeqZ {
	prop := new(PseudoXplusYeqZ)
	prop.x, prop.y, prop.z = x, y, z
	return prop
}
