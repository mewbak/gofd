package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// implementation of the Roots constraint
// signature: Roots({X1,...,Xi}, S, T) iff S = {i | Xi ∈ T }
// if the index of a variable in xi appears in S,
// it must be assigned a value from T
// (note: does not work with setvariables as it should)
type Roots struct {
	xi                 []core.VarId
	s, t               core.VarId
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	xi_Domains         map[core.VarId]core.Domain
	s_Domain, t_Domain core.Domain
	id                 core.PropId
	store              *core.Store
}

func (this *Roots) Clone() core.Constraint {
	prop := new(Roots)
	prop.xi, prop.s, prop.t = this.xi, this.s, this.t
	return prop
}

func (this *Roots) Start(store *core.Store) {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	this.makeConsistent(evt)
	this.sendChangesToStore(evt)
	for changeEntry := range this.inCh {
		if loggerDebug {
			msg := "%s_'Incoming Change for %s'"
			core.GetLogger().Df(msg, this, store.GetName(changeEntry.GetID()))
		}
		evt = core.CreateChangeEvent()
		switch var_id := changeEntry.GetID(); var_id {
		case this.s:
			this.s_Domain.Removes(changeEntry.GetValues())
			this.makeConsistent(evt)
			break
		case this.t:
			this.t_Domain.Removes(changeEntry.GetValues())
			this.makeConsistent(evt)
			break
		default:
			this.xi_Domains[var_id].Removes(changeEntry.GetValues())
			this.makeConsistent(evt)
			break
		}
		this.sendChangesToStore(evt)
	}
}

// sendChangesToStore send the collected changes (stored in the event)
// to the store
func (this *Roots) sendChangesToStore(evt *core.ChangeEvent) {
	if core.GetLogger().DoDebug() {
		msg := "%s_propagate_'communicate change, evt-value: %s'"
		core.GetLogger().Df(msg, this, evt)
	}
	this.outCh <- evt
}

// makeConsistent checks if the constraint is arc consistent
// and makes it arc consistent if possible
func (this *Roots) makeConsistent(evt *core.ChangeEvent) bool {
	//if any domain is empty, propagation can be stopped
	countChangeEntries := 0

	// remove all values from D(S) which are bigger than the amount of variables
	// or the indexes of the Xi having domains that do not intersect with T
	var chEntry *core.ChangeEntry = nil
	for index := range this.s_Domain.Values_asMap() {
		if index > len(this.xi) {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(this.s)
			}
			chEntry.Add(index)
		} else {
			x_domain := this.xi_Domains[this.xi[index-1]]
			if x_domain.Intersection(this.t_Domain).IsEmpty() {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.s)
				}
				chEntry.Add(index)
			}
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
		countChangeEntries += 1
	}
	for var_id, x_domain := range this.xi_Domains {
		chEntry = nil
		domIntersection := x_domain.Intersection(this.t_Domain)
		if this.s_Domain.Contains(int(var_id) + 1) {
			if !x_domain.IsSubset(this.t_Domain) &&
				(!domIntersection.IsEmpty()) {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(var_id)
				}
				// if the index of a variable X in Xi appears in D(S) and
				// its domain contains values from and outside of D(T),
				// all values that are not part of D(T) are removed from D(X)
				newDomain := this.t_Domain.Difference(x_domain)
				chEntry.AddValues(newDomain)
				if chEntry != nil {
					evt.AddChangeEntry(chEntry)
					countChangeEntries += 1
				}
			}
		} else {
			if !domIntersection.IsEmpty() {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.xi[var_id])
				}
				// if the index of a variable X in Xi does not appear in D(S)
				// and its domain contains values from D(T), all values
				// that are part of D(T) are removed from D(X)
				chEntry.AddValues(domIntersection)
				if chEntry != nil {
					evt.AddChangeEntry(chEntry)
					countChangeEntries += 1
				}

			}
		}
	}

	if countChangeEntries != 0 {
		return false
	}
	return true
}

// Register registers the propagator at the store.
func (this *Roots) Register(store *core.Store) {
	var domains []core.Domain
	idSlice := make([]core.VarId, len(this.xi)+2)
	i := 0
	for i = 0; i < len(this.xi); i++ {
		idSlice[i] = this.xi[i]
	}
	idSlice[i] = this.s
	idSlice[i+1] = this.t
	this.inCh, domains, this.outCh =
		store.RegisterPropagator(idSlice, this.id)

	this.xi_Domains = make(map[core.VarId]core.Domain)
	for i = 0; i < (len(idSlice) - 2); i++ {
		this.xi_Domains[idSlice[i]] = domains[i]
	}
	this.s_Domain = domains[len(idSlice)-2]
	this.t_Domain = domains[len(idSlice)-1]
	this.store = store
}

func (this *Roots) SetID(propID core.PropId) {
	this.id = propID
}

func (this *Roots) GetID() core.PropId {
	return this.id
}

func CreateRoots(xi []core.VarId, s core.VarId, t core.VarId) *Roots {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateRoots_propagator")
	}
	prop := new(Roots)
	prop.xi = xi
	prop.s = s
	prop.t = t
	return prop
}

func (this *Roots) String() string {
	var s string
	for i := 0; i < len(this.xi); i++ {
		s += this.store.GetName(this.xi[i])
	}
	return fmt.Sprintf("PROP_ROOTS({%s}, {%s}, %s)",
		s,
		this.store.GetName(this.s),
		this.store.GetName(this.t))
}

func (this *Roots) GetVarIds() []core.VarId {
	varIds := make([]core.VarId, len(this.xi))
	for _, var_id := range this.xi {
		varIds = append(varIds, var_id)
	}
	varIds = append(varIds, this.s)
	varIds = append(varIds, this.t)
	return varIds
}

func (this *Roots) GetDomains() []core.Domain {
	domains := make([]core.Domain, len(this.xi_Domains))
	for _, var_id := range this.xi {
		domains = append(domains, this.xi_Domains[var_id])
	}
	domains = append(domains, this.s_Domain)
	domains = append(domains, this.t_Domain)
	return domains
}

func (this *Roots) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *Roots) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
