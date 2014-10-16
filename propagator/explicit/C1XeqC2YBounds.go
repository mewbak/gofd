package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// C1XeqC2YBounds represents the propagator for the constraint C1*X = C2*Y
// example: X*3 = Y*1
type C1XeqC2YBounds struct {
	x, y               core.VarId
	c1, c2             int
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	x_Domain, y_Domain *core.ExDomain
	id                 core.PropId
}

func (this *C1XeqC2YBounds) Clone() core.Constraint {
	prop := new(C1XeqC2YBounds)
	prop.x, prop.y = this.x, this.y
	prop.c1, prop.c2 = this.c1, this.c2
	return prop
}

func (this *C1XeqC2YBounds) Start(store *core.Store) {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	if !this.makeBoundConsistent(evt) {
		this.sendChangesToStore(evt)
		return
	}
	this.xinYout(evt)
	this.yinXout(evt)
	this.sendChangesToStore(evt)
	for changeEntry := range this.inCh {
		if loggerDebug {
			msg := "%s_'Incoming Change for %s'"
			core.GetLogger().Df(msg, this, core.GetNameRegistry().GetName(changeEntry.GetID()))
		}
		evt = core.CreateChangeEvent()
		switch var_id := changeEntry.GetID(); var_id {
		case this.x:
			this.x_Domain.Removes(changeEntry.GetValues())
			if !this.makeBoundConsistent(evt) {
				this.sendChangesToStore(evt)
				return
			}
			this.xinYout(evt)
		case this.y:
			this.y_Domain.Removes(changeEntry.GetValues())
			if !this.makeBoundConsistent(evt) {
				this.sendChangesToStore(evt)
				return
			}
			this.yinXout(evt)
		}
		this.sendChangesToStore(evt)
	}
}

// sendChangesToStore send the collected changes (stored in the event)
// to the store
func (this *C1XeqC2YBounds) sendChangesToStore(evt *core.ChangeEvent) {
	if core.GetLogger().DoDebug() {
		msg := "%s_propagate_'communicate change, evt-value: %s'"
		core.GetLogger().Df(msg, this, evt)
	}
	this.outCh <- evt
}

// makeBoundConsistent checks if the constraint is bound consistent
// and makes it bound consistent if possible
func (this *C1XeqC2YBounds) makeBoundConsistent(evt *core.ChangeEvent) bool {
	if this.x_Domain.IsEmpty() || this.y_Domain.IsEmpty() {
		return false
	}
	xmin, xmax := this.x_Domain.Min, this.x_Domain.Max
	ymin, ymax := this.y_Domain.Min, this.y_Domain.Max
	possmax := this.c1 * xmax
	possmin := this.c1 * xmin
	var chEntry *core.ChangeEntry = nil
	if possmax < this.c2*ymax || possmin > this.c2*ymin {
		for k := range this.y_Domain.Values {
			// delete all values in Y that are bigger than possmax
			if k*this.c2 > possmax {
				this.y_Domain.Remove(k) // BUG: ToDo: mustn't change
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.y)
				}
				chEntry.Add(k)
			}
			// delete all values in Y that are smaller than possmin
			if k*this.c2 < possmin {
				this.y_Domain.Remove(k)
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.y)
				}
				chEntry.Add(k)
			}
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
	// if there are still elements in the domain of Y,
	// it should be bound consistent by now
	if !this.y_Domain.IsEmpty() {
		return true
	}
	return false
}

func (this *C1XeqC2YBounds) xinYout(evt *core.ChangeEvent) {
	// collect changes X=Z/3
	x_Domain := this.x_Domain
	var chEntry *core.ChangeEntry = nil
	for y_val := range this.y_Domain.Values {
		if this.c1 == 0 {
			if y_val != 0 && this.c2 != 0 {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.y)
				}
				chEntry.Add(y_val)
			}
		} else {
			if y_val*this.c2%this.c1 != 0 ||
				!x_Domain.Contains(y_val*this.c2/this.c1) {
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

func (this *C1XeqC2YBounds) yinXout(evt *core.ChangeEvent) {
	x_Domain := this.x_Domain
	y_Domain := this.y_Domain
	var chEntry *core.ChangeEntry = nil
	for x_val := range x_Domain.Values {
		if this.c2 == 0 {
			if x_val != 0 && this.c1 != 0 {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.x)
				}
				chEntry.Add(x_val)
			}
		} else {
			if x_val*this.c1%this.c2 != 0 ||
				!y_Domain.Contains(x_val*this.c1/this.c2) {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.x)
				}
				chEntry.Add(x_val)
			}
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}
}

// Register registers the propagator at the store.
func (this *C1XeqC2YBounds) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y}, this.id)
	varidToDomainMap := core.GetVaridToExplicitDomainsMap(domains)
	this.x_Domain = varidToDomainMap[this.x]
	this.y_Domain = varidToDomainMap[this.y]

}

func (this *C1XeqC2YBounds) SetID(propID core.PropId) {
	this.id = propID
}

func (this *C1XeqC2YBounds) GetID() core.PropId {
	return this.id
}

//X*C=Z
func CreateC1XeqC2YBounds(c1 int, x core.VarId, c2 int, y core.VarId) *C1XeqC2YBounds {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateCXeqCYBounds_propagator")
	}
	prop := new(C1XeqC2YBounds)
	prop.x = x
	prop.y = y
	prop.c1 = c1
	prop.c2 = c2
	return prop
}

func CreateCXeqYBounds(c1 int, x core.VarId, y core.VarId) *C1XeqC2YBounds {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateCXeqYBounds_propagator")
	}
	prop := new(C1XeqC2YBounds)
	prop.x = x
	prop.y = y
	prop.c1 = c1
	prop.c2 = 1
	return prop
}

func CreateXeqYBounds(x core.VarId, y core.VarId) *C1XeqC2YBounds {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXeqYBounds_propagator")
	}
	prop := new(C1XeqC2YBounds)
	prop.x = x
	prop.y = y
	prop.c1 = 1
	prop.c2 = 1
	return prop
}

func (this *C1XeqC2YBounds) String() string {
	return fmt.Sprintf("PROP_%d %d*%s = %d*%s",
		this.id, this.c1, core.GetNameRegistry().GetName(this.x),
		this.c2, core.GetNameRegistry().GetName(this.y))
}

func (this *C1XeqC2YBounds) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y}
}

func (this *C1XeqC2YBounds) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain}
}

func (this *C1XeqC2YBounds) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *C1XeqC2YBounds) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
