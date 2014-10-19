package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

//XgtY represents the propagator for the constraint X > Y
type XgtY struct {
	x, y               core.VarId
	x_Domain, y_Domain *core.IvDomain
	iColl              *indexical.IndexicalCollection
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	id                 core.PropId
	store              *core.Store
}

func (this *XgtY) GetIndexicalCollection() *indexical.IndexicalCollection {
	return this.iColl
}

func (this *XgtY) Clone() core.Constraint {
	prop := new(XgtY)
	prop.x = this.x
	prop.y = this.y
	return prop
}

func (this *XgtY) Start() {

	indexical.InitProcessConstraint(this, false)
	indexical.ProcessConstraint(this, false)
}

// Register registers the propagator at the store.
func (this *XgtY) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y}, this.id)

	varidToDomainMap := core.GetVaridToIntervalDomains(domains)

	this.x_Domain = varidToDomainMap[this.x]
	this.y_Domain = varidToDomainMap[this.y]
	this.store = store

	this.iColl = indexical.CreateIndexicalCollection()

	//X>Y --> X in min(Y)+1..inf
	//	  --> Y in -inf..max(X)-1

	minY := ixterm.CreateMinTerm(this.y, this.y_Domain) //Übergabe von InterVar bzw. Domain
	s := ixterm.CreateValueTerm(1)
	addT := ixterm.CreateAdditionTerm(minY, s)
	inf := ixterm.CreateValueTerm(ixterm.INFINITY)
	r := ixrange.CreateFromToRange(addT, inf) //Übergabe von Termen
	this.iColl.CreateAndAddIndexical(this.x, this.x_Domain, r, indexical.HIGHEST)

	maxX := ixterm.CreateMaxTerm(this.x, this.x_Domain)
	s = ixterm.CreateValueTerm(1)
	subT := ixterm.CreateSubtractionTerm(maxX, s)
	neg_inf := ixterm.CreateValueTerm(ixterm.NEG_INFINITY)
	r2 := ixrange.CreateFromToRange(neg_inf, subT)
	this.iColl.CreateAndAddIndexical(this.y, this.y_Domain, r2, indexical.HIGHEST)
}

func (this *XgtY) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XgtY) GetID() core.PropId {
	return this.id
}

func (this *XgtY) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y}
}

func (this *XgtY) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain}
}

func (this *XgtY) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XgtY) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}

// CreateXgtY creates an IndexicalConstraint X>Y
func CreateXgtY(x core.VarId, y core.VarId) *XgtY {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXgtY")
	}
	prop := new(XgtY)

	prop.x = x
	prop.y = y
	return prop
}

func (this *XgtY) String() string {
	return fmt.Sprintf("PROP_%d %s > %s",
		this.id, this.store.GetName(this.x), this.store.GetName(this.y))
}
