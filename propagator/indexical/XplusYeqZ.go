package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

type XplusYeqZ struct {
	x, y, z           				core.VarId
	x_Domain, y_Domain, z_Domain 	*core.IvDomain
	outCh             				chan<- *core.ChangeEvent
	inCh              				<-chan *core.ChangeEntry
	id                				core.PropId
	store             				*core.Store
	iColl             				*indexical.IndexicalCollection
	checkingIndexical 				*indexical.CheckingIndexical
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
	var domains []core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagator(allvars, this.id)
	
	this.Init(store, domains)
}

// bounds-consistency: X in min(Z)-max(Y)..max(Z)-min(Y), Y in min(Z)-max(X)..max(Z)-min(X),
// Z in min(Y)+min(X)..max(Y)+max(X)
func (this *XplusYeqZ) MakeXplusYeqZBoundsIndexicals() []*indexical.Indexical {

	indexicals := make([]*indexical.Indexical, 3)
	minXTerm := ixterm.CreateMinTerm(this.x, this.x_Domain)
	maxXTerm := ixterm.CreateMaxTerm(this.x, this.x_Domain)
	minYTerm := ixterm.CreateMinTerm(this.y, this.y_Domain)
	maxYTerm := ixterm.CreateMaxTerm(this.y, this.y_Domain)
	minZTerm := ixterm.CreateMinTerm(this.z, this.z_Domain)
	maxZTerm := ixterm.CreateMaxTerm(this.z, this.z_Domain)

	subT1 := ixterm.CreateSubtractionTerm(minZTerm, maxYTerm)
	subT2 := ixterm.CreateSubtractionTerm(maxZTerm, minYTerm)
	r := ixrange.CreateFromToRange(subT1, subT2)
	ind := indexical.CreateIndexical(this.x, this.x_Domain, r)
	indexicals[0] = ind

	subT1 = ixterm.CreateSubtractionTerm(minZTerm, maxXTerm)
	subT2 = ixterm.CreateSubtractionTerm(maxZTerm, minXTerm)
	r = ixrange.CreateFromToRange(subT1, subT2)
	ind = indexical.CreateIndexical(this.y, this.y_Domain, r)
	indexicals[1] = ind

	addT1 := ixterm.CreateAdditionTerm(minYTerm, minXTerm)
	addT2 := ixterm.CreateAdditionTerm(maxYTerm, maxXTerm)
	r = ixrange.CreateFromToRange(addT1, addT2)
	ind = indexical.CreateIndexical(this.z, this.z_Domain, r)
	indexicals[2] = ind

	return indexicals
}

// arc-consistency: X in dom(Z) - dom(Y), Y in dom(Z) - dom(X), Z in dom(Y) + dom(X)
func (this *XplusYeqZ) MakeXplusYeqZArcIndexicals() []*indexical.Indexical {
	
	indexicals := make([]*indexical.Indexical, 3)
	domXRange := ixrange.CreateDomRange(this.x, this.x_Domain)
	domYRange := ixrange.CreateDomRange(this.y, this.y_Domain)
	domZRange := ixrange.CreateDomRange(this.z, this.z_Domain)
	subR := ixrange.CreateSubRange(domZRange, domYRange)
	ind := indexical.CreateIndexical(this.x, this.x_Domain, subR)
	indexicals[0] = ind
	subR = ixrange.CreateSubRange(domZRange, domXRange)
	ind = indexical.CreateIndexical(this.y, this.y_Domain, subR)
	indexicals[1] = ind
	addR := ixrange.CreateAddRange(domXRange, domYRange)
	ind = indexical.CreateIndexical(this.z, this.z_Domain, addR)
	indexicals[2] = ind
	return indexicals
}

func (this *XplusYeqZ) IsEntailed() bool {
	return this.checkingIndexical.IsEntailed()
}

func (this *XplusYeqZ) Init(store *core.Store, domains []core.Domain) {
	this.store = store
	this.x_Domain = core.GetVaridToIntervalDomain(domains[0])
	this.y_Domain = core.GetVaridToIntervalDomain(domains[1])
	this.z_Domain = core.GetVaridToIntervalDomain(domains[2])
	this.iColl = indexical.CreateIndexicalCollection()
	this.iColl.AddIndexicalsAtPrio(this.MakeXplusYeqZBoundsIndexicals(),
		indexical.HIGHEST)
	arcPropagatingIndexicals := this.MakeXplusYeqZArcIndexicals()
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
	return []core.Domain{this.x_Domain,this.y_Domain,this.z_Domain}
}

func (this *XplusYeqZ) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XplusYeqZ) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
