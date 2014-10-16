package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// XmultYeqZ propagator for X*Y=Z
type XmultYeqZ struct {
	x, y, z                      core.VarId
	outCh                        chan<- *core.ChangeEvent
	inCh                         <-chan *core.ChangeEntry
	x_Domain, y_Domain, z_Domain *core.ExDomain
	id                           core.PropId
}

func (this *XmultYeqZ) Clone() core.Constraint {
	prop := new(XmultYeqZ)
	prop.x, prop.y, prop.z = this.x, this.y, this.z
	return prop
}

func (this *XmultYeqZ) Start(store *core.Store) {
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
		// println("X_MULT_Y_EQ_Z")
		if loggerDebug {
			core.GetLogger().Df("%s_'Incoming Change for %s'",
				this, core.GetNameRegistry().GetName(changeEntry.GetID()))
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
		//println("X_MULT_Y_EQ_Z_SEND_STORE")
	}
}

func (this *XmultYeqZ) xinYinZout(evt *core.ChangeEvent) {
	x_Domain := this.x_Domain
	y_Domain := this.y_Domain
	y_DomainNotEmpty := y_Domain.IsEmpty()
	var chEntry *core.ChangeEntry = nil
	for z_val := range this.z_Domain.Values {
		match := false
		for x_val := range x_Domain.Values {
			if x_val != 0 {
				if z_val%x_val == 0 && y_Domain.Contains(z_val/x_val) {
					match = true
					break
				}
			} else {
				if y_DomainNotEmpty && z_val == 0 {
					match = true
					break
				}
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

func (this *XmultYeqZ) xinZinYout(evt *core.ChangeEvent) {
	x_Domain := this.x_Domain
	z_Domain := this.z_Domain
	var chEntry *core.ChangeEntry = nil
	for y_val := range this.y_Domain.Values {
		match := false
		for x_val := range x_Domain.Values {
			if z_Domain.Contains(x_val * y_val) {
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

func (this *XmultYeqZ) yinZinXout(evt *core.ChangeEvent) {
	y_Domain := this.y_Domain
	z_Domain := this.z_Domain
	var chEntry *core.ChangeEntry = nil
	for x_val := range this.x_Domain.Values {
		match := false
		for y_val := range y_Domain.Values {
			if z_Domain.Contains(x_val * y_val) {
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
func (this *XmultYeqZ) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y, this.z},
			this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]
	this.y_Domain = varidToDomainMap[this.y]
	this.z_Domain = varidToDomainMap[this.z]

}

// SetID is used by the store to set the propagator's ID,
// don't use it yourself or bad things will happen
func (this *XmultYeqZ) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XmultYeqZ) GetID() core.PropId {
	return this.id
}

func CreateXmultYeqZ(x core.VarId, y core.VarId, z core.VarId) *XmultYeqZ {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXmultYeqC-propagator")
	}
	prop := new(XmultYeqZ)
	prop.x, prop.y, prop.z = x, y, z
	return prop
}

func (this *XmultYeqZ) String() string {
	return fmt.Sprintf("PROP_%d %s*%s=%s", this.id,
		core.GetNameRegistry().GetName(this.x), core.GetNameRegistry().GetName(this.y),
		core.GetNameRegistry().GetName(this.z))
}

func (this *XmultYeqZ) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y}
}

func (this *XmultYeqZ) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain}
}

func (this *XmultYeqZ) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XmultYeqZ) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
