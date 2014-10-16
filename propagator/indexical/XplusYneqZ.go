package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

// idea: X+Y+Q=Z, --> X+Y=H1, --> H1+Q=Z

type XplusYneqZ struct {
	x, y, z           				core.VarId
	x_Domain, y_Domain, z_Domain 	*core.IvDomain
	outCh             				chan<- *core.ChangeEvent
	inCh              				<-chan *core.ChangeEntry
	id                				core.PropId
	store             				*core.Store
	iColl             				*indexical.IndexicalCollection
	checkingIndexical 				*indexical.CheckingIndexical
}

func (this *XplusYneqZ) GetIndexicalCollection() *indexical.IndexicalCollection {
	return this.iColl
}

func (this *XplusYneqZ) Start(store *core.Store) {
	core.LogInitConsistency(this)
	iColl := this.GetIndexicalCollection()
	evt := indexical.ProcessIndexicals(iColl, nil, true)
	core.SendChangesToStore(evt, this)
	for changeEntry := range this.inCh {
		indexical.RemoveValues(this, changeEntry)
		iColl := this.GetIndexicalCollection()
		evt = indexical.ProcessIndexicals(iColl, changeEntry, true)
		core.SendChangesToStore(evt, this)
	}
}

func (this *XplusYneqZ) GetAllVars() []core.VarId {
	return []core.VarId{this.x, this.y, this.z}
}

// Register generates auxiliary variables and makes pseudo structs
// and all vars will be registered at store and get domains and channels
func (this *XplusYneqZ) Register(store *core.Store) {
	allvars := this.GetAllVars()
	var domains []core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagator(allvars, this.id)

	this.Init(store, domains)
}

// arc-consistency: X in not(val(Z)-val(Y)), Y in not(val(Z)-val(X)),
// Z in not(val(X)+val(Y))
func (this *XplusYneqZ) MakeXplusYneqZArcIndexicals() []*indexical.Indexical {
	indexicals := make([]*indexical.Indexical, 3)

	// X in not(val(Z)-val(Y))
	valZ := ixterm.CreateValTerm(this.z, this.z_Domain)
	valY := ixterm.CreateValTerm(this.y, this.y_Domain)
	subT := ixterm.CreateSubtractionTerm(valZ, valY)
	notR := ixrange.CreateNotRange(ixrange.CreateSingleValueRange(subT))
	ind := indexical.CreateIndexical(this.x, this.x_Domain, notR)
	indexicals[0] = ind
	//Y in not(val(Z)-val(X))
	valZ = ixterm.CreateValTerm(this.z, this.z_Domain)
	valX := ixterm.CreateValTerm(this.x, this.x_Domain)
	subT = ixterm.CreateSubtractionTerm(valZ, valX)
	notR = ixrange.CreateNotRange(ixrange.CreateSingleValueRange(subT))
	ind = indexical.CreateIndexical(this.y, this.y_Domain, notR)
	indexicals[1] = ind
	//Z in not(val(X)+val(Y))
	valY = ixterm.CreateValTerm(this.y, this.y_Domain)
	valX = ixterm.CreateValTerm(this.x, this.x_Domain)
	addT := ixterm.CreateAdditionTerm(valY, valX)
	notR = ixrange.CreateNotRange(ixrange.CreateSingleValueRange(addT))
	ind = indexical.CreateIndexical(this.z, this.z_Domain, notR)
	indexicals[2] = ind
	return indexicals
}

func (this *XplusYneqZ) IsEntailed() bool {
	return this.checkingIndexical.IsEntailed()
}

func (this *XplusYneqZ) Init(store *core.Store, domains []core.Domain) {
	this.store = store
	this.x_Domain = core.GetVaridToIntervalDomain(domains[0])
	this.y_Domain = core.GetVaridToIntervalDomain(domains[1])
	this.z_Domain = core.GetVaridToIntervalDomain(domains[2])
	this.iColl = indexical.CreateIndexicalCollection()
	arcIndexicals := this.MakeXplusYneqZArcIndexicals()
	this.checkingIndexical = arcIndexicals[0].GetCheckingIndexical()
	this.iColl.AddIndexicalsAtPrio(arcIndexicals, indexical.HIGH)
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen
func (this *XplusYneqZ) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XplusYneqZ) GetID() core.PropId {
	return this.id
}

func CreateXplusYneqZ(x, y, z core.VarId) *XplusYneqZ {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXplusYneqZ-propagator")
	}
	prop := new(XplusYneqZ)
	prop.x = x
	prop.y = y
	prop.z = z
	return prop
}

func (this *XplusYneqZ) GetNegation() indexical.IReifiableConstraint {
	return CreateXplusYeqZ(this.x, this.y, this.z)
}

func (this *XplusYneqZ) Clone() core.Constraint {
	prop := new(XplusYneqZ)
	prop.x = this.x
	prop.y = this.y
	prop.z = this.z
	return prop
}

func (this *XplusYneqZ) String() string {
	return fmt.Sprintf("PROP_%d %s + %s = %s",
		this.id, this.store.GetName(this.x),
		this.store.GetName(this.y), this.store.GetName(this.z))
}

func (this *XplusYneqZ) GetVarIds() []core.VarId {
	return this.GetAllVars()
}

func (this *XplusYneqZ) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain,this.y_Domain,this.z_Domain}
}

func (this *XplusYneqZ) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XplusYneqZ) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
