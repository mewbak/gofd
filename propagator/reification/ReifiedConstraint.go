package reification

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical"
	"fmt"
)

// for reifying a constraint, see constraints in folder "indexicalconstraint"
// and interface "IReifiableConstraint"

type ReifiedConstraint struct {
	c, negC  indexical.IReifiableConstraint
	b        core.VarId
	b_Domain *core.IvDomain
	outCh    chan<- *core.ChangeEvent
	inCh     <-chan *core.ChangeEntry
	id       core.PropId
	store    *core.Store
	state    int
}

func (this *ReifiedConstraint) Clone() core.Constraint {
	prop := new(ReifiedConstraint)
	prop.c = this.c.Clone().(indexical.IReifiableConstraint)
	prop.negC = this.negC.Clone().(indexical.IReifiableConstraint)
	prop.b = this.b
	prop.state = this.state
	return prop
}

func (this *ReifiedConstraint) Process(store *core.Store, changeEntry *core.ChangeEntry) {
	evt := core.CreateChangeEvent()
	//println(this.c.String())
	if this.b_Domain.IsGround() {

		//propagation
		bVal := this.b_Domain.GetAnyElement()
		var iColl *indexical.IndexicalCollection
		if bVal == 1 {
			//println("b ground 1")
			//propagate C
			core.LogPropagationOfConstraint(this.c)
			iColl = this.c.GetIndexicalCollection()
		} else if bVal == 0 {
			//println("b ground 0")
			//propagate !C
			core.LogPropagationOfConstraint(this.negC)
			iColl = this.negC.GetIndexicalCollection()
		} else {
			panic("B has wrong value!")
		}

		if changeEntry == nil {
			this.state = PROPAGATING
			evt = indexical.ProcessIndexicals(iColl, nil, false)
		} else {
			if this.state == ENTAILMENT_CHECKING {
				this.state = PROPAGATING
				evt = indexical.ProcessIndexicals(iColl, nil, false)
			} else {
				evt = indexical.ProcessIndexicals(iColl, changeEntry, false)
			}
		}
	} else {
		//println("entail")
		//entailment-check
		var d *core.IvDomain
		if this.c.IsEntailed() {
			//println("c entailed")
			core.LogEntailmentDetected(this, this.c)
			d = core.CreateIvDomainFromTo(0, 0)
		} else if this.negC.IsEntailed() {
			//println("!c entailed")
			core.LogEntailmentDetected(this, this.negC)
			d = core.CreateIvDomainFromTo(1, 1)
		} else {
			//println("nothing")
		}
		if d != nil {
			chEntry := core.CreateChangeEntry(this.b)
			chEntry.SetValues(d)
			evt.AddChangeEntry(chEntry)
		}
	}
	//println(evt.String())
	core.SendChangesToStore(evt, this)
}

func (this *ReifiedConstraint) Start(store *core.Store) {

	core.LogInitConsistency(this)
	this.Process(store, nil)

	for changeEntry := range this.inCh {

		indexical.RemoveValues(this, changeEntry)
		this.Process(store, changeEntry)
	}
}

// Register registers the propagator at the store.
func (this *ReifiedConstraint) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(this.GetVarIds(), this.id)

	varidToDomainMap := core.GetVaridToIntervalDomains(domains)

	this.b_Domain = varidToDomainMap[this.b]

	if this.b_Domain.GetMax() > 1 || this.b_Domain.GetMin() < 0 {
		panic("Reified Constraint: bool-Domain has invalid values (must be 0 and/or 1)!")
	}

	this.c.Init(store, varidToDomainMap)
	this.negC.Init(store, varidToDomainMap)
	this.store = store
}

func (this *ReifiedConstraint) SetID(propID core.PropId) {
	this.id = propID
}

func (this *ReifiedConstraint) GetID() core.PropId {
	return this.id
}

func (this *ReifiedConstraint) GetVarIds() []core.VarId {
	varids := make([]core.VarId, 0)
	varids = append(varids, this.b)
	varids = append(varids, this.c.GetVarIds()...)
	return varids
}

func (this *ReifiedConstraint) GetDomains() []core.Domain {
	domains := make([]core.Domain, 0)
	domains = append(domains, this.b_Domain)
	domains = append(domains, this.c.GetDomains()...)

	return domains
}

func (this *ReifiedConstraint) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *ReifiedConstraint) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}

func (this *ReifiedConstraint) String() string {
	var processingC string
	if this.b_Domain.IsGround() {
		v := this.b_Domain.GetAnyElement()
		if v == 0 {
			processingC = this.negC.String()
		} else if v == 1 {
			processingC = this.c.String()
		}
	} else {
		processingC = "nothing"
	}

	return fmt.Sprintf("'%s <=> %s, propagating %s'",
		this.c, this.b_Domain, processingC)
}

func CreateReifiedConstraint(c indexical.IReifiableConstraint, b core.VarId) *ReifiedConstraint {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateReifiedConstraint")
	}

	prop := new(ReifiedConstraint)
	prop.c = c
	prop.negC = c.GetNegation()
	prop.b = b
	prop.state = ENTAILMENT_CHECKING

	return prop
}

func (this *ReifiedConstraint) GetBool() core.VarId {
	return this.b
}

const (
	ENTAILMENT_CHECKING = 0
	PROPAGATING         = 1
)
