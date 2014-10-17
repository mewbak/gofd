package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// XltC propagator for X < C
type XltC struct {
	x        core.VarId
	c        int
	outCh    chan<- *core.ChangeEvent
	inCh     <-chan *core.ChangeEntry
	x_Domain *core.ExDomain
	id       core.PropId
	store    *core.Store
}

func (this *XltC) Clone() core.Constraint {
	prop := new(XltC)
	prop.x = this.x
	prop.c = this.c
	return prop
}

func (this *XltC) Start(store *core.Store) {
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
				this, core.GetNameRegistry().GetName(changeEntry.GetID()))
		}
		evt := core.CreateChangeEvent()
		switch var_id := changeEntry.GetID(); var_id {
		case this.x:
			this.x_Domain.Removes(changeEntry.GetValues())
			// impossible, that anything else will be removed
			// why is the channel still open?
		}
		if loggerDebug {
			msg := "%s_propagate_'communicate change, evt-value: %s'"
			core.GetLogger().Df(msg, this, evt)
		}
		this.outCh <- evt // send changes to store
	}
}

func (this *XltC) xin(evt *core.ChangeEvent) {
	var chEntry *core.ChangeEntry = nil
	for v := range this.x_Domain.Values {
		if v >= this.c { // remove all values v in X with v <= C
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

// Register registers the propagator at the store, subscribes
// for reading on X
func (this *XltC) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x}, this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]

}

func (this *XltC) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XltC) GetID() core.PropId {
	return this.id
}

func CreateXltC(x core.VarId, c int) *XltC {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXltC_propagator")
	}
	XltC := new(XltC)
	XltC.x = x
	XltC.c = c
	return XltC
}

// CreateXlteqC creates propagators for X<=C
// with XgtC-propagators. Example:
// User want: 	 X<=6
// System makes: X<6+1
func CreateXlteqC(x core.VarId, c int) *XltC {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXlteqC_propagator")
	}
	xltc := new(XltC)
	xltc.x = x
	xltc.c = c + 1
	return xltc
}

func (this *XltC) String() string {
	return fmt.Sprintf("PROP_%d %s<%d",
		this.id, core.GetNameRegistry().GetName(this.x), this.c)
}

func (this *XltC) GetVarIds() []core.VarId {
	return []core.VarId{this.x}
}

func (this *XltC) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain}
}

func (this *XltC) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XltC) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
