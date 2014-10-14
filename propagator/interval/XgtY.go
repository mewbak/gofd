package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

//XgtY represents the constraint X > Y
type XgtY struct {
	x, y               core.VarId
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	x_Domain, y_Domain *core.IvDomain
	id                 core.PropId
	store              *core.Store
}

func (this *XgtY) Clone() core.Constraint {
	prop := new(XgtY)
	prop.x = this.x
	prop.y = this.y
	return prop
}

func (this *XgtY) Start(store *core.Store) {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	this.xinYout(evt)
	this.yinXout(evt)
	this.outCh <- evt // send changes to store

	for changeEntry := range this.inCh {
		if loggerDebug {
			core.GetLogger().Df("%s_'Incoming Change for %s'", this,
				store.GetName(changeEntry.GetID()))
		}
		evt = core.CreateChangeEvent()
		switch var_id := changeEntry.GetID(); var_id {
		case this.x:
			this.x_Domain.Removes(changeEntry.GetValues())
			this.xinYout(evt)
		case this.y:
			this.y_Domain.Removes(changeEntry.GetValues())
			this.yinXout(evt)
		}
		if loggerDebug {
			msg := "%s_propagate_'communicate change, evt-value: %s'"
			core.GetLogger().Df(msg, this, evt)
		}
		this.outCh <- evt // send changes to store
	}
}

// x \in X.Domain | x>y_val
func (this *XgtY) xinYout(evt *core.ChangeEvent) {
	xmax := this.x_Domain.GetMax()

	if this.y_Domain.GetMax() >= xmax {
		rDom := core.CreateIvDomainFromTo(xmax, this.y_Domain.GetMax())
		chEntry := core.CreateChangeEntryWithValues(this.y, rDom)
		evt.AddChangeEntry(chEntry)
	}
}

// y \in y.Domain | y<x_val
func (this *XgtY) yinXout(evt *core.ChangeEvent) {
	ymin := this.y_Domain.GetMin()

	if this.x_Domain.GetMin() <= ymin {
		rDom := core.CreateIvDomainFromTo(this.x_Domain.GetMin(), ymin)
		chEntry := core.CreateChangeEntryWithValues(this.x, rDom)
		evt.AddChangeEntry(chEntry)
	}
}

// Register registers the propagator at the store.
func (this *XgtY) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y}, this.id)
	this.x_Domain = core.GetVaridToIntervalDomain(domains[this.x])
	this.y_Domain = core.GetVaridToIntervalDomain(domains[this.y])
	this.store = store
}

func (this *XgtY) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XgtY) GetID() core.PropId {
	return this.id
}

// X+C=Z
func CreateXgtY(x core.VarId, y core.VarId) *XgtY {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXgtY_propagator")
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
