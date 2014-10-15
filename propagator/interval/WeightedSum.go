package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
)

// WeightedSum represents the constraint C1*X1+C2*X2+...+Cn*Xn=Z
// Its propagate functions establish arc consistency (with bounds and arc
// algorithms).
// The basic idea of WeightedSum is to substitute the WeightedSum equation to
// many Ci*Xi=Hi, so that a Sum constraint results with H1+H2+...+Hn=Z. This
// Sum constraint is substituted as well with X+Y=Z equations (see Sum
// constraint for more information).
type WeightedSum struct {
	vars             []core.VarId
	hvars            []core.VarId //helper-variables
	cs               []int
	resultVar        core.VarId
	outCh            chan<- *core.ChangeEvent
	inCh             <-chan *core.ChangeEntry
	varidToDomainMap map[core.VarId]*core.IvDomain
	id               core.PropId
	store            *core.Store
	pseudoPropsXCY   []*PseudoXmultCeqY
	pseudoPropsXYZ   []*PseudoXplusYeqZ
}

func (this *WeightedSum) Start(store *core.Store) {

	// initial check
	evt := core.CreateChangeEvent()
	this.ivweightSumInitialCheck(evt)
	core.SendChangesToStore(evt, this)

	for changeEntry := range this.inCh {
		core.LogIncomingChange(this, store, changeEntry)

		evt = core.CreateChangeEvent()
		varidChanged := changeEntry.GetID()
		changedDom := this.varidToDomainMap[varidChanged]
		changedDom.Removes(changeEntry.GetValues())
		this.ivweightedSumPropagate(varidChanged, evt)
		core.SendChangesToStore(evt, this)
	}
}

// propagate check for changes. First look for X*C=Y propagators, then
// X+Y=Z and finally for the final propagator X=Y. Collect changes
func (this *WeightedSum) ivweightedSumPropagate(varid core.VarId, evt *core.ChangeEvent) {

	this.checkXmultCeqY(varid, evt)

	this.ivsumPropagate(varid, evt)
}

func (this *WeightedSum) ivsumPropagate(varid core.VarId, evt *core.ChangeEvent) {
	hvarXCY := make([]core.VarId, len(this.pseudoPropsXCY))
	for i, hxcy := range this.pseudoPropsXCY {
		hvarXCY[i] = hxcy.y
	}

	ivsumBoundsPropagate(varid, this.varidToDomainMap, this.pseudoPropsXYZ, evt)
	ivsumArcPropagate(varid, this.varidToDomainMap, this.pseudoPropsXYZ, evt)
}

// initialCheck check for changes. First look for X*C=Y propagators, then
// X+Y=Z and finally for the final propagator X=Y. Collect changes
func (this *WeightedSum) ivweightSumInitialCheck(evt *core.ChangeEvent) {

	this.checkXmultCeqY(-1, evt)

	hvarXCY := make([]core.VarId, len(this.pseudoPropsXCY))
	for i, hxcy := range this.pseudoPropsXCY {
		hvarXCY[i] = hxcy.y
	}

	ivsumBoundsInitialCheck(this.varidToDomainMap, this.pseudoPropsXYZ, evt)
	ivsumArcInitialCheck(this.varidToDomainMap, this.pseudoPropsXYZ, evt)
}

func (this *WeightedSum) checkXmultCeqY(varid core.VarId, evt *core.ChangeEvent) {
	for i, prop := range this.pseudoPropsXCY {
		xDom := this.varidToDomainMap[prop.x]
		yDom := this.varidToDomainMap[prop.y]
		c := this.cs[i]
		if varid == -1 {
			xDom := this.varidToDomainMap[prop.x]
			yDom := this.varidToDomainMap[prop.y]
			c := this.cs[i]

			ivfirstInMultSecondOutBOUNDS(xDom, c, yDom, prop.y, evt)
			ivsecondInMultFirstOutBOUNDS(yDom, c, xDom, prop.x, evt)
			ivfirstInMultSecondOutARC(xDom, c, yDom, prop.y, evt)
			ivsecondInMultFirstOutARC(yDom, c, xDom, prop.x, evt)

		} else if prop.x == varid {
			ivfirstInMultSecondOutBOUNDS(xDom, c, yDom, prop.y, evt)
			ivfirstInMultSecondOutARC(xDom, c, yDom, prop.y, evt)
		} else if prop.y == varid {
			ivsecondInMultFirstOutBOUNDS(yDom, c, xDom, prop.x, evt)
			ivsecondInMultFirstOutARC(yDom, c, xDom, prop.x, evt)
		}
	}
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
	this.varidToDomainMap = core.GetVaridToIntervalDomains(domains)
	this.store = store
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
	prop.pseudoPropsXCY = make([]*PseudoXmultCeqY, len(prop.vars))
	prop.hvars = make([]core.VarId, 0)
	for i, X := range prop.vars {
		H := core.CreateAuxIntVarIvValues(store,
			core.ScalarSlice(prop.cs[i], store.GetDomain(X).Values_asSlice()))
		prop.pseudoPropsXCY[i] = CreatePseudoXmultCeqY(X, prop.cs[i], H)
		prop.hvars = append(prop.hvars, H)
	}
	prop.pseudoPropsXYZ = make([]*PseudoXplusYeqZ, len(prop.pseudoPropsXCY)-1)
	H := prop.pseudoPropsXCY[0].y
	newHVars := make([]core.VarId, 0)
	for i, p := range prop.pseudoPropsXCY[1 : len(prop.vars)-1] {
		NewH := core.CreateAuxIntVarIvFromTo(store,
			store.GetDomain(H).GetMin()+store.GetDomain(p.y).GetMin(),
			store.GetDomain(H).GetMax()+store.GetDomain(p.y).GetMax())
		prop.pseudoPropsXYZ[i] = CreatePseudoXplusYeqZ(H, p.y, NewH)
		H = NewH
		newHVars = append(newHVars, NewH)
	}
	X := prop.hvars[len(prop.hvars)-1]
	prop.hvars = append(prop.hvars, newHVars...)
	prop.pseudoPropsXYZ[len(prop.pseudoPropsXYZ)-1] = CreatePseudoXplusYeqZ(H, X, prop.resultVar)

	return prop
}

func (this *WeightedSum) String() string {
	vars_str := make([]string, len(this.vars))
	for i, var_id := range this.vars {
		vars_str[i] = fmt.Sprintf("%v*%s",
			this.cs[i], this.store.GetName(var_id))
	}
	return fmt.Sprintf("PROP_%d %s = %s",
		this.id, strings.Join(vars_str, "+"),
		this.store.GetName(this.resultVar))
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
	prop.pseudoPropsXCY = make([]*PseudoXmultCeqY, len(this.pseudoPropsXCY))
	for i, p := range this.pseudoPropsXCY {
		prop.pseudoPropsXCY[i] = p.Clone()
	}
	prop.pseudoPropsXYZ = make([]*PseudoXplusYeqZ, len(this.pseudoPropsXYZ))
	for i, p := range this.pseudoPropsXYZ {
		prop.pseudoPropsXYZ[i] = p.Clone()
	}
	prop.hvars = make([]core.VarId, len(this.hvars))
	for i, single_var := range this.hvars {
		prop.hvars[i] = single_var
	}

	return prop
}

// INFO: the following two propagation-functions do not work on special
// interval-operations. Reason is: the multiplication is always worst case!
// example-CSP:
// X*C=Y
// X:[1,11]
// C:2
// Y:[1,20]

// RESULT:
// X:[1,10]
// Y:[2] [4] [6] ... [20]

// X*C=Y
// ivfirstInMultSecondOutARC collect changes, when first variable has changed
// e.g. X=Y/C, then X is first variable
func ivfirstInMultSecondOutBOUNDS(firstInDomain *core.IvDomain, c int,
	secondOutDomain *core.IvDomain, secondOutVarId core.VarId,
	evt *core.ChangeEvent) {

	if firstInDomain.IsEmpty() {
		return
	}

	removeParts := make([]*core.IvDomPart, 0)

	minX := firstInDomain.GetMin() * c
	maxX := firstInDomain.GetMax() * c

	for i, yPart := range secondOutDomain.GetParts() {
		removeParts = append(removeParts, yPart.DIFFERENCE_MIN_MAX(minX, maxX)...)
		if yPart.GT(maxX) {
			removeParts = append(removeParts, secondOutDomain.GetParts()[i:]...)
			break
		}
	}

	if len(removeParts) != 0 {
		chEntry := core.CreateChangeEntryWithValues(secondOutVarId, core.CreateIvDomainDomParts(removeParts))
		evt.AddChangeEntry(chEntry)
	}
}

// ivsecondInMultFirstOutARC collect changes, when second variable has changed
// e.g. X*C=Y, then Y is second variable
func ivsecondInMultFirstOutBOUNDS(secondInDomain *core.IvDomain, c int,
	firstOutDomain *core.IvDomain, firstOutVarId core.VarId,
	evt *core.ChangeEvent) {
	removeParts := make([]*core.IvDomPart, 0)

	if secondInDomain.IsEmpty() {
		return
	}

	minY := secondInDomain.GetMin() / c
	maxY := secondInDomain.GetMax() / c

	for i, xPart := range firstOutDomain.GetParts() {

		removeParts = append(removeParts, xPart.DIFFERENCE_MIN_MAX(minY, maxY)...)
		if xPart.GT(maxY) {
			removeParts = append(removeParts, firstOutDomain.GetParts()[i:]...)
			break
		}
	}

	if len(removeParts) != 0 {
		chEntry := core.CreateChangeEntryWithValues(firstOutVarId, core.CreateIvDomainDomParts(removeParts))
		evt.AddChangeEntry(chEntry)
	}
}

// X*C=Y
// ivfirstInMultSecondOutARC collect changes, when first variable has changed
// e.g. Y=X/C, then X is first variable
func ivfirstInMultSecondOutARC(firstInDomain *core.IvDomain, c int,
	secondOutDomain *core.IvDomain, secondOutVarId core.VarId,
	evt *core.ChangeEvent) {

	if firstInDomain.IsEmpty() {
		return
	}

	vals := make([]int, 0)
	for _, y_val := range secondOutDomain.GetValues() {
		if c == 0 || !(y_val%c == 0) {
			if y_val != 0 {
				vals = append(vals, y_val)
			}
		} else {
			if !firstInDomain.Contains(y_val / c) {
				vals = append(vals, y_val)
			}
		}
	}

	if len(vals) != 0 {
		chEntry := core.CreateChangeEntryWithIntValues(secondOutVarId, vals)
		evt.AddChangeEntry(chEntry)
	}
}

// ivsecondInMultFirstOutARC collect changes, when second variable has changed
// e.g. X*C=Y, then Y is second variable
func ivsecondInMultFirstOutARC(secondInDomain *core.IvDomain, c int,
	firstOutDomain *core.IvDomain, firstOutVarId core.VarId,
	evt *core.ChangeEvent) {

	if secondInDomain.IsEmpty() {
		return
	}

	vals := make([]int, 0)
	for _, x_val := range firstOutDomain.GetValues() {
		if !secondInDomain.Contains(x_val * c) {
			vals = append(vals, x_val)
		}
	}

	if len(vals) != 0 {
		chEntry := core.CreateChangeEntryWithIntValues(firstOutVarId, vals)
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
	return core.ValuesOfMapVarIdToIvDomain(this.GetAllVars(), this.varidToDomainMap)
}

func (this *WeightedSum) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *WeightedSum) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}

// XplusYneqZ represents the propagator for the constraint X + Y == Z
type PseudoXmultCeqY struct {
	x, y core.VarId
	c    int
}

func (this *PseudoXmultCeqY) Clone() *PseudoXmultCeqY {
	prop := new(PseudoXmultCeqY)
	prop.x, prop.c, prop.y = this.x, this.c, this.y
	return prop
}

func CreatePseudoXmultCeqY(x core.VarId, c int, y core.VarId) *PseudoXmultCeqY {
	prop := new(PseudoXmultCeqY)
	prop.x, prop.c, prop.y = x, c, y
	return prop
}
