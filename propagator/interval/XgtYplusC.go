package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

//XgtYplus represents the constraint X>Y+C
type XgtYplusC struct {
	x, y               core.VarId
	c                  int
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	x_Domain, y_Domain *core.IvDomain
	id                 core.PropId
	store              *core.Store
}

func (this *XgtYplusC) Clone() core.Constraint {
	prop := new(XgtYplusC)
	prop.x, prop.y, prop.c = this.x, this.y, this.c
	return prop
}

func (this *XgtYplusC) Start(store *core.Store) {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	this.xinYout(evt)
	this.yinXout(evt)
	this.outCh <- evt // send changes to store
	for changeEntry := range this.inCh {
		//println("X_GT_Y")
		if loggerDebug {
			core.GetLogger().Df("%s_'Incoming Change for %s'",
				this, store.GetName(changeEntry.GetID()))
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
		//println("X_GT_Y_SEND_STORE")
	}
}

func (this *XgtYplusC) xinYout(evt *core.ChangeEvent) {

	xmaxmc := this.x_Domain.GetMax() - this.c

	if (this.y_Domain.GetMax()) >= xmaxmc {
		rDom := core.CreateIvDomainFromTo(xmaxmc, this.y_Domain.GetMax())
		chEntry := core.CreateChangeEntryWithValues(this.y, rDom)
		evt.AddChangeEntry(chEntry)
	}
}

func (this *XgtYplusC) yinXout(evt *core.ChangeEvent) {
	yminmc := this.y_Domain.GetMin() + this.c
	if (this.x_Domain.GetMin()) <= yminmc {
		rDom := core.CreateIvDomainFromTo(this.x_Domain.GetMin(), yminmc)
		chEntry := core.CreateChangeEntryWithValues(this.x, rDom)
		evt.AddChangeEntry(chEntry)
	}
}

// Register registers the propagator at the store.
func (this *XgtYplusC) Register(store *core.Store) {
	var domains []core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagator([]core.VarId{this.x, this.y}, this.id)
	this.x_Domain = core.GetVaridToIntervalDomain(domains[0])
	this.y_Domain = core.GetVaridToIntervalDomain(domains[1])
	this.store = store
}

func (this *XgtYplusC) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XgtYplusC) GetID() core.PropId {
	return this.id
}

// X > Y+C
func CreateXgtYplusC(x core.VarId, y core.VarId, c int) *XgtYplusC {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXgtYplusC_propagator")
	}
	prop := new(XgtYplusC)
	prop.x, prop.y, prop.c = x, y, c
	return prop
}

// X >= Y
func CreateXgteqY(x core.VarId, y core.VarId) *XgtYplusC {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXgteqYplusC_propagator")
	}
	prop := new(XgtYplusC)
	prop.x, prop.y, prop.c = x, y, -1
	return prop
}

func (this *XgtYplusC) String() string {
	return fmt.Sprintf("PROP_%d %s > %s + %d",
		this.id, this.store.GetName(this.x),
		this.store.GetName(this.y), this.c)
}

func (this *XgtYplusC) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y}
}

func (this *XgtYplusC) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain}
}

func (this *XgtYplusC) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XgtYplusC) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
