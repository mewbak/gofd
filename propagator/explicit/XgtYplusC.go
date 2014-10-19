package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

type XgtYplusC struct {
	x, y               core.VarId
	c                  int
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	x_Domain, y_Domain *core.ExDomain
	id                 core.PropId
	store              *core.Store
}

func (this *XgtYplusC) Clone() core.Constraint {
	prop := new(XgtYplusC)
	prop.x, prop.y, prop.c = this.x, this.y, this.c
	return prop
}

func (this *XgtYplusC) Start() {
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
				this, this.store.GetName(changeEntry.GetID()))
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
	xmaxmc := this.x_Domain.Max - this.c
	var chEntry *core.ChangeEntry = nil
	if xmaxmc > this.y_Domain.Max {
		return
	}
	for y_val := range this.y_Domain.Values {
		if xmaxmc <= y_val {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(this.y)
			}
			chEntry.Add(y_val)
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

func (this *XgtYplusC) yinXout(evt *core.ChangeEvent) {
	yminmc := this.y_Domain.Min + this.c
	var chEntry *core.ChangeEntry = nil
	if yminmc < this.x_Domain.Min {
		return
	}
	for x_val := range this.x_Domain.Values {
		if yminmc >= x_val {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(this.x)
			}
			chEntry.Add(x_val)
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

// Register registers the propagator at the store.
func (this *XgtYplusC) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y}, this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]
	this.y_Domain = varidToDomainMap[this.y]
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
