package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// C1XmultC2YeqC3ZBounds propagator for C1*X * C2*Y = C3*Z
type C1XmultC2YeqC3ZBounds struct {
	x, y, z                      core.VarId
	c1, c2, c3                   int
	outCh                        chan<- *core.ChangeEvent
	inCh                         <-chan *core.ChangeEntry
	x_Domain, y_Domain, z_Domain *core.ExDomain
	id                           core.PropId
	store                        *core.Store
}

func (this *C1XmultC2YeqC3ZBounds) Clone() core.Constraint {
	prop := new(C1XmultC2YeqC3ZBounds)
	prop.x = this.x
	prop.y = this.y
	prop.z = this.z
	prop.c1 = this.c1
	prop.c2 = this.c2
	prop.c3 = this.c3
	return prop
}

func (this *C1XmultC2YeqC3ZBounds) Start() {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	if !this.makeBoundConsistent(evt) {
		this.sendChangesToStore(evt)
		return
	}
	this.xinYinZout(evt)
	this.xinZinYout(evt)
	this.yinZinXout(evt)
	this.sendChangesToStore(evt)

	for changeEntry := range this.inCh {
		if loggerDebug {
			core.GetLogger().Df("%s_'Incoming Change for %s'",
				this, this.store.GetName(changeEntry.GetID()))
		}
		evt = core.CreateChangeEvent()
		switch var_id := changeEntry.GetID(); var_id {
		case this.x:
			this.x_Domain.Removes(changeEntry.GetValues())
			if !this.makeBoundConsistent(evt) {
				this.sendChangesToStore(evt)
				return
			}
			this.xinYinZout(evt)
			this.xinZinYout(evt)
		case this.y:
			this.y_Domain.Removes(changeEntry.GetValues())
			if !this.makeBoundConsistent(evt) {
				this.sendChangesToStore(evt)
				return
			}
			this.xinYinZout(evt)
			this.yinZinXout(evt)
		case this.z:
			this.z_Domain.Removes(changeEntry.GetValues())
			if !this.makeBoundConsistent(evt) {
				this.sendChangesToStore(evt)
				return
			}
			this.xinZinYout(evt)
			this.yinZinXout(evt)
		}
		this.sendChangesToStore(evt)
	}
}

// sendChangesToStore send the collected changes (stored in the event)
// to the store
func (this *C1XmultC2YeqC3ZBounds) sendChangesToStore(evt *core.ChangeEvent) {
	if core.GetLogger().DoDebug() {
		msg := "%s_propagate_'communicate change, evt-value: %s'"
		core.GetLogger().Df(msg, this, evt)
	}
	this.outCh <- evt
}

// makeBoundConsistent checks if the constraint is bound consistent
// and makes it bound consistent if possible
func (this *C1XmultC2YeqC3ZBounds) makeBoundConsistent(evt *core.ChangeEvent) bool {
	if this.x_Domain.IsEmpty() || this.y_Domain.IsEmpty() ||
		this.z_Domain.IsEmpty() {
		return false
	}
	xmin, xmax := this.x_Domain.Min, this.x_Domain.Max
	ymin, ymax := this.y_Domain.Min, this.y_Domain.Max
	zmin, zmax := this.z_Domain.Min, this.z_Domain.Max
	possmax := this.c1 * xmax * this.c2 * ymax
	possmin := this.c1 * xmin * this.c2 * ymin
	var chEntry *core.ChangeEntry = nil

	// if max of domain Z is bigger than possmax (= xmax + ymax)
	// or min of domain Z is smaller than possmin (= xmin + ymin)
	if possmax < this.c3*zmax || possmin > this.c3*zmin {
		for k := range this.z_Domain.Values {
			// delete all values in Z that are bigger than possmax
			if k*this.c3 > possmax {
				this.z_Domain.Remove(k) // BUG: ToDo: shouldn't remove
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.z)
				}
				chEntry.Add(k)
			}
			// delete all values in Z that are smaller than possmin
			if k*this.c3 < possmin {
				this.z_Domain.Remove(k) // BUG: ToDo: shouldn't remove
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.z)
				}
				chEntry.Add(k)
			}
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
	// if there are still elements in the domain of Z,
	// it should be bound consistent by now
	if !this.z_Domain.IsEmpty() {
		return true
	}
	return false
}

func (this *C1XmultC2YeqC3ZBounds) xinYinZout(evt *core.ChangeEvent) {
	//collect changes
	x_Domain := this.x_Domain
	y_Domain := this.y_Domain
	y_DomainNotEmpty := !y_Domain.IsEmpty()
	var chEntry *core.ChangeEntry = nil
	for z_val := range this.z_Domain.Values {
		match := false
		for x_val := range x_Domain.Values {
			if x_val != 0 {
				check_val := (this.c3 * z_val / (this.c1 * x_val)) / this.c2
				if (this.c3*z_val)%(this.c1*x_val) == 0 &&
					y_Domain.Contains(check_val) {
					match = true
					break
				}
			} else {
				if z_val == 0 && y_DomainNotEmpty {
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

func (this *C1XmultC2YeqC3ZBounds) xinZinYout(evt *core.ChangeEvent) {
	x_Domain := this.x_Domain
	z_Domain := this.z_Domain
	var chEntry *core.ChangeEntry = nil
	for y_val := range this.y_Domain.Values {
		match := false
		for x_val := range x_Domain.Values {
			if z_Domain.Contains((this.c1 * x_val * this.c2 * y_val) / this.c3) {
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

func (this *C1XmultC2YeqC3ZBounds) yinZinXout(evt *core.ChangeEvent) {
	y_Domain := this.y_Domain
	z_Domain := this.z_Domain
	var chEntry *core.ChangeEntry = nil
	for x_val := range this.x_Domain.Values {
		match := false
		for y_val := range y_Domain.Values {
			if z_Domain.Contains((this.c1 * x_val * this.c2 * y_val) / this.c3) {
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
func (this *C1XmultC2YeqC3ZBounds) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(
			[]core.VarId{this.x, this.y, this.z}, this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]
	this.y_Domain = varidToDomainMap[this.y]
	this.z_Domain = varidToDomainMap[this.z]
	this.store = store
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen
func (this *C1XmultC2YeqC3ZBounds) SetID(propID core.PropId) {
	this.id = propID
}

func (this *C1XmultC2YeqC3ZBounds) GetID() core.PropId {
	return this.id
}

func CreateC1XmultC2YeqC3ZBounds(c1 int, x core.VarId, c2 int,
	y core.VarId, c3 int, z core.VarId) *C1XmultC2YeqC3ZBounds {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXmultYeqC-propagator")
	}
	prop := new(C1XmultC2YeqC3ZBounds)
	prop.x = x
	prop.y = y
	prop.z = z
	prop.c1 = c1
	prop.c2 = c2
	prop.c3 = c3
	return prop
}

func (this *C1XmultC2YeqC3ZBounds) String() string {
	return fmt.Sprintf("PROP_%d %d*%s*%d*%s=%d*%s",
		this.id, this.c1, this.store.GetName(this.x),
		this.c2, this.store.GetName(this.y),
		this.c3, this.store.GetName(this.z))
}

func (this *C1XmultC2YeqC3ZBounds) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y, this.z}
}

func (this *C1XmultC2YeqC3ZBounds) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain, this.z_Domain}
}

func (this *C1XmultC2YeqC3ZBounds) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *C1XmultC2YeqC3ZBounds) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
