package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// XplusCneqY represents the propagator for the constraint X + C != Y
type XplusCneqY struct {
	x, y               core.VarId
	c                  int
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	x_Domain, y_Domain *core.ExDomain
	id                 core.PropId
	store              *core.Store
}

func (this *XplusCneqY) Clone() core.Constraint {
	prop := new(XplusCneqY)
	prop.x = this.x
	prop.y = this.y
	prop.c = this.c
	return prop
}

func (this *XplusCneqY) Start() {
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
	}
}

// xinYout removes fixed values X+c from Y
// Example with C=1: X:{5}, Y:{6,7} --> X:{5}, Y:{7}
func (this *XplusCneqY) xinYout(evt *core.ChangeEvent) {
	if this.x_Domain.IsGround() {
		fixed_val := this.x_Domain.Min + this.c
		if this.y_Domain.Contains(fixed_val) {
			chEntry := core.CreateChangeEntry(this.y)
			chEntry.Add(fixed_val)
			evt.AddChangeEntry(chEntry)
			return
		}
	}
}

// yinXout see xinYout, but vice versa
func (this *XplusCneqY) yinXout(evt *core.ChangeEvent) {
	if this.y_Domain.IsGround() {
		fixed_val := this.y_Domain.Min - this.c
		if this.x_Domain.Contains(fixed_val) {
			chEntry := core.CreateChangeEntry(this.x)
			chEntry.Add(fixed_val)
			evt.AddChangeEntry(chEntry)
			return
		}
	}
}

// Register registers the propagator at the store.
func (this *XplusCneqY) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y}, this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]
	this.y_Domain = varidToDomainMap[this.y]
	this.store = store
}

func (this *XplusCneqY) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XplusCneqY) GetID() core.PropId {
	return this.id
}

// CreateXplusCneqY creates propagators for the constraint X+C!=Y
func CreateXplusCneqY(x core.VarId, c int, y core.VarId) *XplusCneqY {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXplusCneqZ_propagator")
	}
	prop := new(XplusCneqY)
	prop.x, prop.y, prop.c = x, y, c
	return prop
}

func (this *XplusCneqY) String() string {
	return fmt.Sprintf("PROP_%d %s + %d != %s", this.id,
		this.store.GetName(this.x), this.c, this.store.GetName(this.y))
}

func (this *XplusCneqY) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y}
}

func (this *XplusCneqY) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain}
}

func (this *XplusCneqY) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XplusCneqY) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
