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
	x, y, z           core.VarId
	outCh             chan<- *core.ChangeEvent
	inCh              <-chan *core.ChangeEntry
	varidToDomainMap  map[core.VarId]*core.IvDomain
	id                core.PropId
	store             *core.Store
	iColl             *indexical.IndexicalCollection
	checkingIndexical *indexical.CheckingIndexical
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
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(allvars, this.id)

	varidToDomainMap := core.GetVaridToIntervalDomains(domains)

	this.Init(store, varidToDomainMap)
}

// arc-consistency: X in not(val(Z)-val(Y)), Y in not(val(Z)-val(X)),
// Z in not(val(X)+val(Y))
func (this *XplusYneqZ) MakeXplusYneqZArcIndexicals(varidToDomainMap map[core.VarId]*core.IvDomain) []*indexical.Indexical {
	var xDom, yDom, zDom *core.IvDomain
	indexicals := make([]*indexical.Indexical, 3)
	X := this.x
	Y := this.y
	Z := this.z
	xDom = varidToDomainMap[X]
	yDom = varidToDomainMap[Y]
	zDom = varidToDomainMap[Z]
	// X in not(val(Z)-val(Y))
	valZ := ixterm.CreateValTerm(Z, zDom)
	valY := ixterm.CreateValTerm(Y, yDom)
	subT := ixterm.CreateSubtractionTerm(valZ, valY)
	notR := ixrange.CreateNotRange(ixrange.CreateSingleValueRange(subT))
	ind := indexical.CreateIndexical(X, xDom, notR)
	indexicals[0] = ind
	//Y in not(val(Z)-val(X))
	valZ = ixterm.CreateValTerm(Z, zDom)
	valX := ixterm.CreateValTerm(X, xDom)
	subT = ixterm.CreateSubtractionTerm(valZ, valX)
	notR = ixrange.CreateNotRange(ixrange.CreateSingleValueRange(subT))
	ind = indexical.CreateIndexical(Y, yDom, notR)
	indexicals[1] = ind
	//Z in not(val(X)+val(Y))
	valY = ixterm.CreateValTerm(Y, yDom)
	valX = ixterm.CreateValTerm(X, xDom)
	addT := ixterm.CreateAdditionTerm(valY, valX)
	notR = ixrange.CreateNotRange(ixrange.CreateSingleValueRange(addT))
	ind = indexical.CreateIndexical(Z, zDom, notR)
	indexicals[2] = ind
	return indexicals
}

func (this *XplusYneqZ) IsEntailed() bool {
	return this.checkingIndexical.IsEntailed()
}

func (this *XplusYneqZ) Init(store *core.Store, domains map[core.VarId]*core.IvDomain) {
	this.store = store
	this.varidToDomainMap = domains
	this.iColl = indexical.CreateIndexicalCollection()
	arcIndexicals := this.MakeXplusYneqZArcIndexicals(this.varidToDomainMap)
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
	return core.ValuesOfMapVarIdToIvDomain(this.GetAllVars(), this.varidToDomainMap)
}

func (this *XplusYneqZ) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XplusYneqZ) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
