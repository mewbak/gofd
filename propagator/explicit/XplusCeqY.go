package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// XplusCeqY represents the propagator for the constraint X + C = Y
// example: X + 3 = Y
type XplusCeqY struct {
	x, y               core.VarId
	c                  int
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	x_Domain, y_Domain *core.ExDomain
	id                 core.PropId
}

func (this *XplusCeqY) Clone() core.Constraint {
	prop := new(XplusCeqY)
	prop.x = this.x
	prop.y = this.y
	prop.c = this.c
	return prop
}

func (this *XplusCeqY) Start(store *core.Store) {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	this.xinYout(evt)
	this.yinXout(evt)
	this.outCh <- evt // sendChangestoStore

	for changeEntry := range this.inCh {
		if loggerDebug {
			core.GetLogger().Df("%s_'Incoming Change for %s'",
				this, core.GetNameRegistry().GetName(changeEntry.GetID()))
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
		this.outCh <- evt // sendChangestoStore
	}

}

func (this *XplusCeqY) xinYout(evt *core.ChangeEvent) {
	x_Domain := this.x_Domain
	var chEntry *core.ChangeEntry = nil
	for y_val := range this.y_Domain.Values {
		if !x_Domain.Contains(y_val - this.c) {
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

func (this *XplusCeqY) yinXout(evt *core.ChangeEvent) {
	y_Domain := this.y_Domain
	var chEntry *core.ChangeEntry = nil
	for x_val := range this.x_Domain.Values {
		if !y_Domain.Contains(x_val + this.c) {
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
func (this *XplusCeqY) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y}, this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]
	this.y_Domain = varidToDomainMap[this.y]

}

func (this *XplusCeqY) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XplusCeqY) GetID() core.PropId {
	return this.id
}

//X+C=Z
func CreateXplusCeqY(x core.VarId, c int, y core.VarId) *XplusCeqY {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXplusCeqZ_propagator")
	}
	prop := new(XplusCeqY)
	prop.x = x
	prop.y = y
	prop.c = c
	return prop
}

//X=Z
func CreateXeqY(x core.VarId, y core.VarId) *XplusCeqY {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXplusCeqZ_propagator")
	}
	prop := new(XplusCeqY)
	prop.x = x
	prop.y = y
	prop.c = 0
	return prop
}

func (this *XplusCeqY) String() string {
	return fmt.Sprintf("PROP_%d %s+%d = %s",
		this.id, core.GetNameRegistry().GetName(this.x), this.c,
		core.GetNameRegistry().GetName(this.y))
}

func (this *XplusCeqY) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y}
}

func (this *XplusCeqY) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain}
}

func (this *XplusCeqY) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XplusCeqY) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
