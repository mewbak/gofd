package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// XneqY represents the propagator for the constraint X != Y
type XneqY struct {
	x, y               core.VarId
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	x_Domain, y_Domain *core.ExDomain
	id                 core.PropId
	store              *core.Store
}

func (this *XneqY) Clone() core.Constraint {
	prop := new(XneqY)
	prop.x = this.x
	prop.y = this.y
	return prop
}

func (this *XneqY) Start(store *core.Store) {
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
		// ToDo: If either one of the variables is ground there is the only
		// propagation that can happen at all. Thus, afterwards the propgator
		// can deregister itself, which is not implemented yet.
		// ToDo: Is it possible that *any* propagator that has just one
		// nonground variable left can deregister itself? If yes - do it,
		// if not at least the option during register might be given so that
		// the store can deregister as it is correctly doing if *all* variables
		// are ground.
		if loggerDebug {
			msg := "%s_propagate_'communicate change, evt-value: %s'"
			core.GetLogger().Df(msg, this, evt)
		}
		this.outCh <- evt // send changes to store
	}
}

// xinYout removes fixed values X from Y
// Example: X:{6}, Y:{6,7} --> X:{6}, Y:{7}
func (this *XneqY) xinYout(evt *core.ChangeEvent) {
	if this.x_Domain.IsGround() {
		fixed_val := this.x_Domain.Min
		if this.y_Domain.Contains(fixed_val) {
			chEntry := core.CreateChangeEntry(this.y)
			chEntry.Add(fixed_val)
			evt.AddChangeEntry(chEntry)
			return
		}
	}
}

// yinXout see xinYout, but vice versa
func (this *XneqY) yinXout(evt *core.ChangeEvent) {
	if this.y_Domain.IsGround() {
		fixed_val := this.y_Domain.Min
		if this.x_Domain.Contains(fixed_val) {
			chEntry := core.CreateChangeEntry(this.x)
			chEntry.Add(fixed_val)
			evt.AddChangeEntry(chEntry)
			return
		}
	}
}

// Register registers the propagator at the store.
func (this *XneqY) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y}, this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]
	this.y_Domain = varidToDomainMap[this.y]
	this.store = store
}

func (this *XneqY) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XneqY) GetID() core.PropId {
	return this.id
}

// CreateXneqY creates propagators for the constraint X!=Y
func CreateXneqY(x core.VarId, y core.VarId) *XneqY {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXneqY_propagator")
	}
	prop := new(XneqY)
	prop.x, prop.y = x, y
	return prop
}

func (this *XneqY) String() string {
	return fmt.Sprintf("PROP_%d %s != %s", this.id,
		this.store.GetName(this.x), this.store.GetName(this.y))
}

func (this *XneqY) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y}
}

func (this *XneqY) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain}
}

func (this *XneqY) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XneqY) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
