package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// XgtC represents the propagator for the constraint X > C
type XgtC struct {
	x        core.VarId
	c        int
	outCh    chan<- *core.ChangeEvent
	inCh     <-chan *core.ChangeEntry
	x_Domain *core.ExDomain
	id       core.PropId
	store    *core.Store
}

func (this *XgtC) Clone() core.Constraint {
	prop := new(XgtC)
	prop.x = this.x
	prop.c = this.c
	return prop
}

func (this *XgtC) Start() {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_Start_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	this.xin(evt)
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
		}
		if loggerDebug {
			msg := "%s_propagate_'communicate change, evt-value: %s'"
			core.GetLogger().Df(msg, this, evt)
		}
		this.outCh <- evt // send changes to store
	}
}

func (this *XgtC) xin(evt *core.ChangeEvent) {
	var chEntry *core.ChangeEntry = nil
	for v := range this.x_Domain.Values {
		if v <= this.c { // remove all values v \in X with v <= C
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(this.x)
			}
			chEntry.Add(v)
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

// Register registers the propagator at the store.
func (this *XgtC) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x}, this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]
	this.store = store
}

func (this *XgtC) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XgtC) GetID() core.PropId {
	return this.id
}

func CreateXgtC(x core.VarId, c int) *XgtC {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXgtC_propagator")
	}
	xgtc := new(XgtC)
	xgtc.x = x
	xgtc.c = c
	return xgtc
}

func (this *XgtC) String() string {
	return fmt.Sprintf("PROP_%d %s > %d",
		this.id, this.store.GetName(this.x), this.c)
}

// CreateXgteqC creates propagators for X>=C
// with XgtC-propagators. Example:
// User want: 	 X>=6
// System makes: X>6-1
func CreateXgteqC(x core.VarId, c int) *XgtC {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXgteqC_propagator")
	}
	xgtc := new(XgtC)
	xgtc.x = x
	xgtc.c = c - 1
	return xgtc
}

func (this *XgtC) GetVarIds() []core.VarId {
	return []core.VarId{this.x}
}

func (this *XgtC) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain}
}

func (this *XgtC) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XgtC) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
