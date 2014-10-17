package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// XplusYneqZ represents the propagator for the constraint X + Y == Z
type XPlusYEqZ struct {
	x, y, z                      core.VarId
	outCh                        chan<- *core.ChangeEvent
	inCh                         <-chan *core.ChangeEntry
	x_Domain, y_Domain, z_Domain *core.ExDomain
	id                           core.PropId
	store                        *core.Store
}

func (this *XPlusYEqZ) Clone() core.Constraint {
	prop := new(XPlusYEqZ)
	prop.x, prop.y, prop.z = this.x, this.y, this.z
	return prop
}

func (this *XPlusYEqZ) Start(store *core.Store) {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	this.xinYinZout(evt)
	this.xinZinYout(evt)
	this.yinZinXout(evt)
	this.outCh <- evt // send changes to store
	for changeEntry := range this.inCh {
		// println("X_PLUS_Y_EQ_Z")
		if loggerDebug {
			core.GetLogger().Df("%s_'Incoming Change for %s'",
				this, store.GetName(changeEntry.GetID()))
		}
		evt = core.CreateChangeEvent()
		switch var_id := changeEntry.GetID(); var_id {
		case this.x:
			this.x_Domain.Removes(changeEntry.GetValues())
			this.xinYinZout(evt)
			this.xinZinYout(evt)
		case this.y:
			this.y_Domain.Removes(changeEntry.GetValues())
			this.xinYinZout(evt)
			this.yinZinXout(evt)
		case this.z:
			this.z_Domain.Removes(changeEntry.GetValues())
			this.xinZinYout(evt)
			this.yinZinXout(evt)
		}
		if loggerDebug {
			msg := "%s_propagate_'communicate change, evt-value: %s'"
			core.GetLogger().Df(msg, this, evt)
		}
		this.outCh <- evt // send changes to store
		// println("X_PLUS_Y_EQ_Z_SEND_STORE")
	}
}

func (this *XPlusYEqZ) xinYinZout(evt *core.ChangeEvent) {
	y_Domain := this.y_Domain
	var chEntry *core.ChangeEntry = nil
	for z_val := range this.z_Domain.Values {
		match := false
		for x_val := range this.x_Domain.Values {
			if y_Domain.Contains(z_val - x_val) {
				match = true
				break
			}
		}
		if !match {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(this.z)
			}
			chEntry.Add(z_val)
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

func (this *XPlusYEqZ) xinZinYout(evt *core.ChangeEvent) {
	z_Domain := this.z_Domain
	var chEntry *core.ChangeEntry = nil
	for y_val := range this.y_Domain.Values {
		match := false
		for x_val := range this.x_Domain.Values {
			if z_Domain.Contains(y_val + x_val) {
				match = true
				break
			}
		}
		if !match {
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

func (this *XPlusYEqZ) yinZinXout(evt *core.ChangeEvent) {
	z_Domain := this.z_Domain
	var chEntry *core.ChangeEntry = nil
	for x_val := range this.x_Domain.Values {
		match := false
		for y_val := range this.y_Domain.Values {
			if z_Domain.Contains(x_val + y_val) {
				match = true
				break
			}
		}
		if !match {
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
func (this *XPlusYEqZ) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y, this.z},
			this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]
	this.y_Domain = varidToDomainMap[this.y]
	this.z_Domain = varidToDomainMap[this.z]
	this.store = store
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen
func (this *XPlusYEqZ) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XPlusYEqZ) GetID() core.PropId {
	return this.id
}

func CreateXplusYeqZ(x core.VarId, y core.VarId, z core.VarId) *XPlusYEqZ {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXplusYeqZ-propagator")
	}
	prop := new(XPlusYEqZ)
	prop.x, prop.y, prop.z = x, y, z
	return prop
}

func (this *XPlusYEqZ) String() string {
	return fmt.Sprintf("PROP_%d %s+%s=%s",
		this.id, this.store.GetName(this.x), this.store.GetName(this.y),
		this.store.GetName(this.z))
}

func (this *XPlusYEqZ) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y, this.z}
}

func (this *XPlusYEqZ) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain, this.z_Domain}
}

func (this *XPlusYEqZ) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XPlusYEqZ) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
