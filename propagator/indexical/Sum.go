package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
	"strings"
)

// idea: X+Y+Q=Z, --> X+Y=H1, --> H1+Q=Z

type Sum struct {
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

func (this *Sum) GetIndexicalCollection() *indexical.IndexicalCollection {
	return this.iColl
}

func (this *Sum) Start() {
	indexical.InitProcessConstraint(this, true)
	indexical.ProcessConstraint(this, true)
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

// Register generates auxiliary variables and makes pseudo structs
// and all vars will be registered at store and get domains and channels
func (this *Sum) Register(store *core.Store) {
	allvars := this.GetAllVars()

	var domains map[core.VarId]core.Domain

	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(allvars, this.id)

	this.varidToDomainMap = core.GetVaridToIntervalDomains(domains)

	this.store = store

	this.iColl = indexical.CreateIndexicalCollection()

	this.iColl.AddIndexicalsAtPrio(MakeSumBoundsIndexicals(this.pseudoProps, this.varidToDomainMap), indexical.HIGHEST)
	this.iColl.AddIndexicalsAtPrio(MakeSumArcIndexicals(this.pseudoProps, this.varidToDomainMap), indexical.HIGH)
}

// bounds-consistency: X in min(Z)-max(Y)..max(Z)-min(Y), Y in min(Z)-max(X)..max(Z)-min(X),
// Z in min(Y)+min(X)..max(Y)+max(X)
func MakeSumBoundsIndexicals(pseudoProps []*XplusYeqZ_Rel,
	varidToDomainMap map[core.VarId]*core.IvDomain) []*indexical.Indexical {
	var xDom, yDom, zDom *core.IvDomain
	var minXTerm, maxXTerm, minYTerm, maxYTerm, minZTerm, maxZTerm ixterm.ITerm
	var subT1, subT2 *ixterm.SubtractionTerm
	var addT1, addT2 *ixterm.AdditionTerm
	var r *ixrange.FromToRange
	var ind *indexical.Indexical
	indexicals := make([]*indexical.Indexical, len(pseudoProps)*3)
	i := 0
	for _, p := range pseudoProps {
		X := p.x
		Y := p.y
		Z := p.z
		xDom = varidToDomainMap[X]
		yDom = varidToDomainMap[Y]
		zDom = varidToDomainMap[Z]
		minXTerm = ixterm.CreateMinTerm(X, xDom)
		maxXTerm = ixterm.CreateMaxTerm(X, xDom)
		minYTerm = ixterm.CreateMinTerm(Y, yDom)
		maxYTerm = ixterm.CreateMaxTerm(Y, yDom)
		minZTerm = ixterm.CreateMinTerm(Z, zDom)
		maxZTerm = ixterm.CreateMaxTerm(Z, zDom)
		subT1 = ixterm.CreateSubtractionTerm(minZTerm, maxYTerm)
		subT2 = ixterm.CreateSubtractionTerm(maxZTerm, minYTerm)
		r = ixrange.CreateFromToRange(subT1, subT2)
		ind = indexical.CreateIndexical(X, xDom, r)
		indexicals[i] = ind
		i += 1
		subT1 = ixterm.CreateSubtractionTerm(minZTerm, maxXTerm)
		subT2 = ixterm.CreateSubtractionTerm(maxZTerm, minXTerm)
		r = ixrange.CreateFromToRange(subT1, subT2)
		ind = indexical.CreateIndexical(Y, yDom, r)
		indexicals[i] = ind
		i += 1
		addT1 = ixterm.CreateAdditionTerm(minYTerm, minXTerm)
		addT2 = ixterm.CreateAdditionTerm(maxYTerm, maxXTerm)
		r = ixrange.CreateFromToRange(addT1, addT2)
		ind = indexical.CreateIndexical(Z, zDom, r)
		indexicals[i] = ind
		i += 1
	}
	return indexicals
}

// arc-consistency: X in dom(Z) - dom(Y), Y in dom(Z) - dom(X), Z in dom(Y) + dom(X)
func MakeSumArcIndexicals(pseudoProps []*XplusYeqZ_Rel,
	varidToDomainMap map[core.VarId]*core.IvDomain) []*indexical.Indexical {
	var xDom, yDom, zDom *core.IvDomain
	var domXRange, domYRange, domZRange *ixrange.DomRange
	var subR *ixrange.SubRange
	var addR *ixrange.AddRange
	var ind *indexical.Indexical
	indexicals := make([]*indexical.Indexical, len(pseudoProps)*3)
	i := 0
	for _, P := range pseudoProps {
		X := P.x
		Y := P.y
		Z := P.z
		xDom = varidToDomainMap[X]
		yDom = varidToDomainMap[Y]
		zDom = varidToDomainMap[Z]
		domXRange = ixrange.CreateDomRange(X, xDom)
		domYRange = ixrange.CreateDomRange(Y, yDom)
		domZRange = ixrange.CreateDomRange(Z, zDom)
		subR = ixrange.CreateSubRange(domZRange, domYRange)
		ind = indexical.CreateIndexical(X, xDom, subR)
		indexicals[i] = ind
		i += 1
		subR = ixrange.CreateSubRange(domZRange, domXRange)
		ind = indexical.CreateIndexical(Y, yDom, subR)
		indexicals[i] = ind
		i += 1
		addR = ixrange.CreateAddRange(domXRange, domYRange)
		ind = indexical.CreateIndexical(Z, zDom, addR)
		indexicals[i] = ind
		i += 1
	}
	return indexicals
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

func (this *Sum) Clone() core.Constraint {
	prop := new(Sum)
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

func (this *Sum) String() string {
	vars_str := make([]string, len(this.vars))
	for i, var_id := range this.vars {
		vars_str[i] = this.store.GetName(var_id)
	}
	return fmt.Sprintf("PROP_%d %s = %s",
		this.id, strings.Join(vars_str, "+"),
		this.store.GetName(this.resultVar))
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
type XplusYeqZ_Rel struct {
	x, y, z core.VarId
}

func (this *XplusYeqZ_Rel) Clone() *XplusYeqZ_Rel {
	prop := new(XplusYeqZ_Rel)
	prop.x, prop.y, prop.z = this.x, this.y, this.z
	return prop
}

func CreateXplusYeqZ_Rel(x core.VarId, y core.VarId, z core.VarId) *XplusYeqZ_Rel {
	prop := new(XplusYeqZ_Rel)
	prop.x, prop.y, prop.z = x, y, z
	return prop
}
