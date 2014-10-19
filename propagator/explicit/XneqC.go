package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

type XneqC struct {
	x        core.VarId
	c        int
	outCh    chan<- *core.ChangeEvent
	inCh     <-chan *core.ChangeEntry
	x_Domain *core.ExDomain
	id       core.PropId
	store    *core.Store
}

func (this *XneqC) Start() {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_Start_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	this.xin(evt)     // can fire just initially
	this.outCh <- evt // send changes to store
	for changeEntry := range this.inCh {
		if loggerDebug {
			core.GetLogger().Df("%s_Start_'Incoming Change for %s'",
				this, this.store.GetName(changeEntry.GetID()))
		}
		evt := core.CreateChangeEvent()
		switch var_id := changeEntry.GetID(); var_id {
		case this.x:
			this.x_Domain.Removes(changeEntry.GetValues())
			break
		}
		if loggerDebug {
			msg := "%s_propagate_'communicate change, evt-value: %s'"
			core.GetLogger().Df(msg, this, evt)
		}
		this.outCh <- evt // send changes to store
	}
}

func (this *XneqC) xin(evt *core.ChangeEvent) {
	if this.x_Domain.Contains(this.c) {
		chEntry := core.CreateChangeEntry(this.x)
		chEntry.Add(this.c)
		evt.AddChangeEntry(chEntry)
	}
}

// Register registers the propagator at the store.
func (this *XneqC) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x}, this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]
	this.store = store
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

func CreateXneqC(x core.VarId, c int) *XneqC {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXneqC_propagator")
	}
	prop := new(XneqC)
	prop.x = x
	prop.c = c
	return prop
}

func (this *XneqC) String() string {
	return fmt.Sprintf("PROP_%d %s != %d",
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
