package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

type XplusYeqZ struct {
	x, y, z           core.VarId
	outCh             chan<- *core.ChangeEvent
	inCh              <-chan *core.ChangeEntry
	varidToDomainMap  map[core.VarId]*core.IvDomain
	id                core.PropId
	store             *core.Store
	iColl             *indexical.IndexicalCollection
	checkingIndexical *indexical.CheckingIndexical
}

func (this *XplusYeqZ) GetIndexicalCollection() *indexical.IndexicalCollection {
	return this.iColl
}

func (this *XplusYeqZ) Start(store *core.Store) {
	indexical.InitProcessConstraint(this, false)
	indexical.ProcessConstraint(this, false)
}

func (this *XplusYeqZ) GetAllVars() []core.VarId {
	return []core.VarId{this.x, this.y, this.z}
}

// Register generates auxiliary variables and makes pseudo structs
// and all vars will be registered at store and get domains and channels
func (this *XplusYeqZ) Register(store *core.Store) {
	allvars := this.GetAllVars()
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(allvars, this.id)
	varidToDomainMap := core.GetVaridToIntervalDomains(domains)
	this.Init(store, varidToDomainMap)
}

// bounds-consistency: X in min(Z)-max(Y)..max(Z)-min(Y), Y in min(Z)-max(X)..max(Z)-min(X),
// Z in min(Y)+min(X)..max(Y)+max(X)
func (this *XplusYeqZ) MakeXplusYeqZBoundsIndexicals(varidToDomainMap map[core.VarId]*core.IvDomain) []*indexical.Indexical {
	var xDom, yDom, zDom *core.IvDomain
	var minXTerm, maxXTerm, minYTerm, maxYTerm, minZTerm, maxZTerm ixterm.ITerm
	var subT1, subT2 *ixterm.SubtractionTerm
	var addT1, addT2 *ixterm.AdditionTerm
	var r *ixrange.FromToRange
	var ind *indexical.Indexical

	indexicals := make([]*indexical.Indexical, 3)
	X := this.x
	Y := this.y
	Z := this.z

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
	indexicals[0] = ind

	subT1 = ixterm.CreateSubtractionTerm(minZTerm, maxXTerm)
	subT2 = ixterm.CreateSubtractionTerm(maxZTerm, minXTerm)
	r = ixrange.CreateFromToRange(subT1, subT2)
	ind = indexical.CreateIndexical(Y, yDom, r)
	indexicals[1] = ind

	addT1 = ixterm.CreateAdditionTerm(minYTerm, minXTerm)
	addT2 = ixterm.CreateAdditionTerm(maxYTerm, maxXTerm)
	r = ixrange.CreateFromToRange(addT1, addT2)
	ind = indexical.CreateIndexical(Z, zDom, r)
	indexicals[2] = ind

	return indexicals
}

// arc-consistency: X in dom(Z) - dom(Y), Y in dom(Z) - dom(X), Z in dom(Y) + dom(X)
func (this *XplusYeqZ) MakeXplusYeqZArcIndexicals(varidToDomainMap map[core.VarId]*core.IvDomain) []*indexical.Indexical {
	var xDom, yDom, zDom *core.IvDomain
	var domXRange, domYRange, domZRange *ixrange.DomRange
	var subR *ixrange.SubRange
	var addR *ixrange.AddRange
	var ind *indexical.Indexical
	indexicals := make([]*indexical.Indexical, 3)
	X := this.x
	Y := this.y
	Z := this.z
	xDom = varidToDomainMap[X]
	yDom = varidToDomainMap[Y]
	zDom = varidToDomainMap[Z]
	domXRange = ixrange.CreateDomRange(X, xDom)
	domYRange = ixrange.CreateDomRange(Y, yDom)
	domZRange = ixrange.CreateDomRange(Z, zDom)
	subR = ixrange.CreateSubRange(domZRange, domYRange)
	ind = indexical.CreateIndexical(X, xDom, subR)
	indexicals[0] = ind
	subR = ixrange.CreateSubRange(domZRange, domXRange)
	ind = indexical.CreateIndexical(Y, yDom, subR)
	indexicals[1] = ind
	addR = ixrange.CreateAddRange(domXRange, domYRange)
	ind = indexical.CreateIndexical(Z, zDom, addR)
	indexicals[2] = ind
	return indexicals
}

func (this *XplusYeqZ) IsEntailed() bool {
	return this.checkingIndexical.IsEntailed()
}

func (this *XplusYeqZ) Init(store *core.Store, domains map[core.VarId]*core.IvDomain) {
	this.store = store
	this.varidToDomainMap = domains
	this.iColl = indexical.CreateIndexicalCollection()
	this.iColl.AddIndexicalsAtPrio(this.MakeXplusYeqZBoundsIndexicals(this.varidToDomainMap),
		indexical.HIGHEST)
	arcPropagatingIndexicals := this.MakeXplusYeqZArcIndexicals(this.varidToDomainMap)
	this.iColl.AddIndexicalsAtPrio(arcPropagatingIndexicals, indexical.HIGH)
	this.checkingIndexical = arcPropagatingIndexicals[0].GetCheckingIndexical()
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen
func (this *XplusYeqZ) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XplusYeqZ) GetID() core.PropId {
	return this.id
}

func CreateXplusYeqZ(x, y, z core.VarId) *XplusYeqZ {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXplusYeqZ-propagator")
	}
	prop := new(XplusYeqZ)
	prop.x = x
	prop.y = y
	prop.z = z

	return prop
}

func (this *XplusYeqZ) GetNegation() indexical.IReifiableConstraint {
	return CreateXplusYneqZ(this.x, this.y, this.z)
}

func (this *XplusYeqZ) Clone() core.Constraint {
	prop := new(XplusYeqZ)
	prop.x = this.x
	prop.y = this.y
	prop.z = this.z

	return prop
}

func (this *XplusYeqZ) String() string {
	return fmt.Sprintf("PROP_%d %s + %s = %s",
		this.id, this.store.GetName(this.x),
		this.store.GetName(this.y), this.store.GetName(this.z))
}

func (this *XplusYeqZ) GetVarIds() []core.VarId {
	return this.GetAllVars()
}

func (this *XplusYeqZ) GetDomains() []core.Domain {
	return core.ValuesOfMapVarIdToIvDomain(this.GetAllVars(), this.varidToDomainMap)
}

func (this *XplusYeqZ) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XplusYeqZ) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
