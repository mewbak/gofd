package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

// XeqC represents the propagator for the indexical constraint X = C
type XeqC struct {
	x        core.VarId
	c        int
	x_Domain *core.IvDomain
	iColl    *indexical.IndexicalCollection
	outCh    chan<- *core.ChangeEvent
	inCh     <-chan *core.ChangeEntry
	id       core.PropId
	store    *core.Store
}

func (this *XeqC) GetIndexicalCollection() *indexical.IndexicalCollection {
	return this.iColl
}

func (this *XeqC) Start() {
	indexical.InitProcessConstraint(this, false)
	indexical.ProcessConstraint(this, false)
}

// Register registers the propagator at the store.
func (this *XeqC) Register(store *core.Store) {
	var domains []core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagator([]core.VarId{this.x}, this.id)

	this.Init(store, domains)
	this.store = store
}

func (this *XeqC) Clone() core.Constraint {
	prop := new(XeqC)
	prop.x = this.x
	prop.c = this.c
	return prop
}

func (this *XeqC) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XeqC) GetID() core.PropId {
	return this.id
}

func (this *XeqC) String() string {
	return fmt.Sprintf("XeqC_%d %s = %d",
		this.id, this.store.GetName(this.x), this.c)
}

func (this *XeqC) GetVarIds() []core.VarId {
	return []core.VarId{this.x}
}

func (this *XeqC) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain}
}

func (this *XeqC) Init(store *core.Store, domains []core.Domain) {
	this.store = store

	this.x_Domain = core.GetVaridToIntervalDomain(domains[0])

	this.iColl = indexical.CreateIndexicalCollection()

	cTerm := ixterm.CreateValueTerm(this.c)
	r := ixrange.CreateSingleValueRange(cTerm)
	this.iColl.CreateAndAddIndexical(this.x, this.x_Domain, r, indexical.HIGHEST)
}

func (this *XeqC) IsEntailed() bool {
	// ToDo: use checking indexical
	if this.x_Domain.IsGround() {
		return this.x_Domain.Contains(this.c)
	}
	return false
}

func (this *XeqC) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XeqC) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}

func (this *XeqC) GetNegation() indexical.IReifiableConstraint {
	return CreateXneqC(this.x, this.c)
}

func CreateXeqC(x core.VarId, c int) *XeqC {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXeqC")
	}
	prop := new(XeqC)
	prop.x = x
	prop.c = c
	return prop
}
