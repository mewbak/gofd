package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

// XneqC represents the propagator for the indexical constraint X = C
type XneqC struct {
	x        core.VarId
	c        int
	x_Domain *core.IvDomain
	iColl    *indexical.IndexicalCollection
	outCh    chan<- *core.ChangeEvent
	inCh     <-chan *core.ChangeEntry
	id       core.PropId
	store    *core.Store
}

func (this *XneqC) GetIndexicalCollection() *indexical.IndexicalCollection {
	return this.iColl
}

func (this *XneqC) Start() {
	indexical.InitProcessConstraint(this, false)
	indexical.ProcessConstraint(this, false)
}

// Register registers the propagator at the store.
func (this *XneqC) Register(store *core.Store) {
	var domains []core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagator([]core.VarId{this.x}, this.id)

	this.Init(store, domains)
	this.store = store
}

func (this *XneqC) Init(store *core.Store, domains []core.Domain) {
	this.store = store

	this.x_Domain = core.GetVaridToIntervalDomain(domains[0])

	this.iColl = indexical.CreateIndexicalCollection()

	cTerm := ixterm.CreateValueTerm(this.c)
	r := ixrange.CreateSingleValueRange(cTerm)
	negR := ixrange.CreateNotRange(r)
	this.iColl.CreateAndAddIndexical(this.x, this.x_Domain, negR, indexical.HIGHEST)
}

func (this *XneqC) Clone() core.Constraint {
	prop := new(XneqC)
	prop.x = this.x
	prop.c = this.c
	return prop
}

func (this *XneqC) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XneqC) GetID() core.PropId {
	return this.id
}

func (this *XneqC) String() string {
	return fmt.Sprintf("IC_%d %s != %d",
		this.id, this.store.GetName(this.x), this.c)
}

func (this *XneqC) GetVarIds() []core.VarId {
	return []core.VarId{this.x}
}

func (this *XneqC) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain}
}

func (this *XneqC) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XneqC) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}

func (this XneqC) IsEntailed() bool {
	// ToDo: use checking indexical
	return !this.x_Domain.Contains(this.c)
}

func CreateXneqC(x core.VarId, c int) *XneqC {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXneqC")
	}
	prop := new(XneqC)
	prop.x = x
	prop.c = c
	return prop
}

func (this *XneqC) GetNegation() indexical.IReifiableConstraint {
	return CreateXeqC(this.x, this.c)
}
