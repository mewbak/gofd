package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// XeqC represents the constraint X = C
type XeqC struct {
	x        core.VarId
	c        int
	outCh    chan<- *core.ChangeEvent
	inCh     <-chan *core.ChangeEntry
	x_Domain *core.IvDomain
	id       core.PropId
	store    *core.Store
}

func (this *XeqC) Start(store *core.Store) {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_Start_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	this.xin(evt)
	this.outCh <- evt // send changes to store

	for changeEntry := range this.inCh {
		if loggerDebug {
			msg := "%s_Start_'Incoming Change for %s'"
			core.GetLogger().Df(msg, this, store.GetName(changeEntry.GetID()))
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

func (this *XeqC) xin(evt *core.ChangeEvent) {
	var chEntry *core.ChangeEntry = nil
	for x_val := range this.x_Domain.Values_asMap() {
		if x_val != this.c {
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
func (this *XeqC) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x}, this.id)

	varidToDomainMap := core.GetVaridToIntervalDomains(domains)

	this.x_Domain = varidToDomainMap[this.x]
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
	return fmt.Sprintf("PROP_%d %s = %d",
		this.id, this.store.GetName(this.x), this.c)
}

func CreateXeqC(x core.VarId, c int) *XeqC {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXeqC_propagator")
	}
	prop := new(XeqC)
	prop.x = x
	prop.c = c
	return prop
}

func (this *XeqC) GetVarIds() []core.VarId {
	return []core.VarId{this.x}
}

func (this *XeqC) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain}
}

func (this *XeqC) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XeqC) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
