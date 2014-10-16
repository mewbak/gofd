package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// XmultCeqY propagator for X*C = Y
type XmultCeqY struct {
	x, y               core.VarId
	c                  int
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	x_Domain, y_Domain *core.ExDomain
	id                 core.PropId
}

func (this *XmultCeqY) Clone() core.Constraint {
	prop := new(XmultCeqY)
	prop.x = this.x
	prop.y = this.y
	prop.c = this.c
	return prop
}

func (this *XmultCeqY) Start(store *core.Store) {
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
				core.GetNameRegistry().GetName(changeEntry.GetID()))
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

func (this *XmultCeqY) xinYout(evt *core.ChangeEvent) {
	x_Domain := this.x_Domain
	var chEntry *core.ChangeEntry = nil
	for y_val := range this.y_Domain.Values {
		if this.c == 0 {
			if y_val != 0 {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.y)
				}
				chEntry.Add(y_val)
			}
		} else {
			if y_val%this.c != 0 || !x_Domain.Contains(y_val/this.c) {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.y)
				}
				chEntry.Add(y_val)
			}
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

func (this *XmultCeqY) yinXout(evt *core.ChangeEvent) {
	y_Domain := this.y_Domain
	var chEntry *core.ChangeEntry = nil
	for x_val := range this.x_Domain.Values {
		if !y_Domain.Contains(x_val * this.c) {
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
func (this *XmultCeqY) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y}, this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]
	this.y_Domain = varidToDomainMap[this.y]

}

func (this *XmultCeqY) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XmultCeqY) GetID() core.PropId {
	return this.id
}

func CreateXmultCeqY(x core.VarId, c int, y core.VarId) *XmultCeqY {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXmultCeqZ_propagator")
	}
	prop := new(XmultCeqY)
	prop.x = x
	prop.y = y
	prop.c = c
	return prop
}

func (this *XmultCeqY) String() string {
	return fmt.Sprintf("PROP_%d  %s*%d = %s",
		this.id, core.GetNameRegistry().GetName(this.x), this.c,
		core.GetNameRegistry().GetName(this.y))
}

func (this *XmultCeqY) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y}
}

func (this *XmultCeqY) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain}
}

func (this *XmultCeqY) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XmultCeqY) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
