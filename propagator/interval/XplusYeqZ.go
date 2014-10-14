package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// XplusYneqZ represents the constraint X + Y = Z
type XplusYeqZ struct {
	x, y, z                      core.VarId
	outCh                        chan<- *core.ChangeEvent
	inCh                         <-chan *core.ChangeEntry
	x_Domain, y_Domain, z_Domain *core.IvDomain
	id                           core.PropId
	store                        *core.Store
}

func (this *XplusYeqZ) Clone() core.Constraint {
	prop := new(XplusYeqZ)
	prop.x, prop.y, prop.z = this.x, this.y, this.z
	return prop
}

func (this *XplusYeqZ) Start(store *core.Store) {
	// initial check
	evt := core.CreateChangeEvent()
	ivsecondInResultInFirstOut(this.y_Domain, this.z_Domain, this.x_Domain,
		this.x, evt)
	ivfirstInResultInSecondOut(this.x_Domain, this.z_Domain, this.y_Domain,
		this.y, evt)
	ivfirstInSecondInResultOut(this.x_Domain, this.y_Domain, this.z_Domain,
		this.z, evt)
	core.SendChangesToStore(evt, this)
	for changeEntry := range this.inCh {
		core.LogIncomingChange(this, store, changeEntry)
		evt = core.CreateChangeEvent()
		switch var_id := changeEntry.GetID(); var_id {
		case this.x:
			this.x_Domain.Removes(changeEntry.GetValues())
			ivfirstInResultInSecondOut(this.x_Domain, this.y_Domain,
				this.z_Domain, this.y, evt)
			ivfirstInSecondInResultOut(this.x_Domain, this.y_Domain,
				this.z_Domain, this.z, evt)
		case this.y:
			this.y_Domain.Removes(changeEntry.GetValues())
			ivsecondInResultInFirstOut(this.x_Domain, this.y_Domain,
				this.z_Domain, this.x, evt)
			ivfirstInSecondInResultOut(this.x_Domain, this.y_Domain,
				this.z_Domain, this.z, evt)
		case this.z:
			this.z_Domain.Removes(changeEntry.GetValues())
			ivsecondInResultInFirstOut(this.x_Domain, this.y_Domain,
				this.z_Domain, this.x, evt)
			ivfirstInResultInSecondOut(this.x_Domain, this.y_Domain,
				this.z_Domain, this.y, evt)
		}
		core.SendChangesToStore(evt, this)
	}
}

// Register registers the propagator at the store.
func (this *XplusYeqZ) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y, this.z},
			this.id)
	this.x_Domain = core.GetVaridToIntervalDomain(domains[this.x])
	this.y_Domain = core.GetVaridToIntervalDomain(domains[this.y])
	this.z_Domain = core.GetVaridToIntervalDomain(domains[this.z])
	this.store = store
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen
func (this *XplusYeqZ) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XplusYeqZ) GetID() core.PropId {
	return this.id
}

func CreateXplusYeqZ(x core.VarId, y core.VarId, z core.VarId) *XplusYeqZ {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXplusYeqZ_intervals-propagator")
	}
	prop := new(XplusYeqZ)
	prop.x, prop.y, prop.z = x, y, z
	return prop
}

func (this *XplusYeqZ) String() string {
	return fmt.Sprintf("PROP_%d %s+%s=%s",
		this.id, this.store.GetName(this.x), this.store.GetName(this.y),
		this.store.GetName(this.z))
}

func (this *XplusYeqZ) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y, this.z}
}

func (this *XplusYeqZ) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain, this.z_Domain}
}

func (this *XplusYeqZ) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XplusYeqZ) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
