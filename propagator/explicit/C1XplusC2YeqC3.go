package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// C1XplusC2YeqC3 implementierung
// <=>   c1*x + c2*y = c3
// <=>   c1*x *(-1)*c3 = (-1)*c2*y
// --->  c1*x + -c3 = -c2 *y
type C1XplusC2YeqC3 struct {
	x, y               core.VarId
	c1, c2, c3         int
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	x_Domain, y_Domain *core.ExDomain
	id                 core.PropId
	store              *core.Store
}

func (this *C1XplusC2YeqC3) Clone() core.Constraint {
	prop := new(C1XplusC2YeqC3)
	prop.x, prop.y = this.x, this.y
	prop.c1, prop.c2, prop.c3 = this.c1, this.c2, this.c3
	return prop
}

// Start starts propagation with y as output and x as input variable and
// listens to channels for incoming domain updates
func (this *C1XplusC2YeqC3) Start(store *core.Store) {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%v_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	this.xinYout(evt)
	this.yinXout(evt)
	this.sendChangesToStore(evt)

	for changeEntry := range this.inCh {
		if loggerDebug {
			msg := "%v_'Incoming Change for %s'"
			core.GetLogger().Df(msg, this, store.GetName(changeEntry.GetID()))
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
		this.sendChangesToStore(evt)
	}
}

// sendChangesToStore send the collected changes (stored in the event)
// to the store
func (this *C1XplusC2YeqC3) sendChangesToStore(evt *core.ChangeEvent) {
	if core.GetLogger().DoDebug() {
		msg := "%s_propagate_'communicate change, evt-value: %s'"
		core.GetLogger().Df(msg, this, evt)
	}
	this.outCh <- evt
}

func (this *C1XplusC2YeqC3) xinYout(evt *core.ChangeEvent) {
	//collect changes
	x_Domain := this.x_Domain
	var chEntry *core.ChangeEntry = nil
	for y_val := range this.y_Domain.Values {
		numerator := this.c3 - this.c2*y_val
		if (numerator%this.c1) != 0 ||
			!x_Domain.Contains(numerator/this.c1) {
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

func (this *C1XplusC2YeqC3) yinXout(evt *core.ChangeEvent) {
	y_Domain := this.y_Domain
	var chEntry *core.ChangeEntry = nil
	for x_val := range this.x_Domain.Values {
		numerator := this.c3 - this.c1*x_val
		if (numerator%this.c2) != 0 ||
			!y_Domain.Contains(numerator/this.c2) {
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
func (this *C1XplusC2YeqC3) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y}, this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]
	this.y_Domain = varidToDomainMap[this.y]
	this.store = store
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen
func (this *C1XplusC2YeqC3) SetID(propID core.PropId) {
	this.id = propID
}

func (this *C1XplusC2YeqC3) GetID() core.PropId {
	return this.id
}

func CreateC1XplusC2YeqC3(c1 int, x core.VarId, c2 int, y core.VarId,
	c3 int) *C1XplusC2YeqC3 {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateC1XplusC2YeqC3_propagator")
	}
	prop := new(C1XplusC2YeqC3)
	prop.x = x
	prop.y = y
	prop.c1 = c1
	prop.c2 = c2
	prop.c3 = c3
	return prop
}

func (this *C1XplusC2YeqC3) Equals(prop *C1XplusC2YeqC3) bool {
	return (this.id == prop.id) && (this.x == prop.x) && (this.y == prop.y) &&
		(this.c1 == prop.c1) && (this.c2 == prop.c2) && (this.c3 == prop.c3)
}

func (this *C1XplusC2YeqC3) String() string {
	return fmt.Sprintf("PROP_%d %d*%s+%d*%s=%d",
		this.id, this.c1, this.store.GetName(this.x),
		this.c2, this.store.GetName(this.y), this.c3)
}

func CreateXplusYeqC(x core.VarId, y core.VarId, c int) *C1XplusC2YeqC3 {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateC1XplusC2YeqC3_propagator")
	}
	prop := new(C1XplusC2YeqC3)
	prop.x = x
	prop.y = y
	prop.c1 = 1
	prop.c2 = 1
	prop.c3 = c
	return prop
}

func (this *C1XplusC2YeqC3) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y}
}

func (this *C1XplusC2YeqC3) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain}
}

func (this *C1XplusC2YeqC3) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *C1XplusC2YeqC3) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
